package OrgSubscription

//组织订阅处理
/**
1、用户在本模块内发起请求，本模块将触发新的订单；
2、用户完成订单支付后，本系统将巡逻发现完成状态；
3、如果完成，将为用户指定组织开通指定的服务项目。
4、【后续】检查组织订阅状态，如果没有续约，将自动剥离组织的服务项目。如果达到N天没有续约，将标记删除组织。
*/

var (
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() {
	if OpenSub {
		//消息列队
		subNats()
	}
}
