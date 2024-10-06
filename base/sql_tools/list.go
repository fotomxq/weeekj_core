package BaseSQLTools

import (
	"errors"
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	"reflect"
)

// QuickList 获取列表
type QuickList struct {
	//Quick
	quickClient *Quick
}

// ArgsGetListSimple 获取简单列表参数
type ArgsGetListSimple struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//检索字段
	// 仅支持field_list:"true"的字段
	ConditionFields []ArgsGetListSimpleConditionID
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

type ArgsGetListSimpleConditionID struct {
	//字段名称
	Name string `json:"name"`
	//值
	Val any `json:"val"`
}

// GetListSimple 获取简单列表
func (c *QuickList) GetListSimple(args *ArgsGetListSimple, result any) (dataCount int64, err error) {
	//组装条件
	ctx := c.quickClient.client.Select().SetDefaultListFields().SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove)
	if args.ConditionFields != nil && len(args.ConditionFields) > 0 {
		conditionFields := c.quickClient.getFieldsNameByConditionBoolTrue("field_list")
		for _, v := range args.ConditionFields {
			isFind := false
			for _, v2 := range conditionFields {
				if v.Name == v2 {
					isFind = true
					break
				}
			}
			if !isFind {
				err = errors.New(fmt.Sprintf("no support field: %s", v.Name))
			}
			switch v.Val.(type) {
			case int:
				if v.Val.(int) < 0 {
					break
				}
				ctx = ctx.SetIntQuery(v.Name, v.Val.(int))
			case int64:
				if v.Val.(int64) < 0 {
					break
				}
				ctx = ctx.SetIDQuery(v.Name, v.Val.(int64))
			case string:
				if v.Val.(string) == "" {
					break
				}
				ctx = ctx.SetStringQuery(v.Name, v.Val.(string))
			default:
				err = errors.New(fmt.Sprintf("no support type: %s(%s)", v.Name, reflect.TypeOf(v.Val).String()))
				return
			}
		}
	}
	if args.Search != "" {
		searchFields := c.quickClient.getFieldsNameByConditionBoolTrue("field_search")
		if len(searchFields) > 0 {
			ctx = ctx.SetSearchQuery(searchFields, args.Search)
		}
	}
	//获取数据
	dataCount, err = ctx.ResultAndCount(result)
	if err != nil {
		return
	}
	//result转化为数组
	//paramsType := reflect.TypeOf(result).Elem()
	//valueType := reflect.ValueOf(result).Elem()
	//step := 0
	//for step < paramsType.NumField()-1 {
	//	//捕捉结构
	//	vField := paramsType.Field(step)
	//	vValueType := valueType.Field(step)
	//	//下一步
	//	step += 1
	//	//找到ID
	//	if vField.Name == "ID" || vField.Tag.Get("db") == "id" {
	//		//获取ID
	//		id := vValueType.Int()
	//		if id < 1 {
	//			err = errors.New("id error")
	//			return
	//		}
	//		//跳出
	//		break
	//	}
	//}
	//反馈
	return
}

// ArgsGetListPrev 获取高级列表参数
type ArgsGetListPrev struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetListPrev 获取高级列表
func (c *QuickList) GetListPrev(args *ArgsGetListPrev) (ctx *CoreSQL2.ClientListCtx) {
	ctx = c.quickClient.client.Select().SetDefaultListFields().SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove)
	if args.Search != "" {
		searchFields := c.quickClient.getFieldsNameByConditionBoolTrue("field_search")
		if len(searchFields) > 0 {
			ctx = ctx.SetSearchQuery(searchFields, args.Search)
		}
	}
	return
}

func (c *QuickList) GetListDo(ctx *CoreSQL2.ClientListCtx, result any) (dataCount int64, err error) {
	//获取数据
	dataCount, err = ctx.ResultAndCount(result)
	if err != nil {
		return
	}
	//反馈
	return
}
