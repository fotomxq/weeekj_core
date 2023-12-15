package ToolsHelpContent

import (
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	ClassTag "github.com/fotomxq/weeekj_core/v5/class/tag"
)

//帮助内容系统
// 提供自助内容服务的帮助系统
// 1、界面提供帮助说明
// 2、提供简单的用户互动帮助服务
// 3、可以直接和系统开发商发起沟通服务

var (
	//Sort 分类系统
	Sort = ClassSort.Sort{
		SortTableName: "tools_help_content_sort",
	}
	//Tag 标签系统
	Tag = ClassTag.Tag{
		TagTableName: "tools_help_content_tag",
	}
)
