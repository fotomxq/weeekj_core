package BlogCore

import (
	ClassComment "gitee.com/weeekj/weeekj_core/v5/class/comment"
	ClassSort "gitee.com/weeekj/weeekj_core/v5/class/sort"
	ClassTag "gitee.com/weeekj/weeekj_core/v5/class/tag"
	"sync"
)

//内容服务核心
// 用于对外提供一整套的内容，方便外部访问管理。注意，提供方式需分为API方式和静态网站方式，API方式采用传统形式，例如JS-API获取数据并展示；静态网站方式为静态go模版形式，由后端渲染后输出即可。

var (
	//Sort 分类系统
	Sort = ClassSort.Sort{
		SortTableName: "blog_core_sort",
	}
	//Tag 标签系统
	Tag = ClassTag.Tag{
		TagTableName: "blog_core_tags",
	}
	//Comment 评论
	Comment = ClassComment.Comment{
		TableName:         "blog_core_comment",
		UserMoreComment:   false,
		UserEditComment:   false,
		UserDeleteComment: false,
		OrgDeleteComment:  false,
		System:            "blog_core_content",
	}
	//OpenSub 是否启动订阅
	OpenSub = false
	//key锁定机制
	makeKeyLock sync.Mutex
)

// Init 初始化
func Init() {
	//nats
	if OpenSub {
		subNats()
	}
}
