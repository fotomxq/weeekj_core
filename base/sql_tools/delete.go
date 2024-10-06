package SQLTools

import "errors"

// QuickDelete 删除信息
type QuickDelete struct {
	//Quick
	quickClient *Quick
}

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
