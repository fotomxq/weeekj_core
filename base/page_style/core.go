package BasePageStyle

import (
	ClassSort "gitee.com/weeekj/weeekj_core/v5/class/sort"
	ClassTag "gitee.com/weeekj/weeekj_core/v5/class/tag"
)

/**
页面样式设计模块
*/

var (
	//ComponentSort 组件分类
	ComponentSort = ClassSort.Sort{
		SortTableName: "core_page_style_component_sort",
	}
	//ComponentTags 组件标签
	ComponentTags = ClassTag.Tag{
		TagTableName: "core_page_style_component_tags",
	}
	//TemplateSort 模版分类
	TemplateSort = ClassSort.Sort{
		SortTableName: "core_page_style_template_sort",
	}
	//TemplateTags 模版标签
	TemplateTags = ClassTag.Tag{
		TagTableName: "core_page_style_template_tags",
	}
)
