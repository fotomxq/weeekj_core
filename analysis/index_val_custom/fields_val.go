package AnalysisIndexValCustom

import "time"

// FieldsVal 数据集合
type FieldsVal struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0" index:"true"`
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
	/////////////////////////////////////////////////////////////////////////////////////////////////
	//数据值集合
	// 根据项目需求，可赋予值具体的定义和内容
	// 不应该对具体日期等其他类型数据记录，因为统计分析中，主要围绕的是值的变化，而不是具体的日期或其他类型的数据
	// 如存在其他数据的，建议归一化后存储
	/////////////////////////////////////////////////////////////////////////////////////////////////
	//值结果1
	Val1 float64 `db:"val1" json:"val1"`
	//值结果2
	Val2 float64 `db:"val2" json:"val2"`
	//值结果3
	Val3 float64 `db:"val3" json:"val3"`
	//值结果4
	Val4 float64 `db:"val4" json:"val4"`
	//值结果5
	Val5 float64 `db:"val5" json:"val5"`
	//值结果6
	Val6 float64 `db:"val6" json:"val6"`
	//值结果7
	Val7 float64 `db:"val7" json:"val7"`
	//值结果8
	Val8 float64 `db:"val8" json:"val8"`
	//值结果9
	Val9 float64 `db:"val9" json:"val9"`
	//值结果10
	Val10 float64 `db:"val10" json:"val10"`
	//值结果11
	Val11 float64 `db:"val11" json:"val11"`
	//值结果12
	Val12 float64 `db:"val12" json:"val12"`
	//值结果13
	Val13 float64 `db:"val13" json:"val13"`
	//值结果14
	Val14 float64 `db:"val14" json:"val14"`
	//值结果15
	Val15 float64 `db:"val15" json:"val15"`
	//值结果16
	Val16 float64 `db:"val16" json:"val16"`
	//值结果17
	Val17 float64 `db:"val17" json:"val17"`
	//值结果18
	Val18 float64 `db:"val18" json:"val18"`
	//值结果19
	Val19 float64 `db:"val19" json:"val19"`
	//值结果20
	Val20 float64 `db:"val20" json:"val20"`
	//值结果21
	Val21 float64 `db:"val21" json:"val21"`
	//值结果22
	Val22 float64 `db:"val22" json:"val22"`
	//值结果23
	Val23 float64 `db:"val23" json:"val23"`
	//值结果24
	Val24 float64 `db:"val24" json:"val24"`
	//值结果25
	Val25 float64 `db:"val25" json:"val25"`
	//值结果26
	Val26 float64 `db:"val26" json:"val26"`
	//值结果27
	Val27 float64 `db:"val27" json:"val27"`
	//值结果28
	Val28 float64 `db:"val28" json:"val28"`
	//值结果29
	Val29 float64 `db:"val29" json:"val29"`
	//值结果30
	Val30 float64 `db:"val30" json:"val30"`
	//值结果31
	Val31 float64 `db:"val31" json:"val31"`
	//值结果32
	Val32 float64 `db:"val32" json:"val32"`
	//值结果33
	Val33 float64 `db:"val33" json:"val33"`
	//值结果34
	Val34 float64 `db:"val34" json:"val34"`
	//值结果35
	Val35 float64 `db:"val35" json:"val35"`
	//值结果36
	Val36 float64 `db:"val36" json:"val36"`
	//值结果37
	Val37 float64 `db:"val37" json:"val37"`
	//值结果38
	Val38 float64 `db:"val38" json:"val38"`
	//值结果39
	Val39 float64 `db:"val39" json:"val39"`
	//值结果40
	Val40 float64 `db:"val40" json:"val40"`
	//值结果41
	Val41 float64 `db:"val41" json:"val41"`
	//值结果42
	Val42 float64 `db:"val42" json:"val42"`
	//值结果43
	Val43 float64 `db:"val43" json:"val43"`
	//值结果44
	Val44 float64 `db:"val44" json:"val44"`
	//值结果45
	Val45 float64 `db:"val45" json:"val45"`
	//值结果46
	Val46 float64 `db:"val46" json:"val46"`
	//值结果47
	Val47 float64 `db:"val47" json:"val47"`
	//值结果48
	Val48 float64 `db:"val48" json:"val48"`
	//值结果49
	Val49 float64 `db:"val49" json:"val49"`
	//值结果50
	Val50 float64 `db:"val50" json:"val50"`
}
