package BaseDistribution

type DataType struct {
	Service  []FieldsDistribution
	Child    []FieldsDistributionChild
	ChildRun []FieldsDistributionChildRun
}

//顶层mongodb配置
type DataGlobMongodb struct {
	//数据库连接地址
	URL string
	//数据库名称
	DBName string
}

//顶层postgres配置
type DataGlobPostgres struct {
	//数据库连接地址
	// 地址同时包含其他内容
	URL string
}

//通用配置设计
type DataServiceConfig struct {
	//服务器广播信息
	ServerName string `json:"serverName"`
	ServerMark string `json:"serverMark"`
	//debug模式
	Debug bool `json:"debug"`
	//是否需要全局debug覆盖
	NeedGlobDebug bool `json:"needGlobDebug"`
}

//配置和附带内部设计
type DataService struct {
	//配置
	Config DataServiceConfig
	//过期时间
	ExpireTime int64
}
