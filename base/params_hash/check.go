package BaseParamsHash

import (
	"encoding/json"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// Check 检查参数是否重复提交
func Check(tokenID int64, userID int64, params any) (b bool) {
	paramsRaw, _ := json.Marshal(params)
	hash := CoreFilter.GetSha1Str(string(paramsRaw))
	return CheckHash(tokenID, userID, hash)
}

// CheckHash 检查指定的hash是否重复
func CheckHash(tokenID int64, userID int64, hash string) (b bool) {
	hashP := CoreFilter.GetSha1Str(hash)
	mark := fmt.Sprint("base.params.hash.", tokenID, ".", userID, ".", hashP)
	var rawData DataHash
	if err := Router2SystemConfig.MainCache.GetStruct(mark, &rawData); err == nil && rawData.Hash != "" {
		return
	}
	rawData = DataHash{
		Hash:    hashP,
		TokenID: tokenID,
		UserID:  userID,
	}
	Router2SystemConfig.MainCache.SetStruct(mark, rawData, 60)
	b = true
	return
}
