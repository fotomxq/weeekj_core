package BaseDistribution

import (
	"context"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	RPCXClient "github.com/smallnest/rpcx/client"
	"sync"
	"time"
)

// 链接持数据集合
type runTestAvgDataType struct {
	//基础配对信息
	// 该信息结构将用于匹配链接
	Host string
	//链接基础
	Service RPCXClient.ServiceDiscovery
	Client  RPCXClient.XClient
	//创建时间
	CreateTime int64
	//最后一次使用时间
	// 用于判定该链接是否有效
	LastTime int64
}

// 本服务将建立维护持久化链接持，一般情况下不会反复重新链接
var (
	runTestAvgClientLock sync.Mutex
	runTestAvgClients    []runTestAvgDataType
)

// 子服务，自动连接默认方法，测试效率。如果方法不可用，则自动按照-1秒超延迟记录
func runTestAvg() error {
	//获取列表
	serviceList := cacheData
	if len(serviceList) < 1 {
		//关闭上下文
		return nil
	}
	//遍历服务
	for _, vService := range serviceList {
		//如果方法为空则不需要测试
		if vService.DefaultAction == "" {
			continue
		}
		//查询所有负载
		vChildes, err := GetChildAll(&ArgsGetChildAll{
			Mark: vService.Mark,
		})
		if err != nil || len(vChildes) < 1 {
			continue
		}
		//遍历负载
		for _, vChild := range vChildes {
			//已经过期的不进行负载处理
			if vChild.ExpireTime < CoreFilter.GetNowTime().Unix() {
				continue
			}
			//连接服务
			clientHost, clientService, client, err := clientConnectByData("tcp", vChild.ServerIP, vChild.ServerPort, vService.DefaultAction)
			if err != nil {
				//标记为连接失败
				runSetAddCount(vChild, -1)
				//跳过
				continue
			}
			//上下文设置
			ctx3, cancel3 := context.WithTimeout(context.Background(), time.Second*60)
			//测试方法连接时间
			nowTime := CoreFilter.GetNowTime().UnixNano()
			if err := clientCall(ctx3, client, vService.DefaultFunc, nil, nil); err != nil {
				CoreLog.Error("run test client, call failed to mark: ", vService.Mark, ", err: ", err)
				//关闭上下文
				cancel3()
				//呼叫失败时，关闭链接
				if err := client.Close(); err != nil {
					CoreLog.Error("run test close client, ", err)
				}
				clientService.Close()
				//将该链接从共享库清除
				clearRunTestAvgClient(clientHost)
				//跳过
				continue
			}
			//递交新的时间节点
			lastTime := CoreFilter.GetNowTime().UnixNano()
			runSetAddCount(vChild, lastTime-nowTime)
			//关闭上下文
			cancel3()
		}
	}
	return nil
}

// 负载测试失败处理方法
func runSetAddCount(data FieldsDistributionChild, addCount int64) {
	var newLastCount []int64
	newLastCount = append(newLastCount, addCount)
	max := 10
	var lastCount int64 = 0
	for _, v := range data.LastCount {
		if len(newLastCount) >= max {
			break
		}
		newLastCount = append(newLastCount, v)
		lastCount += v
	}
	data.LastCount = newLastCount
	data.LastCountAvg = lastCount / int64(len(data.LastCount))
	saveChildCache(data)
}

// 连接某个负载
// 必须开放，部分场景可能需要连接所有负载进行全局通知
func clientConnectByData(serverType, ip, port, serverAction string) (string, RPCXClient.ServiceDiscovery, RPCXClient.XClient, error) {
	//锁机制，避免一个时间重建多个重复链接
	runTestAvgClientLock.Lock()
	defer runTestAvgClientLock.Unlock()
	//建立基本握手信息
	serverHost := serverType + "@" + ip + ":" + port
	//查询该连接是否已经存在？
	// 默认情况下，采用了自定义配置项，超时时间为10秒，如果超时也会反馈连接，但会造成短暂的失效异常
	// 如果发生异常，应强制关闭链接
	for k, v := range runTestAvgClients {
		if v.Host == serverHost {
			//找到数据开始更新
			runTestAvgClients[k].LastTime = CoreFilter.GetNowTime().Unix()
			//反馈数据
			return v.Host, v.Service, v.Client, nil
		}
	}
	//如果没有找到，则建立链接
	serverDiscovery := RPCXClient.NewPeer2PeerDiscovery(serverHost, "")
	client := RPCXClient.NewXClient(serverAction, RPCXClient.Failtry, RPCXClient.RandomSelect, serverDiscovery, RPCXClient.DefaultOption)
	//将链接信息放入共享库
	isFind := false
	for k, v := range runTestAvgClients {
		if v.Host == serverHost {
			isFind = true
			runTestAvgClients[k] = runTestAvgDataType{
				Host:       serverHost,
				Service:    serverDiscovery,
				Client:     client,
				CreateTime: CoreFilter.GetNowTime().Unix(),
				LastTime:   CoreFilter.GetNowTime().Unix(),
			}
			break
		}
	}
	if !isFind {
		runTestAvgClients = append(runTestAvgClients, runTestAvgDataType{
			Host:       serverHost,
			Service:    serverDiscovery,
			Client:     client,
			CreateTime: CoreFilter.GetNowTime().Unix(),
			LastTime:   CoreFilter.GetNowTime().Unix(),
		})
	}
	//反馈信息
	return serverHost, serverDiscovery, client, nil
}

// 执行某个方法
func clientCall(ctx context.Context, client RPCXClient.XClient, serviceMethod string, args interface{}, reply interface{}) (err error) {
	//触发器
	err = client.Call(ctx, serviceMethod, args, reply)
	return
}

// 清除某个链接
func clearRunTestAvgClient(host string) {
	runTestAvgClientLock.Lock()
	defer runTestAvgClientLock.Unlock()
	var newData []runTestAvgDataType
	for _, v := range runTestAvgClients {
		if v.Host == host {
			continue
		}
		newData = append(newData, v)
	}
	runTestAvgClients = newData
}
