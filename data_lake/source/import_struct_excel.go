package DataLakeSource

import (
	"encoding/csv"
	"errors"
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"github.com/mozillazg/go-pinyin"
	"github.com/xuri/excelize/v2"
	"io"
	"os"
	"strings"
)

// ArgsImportStructExcel 通用导入Excel文件快速建立表结构参数
type ArgsImportStructExcel struct {
	//表名称
	TableName string `db:"table_name" json:"tableName" field_search:"true"`
	//表描述
	TableDesc string `db:"table_desc" json:"tableDesc" field_search:"true"`
	//提示名称
	TipName string `db:"tip_name" json:"tipName" field_search:"true"`
	//数据唯一渠道名称
	// 如果是多处来源，应拆分表
	ChannelName string `db:"channel_name" json:"channelName" field_search:"true"`
	//数据唯一渠道提示名称
	ChannelTipName string `db:"channel_tip_name" json:"channelTipName" field_search:"true"`
	//源文件路径
	// 仅支持csv/xlsx文件
	Src string `json:"src"`
}

// ImportStructExcel 通用导入Excel文件快速建立表结构
func ImportStructExcel(args *ArgsImportStructExcel) (tableID int64, errCode string, err error) {
	//检查表名称的唯一性
	var tableData FieldsTable
	tableData, err = GetTableDetailByName(args.TableName)
	if err == nil && tableData.ID > 0 && CoreFilter.CheckHaveTime(tableData.DeleteAt) {
		errCode = "err_have_replace"
		err = errors.New(fmt.Sprint("table is exist, table name: ", args.TableName))
		return
	}
	//创建表
	tableID, err = CreateTable(&ArgsCreateTable{
		TableName:      args.TableName,
		TableDesc:      args.TableDesc,
		TipName:        args.TipName,
		ChannelName:    args.ChannelName,
		ChannelTipName: args.ChannelTipName,
	})
	if err != nil {
		errCode = "report_create_failed"
		err = errors.New(fmt.Sprint("create table failed: ", err))
		return
	}
	//重新获取表
	tableData, err = GetTableDetail(tableID)
	if err != nil {
		errCode = "report_data_empty"
		err = errors.New(fmt.Sprint("get table failed: ", err))
		return
	}
	//列头清单
	var head []string
	//参考数据表
	var params []string
	//识别文件格式
	// 仅支持csv/xlsx文件
	switch CoreFile.GetFileType(args.Src) {
	case "csv":
		//导入csv文件
		var fs *os.File
		fs, err = os.Open(args.Src)
		if err != nil {
			errCode = "err_io"
			return
		}
		defer fs.Close()
		r := csv.NewReader(fs)
		//读取数据
		step := 0
		for {
			if step == 0 {
				head, err = r.Read()
			}
			if step == 1 {
				params, err = r.Read()
			}
			if err != nil && err != io.EOF {
				errCode = "err_io"
				return
			}
			if step == 2 {
				break
			}
			if err == io.EOF {
				break
			}
			step += 1
		}
	case "xlsx":
		//导入xlsx文件
		//读取excel文件
		var fs *excelize.File
		fs, err = excelize.OpenFile(args.Src)
		if err != nil {
			errCode = "err_io"
			return
		}
		defer fs.Close()
		//获取第一张数据表
		sheetName := fs.GetSheetName(0)
		var rows *excelize.Rows
		rows, err = fs.Rows(sheetName)
		if err != nil {
			errCode = "err_io"
			return
		}
		step := 0
		for rows.Next() {
			if step == 0 {
				head, err = rows.Columns()
			}
			if step == 1 {
				params, err = rows.Columns()
			}
			if step == 2 {
				break
			}
			if err != nil {
				errCode = "err_io"
				return
			}
			step += 1
		}
	default:
		//不支持的文件格式
		errCode = "err_file_type"
		err = errors.New(fmt.Sprint("not support file type: ", args.Src))
		return
	}
	//遍历列头，并创建
	for k := 0; k < len(head); k++ {
		v := head[k]
		//分析数据类型
		var inputType string
		var fieldType string
		paramType := CoreFilter.DetermineType(params[k])
		switch paramType {
		case "int":
			inputType = FIELDS_INPUT_TYPE_ENUM_NUMBER
			fieldType = FIELDS_DATA_TYPE_ENUM_INT
		case "int64":
			inputType = FIELDS_INPUT_TYPE_ENUM_NUMBER
			fieldType = FIELDS_DATA_TYPE_ENUM_INT64
		case "float64":
			inputType = FIELDS_INPUT_TYPE_ENUM_NUMBER
			fieldType = FIELDS_DATA_TYPE_ENUM_FLOAT
		case "bool":
			inputType = FIELDS_INPUT_TYPE_ENUM_RADIO
			fieldType = FIELDS_DATA_TYPE_ENUM_BOOL
		case "string":
			inputType = FIELDS_INPUT_TYPE_ENUM_TEXTAREA
			fieldType = FIELDS_DATA_TYPE_ENUM_TEXT
		case "date":
			inputType = FIELDS_INPUT_TYPE_ENUM_DATE
			fieldType = FIELDS_DATA_TYPE_ENUM_DATE
		case "datetime":
			inputType = FIELDS_INPUT_TYPE_ENUM_DATETIME
			fieldType = FIELDS_DATA_TYPE_ENUM_DATETIME
		default:
			inputType = FIELDS_INPUT_TYPE_ENUM_TEXTAREA
			fieldType = "text"
		}
		//修订字段名称
		var vFieldName = v
		//检查v.FieldName如果为非英文
		if !CoreFilter.CheckMark(v) {
			//转为拼音
			strArr := pinyin.LazyPinyin(v, pinyin.NewArgs())
			for strK, strV := range strArr {
				strArr[strK] = strings.ToUpper(strV)
			}
			vFieldName = strings.Join(strArr, "")
		} else {
			//转为大写
			vFieldName = strings.ToUpper(vFieldName)
		}
		//如果vFieldName为空，则按照序列给值
		if vFieldName == "" {
			vFieldName = fmt.Sprint("F", k)
		}
		//创建字段
		_, err = CreateFields(&ArgsCreateFields{
			TableID:       tableData.ID,
			InputName:     v,
			InputType:     inputType,
			InputLength:   0,
			InputDefault:  "",
			InputRequired: false,
			InputPattern:  "",
			FieldName:     vFieldName,
			FieldLabel:    v,
			IsPrimary:     false,
			IsIndex:       false,
			IsSystem:      false,
			IsSearch:      false,
			DataType:      fieldType,
			FieldDesc:     v,
		})
		if err != nil {
			errCode = "report_create_failed"
			err = errors.New(fmt.Sprint("create fields failed: ", err, ", field name: ", v, ", field type: ", fieldType, ", field desc: ", v))
			return
		}
	}
	//如果导入失败，则清理数据表和字段记录
	if err != nil || len(head) < 1 {
		//删除表
		if tableData.ID > 0 {
			err = ClearFields(tableData.ID)
			if err != nil {
				errCode = "report_delete_failed"
				err = errors.New(fmt.Sprint("clear fields failed: ", err))
				return
			}
			err = DeleteTable(tableData.ID)
			if err != nil {
				errCode = "report_delete_failed"
				err = errors.New(fmt.Sprint("delete table failed: ", err))
				return
			}
		}
	}
	//构建实体表
	errCode, err = importStructTableRelation(tableID)
	if err != nil {
		err = errors.New(fmt.Sprint("import struct table relation failed: ", err))
		return
	}
	//反馈
	return
}
