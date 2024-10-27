package DataLakeSource

import (
	"encoding/csv"
	"errors"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"github.com/xuri/excelize/v2"
	"io"
	"os"
)

// ArgsImportStructExcel 通用导入Excel文件快速建立表结构参数
type ArgsImportStructExcel struct {
	//表名称
	TableName string `json:"tableName"`
	//源文件路径
	// 仅支持csv/xlsx文件
	Src string `json:"src"`
}

// ImportStructExcel 通用导入Excel文件快速建立表结构
func ImportStructExcel(args *ArgsImportStructExcel) (errCode string, err error) {
	//检查表名称的唯一性
	var tableData FieldsTable
	tableData, err = GetTableDetailByName(args.TableName)
	if err == nil && tableData.ID > 0 && CoreFilter.CheckHaveTime(tableData.DeleteAt) {
		errCode = "err_have_replace"
		err = errors.New("table is exist")
		return
	}
	//创建表
	tableData, err = GetTableDetailByName(args.TableName)
	if err != nil {
		errCode = "report_data_empty"
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
		err = errors.New("not support file type")
		return
	}
	//遍历列头，并创建
	for k := 0; k < len(head); k++ {
		//分析数据类型
		var dataType string
		paramType := CoreFilter.DetermineType(params[k])
		switch paramType {
		case "int":
			dataType = "integer"
		case "int64":
			dataType = "bigint"
		case "float64":
			dataType = "float"
		case "string":
			dataType = "text"
		default:
			dataType = "text"
		}
		//创建字段
		_, err = CreateFields(&ArgsCreateFields{
			TableID:       tableData.ID,
			FieldName:     head[k],
			FieldLabel:    head[k],
			InputType:     "",
			InputLength:   0,
			InputDefault:  "",
			InputRequired: false,
			InputPattern:  "",
			DataType:      dataType,
			FieldDesc:     head[k],
		})
		if err != nil {
			errCode = "report_create_failed"
			return
		}
	}
	//如果导入失败，则清理数据表和字段记录
	if err != nil || len(head) < 1 {
		//删除表
		if tableData.ID > 0 {
			err = DeleteTable(tableData.ID)
			if err != nil {
				errCode = "report_delete_failed"
				return
			}
		}
	}
	//反馈
	return
}
