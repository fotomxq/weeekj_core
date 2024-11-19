package AnalysisIndexDimensions

import (
	"errors"
	"fmt"
)

// ArgsGetVals 根据表关系获取维度的全部可选值参数
type ArgsGetVals struct {
	//编码
	// 维度编码，用于程序内部识别
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50"`
	//是否限制长度
	Limit int64 `db:"limit" json:"limit"`
}

// DataGetVals 根据表关系获取维度的全部可选值数据
type DataGetVals struct {
	//维度值
	Val string `db:"val" json:"val"`
}

// GetVals 根据表关系获取维度的全部可选值
func GetVals(args *ArgsGetVals) (dataList []DataGetVals, err error) {
	//获取维度定义
	var dimensionsData FieldsDimensions
	dimensionsData, err = GetDimensionsByCode(args.Code)
	if err != nil {
		err = errors.New(fmt.Sprint("dimensions not found, ", err))
		return
	}
	//必须定义过维度的表和字段
	if dimensionsData.TableName == "" || dimensionsData.FieldName == "" {
		err = errors.New("table or field not defined")
		return
	}
	//根据维度的定义，获取指定表和字段的数据集
	_ = dimensionsDB.GetClient().DB.GetPostgresql().Select(&dataList, fmt.Sprint("select ", dimensionsData.FieldName, " as val from ", dimensionsData.TableName, " group by ", dimensionsData.FieldName))
	//反馈
	return
}
