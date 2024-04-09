package BaseService

// 永远不要调用此方法
// 该方法会再coreNats模块中复现
func pushRequest(serviceCode string, mark string) {
	//CoreNats.PushDataNoErr("/base/service/request", serviceCode, 0, mark, nil)
}
