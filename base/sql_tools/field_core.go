package SQLTools

import "reflect"

// GetFields 获取结构体字段列表
func (c *Quick) getFields() []string {
	return c.client.GetFields()
}

// resultGetFieldsByCondition 获取满足条件的字段列表
type resultGetFieldsByCondition struct {
	//排序
	Index int
	//字段名称
	FieldName string
	//条件名称
	ConditionName string
	//条件值
	ConditionValue string
}

// getFieldsByCondition 获取满足条件的字段列表
func (c *Quick) getFieldsByCondition(conditionName string) (result []resultGetFieldsByCondition) {
	//获取结构体
	paramsType := reflect.TypeOf(c.client.StructData).Elem()
	step := 0
	for step < paramsType.NumField() {
		//获取当前节点值
		vField := paramsType.Field(step)
		//获取值
		vVal := vField.Tag.Get(conditionName)
		if vVal != "" {
			result = append(result, resultGetFieldsByCondition{
				Index:          step,
				FieldName:      vField.Name,
				ConditionName:  conditionName,
				ConditionValue: vVal,
			})
		}
		//下一步
		step += 1
	}
	return
}

// getFieldsNameByConditionBoolTrue 获取满足条件的字段名称列表
func (c *Quick) getFieldsNameByConditionBoolTrue(conditionName string) (result []string) {
	//获取满足条件的字段列表
	fields := c.getFieldsByCondition(conditionName)
	for _, v := range fields {
		if v.ConditionValue == "true" {
			result = append(result, v.FieldName)
		}
	}
	return
}
