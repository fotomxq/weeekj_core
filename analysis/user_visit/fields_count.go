package AnalysisUserVisit

import "time"

//FieldsCount 用户访问统计
// 记录用户总的进入人次
type FieldsCount struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	// 每隔1小时统计一次
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// 如果存在数据，则表明该数据隶属于指定组织
	// 组织依可查看该数据
	OrgID int64 `db:"org_id" json:"orgID"`
	//行为类型
	// 0 注册新用户; 1 登陆人次(活跃度)
	// 2 购物车行为次数; 3 实际购物行为次数; 4 进入下单页面次数
	// 5 手机号注册; 6 微信小程序注册; 7 后台强制创建用户; 8 邮箱注册; 9 微信APP授权注册
	// 100 账户密码登陆; 101 手机短信登陆; 102 微信小程序登陆; 103 手机扫码登陆; 104 email登陆
	// 200 进入平台人次; 201 点击重要按钮次数; 202 点击购物按钮次数
	Mark int `db:"mark" json:"mark"`
	//统计数量
	Count int64 `db:"count" json:"count"`
}
