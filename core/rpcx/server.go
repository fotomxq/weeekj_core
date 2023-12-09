package CoreRPCX

import (
	BaseDistribution "gitee.com/weeekj/weeekj_core/v5/base/distribution"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	PRCXServer "github.com/smallnest/rpcx/server"
	"time"
)

//本模块用于微服务构建和初始化

// ArgsInitService 服务基础信息
type ArgsInitService struct {
	ServiceMark    string
	ServiceName    string
	ExpireInterval int64
	DefaultAction  string
	DefaultFunc    string
}

// ArgsInitServer 负载信息
type ArgsInitServer struct {
	ServerName string
	ServerIP   string
	ServerPort string
}

type Server struct {
	//参数
	argsService ArgsInitService
	argsServer  ArgsInitServer
	//心跳周期
	runTime time.Duration
	//微服务组件
	Server *PRCXServer.Server
	//本服务广播地址
	host string
}

// Init 初始化设置
func (t *Server) Init(argsService ArgsInitService, argsServer ArgsInitServer, host string) error {
	//设置服务信息
	t.argsService = argsService
	//设置负载信息
	t.argsServer = argsServer
	//初始化配置模块
	Router2SystemConfig.ServerName = t.argsServer.ServerName
	Router2SystemConfig.ServerIP = t.argsServer.ServerIP
	Router2SystemConfig.ServerPort = t.argsServer.ServerPort
	if err := Router2SystemConfig.Init(); err != nil {
		return err
	}
	//创建rpcx
	t.Server = PRCXServer.NewServer()
	//设置巡逻时间
	t.runTime = time.Second * 10
	//设置广播地址
	t.host = host
	//反馈
	return nil
}

// Run 心跳处理和负载情况上报
func (t *Server) Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("rpcx server run, ", r)
		}
	}()
	//基础日志处理模块，必备
	go CoreLog.Run()
	//维护主体
	for {
		var err error
		//设置主服务
		err = BaseDistribution.SetService(&BaseDistribution.ArgsSetService{
			Mark: t.argsService.ServiceMark, Name: t.argsService.ServiceName, ExpireInterval: t.argsService.ExpireInterval, DefaultAction: t.argsService.DefaultAction, DefaultFunc: t.argsService.DefaultFunc,
		})
		if err != nil {
			CoreLog.Error("rpcx server run, distribution set service, service args: ", t.argsService, ", err: ", err)
		} else {
			//设置负载服务
			err = BaseDistribution.SetChild(&BaseDistribution.ArgsSetChild{
				Mark: t.argsService.ServiceMark, Name: t.argsServer.ServerName, IP: t.argsServer.ServerIP, Port: t.argsServer.ServerPort,
			})
			if err != nil {
				CoreLog.Error("rpcx server run error, distribution set service child, server args: ", t.argsService, ", err: ", err)
			}
		}
		//10秒发送1次心跳和负载状态
		time.Sleep(t.runTime)
	}
}

// 简化run服务包装
func (t *Server) SetChildRun(runMark string, expireAddTime int64) {
	if err := BaseDistribution.SetChildRun(&BaseDistribution.ArgsSetChildRun{
		Mark: t.argsService.ServiceMark, IP: t.argsServer.ServerIP, Port: t.argsServer.ServerPort, RunMark: runMark, ExpireAddTime: expireAddTime,
	}); err != nil {
		//设置错误
		CoreLog.Error("distribution set service child run, server args: ", t.argsServer, ", err: ", err)
	}
}

// 启动tcp广播
func (t *Server) RunServer() error {
	if err := t.Server.Serve("tcp", t.host); err != nil {
		return err
	}
	return nil
}
