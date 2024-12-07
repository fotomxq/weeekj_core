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

// InsertOrUpdateRowByID 融合插入数据
// 自动更新或插入，根据 id 判断
func (c *QuickInsert) InsertOrUpdateRowByID(args any) (err error) {
	//找到ID
	var findID int64
	//获取参数
	paramsType := reflect.TypeOf(args).Elem()
	valueType := reflect.ValueOf(args).Elem()
	//开始遍历
	step := 0
	for step < paramsType.NumField() {
		//捕捉结构
		vField := paramsType.Field(step)
		vValueType := valueType.Field(step)
		vTagDB := vField.Tag.Get("db")
		//下一步
		step += 1
		//内置字段禁止设置，主要用于其他操作
		switch vTagDB {
		case "id":
			findID = vValueType.Int()
		}
	}
	//如果存在ID，则更新数据
	if findID > 0 {
		err = c.quickClient.GetUpdate().UpdateByID(args)
	} else {
		_, err = c.InsertRow(args)
	}
	//反馈
	return
}

// InsertOrUpdateRowByField 融合插入数据
// 自动更新或插入，根据 findField 判断
func (c *QuickInsert) InsertOrUpdateRowByField(args any, findField string, findVal any, haveDelete bool) (err error) {
	//找到数据
	var findID int64
	findID, err = c.quickClient.GetInfo().GetInfoByFieldToID(findField, findVal, haveDelete)
	//TODO: 需优化，找到ID后赋予 args
	//如果存在ID，则更新数据
	if err == nil && findID > 0 {
		err = c.quickClient.GetUpdate().UpdateByID(args)
	} else {
		_, err = c.InsertRow(args)
	}
	//反馈
	return
}
