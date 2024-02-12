package BaseToken2

import CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"

// GetTokenByS 获取对应值的会话
func GetTokenByS(val string) (data FieldsToken) {
	var sData FieldsTokenS
	_ = baseTokenS.Get().SetFieldsOne([]string{"id", "create_at", "token_id", "val"}).AppendWhere("val = $1", val).NeedLimit().Result(&sData)
	if sData.ID < 1 {
		return
	}
	data = GetByID(sData.TokenID)
	return
}

// CheckTokenS 检查是否合法
func CheckTokenS(val string) (b bool) {
	data := GetTokenByS(val)
	b = data.ID > 0
	return
}

// CreateTokenS 生成一个短会话地址
func CreateTokenS(tokenID int64) (val string) {
	//生成短地址
	valRand := CoreFilter.GetRandStr4(50)
	if valRand == "" {
		return
	}
	val = CoreFilter.GetSha1Str(valRand)
	if val == "" {
		return
	}
	//写入短地址
	err := baseTokenS.Insert().SetFields([]string{"token_id", "val"}).Add(map[string]interface{}{
		"token_id": tokenID,
		"val":      val,
	}).ExecAndCheckID()
	if err != nil {
		val = ""
		return
	}
	//反馈
	return
}

// DeleteTokenSByTokenID 删除tokenID对应的所有短地址
func DeleteTokenSByTokenID(tokenID int64) (err error) {
	err = baseTokenS.Delete().NeedSoft(false).SetWhereAnd("token_id", tokenID).ExecNamed(map[string]interface{}{
		"token_id": tokenID,
	})
	return
}
