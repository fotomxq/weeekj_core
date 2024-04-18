package EAMRepair

import "time"

// FieldsRepair 维修工单
type FieldsRepair struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	///////////////////////////////////////////////////////////////////////////////////////////////////////
	//创建时
	///////////////////////////////////////////////////////////////////////////////////////////////////////
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//提交人组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//提交人用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//维修地点
	LocationDes string `db:"location_des" json:"locationDes" check:"des" min:"1" max:"300"`
	//库存产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
	//维修数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//维修状态
	// 0: 未维修; 1: 维修中; 2: 维修完成; 3: 维修失败
	RepairStatus int `db:"repair_status" json:"repairStatus"`
	//EAM ID
	EAMID int64 `db:"eam_id" json:"eamID" check:"id"`
	//EAM编码
	EAMCode string `db:"eam_code" json:"eamCode" check:"des" min:"1" max:"50"`
	///////////////////////////////////////////////////////////////////////////////////////////////////////
	//审批时
	// TODO: 需采用BPM替代，审批流无法针对特殊内容进行指定
	///////////////////////////////////////////////////////////////////////////////////////////////////////
	//指派维修供应商ID
	// 集团审批时指定
	RepairCompanyID int64 `db:"repair_company_id" json:"repairCompanyID" check:"id" empty:"true"`
	//指派维修用户ID
	// 集团审批时指定
	RepairUserID int64 `db:"repair_user_id" json:"repairUserID" check:"id" empty:"true"`
	//指派维修人姓名
	// 集团审批时指定
	RepairUserName string `db:"repair_user_name" json:"repairUserName" check:"des" min:"1" max:"50" empty:"true"`
	//维修预估费用
	// 供应商填写的维修费用预估
	Estimate int64 `db:"estimate" json:"estimate" check:"int64Than0"`
	///////////////////////////////////////////////////////////////////////////////////////////////////////
	//维修后
	///////////////////////////////////////////////////////////////////////////////////////////////////////
	//维修时间
	RepairAt time.Time `db:"repair_at" json:"repairAt" empty:"true"`
	//维修备注
	RepairRemark string `db:"repair_remark" json:"repairRemark" check:"des" min:"1" max:"300" empty:"true"`
	//维修费用
	RepairTotal int64 `db:"repair_total" json:"repairTotal" check:"int64Than0" empty:"true"`
	//维修前价值
	// 维修前产品价值
	BeforePrice int64 `db:"before_price" json:"beforePrice" check:"int64Than0" empty:"true"`
	//维修后价值
	// 维修后产品价值
	AfterPrice int64 `db:"after_price" json:"afterPrice" check:"int64Than0" empty:"true"`
	///////////////////////////////////////////////////////////////////////////////////////////////////////
	//验收后
	///////////////////////////////////////////////////////////////////////////////////////////////////////
	//验收人组织成员ID
	AcceptBindID int64 `db:"accept_bind_id" json:"acceptBindID" check:"id" empty:"true"`
	//验收时间
	AcceptAt time.Time `db:"accept_at" json:"acceptAt" empty:"true"`
	//验收备注
	AcceptRemark string `db:"accept_remark" json:"acceptRemark" check:"des" min:"1" max:"300" empty:"true"`
}
