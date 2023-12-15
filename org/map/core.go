package OrgMap

import (
	ClassComment "github.com/fotomxq/weeekj_core/v5/class/comment"
)

var (
	//OpenSub 是否启动订阅
	OpenSub = false
	//Comment 评论
	Comment = ClassComment.Comment{
		TableName:         "org_map_comment",
		UserMoreComment:   false,
		UserEditComment:   false,
		UserDeleteComment: false,
		OrgDeleteComment:  false,
		System:            "org_map_comment",
	}
	//OpenAnalysis 是否启动analysis
	OpenAnalysis = false
)

// Init 初始化
func Init() {
	//统计
	if OpenAnalysis {
	}
	//中间件处理
	if OpenSub {
		//消息列队
		subNats()
	}
}
