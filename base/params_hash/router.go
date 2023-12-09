package BaseParamsHash

import Router2Mid "gitee.com/weeekj/weeekj_core/v5/router2/mid"

// CheckParamsHash 检查随机掩码
func CheckParamsHash(c any, hash string) (b bool) {
	if hash == "" {
		b = true
		return
	}
	ctx := Router2Mid.GetContext(c)
	tokenID := Router2Mid.GetTokenID(ctx)
	userID, _ := Router2Mid.TryGetUserID(ctx)
	b = CheckHash(tokenID, userID, hash)
	if !b {
		Router2Mid.ReportBaseError(c, "err_params_hash")
		return
	}
	return
}
