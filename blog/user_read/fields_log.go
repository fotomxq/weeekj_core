package BlogUserRead

import "time"

//FieldsLog 阅读记录
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//子组织ID
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID"`
	//用户
	UserID int64 `db:"user_id" json:"userID"`
	//阅读渠道
	// 访问渠道的特征码
	FromMark string `db:"from_mark" json:"fromMark"`
	FromName string `db:"from_name" json:"fromName"`
	//姓名
	Name string `db:"name" json:"name"`
	//IP
	IP string `db:"ip" json:"ip"`
	//文章ID
	ContentID int64 `db:"content_id" json:"contentID"`
	//文章分类
	// 每个分类会构建一条统计记录
	SortID int64 `db:"sort_id" json:"sortID"`
	//离开时间
	LeaveAt time.Time `db:"leave_at" json:"leaveAt"`
	//阅读时间
	// 进入和离开时间的秒差值，如果离开没记录则不会记录本数据
	ReadTime int64 `db:"read_time" json:"readTime"`
}
