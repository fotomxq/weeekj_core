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
	if err = c.quickClient.GetCacheInfoByID(id, result); err != nil {
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
func (c *QuickInfo) GetInfoByField(fieldName string, fieldVal any, result any) (err error) {
	//获取数据
	switch fieldVal.(type) {
	case int:
		err = c.quickClient.client.Get().SetDefaultFields().SetIntQuery(fieldName, fieldVal.(int)).Result(result)
	case int64:
		err = c.quickClient.client.Get().SetDefaultFields().SetInt64Query(fieldName, fieldVal.(int64)).Result(result)
	case string:
		err = c.quickClient.client.Get().SetDefaultFields().SetStringQuery(fieldName, fieldVal.(string)).Result(result)
	default:
		return errors.New("field type error")
	}
	if err != nil {
		return
	}
	//反馈
	return
}
