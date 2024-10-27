package DataLakeSource

import (
	"errors"
	"fmt"
	BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"reflect"
	"strings"
)

// importStructTableRelation 通过表格关系构建实体表
// 也可以用于同步表格结构，即已经具备实体表，但需要对表结构进行更新
func importStructTableRelation(tableID int64) (errCode string, err error) {
	//获取表和结构体
	var tableData FieldsTable
	tableData, err = GetTableDetail(tableID)
	if err != nil || tableData.ID < 1 || CoreFilter.CheckHaveTime(tableData.DeleteAt) {
		errCode = "report_data_empty"
		err = errors.New(fmt.Sprint("table not found, ", err))
		return
	}
	var fieldList []FieldsFields
	fieldList, _, err = GetFieldsListByTableID(tableID)
	if err != nil {
		errCode = "report_data_empty"
		err = errors.New(fmt.Sprint("field not found, ", err))
		return
	}
	//构建结构体
	fields := make([]reflect.StructField, 0)
	//构建插入sql
	for k := 0; k < len(fieldList); k++ {
		v := fieldList[k]
		switch v.FieldName {
		case "id":
			fields = append(fields, reflect.StructField{
				Name:      "ID",
				PkgPath:   "",
				Type:      reflect.TypeOf(""),
				Tag:       reflect.StructTag(fmt.Sprintf("`db:\"%s\" json:\"%s\"` unique:\"true\"", "id", "id")),
				Offset:    0,
				Index:     nil,
				Anonymous: false,
			})
		case "create_at":
			fields = append(fields, reflect.StructField{
				Name:      "CreateAt",
				PkgPath:   "",
				Type:      reflect.TypeOf(""),
				Tag:       reflect.StructTag(fmt.Sprintf("`db:\"%s\" json:\"%s\"` default:\"now()\"", "create_at", "createAt")),
				Offset:    0,
				Index:     nil,
				Anonymous: false,
			})
		case "update_at":
			fields = append(fields, reflect.StructField{
				Name:      "UpdateAt",
				PkgPath:   "",
				Type:      reflect.TypeOf(""),
				Tag:       reflect.StructTag(fmt.Sprintf("`db:\"%s\" json:\"%s\"` default:\"now()\"", "update_at", "updateAt")),
				Offset:    0,
				Index:     nil,
				Anonymous: false,
			})
		case "delete_at":
			fields = append(fields, reflect.StructField{
				Name:      "DeleteAt",
				PkgPath:   "",
				Type:      reflect.TypeOf(""),
				Tag:       reflect.StructTag(fmt.Sprintf("`db:\"%s\" json:\"%s\"` default:\"0\"", "delete_at", "deleteAt")),
				Offset:    0,
				Index:     nil,
				Anonymous: false,
			})
		default:
			//自动根据字段构建数据
			vJsonStr := fmt.Sprintf("`db:\"%s\" json:\"%s\"`", v.FieldName, v.FieldName)
			if v.IsIndex {
				vJsonStr += " index:\"true\""
			}
			if v.IsSearch {
				vJsonStr += " field_search:\"true\""
			}
			fields = append(fields, reflect.StructField{
				Name:      strings.ToUpper(v.FieldName),
				PkgPath:   "",
				Type:      reflect.TypeOf(""),
				Tag:       reflect.StructTag(vJsonStr),
				Offset:    0,
				Index:     nil,
				Anonymous: false,
			})
		}
	}
	structType := reflect.StructOf(fields)
	//创建表结构
	var newTableStruct BaseSQLTools.Quick
	err = newTableStruct.Init(tableData.TableName, reflect.New(structType).Interface())
	if err != nil {
		errCode = "report_create_failed"
		err = errors.New(fmt.Sprint("create table struct failed, ", err))
		return
	}
	//反馈
	return
}
