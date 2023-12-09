package BaseWeixinPayClient

//ResponseBase 请求微信返回基础数据
type ResponseBase struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}
