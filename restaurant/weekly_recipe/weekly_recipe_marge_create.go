package RestaurantWeeklyRecipeMarge

import (
	RestaurantRecipe "github.com/fotomxq/weeekj_core/v5/restaurant/recipe"
	"time"
)

// ArgsCreateWeeklyRecipe 创建WeeklyRecipe参数
type ArgsCreateWeeklyRecipe struct {
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
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"1000" empty:"true"`
	//数据包
	RawData FieldsWeeklyRecipeHeaders `db:"raw_data" json:"rawData"`
}

// CreateWeeklyRecipe 创建WeeklyRecipe
func CreateWeeklyRecipe(args *ArgsCreateWeeklyRecipe) (id int64, err error) {
	//添加菜品名称
	var newRawData FieldsWeeklyRecipeHeaders
	if len(args.RawData) > 0 {
		for _, v := range args.RawData {
			var newV FieldsWeeklyRecipeHeader
			newV.DiningDate = v.DiningDate
			for _, v2 := range v.Breakfast {
				if v2.RecipeID > 0 {
					v2.Name = RestaurantRecipe.GetRecipeNameByID(v2.RecipeID)
				}
				newV.Breakfast = append(newV.Breakfast, v2)
			}
			for _, v2 := range v.Lunch {
				if v2.RecipeID > 0 {
					v2.Name = RestaurantRecipe.GetRecipeNameByID(v2.RecipeID)
				}
				newV.Lunch = append(newV.Lunch, v2)
			}
			for _, v2 := range v.Dinner {
				if v2.RecipeID > 0 {
					v2.Name = RestaurantRecipe.GetRecipeNameByID(v2.RecipeID)
				}
				newV.Dinner = append(newV.Dinner, v2)
			}
			newRawData = append(newRawData, newV)
		}
	}
	//创建数据
	id, err = weeklyRecipeMargeDB.Insert().SetFields([]string{"org_id", "store_id", "submit_org_bind_id", "submit_user_id", "submit_user_name", "audit_at", "audit_status", "audit_org_bind_id", "audit_user_id", "audit_user_name", "name", "remark", "raw_data"}).Add(map[string]any{
		"org_id":             args.OrgID,
		"store_id":           args.StoreID,
		"submit_org_bind_id": args.SubmitOrgBindID,
		"submit_user_id":     args.SubmitUserID,
		"submit_user_name":   args.SubmitUserName,
		"audit_at":           time.Time{},
		"audit_status":       0,
		"audit_org_bind_id":  0,
		"audit_user_id":      0,
		"audit_user_name":    "",
		"name":               args.Name,
		"remark":             args.Remark,
		"raw_data":           newRawData,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}
