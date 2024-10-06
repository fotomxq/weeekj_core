package SQLTools

import (
	"errors"
	"reflect"
)

// QuickUpdate 更新模块
type QuickUpdate struct {
	//Quick
	quickClient *Quick
}

// UpdateByID 根据ID更新数据
func (c *QuickUpdate) UpdateByID(args any) (err error) {
	//获取当前结构体
	fieldList := c.quickClient.getFields()
	//获取写入参数
	var argID int64
	var setFields []string
	var setVals map[string]any
	paramsType := reflect.TypeOf(args).Elem()
	valueType := reflect.ValueOf(args).Elem()
	step := 0
	for step < paramsType.NumField()-1 {
		//捕捉结构
		vField := paramsType.Field(step)
		vValueType := valueType.Field(step)
		//下一步
		step += 1
		//检查参数是否存在
		isFind := false
		for _, v := range fieldList {
			if v == vField.Tag.Get("db") {
				isFind = true
				break
			}
		}
		if !isFind {
			err = errors.New("no support field: " + vField.Tag.Get("db"))
			return
		}
		//找到ID
		if vField.Tag.Get("db") == "id" {
			argID = vValueType.Int()
		}
		//找到更新字段
		if vField.Tag.Get("db") != "id" {
			setFields = append(setFields, vField.Tag.Get("db"))
			if setVals == nil {
				setVals = make(map[string]any)
			}
			setVals[vField.Tag.Get("db")] = vValueType.Interface()
		}
	}
	//执行更新
	ctx := c.quickClient.client.Update().NeedSoft(c.quickClient.openSoftDelete).AddWhereID(argID).SetFields(setFields)
	err = ctx.NamedExec(setVals)
	if err != nil {
		return
	}
	//删除缓冲
	c.quickClient.DeleteCacheByID(argID)
	//反馈
	return
}
