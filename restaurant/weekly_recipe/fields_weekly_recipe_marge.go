package RestaurantWeeklyRecipeMarge

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// FieldsWeeklyRecipe 每周提交菜谱表头
type FieldsWeeklyRecipe struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id"`
	//提交组织成员ID
	SubmitOrgBindID int64 `db:"submit_org_bind_id" json:"submitOrgBindID" check:"id" empty:"true"`
	//提交用户ID
	// 与组织ID二选一，如果组织成员ID为空，则使用用户ID；如果组织ID不为空，则使用组织成员ID+用户ID
	SubmitUserID int64 `db:"submit_user_id" json:"submitUserID" check:"id" empty:"true"`
	//提交人姓名
	SubmitUserName string `db:"submit_user_name" json:"submitUserName" check:"des" min:"1" max:"300" empty:"true"`
	//审核时间
	AuditAt time.Time `db:"audit_at" json:"auditAt"`
	//审核状态
	// 0 未审核; 1 审核通过; 2 审核不通过
	AuditStatus int `db:"audit_status" json:"auditStatus" check:"intThan0" empty:"true"`
	//审核人ID
	AuditOrgBindID int64 `db:"audit_org_bind_id" json:"auditOrgBindID" check:"id" empty:"true"`
	//审核用户ID
	// 与组织ID二选一，如果组织成员ID为空，则使用用户ID；如果组织ID不为空，则使用组织成员ID+用户ID
	AuditUserID int64 `db:"audit_user_id" json:"auditUserID" check:"id" empty:"true"`
	//审核人姓名
	AuditUserName string `db:"audit_user_name" json:"auditUserName" check:"des" min:"1" max:"300" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"1000" empty:"true"`
	//数据包
	RawData FieldsWeeklyRecipeHeaders `db:"raw_data" json:"rawData"`
}

type FieldsWeeklyRecipeHeaders []FieldsWeeklyRecipeHeader

// Value sql底层处理器
func (t FieldsWeeklyRecipeHeaders) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsWeeklyRecipeHeaders) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsWeeklyRecipeHeader struct {
	// 用餐日期
	// 例如：2021-01-01
	DiningDate string `db:"dining_date" json:"diningDate"`
	//早餐
	Breakfast FieldsWeeklyRecipeItemList `json:"breakfast"`
	//中餐
	Lunch FieldsWeeklyRecipeItemList `json:"lunch"`
	//晚餐
	Dinner FieldsWeeklyRecipeItemList `json:"dinner"`
}

// Value sql底层处理器
func (t FieldsWeeklyRecipeHeader) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsWeeklyRecipeHeader) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsWeeklyRecipeItemList []FieldsWeeklyRecipeItem

// Value sql底层处理器
func (t FieldsWeeklyRecipeItemList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsWeeklyRecipeItemList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsWeeklyRecipeItem struct {
	//菜品ID
	RecipeID int64 `db:"recipe_id" json:"recipeID" check:"id"`
	//菜品名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//售价
	Price int64 `db:"price" json:"price" check:"int64Than0"`
	//数量
	Count int `db:"count" json:"count" check:"intThan0"`
	//单位
	Unit string `db:"unit" json:"unit" check:"des" min:"1" max:"300" empty:"true"`
}

// Value sql底层处理器
func (t FieldsWeeklyRecipeItem) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsWeeklyRecipeItem) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
