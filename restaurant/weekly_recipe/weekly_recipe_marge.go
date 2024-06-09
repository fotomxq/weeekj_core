package RestaurantWeeklyRecipeMarge

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
)

// DataWeeklyRecipeMarge 聚合数据包
type DataWeeklyRecipeMarge struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//创建时间
	CreateAt string `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt string `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt string `db:"delete_at" json:"deleteAt" default:"0"`
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
	AuditAt string `db:"audit_at" json:"auditAt" index:"true"`
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
	//日数据
	DayList []DataGetWeeklyRecipeMargeDay `json:"dayList"`
}

type DataGetWeeklyRecipeMargeDay struct {
	// 用餐日期
	// 例如：20210101
	DiningDate int `db:"dining_date" json:"diningDate" index:"true"`
	//早餐
	Breakfast []DataGetWeeklyRecipeMargeDayItem `json:"breakfast"`
	//午餐
	Lunch []DataGetWeeklyRecipeMargeDayItem `json:"lunch"`
	//晚餐
	Dinner []DataGetWeeklyRecipeMargeDayItem `json:"dinner"`
}

type DataGetWeeklyRecipeMargeDayItem struct {
	//唯一ID
	ID int64 `db:"id" json:"id"`
	//菜品ID
	RecipeID int64 `db:"recipe_id" json:"recipeID" check:"id" index:"true"`
	//菜品名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//售价
	Price int64 `db:"price" json:"price" check:"int64Than0" empty:"true"`
	//数量
	RecipeCount int `db:"recipe_count" json:"recipeCount" check:"intThan0" empty:"true"`
	//单位
	Unit string `db:"unit" json:"unit" check:"des" min:"1" max:"300" empty:"true"`
	//单位ID
	UnitID int64 `db:"unit_id" json:"unitID" index:"true" check:"id" empty:"true"`
	//上周同时间段是否出现过
	IsRepeat bool `db:"is_repeat" json:"isRepeat" default:"false"`
	//上周早中晚是否全部出现过
	IsRepeatAll bool `db:"is_repeat_all" json:"isRepeatAll" default:"false"`
}

//GetWeeklyRecipeMarge 获取周菜谱聚合数据
/**
1. 底层存储改为分表结构
2. 对外输出采用rawData结构
*/
func GetWeeklyRecipeMarge(weeklyRecipeID int64) (data DataWeeklyRecipeMarge, err error) {
	//获取周菜谱数据
	weeklyRecipeData := getWeeklyRecipeByID(weeklyRecipeID)
	if weeklyRecipeData.ID < 1 {
		err = errors.New("no data")
		return
	}
	//根据本周数据，获取上周数据
	var beforeData []DataGetWeeklyRecipeMargeDay
	var beforeList []FieldsWeeklyRecipeDay
	_ = weeklyRecipeDayDB.Select().SetFieldsSort([]string{"create_at"}).SetFieldsAll().SetIDQuery("org_id", weeklyRecipeData.OrgID).SetIDQuery("store_id", weeklyRecipeData.StoreID).SetIntQuery("audit_status", 1).SetIDQuery("recipe_type_id", weeklyRecipeData.RecipeTypeID).SetDeleteQuery("delete_at", false).SetPages(CoreSQL2.ArgsPages{
		Page: 1,
		Max:  1,
		Sort: "create_at",
		Desc: true,
	}).Result(&beforeList)
	if len(beforeList) > 0 {
		beforeData, _ = GetWeeklyRecipeBeforeMarge(beforeList[0].ID)
	}
	//获取菜谱下所有数据
	var rawList []FieldsWeeklyRecipeDay
	err = weeklyRecipeDayDB.Select().SetFieldsAll().SetIDQuery("weekly_recipe_id", weeklyRecipeID).SetDeleteQuery("delete_at", false).Result(&rawList)
	if err != nil {
		return
	}
	var rawList2 []FieldsWeeklyRecipeChild
	_ = weeklyRecipeChildDB.Select().SetFieldsAll().SetIDQuery("weekly_recipe_id", weeklyRecipeID).SetDeleteQuery("delete_at", false).Result(&rawList2)
	//整理数据
	data = DataWeeklyRecipeMarge{
		ID:              weeklyRecipeData.ID,
		CreateAt:        CoreFilter.GetTimeToDefaultTime(weeklyRecipeData.CreateAt),
		UpdateAt:        CoreFilter.GetTimeToDefaultTime(weeklyRecipeData.UpdateAt),
		DeleteAt:        CoreFilter.GetTimeToDefaultTime(weeklyRecipeData.DeleteAt),
		OrgID:           weeklyRecipeData.OrgID,
		StoreID:         weeklyRecipeData.StoreID,
		SubmitOrgBindID: weeklyRecipeData.SubmitOrgBindID,
		SubmitUserID:    weeklyRecipeData.SubmitUserID,
		SubmitUserName:  weeklyRecipeData.SubmitUserName,
		AuditAt:         CoreFilter.GetTimeToDefaultTime(weeklyRecipeData.AuditAt),
		AuditStatus:     weeklyRecipeData.AuditStatus,
		AuditOrgBindID:  weeklyRecipeData.AuditOrgBindID,
		AuditUserID:     weeklyRecipeData.AuditUserID,
		AuditUserName:   weeklyRecipeData.AuditUserName,
		Name:            weeklyRecipeData.Name,
		Remark:          weeklyRecipeData.Remark,
		RecipeTypeID:    weeklyRecipeData.RecipeTypeID,
		RecipeTypeName:  weeklyRecipeData.RecipeTypeName,
		DayList:         []DataGetWeeklyRecipeMargeDay{},
	}
	for k := 0; k < len(rawList); k++ {
		v := rawList[k]
		//构建数据
		appendData := DataGetWeeklyRecipeMargeDay{
			DiningDate: v.DiningDate,
			Breakfast:  []DataGetWeeklyRecipeMargeDayItem{},
			Lunch:      []DataGetWeeklyRecipeMargeDayItem{},
			Dinner:     []DataGetWeeklyRecipeMargeDayItem{},
		}
		//找出子集合
		for k2 := 0; k2 < len(rawList2); k2++ {
			v2 := rawList2[k2]
			//构建子数据
			if v.ID != v2.WeeklyRecipeDayID || v.WeeklyRecipeID != v2.WeeklyRecipeID {
				continue
			}
			//检查上周是否重复出现
			var isRepeat, isRepeatAll bool
			for _, v3 := range beforeData {
				if v3.DiningDate != v.DiningDate {
					continue
				}
				for _, v4 := range v3.Breakfast {
					if v4.RecipeID != v2.RecipeID {
						continue
					}
					isRepeatAll = true
					break
				}
				if isRepeatAll {
					break
				}
				for _, v4 := range v3.Lunch {
					if v4.RecipeID != v2.RecipeID {
						continue
					}
					isRepeatAll = true
					break
				}
				if isRepeatAll {
					break
				}
				for _, v4 := range v3.Dinner {
					if v4.RecipeID != v2.RecipeID {
						continue
					}
					isRepeatAll = true
					break
				}
				if isRepeatAll {
					break
				}
			}
			//找到对应的数据
			switch v2.DayType {
			case 1:
				for _, v3 := range beforeData {
					if v3.DiningDate != v.DiningDate {
						continue
					}
					for _, v4 := range v3.Breakfast {
						if v4.RecipeID != v2.RecipeID {
							continue
						}
						isRepeat = true
						break
					}
					if isRepeat {
						break
					}
				}
				appendData.Breakfast = append(appendData.Lunch, DataGetWeeklyRecipeMargeDayItem{
					ID:          v2.ID,
					RecipeID:    v2.RecipeID,
					Name:        v2.Name,
					Price:       v2.Price,
					RecipeCount: v2.RecipeCount,
					Unit:        v2.Unit,
					UnitID:      v2.UnitID,
					IsRepeat:    isRepeat,
					IsRepeatAll: isRepeatAll,
				})
			case 2:
				for _, v3 := range beforeData {
					if v3.DiningDate != v.DiningDate {
						continue
					}
					for _, v4 := range v3.Lunch {
						if v4.RecipeID != v2.RecipeID {
							continue
						}
						isRepeat = true
						break
					}
					if isRepeat {
						break
					}
				}
				appendData.Lunch = append(appendData.Lunch, DataGetWeeklyRecipeMargeDayItem{
					ID:          v2.ID,
					RecipeID:    v2.RecipeID,
					Name:        v2.Name,
					Price:       v2.Price,
					RecipeCount: v2.RecipeCount,
					Unit:        v2.Unit,
					UnitID:      v2.UnitID,
					IsRepeat:    isRepeat,
					IsRepeatAll: isRepeatAll,
				})
			case 3:
				for _, v3 := range beforeData {
					if v3.DiningDate != v.DiningDate {
						continue
					}
					for _, v4 := range v3.Dinner {
						if v4.RecipeID != v2.RecipeID {
							continue
						}
						isRepeat = true
						break
					}
					if isRepeat {
						break
					}
				}
				appendData.Dinner = append(appendData.Dinner, DataGetWeeklyRecipeMargeDayItem{
					ID:          v2.ID,
					RecipeID:    v2.RecipeID,
					Name:        v2.Name,
					Price:       v2.Price,
					RecipeCount: v2.RecipeCount,
					Unit:        v2.Unit,
					UnitID:      v2.UnitID,
					IsRepeat:    isRepeat,
					IsRepeatAll: isRepeatAll,
				})
			}
		}
		//追加数据
		data.DayList = append(data.DayList, appendData)
	}
	//反馈
	return
}

func GetWeeklyRecipeBeforeMarge(weeklyRecipeID int64) (dayList []DataGetWeeklyRecipeMargeDay, err error) {
	//获取菜谱下所有数据
	var rawList []FieldsWeeklyRecipeDay
	err = weeklyRecipeDayDB.Select().SetFieldsAll().SetIDQuery("weekly_recipe_id", weeklyRecipeID).SetDeleteQuery("delete_at", false).Result(&rawList)
	if err != nil {
		return
	}
	var rawList2 []FieldsWeeklyRecipeChild
	_ = weeklyRecipeChildDB.Select().SetFieldsAll().SetIDQuery("weekly_recipe_id", weeklyRecipeID).SetDeleteQuery("delete_at", false).Result(&rawList2)

	for k := 0; k < len(rawList); k++ {
		v := rawList[k]
		//构建数据
		appendData := DataGetWeeklyRecipeMargeDay{
			DiningDate: v.DiningDate,
			Breakfast:  []DataGetWeeklyRecipeMargeDayItem{},
			Lunch:      []DataGetWeeklyRecipeMargeDayItem{},
			Dinner:     []DataGetWeeklyRecipeMargeDayItem{},
		}
		//找出子集合
		for k2 := 0; k2 < len(rawList2); k2++ {
			v2 := rawList2[k2]
			//构建子数据
			if v.ID != v2.WeeklyRecipeDayID || v.WeeklyRecipeID != v2.WeeklyRecipeID {
				continue
			}
			//找到对应的数据
			switch v2.DayType {
			case 1:
				appendData.Breakfast = append(appendData.Lunch, DataGetWeeklyRecipeMargeDayItem{
					RecipeID:    v2.RecipeID,
					Name:        v2.Name,
					Price:       v2.Price,
					RecipeCount: v2.RecipeCount,
					Unit:        v2.Unit,
					UnitID:      v2.UnitID,
					IsRepeat:    false,
					IsRepeatAll: false,
				})
			case 2:
				appendData.Lunch = append(appendData.Lunch, DataGetWeeklyRecipeMargeDayItem{
					RecipeID:    v2.RecipeID,
					Name:        v2.Name,
					Price:       v2.Price,
					RecipeCount: v2.RecipeCount,
					Unit:        v2.Unit,
					UnitID:      v2.UnitID,
					IsRepeat:    false,
					IsRepeatAll: false,
				})
			case 3:
				appendData.Dinner = append(appendData.Dinner, DataGetWeeklyRecipeMargeDayItem{
					RecipeID:    v2.RecipeID,
					Name:        v2.Name,
					Price:       v2.Price,
					RecipeCount: v2.RecipeCount,
					Unit:        v2.Unit,
					UnitID:      v2.UnitID,
					IsRepeat:    false,
					IsRepeatAll: false,
				})
			}
		}
		//追加数据
		dayList = append(dayList, appendData)
	}
	//反馈
	return
}

// ArgsCreateWeeklyRecipeMarge 聚合创建数据参数
type ArgsCreateWeeklyRecipeMarge struct {
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
	//菜谱类型ID
	RecipeTypeID int64 `db:"recipe_type_id" json:"recipeTypeID" check:"id" index:"true"`
	//日数据
	DayList []DataGetWeeklyRecipeMargeDay `json:"dayList"`
}

// CreateWeeklyRecipeMarge 聚合创建数据
func CreateWeeklyRecipeMarge(args *ArgsCreateWeeklyRecipeMarge) (weeklyRecipeID int64, err error) {
	//创建周数据
	weeklyRecipeID, err = CreateWeeklyRecipe(&ArgsCreateWeeklyRecipe{
		OrgID:           args.OrgID,
		StoreID:         args.StoreID,
		SubmitOrgBindID: args.SubmitOrgBindID,
		SubmitUserID:    args.SubmitUserID,
		SubmitUserName:  args.SubmitUserName,
		Name:            args.Name,
		Remark:          args.Remark,
		RecipeTypeID:    args.RecipeTypeID,
	})
	if err != nil {
		return
	}
	//创建日数据
	_, err = SetWeeklyRecipeDay(weeklyRecipeID, args.DayList)
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateWeeklyRecipeMarge 聚合修改数据参数
type ArgsUpdateWeeklyRecipeMarge struct {
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
	//菜谱类型ID
	RecipeTypeID int64 `db:"recipe_type_id" json:"recipeTypeID" check:"id" index:"true"`
	//日数据
	DayList []DataGetWeeklyRecipeMargeDay `json:"dayList"`
}

// UpdateWeeklyRecipeMarge 聚合修改数据
func UpdateWeeklyRecipeMarge(args *ArgsUpdateWeeklyRecipeMarge) (err error) {
	//修改周数据
	err = UpdateWeeklyRecipe(&ArgsUpdateWeeklyRecipe{
		ID:           args.ID,
		OrgID:        args.OrgID,
		StoreID:      args.StoreID,
		Name:         args.Name,
		Remark:       args.Remark,
		RecipeTypeID: args.RecipeTypeID,
	})
	if err != nil {
		return
	}
	//修改日数据
	_, err = SetWeeklyRecipeDay(args.ID, args.DayList)
	if err != nil {
		return
	}
	//反馈
	return
}
