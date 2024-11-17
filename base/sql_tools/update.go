package BaseSQLTools

import (
	"errors"
	"fmt"
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
		vTagDB := vField.Tag.Get("db")
		//下一步
		step += 1
		//检查参数是否存在
		isFind := false
		for _, v := range fieldList {
			if v == vTagDB {
				isFind = true
				break
			}
		}
		if !isFind {
			err = errors.New("no support field: " + vField.Tag.Get("db"))
			return
		}
		//内置字段禁止设置，主要用于其他操作
		switch vTagDB {
		case "id":
			//找到ID
			argID = vValueType.Int()
		case "create_at":
			//禁止更新
		case "delete_at":
			//软删除，不能采用此方法操作
		default:
			//找到更新字段
			setFields = append(setFields, vTagDB)
			if setVals == nil {
				setVals = make(map[string]any)
			}
			setVals[vField.Tag.Get("db")] = vValueType.Interface()
		}
	}
	//执行更新
	if len(setFields) < 1 {
		err = errors.New(fmt.Sprint("no update field, id: ", argID, ", fields: ", setFields))
		return
	}
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
