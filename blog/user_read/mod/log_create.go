package BlogUserReadMod

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
)

// ArgsCreateLog 添加日志参数
type ArgsCreateLog struct {
	//子组织ID
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//阅读渠道
	// 访问渠道的特征码
	FromMark string `db:"from_mark" json:"fromMark" check:"mark"`
	FromName string `db:"from_name" json:"fromName"`
	//姓名
	Name string `db:"name" json:"name" check:"name" empty:"true"`
	//IP
	IP string `db:"ip" json:"ip" check:"ip"`
	//文章ID
	ContentID int64 `db:"content_id" json:"contentID" check:"id"`
	//进入时间
	CreateAt string `db:"create_at" json:"createAt" check:"isoTime"`
	//离开时间
	LeaveAt string `db:"leave_at" json:"leaveAt" check:"isoTime" empty:"true"`
}

// CreateLog 添加日志
func CreateLog(args ArgsCreateLog) {
	CoreNats.PushDataNoErr("/blog/user_read/new", "", 0, "", args)
}
