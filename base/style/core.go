package BaseStyle

import (
	ClassSort "gitee.com/weeekj/weeekj_core/v5/class/sort"
	ClassTag "gitee.com/weeekj/weeekj_core/v5/class/tag"
)

/*
*
样式库设计
本模块用于页面展示时，对布局、风格进行的调整设计。
外部任意模块可引用该内容，作为页面表现形式。
*/
var (
	//ComponentSort 组件库分类
	ComponentSort = ClassSort.Sort{
		SortTableName: "core_style_component_sort",
	}
	//StyleSort 页面样式库分类
	StyleSort = ClassSort.Sort{
		SortTableName: "core_style_sort",
	}
	//ComponentTag 组件库标签
	ComponentTag = ClassTag.Tag{
		TagTableName: "core_style_component_tag",
	}
	//StyleTag 页面样式库标签
	StyleTag = ClassTag.Tag{
		TagTableName: "core_style_tag",
	}
)
