package CoreNats

// 方法复现于base/service/nats.go
// action 服务code
// mark 订阅和推送类型: sub订阅; pub发布
func pushRequest(serviceCode string, mark string) {
	PushDataNoErr("base_service_request", "/base/service/request", serviceCode, 0, mark, nil)
}
