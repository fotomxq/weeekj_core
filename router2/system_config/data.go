package Router2SystemConfig

// GlobConfigType 配置数据结构
type GlobConfigType struct {
	//路由
	Router GlobConfigRouterType
	//第三方平台
	OtherAPI GLobConfigOtherAPI
	//用户
	User GlobConfigUserType
	//财务
	Finance GlobConfigFinanceType
	//安全
	Safe GlobConfigSafeType
	//组织
	Org GlobConfigOrg
	//ERP
	ERP GlobConfigERP
	//地图
	Map GlobConfigMap
	//服务
	Service GlobConfigService
}

type GlobConfigRouterType struct {
	//是否需要会话日志
	NeedTokenLog bool
}

type GLobConfigOtherAPI struct {
	//是否启动同步天气预报
	OpenSyncWeather bool
	//是否启动同步假期
	OpenSyncHolidaySeason bool
}

type GlobConfigUserType struct {
	//LoginViewPhone 登录用户是否可以看到手机号
	LoginViewPhone bool
	//LoginViewEmail 登录用户是否可以看到email
	LoginViewEmail bool
	//SyncUserPhoneUsername 修改绑定手机号后同步修改用户登陆名
	SyncUserPhoneUsername bool
	//前端是否可以消费用户积分
	ClientCostIntegral bool
	//允许公开用户数据
	GlobShowUser bool
}

type GlobConfigFinanceType struct {
	//NeedNoCheck 是否不校验
	NeedNoCheck bool
}

type GlobConfigSafeType struct {
	// SafeRouterTimeBlocker 是否启动中间人攻击拦击拦截机制
	SafeRouterTimeBlocker bool
}

type GlobConfigOrg struct {
	//默认开通的功能
	DefaultOpenFunc []string
}

type GlobConfigERP struct {
	//仓储
	Warehouse GlobConfigERPWarehouse
}

type GlobConfigERPWarehouse struct {
	//库存是否可以为负数
	StoreLess0 bool
}

type GlobConfigMap struct {
	//是否打开地图
	MapOpen bool
}

type GlobConfigService struct {
	//信息交互订单完成后自动下架交互内容
	InfoExchangeOrderFinishAutoDown bool
}
