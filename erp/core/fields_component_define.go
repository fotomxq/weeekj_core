package ERPCore

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
)

// FieldsComponentDefineList 节点组件列
type FieldsComponentDefineList []FieldsComponentDefine

// Value sql底层处理器
func (t FieldsComponentDefineList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsComponentDefineList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// Len 排序支持
func (t FieldsComponentDefineList) Len() int {
	return len(t)
}
func (t FieldsComponentDefineList) Less(i, j int) bool {
	return t[i].Sort < t[j].Sort
}

func (t FieldsComponentDefineList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// FieldsComponentDefine 节点组
type FieldsComponentDefine struct {
	//组件key
	// 单个节点内必须唯一
	Key string `db:"key" json:"key"`
	//展示顺序
	Sort int `db:"sort" json:"sort"`
	//组件类型
	// 常规组件: text_input 文本输入框; des_input 长文本输入框; md MD_input 文本框; number_int 数字; number_float 浮点数; bool_open 布尔值开关; number_price 价格数据(1.00=100); number_p 百分比数据(100%=10000);
	// 富组件: input_email 邮箱; input_phone 手机号; input_ip IP地址;
	// 时间: time_default 默认时间结构体0000-00-00 00:00:00; time_date 日期0000-00-00;
	// 文件类: file_id 文件ID(上传文件组件); file_ids 一组文件ID列(批量上传文件组件);
	// 地理位置: map_point 选择地图的定位数据; map_marge 聚合地址信息结构（类似订单内部的地址信息结构体）(值会被记录到扩展参数，而不是val中); map_address 输入地址框，自动联动系统的地图组件获取地址的信息并填入扩展参数; map_city 城市选择器;
	// 文档联动: erp_doc_id 文档选择器(扩展参数config_id约定文档配置); service_user_info_id 信息档案数据ID;
	// 用户: user_id 用户ID;
	// 组织相关模块: org_cert_id 组织证件ID(扩展参数config_id约定配置ID);
	// 选择器: customize_select 自定义选择器(扩展参数中约定对应的值和名称，默认值和名称一致);
	//        org_bind_id 组织成员ID; org_bind_ids 一组组织成员ID列; org_group_id 组织成员分组ID; org_group_ids 一组组织成员分组ID;
	//        erp_product_id ERP产品ID; erp_product_ids 一组ERP产品ID列;
	//        erp_company_id ERP公司ID(扩展参数company_type约定公司类型);
	//        mall_core_product_id 商城产品ID; mall_core_product_ids 一组商城产品ID;
	ComponentType string `db:"component_type" json:"componentType"`
	//组件名称
	Name string `db:"name" json:"name"`
	//帮助描述
	HelpDes string `db:"help_des" json:"helpDes"`
	//组件默认值
	Val string `db:"val" json:"val"`
	//验证用的正则表达式
	CheckVal string `db:"check_val" json:"checkVal"`
	//是否必填
	IsRequire bool `db:"is_require" json:"isRequire"`
	//扩展参数
	// open_analysis_count 是否启动对发生次数的统计;
	// open_analysis_sum 是否启动对数据的统计，仅支持number_int/number_float/number_price组件;
	// open_analysis_avg 是否启动对数据的平均数统计，仅支持number_int/number_float/number_price/number_p组件;
	// open_analysis_sort 是否启动对数据的排名统计，仅支持number_int/number_float/number_price/number_p组件;
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// Value sql底层处理器
func (t FieldsComponentDefine) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsComponentDefine) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
