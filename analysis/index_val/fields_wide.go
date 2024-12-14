package AnalysisIndexVal

import "time"

// FieldsWide 超宽表记录
// 该结构和val存在一些差异，但本质上是一致的，可以方便记录同一个时间下、同一个维度的大量数据集
type FieldsWide struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0" index:"true"`
	//年月日
	// 可任意持续，如年，或仅年月
	// 不建议构建小时及以下级别的指标
	// 同一个维度和时间范围，仅会存在一个数据，否则将覆盖
	YearMD string `db:"year_md" json:"yearMD" index:"true" field_list:"true"`
	//扩展维度1
	// 可建议特别的维度关系，例如特定供应商的数据、特定地区的数据等
	Extend1 string `db:"extend1" json:"extend1" index:"true" field_list:"true"`
	//扩展维度2
	Extend2 string `db:"extend2" json:"extend2" index:"true" field_list:"true"`
	//扩展维度3
	Extend3 string `db:"extend3" json:"extend3" index:"true" field_list:"true"`
	//扩展维度4
	Extend4 string `db:"extend4" json:"extend4" index:"true" field_list:"true"`
	//扩展维度5
	Extend5 string `db:"extend5" json:"extend5" index:"true" field_list:"true"`
	//是否为预测值
	IsForecast bool `db:"is_forecast" json:"isForecast"`
	//////////////////////////////////////////////////////////////////////////////////
	//值
	//////////////////////////////////////////////////////////////////////////////////
	//指标编码
	Code1  string  `db:"code1" json:"code1" index:"true" field_list:"true"`
	Val1   float64 `db:"val1" json:"val1"`
	Code2  string  `db:"code2" json:"code2" index:"true" field_list:"true"`
	Val2   float64 `db:"val2" json:"val2"`
	Code3  string  `db:"code3" json:"code3" index:"true" field_list:"true"`
	Val3   float64 `db:"val3" json:"val3"`
	Code4  string  `db:"code4" json:"code4" index:"true" field_list:"true"`
	Val4   float64 `db:"val4" json:"val4"`
	Code5  string  `db:"code5" json:"code5" index:"true" field_list:"true"`
	Val5   float64 `db:"val5" json:"val5"`
	Code6  string  `db:"code6" json:"code6" index:"true" field_list:"true"`
	Val6   float64 `db:"val6" json:"val6"`
	Code7  string  `db:"code7" json:"code7" index:"true" field_list:"true"`
	Val7   float64 `db:"val7" json:"val7"`
	Code8  string  `db:"code8" json:"code8" index:"true" field_list:"true"`
	Val8   float64 `db:"val8" json:"val8"`
	Code9  string  `db:"code9" json:"code9" index:"true" field_list:"true"`
	Val9   float64 `db:"val9" json:"val9"`
	Code10 string  `db:"code10" json:"code10" index:"true" field_list:"true"`
	Val10  float64 `db:"val10" json:"val10"`
	Code11 string  `db:"code11" json:"code11" index:"true" field_list:"true"`
	Val11  float64 `db:"val11" json:"val11"`
	Code12 string  `db:"code12" json:"code12" index:"true" field_list:"true"`
	Val12  float64 `db:"val12" json:"val12"`
	Code13 string  `db:"code13" json:"code13" index:"true" field_list:"true"`
	Val13  float64 `db:"val13" json:"val13"`
	Code14 string  `db:"code14" json:"code14" index:"true" field_list:"true"`
	Val14  float64 `db:"val14" json:"val14"`
	Code15 string  `db:"code15" json:"code15" index:"true" field_list:"true"`
	Val15  float64 `db:"val15" json:"val15"`
	Code16 string  `db:"code16" json:"code16" index:"true" field_list:"true"`
	Val16  float64 `db:"val16" json:"val16"`
	Code17 string  `db:"code17" json:"code17" index:"true" field_list:"true"`
	Val17  float64 `db:"val17" json:"val17"`
	Code18 string  `db:"code18" json:"code18" index:"true" field_list:"true"`
	Val18  float64 `db:"val18" json:"val18"`
	Code19 string  `db:"code19" json:"code19" index:"true" field_list:"true"`
	Val19  float64 `db:"val19" json:"val19"`
	Code20 string  `db:"code20" json:"code20" index:"true" field_list:"true"`
	Val20  float64 `db:"val20" json:"val20"`
	Code21 string  `db:"code21" json:"code21" index:"true" field_list:"true"`
	Val21  float64 `db:"val21" json:"val21"`
	Code22 string  `db:"code22" json:"code22" index:"true" field_list:"true"`
	Val22  float64 `db:"val22" json:"val22"`
	Code23 string  `db:"code23" json:"code23" index:"true" field_list:"true"`
	Val23  float64 `db:"val23" json:"val23"`
	Code24 string  `db:"code24" json:"code24" index:"true" field_list:"true"`
	Val24  float64 `db:"val24" json:"val24"`
	Code25 string  `db:"code25" json:"code25" index:"true" field_list:"true"`
	Val25  float64 `db:"val25" json:"val25"`
	Code26 string  `db:"code26" json:"code26" index:"true" field_list:"true"`
	Val26  float64 `db:"val26" json:"val26"`
	Code27 string  `db:"code27" json:"code27" index:"true" field_list:"true"`
	Val27  float64 `db:"val27" json:"val27"`
	Code28 string  `db:"code28" json:"code28" index:"true" field_list:"true"`
	Val28  float64 `db:"val28" json:"val28"`
	Code29 string  `db:"code29" json:"code29" index:"true" field_list:"true"`
	Val29  float64 `db:"val29" json:"val29"`
	Code30 string  `db:"code30" json:"code30" index:"true" field_list:"true"`
	Val30  float64 `db:"val30" json:"val30"`
	Code31 string  `db:"code31" json:"code31" index:"true" field_list:"true"`
	Val31  float64 `db:"val31" json:"val31"`
	Code32 string  `db:"code32" json:"code32" index:"true" field_list:"true"`
	Val32  float64 `db:"val32" json:"val32"`
	Code33 string  `db:"code33" json:"code33" index:"true" field_list:"true"`
	Val33  float64 `db:"val33" json:"val33"`
	Code34 string  `db:"code34" json:"code34" index:"true" field_list:"true"`
	Val34  float64 `db:"val34" json:"val34"`
	Code35 string  `db:"code35" json:"code35" index:"true" field_list:"true"`
	Val35  float64 `db:"val35" json:"val35"`
	Code36 string  `db:"code36" json:"code36" index:"true" field_list:"true"`
	Val36  float64 `db:"val36" json:"val36"`
	Code37 string  `db:"code37" json:"code37" index:"true" field_list:"true"`
	Val37  float64 `db:"val37" json:"val37"`
	Code38 string  `db:"code38" json:"code38" index:"true" field_list:"true"`
	Val38  float64 `db:"val38" json:"val38"`
	Code39 string  `db:"code39" json:"code39" index:"true" field_list:"true"`
	Val39  float64 `db:"val39" json:"val39"`
	Code40 string  `db:"code40" json:"code40" index:"true" field_list:"true"`
	Val40  float64 `db:"val40" json:"val40"`
	Code41 string  `db:"code41" json:"code41" index:"true" field_list:"true"`
	Val41  float64 `db:"val41" json:"val41"`
	Code42 string  `db:"code42" json:"code42" index:"true" field_list:"true"`
	Val42  float64 `db:"val42" json:"val42"`
	Code43 string  `db:"code43" json:"code43" index:"true" field_list:"true"`
	Val43  float64 `db:"val43" json:"val43"`
	Code44 string  `db:"code44" json:"code44" index:"true" field_list:"true"`
	Val44  float64 `db:"val44" json:"val44"`
	Code45 string  `db:"code45" json:"code45" index:"true" field_list:"true"`
	Val45  float64 `db:"val45" json:"val45"`
	Code46 string  `db:"code46" json:"code46" index:"true" field_list:"true"`
	Val46  float64 `db:"val46" json:"val46"`
	Code47 string  `db:"code47" json:"code47" index:"true" field_list:"true"`
	Val47  float64 `db:"val47" json:"val47"`
	Code48 string  `db:"code48" json:"code48" index:"true" field_list:"true"`
	Val48  float64 `db:"val48" json:"val48"`
	Code49 string  `db:"code49" json:"code49" index:"true" field_list:"true"`
	Val49  float64 `db:"val49" json:"val49"`
	Code50 string  `db:"code50" json:"code50" index:"true" field_list:"true"`
	Val50  float64 `db:"val50" json:"val50"`
}
