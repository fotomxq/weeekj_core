package MallCore

import (
	ClassComment "github.com/fotomxq/weeekj_core/v5/class/comment"
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	ClassTag "github.com/fotomxq/weeekj_core/v5/class/tag"
)

//商城核心模块
// 记录商品基本信息、分类、标签信息

var (
	//Sort 分类
	Sort = ClassSort.Sort{
		SortTableName: "mall_core_sort",
	}
	//Tags 标签
	Tags = ClassTag.Tag{
		TagTableName: "mall_core_tags",
	}
	//Comment 评论
	Comment = ClassComment.Comment{
		TableName:         "mall_core_comment",
		UserMoreComment:   false,
		UserEditComment:   false,
		UserDeleteComment: false,
		OrgDeleteComment:  false,
		System:            "mall_core_product",
	}
	//OpenSub 是否启动订阅
	OpenSub = false
)

// Init 初始化
func Init() {
	if OpenSub {
		//消息列队
		subNats()
	}
}
