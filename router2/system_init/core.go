package Router2SystemInit

//系统初始化模块组
/**
需调用该模块实现对应核心功能组件，独立服务后续对接额外的内容即可
*/

var (
	//AppName 应用基本信息
	AppName = "weeekj"
	//AppVersion 应用版本号
	AppVersion = "2.0.0"
	//OpenSub 是否启动通用框架体系的订阅消息
	OpenSub = false
	//OpenAnalysis 是否启动通用框架体系的统计支持
	OpenAnalysis = false
	//OpenAPIDefaultSet 是否加载API默认静态路由
	OpenAPIDefaultSet = true
	//OpenSystemReg 是否启动注册机机制
	OpenSystemReg = true
)
