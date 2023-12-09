package BaseParamsHash

type DataHash struct {
	//随机掩码
	Hash string `json:"hash"`
	//token
	TokenID int64 `json:"tokenID"`
	//用户ID
	UserID int64 `json:"userID"`
}
