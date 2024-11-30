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
	valueType := reflect.ValueOf(args).Elem()
	//获取写入参数
	var setVals map[string]any
	setVals = make(map[string]any)
	//开始遍历
	step := 0
	for step < paramsType.NumField() {
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
			//禁止设置
		case "create_at":
			//禁止标记创造
		case "update_at":
			//禁止标记更新
		case "delete_at":
			//禁止标记删除
		default:
			//找到写入字段
			setVals[vField.Tag.Get("db")] = vValueType.Interface()
		}
	}
	//写入准备
	ctx := c.quickClient.client.Insert().SetDefaultInsertFields().Add(setVals)
	//写入执行
	newID, err = ctx.ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}
