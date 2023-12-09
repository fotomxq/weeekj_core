package CoreNats

// dataType 通用消息结构体
type dataType struct {
	//行为
	Action string `json:"action"`
	//影响ID
	ID int64 `json:"id"`
	//影响Mark
	Mark string `json:"mark"`
	//数据包
	Data interface{} `json:"data"`
}
