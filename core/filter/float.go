package CoreFilter

import "reflect"

// CorrectFloat64ByStruct 修正特定结构体内所有浮点数长度
/**
用途：
1. 修正浮点数异常数据，该数据主要是由于程序上做了 /0 动作导致的
2. 浮点数最多保留后6位

使用方法：
CorrectFloat64ByStruct(&data)
由于是引用调用关系，所以会直接修改原值，无需捕捉反馈
*/
func CorrectFloat64ByStruct(data interface{}) {
	//获取结构体
	valueType := reflect.ValueOf(data).Elem()
	//开始遍历
	step := 0
	for step < valueType.NumField() {
		//捕捉结构
		vValueType := valueType.Field(step)
		//下一步
		step += 1
		//找到float64数据类型
		if vValueType.Kind() == reflect.Float64 {
			//修正数据
			vValueType.SetFloat(RoundTo6DecimalPlaces(vValueType.Float()))
		}
	}
}
