package BaseWeixinWXXClient

const (
	// WeChatServerError 微信服务器错误时返回返回消息
	WeChatServerError = "微信服务器发生错误"
)

// Response 请求微信返回基础数据
type ResponseBase struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}
