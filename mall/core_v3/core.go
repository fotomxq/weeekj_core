package MallCoreV3

import (
	ClassComment "gitee.com/weeekj/weeekj_core/v5/class/comment"
	ClassSort "gitee.com/weeekj/weeekj_core/v5/class/sort"
	ClassTag "gitee.com/weeekj/weeekj_core/v5/class/tag"
)

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
)
