package FinanceReturnedMoney

import (
	"time"
)

// FieldsMarge 回款记录表
type FieldsMarge struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	// 每个月固定时间创建一条记录
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//回款公司ID
	CompanyID int64 `db:"company_id" json:"companyID"`
	//开始催款时间
	StartAt time.Time `db:"start_at" json:"startAt"`
	//应该回款金额
	NeedPrice int64 `db:"need_price" json:"needPrice"`
	//应该回款时间
	NeedAt time.Time `db:"need_at" json:"needAt"`
	//已经回款金额
	HavePrice int64 `db:"have_price" json:"havePrice"`
	//实际回款时间，最终回款时间
	HaveAt time.Time `db:"have_at" json:"haveAt"`
	//以下数据根据公司设置填入，该设计主要为保留历史数据记录
	//销售人员
	SellOrgBindID int64 `db:"sell_org_bind_id" json:"sellOrgBindID"`
	//催收负责人
	ReturnOrgBindID int64 `db:"return_org_bind_id" json:"returnOrgBindID"`
	//催款是否确认
	ReturnConfirmAt time.Time `db:"return_confirm_at" json:"returnConfirmAt"`
}
