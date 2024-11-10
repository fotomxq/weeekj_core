package BaseSQLTools

import (
	"errors"
	"time"
)

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
	err = c.quickClient.client.Get().NeedLimit().SetDefaultFields().GetByID(id).Result(result)
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
	ctx := c.quickClient.client.Get().NeedLimit().SetDefaultFields()
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

// GetInfoByFields 多个条件获取数据
func (c *QuickInfo) GetInfoByFields(fields map[string]any, haveDelete bool, result any) (err error) {
	//获取数据
	ctx := c.quickClient.client.Get().NeedLimit().SetDefaultFields()
	if haveDelete {
		ctx = ctx.SetDeleteQuery("delete_at", false)
	}
	for fieldName, fieldVal := range fields {
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
	}
	err = ctx.Result(result)
	if err != nil {
		return
	}
	//反馈
	return
}

// CheckInfoByFields 检查多条件是否存在数据
func (c *QuickInfo) CheckInfoByFields(fields map[string]any, haveDelete bool) (b bool, err error) {
	//获取数据
	ctx := c.quickClient.client.Get().SetFieldsOne([]string{"id"})
	if haveDelete {
		ctx = ctx.SetDeleteQuery("delete_at", false)
	}
	for fieldName, fieldVal := range fields {
		switch fieldVal.(type) {
		case int:
			ctx = ctx.SetIntQuery(fieldName, fieldVal.(int))
		case int64:
			ctx = ctx.SetInt64Query(fieldName, fieldVal.(int64))
		case string:
			ctx = ctx.SetStringQuery(fieldName, fieldVal.(string))
		default:
			err = errors.New("field type error")
			return
		}
	}
	//解析结构体
	var id int64
	err = ctx.Result(&id)
	if err != nil {
		return
	}
	if id < 1 {
		return false, nil
	}
	b = true
	//反馈
	return
}

// CheckInfoByFieldsAndTime 检查多条件是否存在数据，且判断时间是否在有效期内
func (c *QuickInfo) CheckInfoByFieldsAndTime(fields map[string]any, haveDelete bool, timeField string, timeMin, timeMax time.Time) (b bool, err error) {
	//获取数据
	ctx := c.quickClient.client.Get().SetFieldsOne([]string{"id"})
	if haveDelete {
		ctx = ctx.SetDeleteQuery("delete_at", false)
	}
	for fieldName, fieldVal := range fields {
		switch fieldVal.(type) {
		case int:
			ctx = ctx.SetIntQuery(fieldName, fieldVal.(int))
		case int64:
			ctx = ctx.SetInt64Query(fieldName, fieldVal.(int64))
		case string:
			ctx = ctx.SetStringQuery(fieldName, fieldVal.(string))
		default:
			err = errors.New("field type error")
			return
		}
	}
	//检查时间范围
	if timeMin != (time.Time{}) {
		ctx = ctx.SetTimeMinQuery(timeField, timeMin)
	}
	if timeMax != (time.Time{}) {
		ctx = ctx.SetTimeMaxQuery(timeField, timeMax)
	}
	//解析结构体
	var id int64
	err = ctx.Result(&id)
	if err != nil {
		return
	}
	if id < 1 {
		return false, nil
	}
	b = true
	//反馈
	return
}
