package AnalysisIndexValCustom

import (
	"fmt"
	BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	"time"
)

// GetValLastTime 获取最后导入的时间
func GetValLastTime() (result time.Time, dataCount int64) {
	var dataList []FieldsVal
	dataCount, _ = indexValCustomDB.GetList().GetListSimple(&BaseSQLTools.ArgsGetListSimple{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  1,
			Sort: "id",
			Desc: true,
		},
		ConditionFields: nil,
		IsRemove:        false,
		Search:          "",
	}, &dataList)
	if len(dataList) < 1 {
		return
	}
	result = dataList[0].CreateAt
	return
}

// ArgsGetVal 获取具体的数据参数
type ArgsGetVal struct {
	//编码
	// 注意和指标编码可以是不同的，主要用于程序内部识别
	// 例如约定指标集合为履约合同数据集，那么此处可约定为一个缩写，方便程序寻找对应数据
	// 如果维度关系太多，建议拆分成不同的code，以便于存储、使用
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
	//年月日
	// 可任意持续，如年，或仅年月
	// 不建议构建小时及以下级别的指标
	// 同一个维度和时间范围，仅会存在一个数据，否则将覆盖
	// 如果存储具体的值，也可以是实际发生的内容，为了统计的便利性，建议使用年月日或年月，以减少数据的复杂性
	YearMD string `db:"year_md" json:"yearMD" index:"true" field_list:"true"`
	/////////////////////////////////////////////////////////////////////////////////////////////////
	//维度关系
	// 维度关系层可依赖于实施数据的切分逻辑，如地区、行为标记等，以方便筛选数据
	// 例如，如果是履约合同，可建议维度关系为供应商、采购商、地区等
	// 也可以直接和维度关系模块进行关联
	/////////////////////////////////////////////////////////////////////////////////////////////////
	//扩展维度1
	// 可建议特别的维度关系，例如特定供应商的数据、特定地区的数据等
	// 如果是履约合同，也可以是采购方式等维度
	Extend1 string `db:"extend1" json:"extend1" index:"true" field_list:"true"`
	//扩展维度2
	Extend2 string `db:"extend2" json:"extend2" index:"true" field_list:"true"`
	//扩展维度3
	Extend3 string `db:"extend3" json:"extend3" index:"true" field_list:"true"`
	//扩展维度4
	Extend4 string `db:"extend4" json:"extend4" index:"true" field_list:"true"`
	//扩展维度5
	Extend5 string `db:"extend5" json:"extend5" index:"true" field_list:"true"`
	//扩展维度6
	Extend6 string `db:"extend6" json:"extend6" index:"true" field_list:"true"`
	//扩展维度7
	Extend7 string `db:"extend7" json:"extend7" index:"true" field_list:"true"`
	//扩展维度8
	Extend8 string `db:"extend8" json:"extend8" index:"true" field_list:"true"`
	//扩展维度9
	Extend9 string `db:"extend9" json:"extend9" index:"true" field_list:"true"`
}

// GetVal 获取具体的数据
func GetVal(args *ArgsGetVal) (data FieldsVal, err error) {
	err = indexValCustomDB.GetClient().DB.GetPostgresql().Get(&data, "SELECT * FROM "+indexValCustomDB.GetClient().TableName+" WHERE delete_at < to_timestamp(1000000) AND code = $1 AND year_md = $2 AND extend1 = $3 AND extend2 = $4 AND extend3 = $5 AND extend4 = $6 AND extend5 = $7 AND extend6 = $8 AND extend7 = $9 AND extend8 = $10 AND extend9 = $11 AND is_remove = false", args.Code, args.YearMD, args.Extend1, args.Extend2, args.Extend3, args.Extend4, args.Extend5, args.Extend6, args.Extend7, args.Extend8, args.Extend9)
	if err != nil {
		return
	}
	return
}

// ArgsGetValListParams 获取指标值列表参数结构体
type ArgsGetValListParams struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetValList 获取指标值列表
func GetValList(args *ArgsGetValListParams) (dataList []FieldsVal, dataCount int64, err error) {
	//获取数据
	dataCount, err = indexValCustomDB.GetList().GetListSimple(&BaseSQLTools.ArgsGetListSimple{
		Pages:           args.Pages,
		ConditionFields: nil,
		IsRemove:        args.IsRemove,
		Search:          args.Search,
	}, &dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		var vData FieldsVal
		err = indexValCustomDB.GetInfo().GetInfoByID(v.ID, &vData)
		if err != nil || vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	//反馈
	return
}

// GetIndexExtendDistinctList 获取指定维度的所有可选值
func GetIndexExtendDistinctList(extendNum int) (dataList []string, err error) {
	//获取数据
	dataList, err = indexValCustomDB.GetList().GetDistinctList(fmt.Sprintf("extend%d", extendNum))
	if err != nil {
		return
	}
	//反馈
	return
}

func getValByID(id int64) (data FieldsVal) {
	_ = indexValCustomDB.GetInfo().GetInfoByID(id, &data)
	//反馈
	return
}
