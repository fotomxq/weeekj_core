package RestaurantWeeklyRecipeMarge

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	RestaurantRecipe "github.com/fotomxq/weeekj_core/v5/restaurant/recipe"
	"time"
)

// ArgsUpdateWeeklyRecipe 修改WeeklyRecipe参数
type ArgsUpdateWeeklyRecipe struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"1000" empty:"true"`
	//数据包
	RawData FieldsWeeklyRecipeHeaders `db:"raw_data" json:"rawData"`
}

// UpdateWeeklyRecipe 修改WeeklyRecipe
func UpdateWeeklyRecipe(args *ArgsUpdateWeeklyRecipe) (err error) {
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
	//更新数据
	err = weeklyRecipeMargeDB.Update().SetFields([]string{"name", "remark", "raw_data"}).NeedUpdateTime().AddWhereID(args.ID).AddWhereOrgID(args.OrgID).SetWhereOrThan("store_id", args.StoreID).NamedExec(map[string]any{
		"name":     args.Name,
		"remark":   args.Remark,
		"raw_data": newRawData,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteWeeklyRecipeCache(args.ID)
	//反馈
	return
}

// ArgsAuditWeeklyRecipe 审核每周菜谱上报参数
type ArgsAuditWeeklyRecipe struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//当前组织ID
	// 用于验证数据是否属于当前组织
	RawOrgID int64 `db:"raw_org_id" json:"rawOrgID" check:"id" empty:"true"`
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
}

// AuditWeeklyRecipe 审核每周菜谱上报
func AuditWeeklyRecipe(args *ArgsAuditWeeklyRecipe) (err error) {
	var auditAt time.Time
	if args.AuditStatus == 1 {
		auditAt = CoreFilter.GetNowTime()
	}
	if args.AuditStatus == 2 {
		auditAt = time.Time{}
	}
	err = weeklyRecipeMargeDB.Update().SetFields([]string{"audit_at", "audit_status", "audit_org_bind_id", "audit_user_id", "audit_user_name"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"audit_at":          auditAt,
		"audit_status":      args.AuditStatus,
		"audit_org_bind_id": args.AuditOrgBindID,
		"audit_user_id":     args.AuditUserID,
		"audit_user_name":   args.AuditUserName,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteWeeklyRecipeCache(args.ID)
	//反馈
	return
}
