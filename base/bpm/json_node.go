package BaseBPM

/**
采用gjson模块解析json数据
eg:
{
	"name": "bpm_name", //bpm名称
	"version": "1.0", //bpm版本
	"bpm_id": 1, //bpm ID
	"theme_category_id": 0, //bpm主题分类ID
	"theme_category_name": "default", //bpm主题分类名称
	"theme_id": 0, //bpm主题ID
	"theme_name": "default", //bpm主题名称
	"nodes": [
		{
			"num": 1, //节点序号，同时用于排序
			"id": "ABC", //节点ID，采用节点HASH值随机码自动生成
			"name": "start", //节点名称
			"type": "form", //节点类型，form/event/condition。form代表表单；condition代表判断条件
			"next": ["DEF", "GHI"], //下一节点ID列
			"prev": [], //上一节点ID列
			"form": [
				{
					"num": 1, //表单序号，同时用于排序
					"id": "form_1", //表单字段ID，本nodes及json内必须唯一
					"name": "form_1", //表单槽名称
					"slot_id": 0, //表单槽ID，对应slot插槽
					"slot_value": "default", //表单槽默认值，值类型根据插槽定义决定，可以为数字、浮点数、字符串等内容
					"params": "", //扩展参数，可以为空；根据实际需求填入，数据类型任意，包含字符串、数字、浮点数、布尔值、json等
					"value": "", //最终存储的值，可以为空；根据实际需求填入，数据类型任意，包含字符串、数字、浮点数、布尔值、json等
					"data_from": "", //数据来源，对应表名称
					"data_field": "", //数据字段，对应表字段名称
					"data_id": 0, //数据ID，对应表数据ID
				},
			], //如果为form类型，则包含该内容，代表包含的表单内容
			"event": [
				{
					"num": 1, //事件序号，同时用于排序
					"id": "event_1", //事件ID
					"name": "event_1", //事件名称
					"event_id": 0, //事件ID，对应event事件的ID
					"fields": [
						{
							"field": "", //判断字段，必须是之前form内出现的值，对应表单字段ID
						},
					], //传递参数值
					"params": "", //扩展参数，可以为空；根据实际需求填入，数据类型任意，包含字符串、数字、浮点数、布尔值、json等
				},
			], //如果为event类型，则包含该内容，代表包含的事件内容
			"condition": [
				{
					"type": "", //判断类型 >< 不等于; = 等于; >= 大于等于; <= 小于等于; > 大于; < 小于; in 在范围内; not in 不在范围内; like 包含; not like 不包含; is null 为空; is not null 不为空
					"field": "", //判断字段，必须是之前form内出现的值，对应表单字段ID。系统会判断该字段的最终值
					"params": "", //扩展参数，可以为空；根据实际需求填入，数据类型任意，包含字符串、数字、浮点数、布尔值、json等
					"value": false, //最终结果，必须是bool值
				},
			], //如果为condition类型，则包含该内容，代表包含的条件内容
		},
	], //节点列
}
*/
