package BaseSQLTools

import (
	"errors"
	"reflect"
)

// QuickInsert 插入模块
type QuickInsert struct {
	//Quick
	quickClient *Quick
}

// InsertRow 快速写入单行
func (c *QuickInsert) InsertRow(args any) (newID int64, err error) {
	//获取当前结构体
	fieldList := c.quickClient.getFields()
	//获取写入参数
	paramsType := reflect.TypeOf(args).Elem()
	//valueType := reflect.ValueOf(args).Elem()
	step := 0
	for step < paramsType.NumField() {
		//捕捉结构
		vField := paramsType.Field(step)
		vDBVal := vField.Tag.Get("db")
		//下一步
		step += 1
		//检查参数是否存在
		isFind := false
		for _, v := range fieldList {
			if v == vDBVal {
				isFind = true
				break
			}
		}
		if !isFind {
			err = errors.New("no support field: " + vField.Tag.Get("db"))
			return
		}
	}
	//写入准备
	ctx := c.quickClient.client.Insert().SetDefaultInsertFields().Add(args)
	//写入执行
	newID, err = ctx.ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}
