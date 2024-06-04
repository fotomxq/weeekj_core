package ClassSort2

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
)

// Sort 对象结构
type Sort struct {
	//排序表名
	SortTableName string
	//数据库句柄
	sortDB CoreSQL2.Client
}

func (t *Sort) Init(mainDB *CoreSQL2.SQLClient) (err error) {
	_, err = t.sortDB.Init2(mainDB, t.SortTableName, &FieldsSort{})
	if err != nil {
		return
	}
	return
}
