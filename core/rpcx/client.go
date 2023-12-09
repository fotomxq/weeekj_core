package CoreRPCX

import (
	"context"
	"errors"
	BaseDistribution "gitee.com/weeekj/weeekj_core/v5/base/distribution"
	RPCXClient "github.com/smallnest/rpcx/client"
	"sync"
	"time"
)

//本机微服务套件
// * 本服务禁止BaseDistribution调用

type Client struct {
	//连接方法
	serverType string
	//分布式服务地址
	serverHost string
	//RPCX Client
	client RPCXClient.XClient
	//PRCX Discovery
	serverDiscovery RPCXClient.ServiceDiscovery
	//服务标识码
	serverMark string
	//服务集
	serverAction string
	//重新连接阻塞机制
	connectLock sync.Mutex
	//是否强制总是重试
	allowRetryConnect bool
	//强制重连等待时间 秒
	retryConnectTime time.Duration
	//强制重试最大上限
	// 如果为0则不限制
	retryCountMax int
	//强制重试几次 记录重试了第几次
	retryCount int
	//强制指定一个分布式服务
	forceDistribution BaseDistribution.FieldsDistributionChild
	//是否初始化连接
	isConnectStart bool
}

// 初始化
func (t *Client) Init(mark, action string, allowRetryConnect bool) {
	t.SetServerType("tcp")
	t.SetMark(mark, action)
	t.allowRetryConnect = allowRetryConnect
	t.retryConnectTime = 10
	t.retryCountMax = 0
	t.isConnectStart = false
}

// 设置mark
func (t *Client) SetMark(mark, action string) {
	t.serverMark = mark
	t.serverAction = action
}

// 设置方法
func (t *Client) SetServerType(serverType string) {
	t.serverType = serverType
}

// 设置强制重连机制
func (t *Client) SetRetryConnect(allowRetryConnect bool, retryConnectTime int64, retryCountMax int) {
	t.allowRetryConnect = allowRetryConnect
	t.retryConnectTime = time.Second * time.Duration(retryConnectTime)
	t.retryCountMax = retryCountMax
}

// 强制约束到一个分布式去连接处理
func (t *Client) SetForceDistribution(data BaseDistribution.FieldsDistributionChild) {
	t.forceDistribution = data
}

// 清理掉强制约束负载
func (t *Client) ClearForceDistribution() {
	t.forceDistribution = BaseDistribution.FieldsDistributionChild{}
}

// 封装call处理
// 封装后自动进行可用治理，发生不可用后将自动重连
func (t *Client) Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error {
	//检查是否已经连接，如果未连接则自动初始化
	if !t.isConnectStart {
		if err := t.connect(); err != nil {
			return errors.New("connect service failed, " + err.Error())
		}
		t.isConnectStart = true
	}
	//触发器
	err := t.client.Call(ctx, serviceMethod, args, reply)
	if err == nil {
		return nil
	}
	//服务掉线
	// XClient的二参数设置，说明已经重试了3次
	// 可能需要更换服务地址
	if err == RPCXClient.ErrXClientShutdown {
		//无限循环处理
		//直到触发特定条件，自动跳出
		for {
			//如果重试超出最大限制条件，则跳出
			if t.retryCountMax >= 1 && t.retryCount >= t.retryCountMax {
				break
			}
			if t.retryCount > 0 && !t.allowRetryConnect {
				break
			}
			if err2 := t.retryConnect(); err2 != nil {
				//重新连接失败，说明服务掉线
				// 造成级联错误
				//重试次数递增1
				t.retryCount += 1
				//否则等待N时间后重试
				time.Sleep(t.retryConnectTime)
				continue
			}
			//连接成功后，重置计数器
			t.retryCount = 0
			//再次触发call
			err := t.client.Call(ctx, serviceMethod, args, reply)
			if err == nil {
				return nil
			}
			//如果还是失败，则继续执行机制
			if err == RPCXClient.ErrXClientShutdown {
				t.retryCount += 1
			}
			// 注意，此处不需要递增计数器，因为错误不一定是该错误问题
			// 等待时间N后重试
			time.Sleep(t.retryConnectTime)
			//continue
		}
	}
	return err
}

// 发现可用的微服务并连接
func (t *Client) connect() error {
	t.connectLock.Lock()
	//初始化负载
	var data BaseDistribution.FieldsDistributionChild
	var err error
	//如果不存在约束负载，则自动获取一个
	if t.forceDistribution.Mark == "" {
		data, err = BaseDistribution.GetBalancing(&BaseDistribution.ArgsGetBalancing{
			Mark: t.serverMark,
		})
		if err != nil {
			t.connectLock.Unlock()
			return err
		}
	} else {
		//否则强制约束设置
		data = t.forceDistribution
	}
	t.serverHost, t.serverDiscovery, t.client, err = t.connectByData(&data, t.serverAction)
	t.connectLock.Unlock()
	return nil
}

// 连接某个负载
// 必须开放，部分场景可能需要连接所有负载进行全局通知
func (t *Client) connectByData(data *BaseDistribution.FieldsDistributionChild, serverAction string) (string, RPCXClient.ServiceDiscovery, RPCXClient.XClient, error) {
	serverHost := t.serverType + "@" + data.ServerIP + ":" + data.ServerPort
	serverDiscovery := RPCXClient.NewPeer2PeerDiscovery(serverHost, "")
	client := RPCXClient.NewXClient(serverAction, RPCXClient.Failtry, RPCXClient.RandomSelect, serverDiscovery, RPCXClient.DefaultOption)
	return serverHost, serverDiscovery, client, nil
}

// 自动重新连接服务
func (t *Client) retryConnect() error {
	_ = t.client.Close()
	if err := t.connect(); err != nil {
		return err
	}
	return nil
}

// 强制关闭连接
func (t *Client) Close() error {
	return t.client.Close()
}
