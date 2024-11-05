package BaseSQLTools

import "errors"

// QuickDelete 删除信息
type QuickDelete struct {
	//Quick
	quickClient *Quick
}

// DeleteByID 根据ID删除
func (c *QuickDelete) DeleteByID(id int64) (err error) {
	//检查ID
	if id < 1 {
		err = errors.New("id error")
		return
	}
	//执行删除
	err = c.quickClient.client.Delete().NeedSoft(c.quickClient.openSoftDelete).AddWhereID(id).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	c.quickClient.DeleteCacheByID(id)
	//返回
	return
}

// DeleteByField 根据字段删除
func (c *QuickDelete) DeleteByField(fieldName string, val any) (err error) {
	//执行删除
	err = c.quickClient.client.Delete().NeedSoft(c.quickClient.openSoftDelete).AddQuery(fieldName, val).ExecNamed(nil)
	if err != nil {
		return
	}
	//返回
	return
}