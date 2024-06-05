package RestaurantWeeklyRecipeMarge

import (
	"time"
)

// FieldsWeeklyRecipe 每周提交菜谱表头
type FieldsWeeklyRecipe struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" index:"true"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" index:"true"`
	//提交组织成员ID
	SubmitOrgBindID int64 `db:"submit_org_bind_id" json:"submitOrgBindID" check:"id" empty:"true" index:"true"`
	//提交用户ID
	// 与组织ID二选一，如果组织成员ID为空，则使用用户ID；如果组织ID不为空，则使用组织成员ID+用户ID
	SubmitUserID int64 `db:"submit_user_id" json:"submitUserID" check:"id" empty:"true" index:"true"`
	//提交人姓名
	SubmitUserName string `db:"submit_user_name" json:"submitUserName" check:"des" min:"1" max:"300" empty:"true"`
	//审核时间
	AuditAt time.Time `db:"audit_at" json:"auditAt" index:"true"`
	//审核状态
	// 0 未审核; 1 审核通过; 2 审核不通过
	AuditStatus int `db:"audit_status" json:"auditStatus" check:"intThan0" empty:"true" index:"true"`
	//审核人ID
	AuditOrgBindID int64 `db:"audit_org_bind_id" json:"auditOrgBindID" check:"id" empty:"true" index:"true"`
	//审核用户ID
	// 与组织ID二选一，如果组织成员ID为空，则使用用户ID；如果组织ID不为空，则使用组织成员ID+用户ID
	AuditUserID int64 `db:"audit_user_id" json:"auditUserID" check:"id" empty:"true" index:"true"`
	//审核人姓名
	AuditUserName string `db:"audit_user_name" json:"auditUserName" check:"des" min:"1" max:"300" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"1000" empty:"true"`
	//菜谱类型ID
	RecipeTypeID int64 `db:"recipe_type_id" json:"recipeTypeID" check:"id" index:"true"`
	//菜谱类型名称
	RecipeTypeName string `db:"recipe_type_name" json:"recipeTypeName" check:"des" min:"1" max:"300" empty:"true"`
}
