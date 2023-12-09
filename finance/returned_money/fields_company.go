package FinanceReturnedMoney

import "time"

type FieldsCompany struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//回款公司ID
	CompanyID int64 `db:"company_id" json:"companyID"`
	//催款间隔月份
	NeedTakeAddMonth int `db:"need_take_add_month" json:"needTakeAddMonth"`
	//标记的应回款开始日
	NeedTakeStartDay int `db:"need_take_start_day" json:"needTakeStartDay"`
	//每个月几号回款
	// 支持: 0 月初、1-28对应日、-1 月底模式
	NeedTakeDay int `db:"need_take_day" json:"needTakeDay"`
	//销售人员
	SellOrgBindID int64 `db:"sell_org_bind_id" json:"sellOrgBindID"`
	//催收负责人
	ReturnOrgBindID int64 `db:"return_org_bind_id" json:"returnOrgBindID"`
	//是否坏账
	IsBan bool `db:"is_ban" json:"isBan"`
	//催款路线
	ReturnLocation string `db:"return_location" json:"returnLocation"`
	//当前超期状态
	// 0 没有应收; 1 存在应收尚未逾期; 2 预留选项; 3 已经完成回款；4 存在逾期; 5 严重逾期30天; 6 违约60天; 7 违约90天; 8 违约365天
	ReturnStatus int `db:"return_status" json:"returnStatus"`
}
