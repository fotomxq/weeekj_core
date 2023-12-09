package BaseWeixinPayNotify

// 收到退款和支付通知后返回给微信服务器的消息
type replay struct {
	Code string `xml:"return_code"` // 返回状态码: SUCCESS/FAIL
	Msg  string `xml:"return_msg"`  // 返回信息: 返回信息，如非空，为错误原因
}

// 根据结果创建返回数据
//
// ok 是否处理成功
// msg 处理不成功原因
func newReplay(ok bool, msg string) replay {
	ret := replay{Msg: msg}
	if ok {
		ret.Code = "SUCCESS"
	} else {
		ret.Code = "FAIL"
	}
	return ret
}
