package BaseSQLTools

import "errors"

// QuickInfo 获取信息
type QuickInfo struct {
	//Quick
	quickClient *Quick
}

// GetInfoByID 获取指定ID的信息
func (c *QuickInfo) GetInfoByID(id int64, result any) (err error) {
	//检查ID
	if id < 1 {
		return errors.New("id error")
	}
	//获取缓冲
	if err = c.quickClient.GetCacheInfoByID(id, result); err == nil {
		return
	}
	//获取数据
	err = c.quickClient.client.Get().SetDefaultFields().GetByID(id).Result(result)
	if err != nil {
		return
	}
	//保存缓冲
	c.quickClient.SetCacheInfoByID(id, result)
	//反馈
	return
}

// GetInfoByField 通过指定字段获取信息（必须是唯一的字段）
// 注意，当字段存在软删除时，请务必启用haveDelete，否则将出现异常
func (c *QuickInfo) GetInfoByField(fieldName string, fieldVal any, haveDelete bool, result any) (err error) {
	//获取数据
	ctx := c.quickClient.client.Get().SetDefaultFields()
	if haveDelete {
		ctx = ctx.SetDeleteQuery("delete_at", false)
	}
	switch fieldVal.(type) {
	case int:
		ctx = ctx.SetIntQuery(fieldName, fieldVal.(int))
	case int64:
		ctx = ctx.SetInt64Query(fieldName, fieldVal.(int64))
	case string:
		ctx = ctx.SetStringQuery(fieldName, fieldVal.(string))
	default:
		return errors.New("field type error")
	}
	err = ctx.Result(result)
	if err != nil {
		return
	}
	//反馈
	return
}
