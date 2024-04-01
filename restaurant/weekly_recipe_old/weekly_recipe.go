package RestaurantWeeklyRecipe

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	RestaurantRecipe "github.com/fotomxq/weeekj_core/v5/restaurant/recipe"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetWeeklyRecipeList 获取WeeklyRecipe列表参数
type ArgsGetWeeklyRecipeList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分公司ID组
	OrgIDs []int64 `db:"org_ids" json:"orgIDs" check:"ids" empty:"true"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true"`
	//门店ID列
	StoreIDs []int64 `db:"store_ids" json:"storeIDs" check:"ids" empty:"true"`
	//提交组织成员ID
	SubmitOrgBindID int64 `db:"submit_org_bind_id" json:"submitOrgBindID" check:"id" empty:"true"`
	//提交用户ID
	// 与组织ID二选一，如果组织成员ID为空，则使用用户ID；如果组织ID不为空，则使用组织成员ID+用户ID
	SubmitUserID int64 `db:"submit_user_id" json:"submitUserID" check:"id" empty:"true"`
	//审核状态
	// 0 未审核; 1 审核通过; 2 审核不通过
	AuditStatus int `db:"audit_status" json:"auditStatus" check:"intThan0" empty:"true"`
	//审核人ID
	AuditOrgBindID int64 `db:"audit_org_bind_id" json:"auditOrgBindID" check:"id" empty:"true"`
	//审核用户ID
	// 与组织ID二选一，如果组织成员ID为空，则使用用户ID；如果组织ID不为空，则使用组织成员ID+用户ID
	AuditUserID int64 `db:"audit_user_id" json:"auditUserID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetWeeklyRecipeList 获取WeeklyRecipe列表
func GetWeeklyRecipeList(args *ArgsGetWeeklyRecipeList) (dataList []FieldsWeeklyRecipe, dataCount int64, err error) {
	dataCount, err = weeklyRecipeDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "dining_date", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("org_id", args.OrgID).SetIDsQuery("org_ids", args.OrgIDs).SetIDQuery("store_id", args.StoreID).SetIDsQuery("store_id", args.StoreIDs).SetIDQuery("submit_org_bind_id", args.SubmitOrgBindID).SetIDQuery("submit_user_id", args.SubmitUserID).SetIntQuery("audit_status", args.AuditStatus).SetIDQuery("audit_org_bind_id", args.AuditOrgBindID).SetIDQuery("audit_user_id", args.AuditUserID).SetSearchQuery([]string{"name", "remark"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getWeeklyRecipeByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetWeeklyRecipeByID 获取WeeklyRecipe数据包参数
type ArgsGetWeeklyRecipeByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true"`
}

// GetWeeklyRecipeByID 获取WeeklyRecipe数
func GetWeeklyRecipeByID(args *ArgsGetWeeklyRecipeByID) (data FieldsWeeklyRecipe, err error) {
	data = getWeeklyRecipeByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) || !CoreFilter.EqID2(args.StoreID, data.StoreID) {
		err = errors.New("no data")
		return
	}
	return
}

// GetWeeklyRecipeNameByID 获取菜品名称
func GetWeeklyRecipeNameByID(id int64) (name string) {
	data := getWeeklyRecipeByID(id)
	if data.ID < 1 {
		return
	}
	return data.Name
}

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
	// 用餐时间
	//0 早餐; 1 午餐; 2 晚餐
	DiningTime int `db:"dining_time" json:"diningTime" check:"intThan0" empty:"true"`
	// 用餐日期
	DiningDate time.Time `db:"dining_date" json:"diningDate"`
}

// CreateWeeklyRecipe 创建WeeklyRecipe
func CreateWeeklyRecipe(args *ArgsCreateWeeklyRecipe) (id int64, err error) {
	//检查当前是否存在数据
	checkDataList, _ := getWeeklyRecipeBetweenDate(args.OrgID, args.StoreID, CoreFilter.GetTimeToDefaultDate(CoreFilter.GetCarbonByTime(args.DiningDate).Time))
	if len(checkDataList) > 0 {
		for _, v := range checkDataList {
			if v.DiningTime == args.DiningTime {
				err = errors.New("data already exists")
				return
			}
		}
	}
	//创建数据
	id, err = weeklyRecipeDB.Insert().SetFields([]string{"org_id", "store_id", "submit_org_bind_id", "submit_user_id", "submit_user_name", "audit_at", "audit_status", "audit_org_bind_id", "audit_user_id", "audit_user_name", "name", "remark", "dining_time", "dining_date"}).Add(map[string]any{
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
		"dining_time":        args.DiningTime,
		"dining_date":        args.DiningDate,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsCreateWeeklyRecipeMargeByDay 打包创建某一天的菜谱参数
type ArgsCreateWeeklyRecipeMargeByDay struct {
	//头部关键信息
	HeaderData ArgsCreateWeeklyRecipe `json:"headerData"`
	//行信息列
	RowData []ArgsCreateWeeklyRecipeMargeByDayItem `json:"rowData"`
}

type ArgsCreateWeeklyRecipeMargeByDayItem struct {
	//菜品ID
	RecipeID int64 `db:"recipe_id" json:"recipeID" check:"id"`
	//售价
	Price int64 `db:"price" json:"price" check:"int64Than0"`
}

// CreateWeeklyRecipeMargeByDay 打包创建某一天的菜谱
func CreateWeeklyRecipeMargeByDay(args *ArgsCreateWeeklyRecipeMargeByDay) (id int64, err error) {
	//创建头部
	id, err = CreateWeeklyRecipe(&args.HeaderData)
	if err != nil {
		return
	}
	//批量创建行
	for _, v := range args.RowData {
		_, err = CreateWeeklyRecipeItem(&ArgsCreateWeeklyRecipeItem{
			OrgID:          args.HeaderData.OrgID,
			StoreID:        args.HeaderData.StoreID,
			WeeklyRecipeID: id,
			RecipeID:       v.RecipeID,
			Name:           RestaurantRecipe.GetRecipeNameByID(v.RecipeID),
			Price:          v.Price,
		})
		if err != nil {
			return
		}
	}
	//反馈
	return
}

// ArgsUpdateWeeklyRecipe 修改WeeklyRecipe参数
type ArgsUpdateWeeklyRecipe struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"1000" empty:"true"`
	// 用餐时间
	//0 早餐; 1 午餐; 2 晚餐
	DiningTime int `db:"dining_time" json:"diningTime" check:"intThan0" empty:"true"`
	// 用餐日期
	DiningDate time.Time `db:"dining_date" json:"diningDate"`
}

// UpdateWeeklyRecipe 修改WeeklyRecipe
func UpdateWeeklyRecipe(args *ArgsUpdateWeeklyRecipe) (err error) {
	//更新数据
	err = weeklyRecipeDB.Update().SetFields([]string{"name", "remark", "dining_time", "dining_date"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"name":        args.Name,
		"remark":      args.Remark,
		"dining_time": args.DiningTime,
		"dining_date": args.DiningDate,
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
	RawOrgID int64 `db:"raw_org_id" json:"rawOrgID" check:"id"`
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
	err = weeklyRecipeDB.Update().SetFields([]string{"audit_at", "audit_status", "audit_org_bind_id", "audit_user_id", "audit_user_name"}).NeedUpdateTime().AddWhereID(args.ID).SetWhereAnd("(raw_org_id = :raw_org_id OR :raw_org_id < 0)", map[string]any{
		"raw_org_id": args.RawOrgID,
	}).NamedExec(map[string]any{
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

// ArgsDeleteWeeklyRecipe 删除WeeklyRecipe参数
type ArgsDeleteWeeklyRecipe struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteWeeklyRecipe 删除WeeklyRecipe
func DeleteWeeklyRecipe(args *ArgsDeleteWeeklyRecipe) (err error) {
	//删除数据
	err = weeklyRecipeDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteWeeklyRecipeCache(args.ID)
	//反馈
	return
}

// getWeeklyRecipeBetweenDate 获取指定时间段内的菜谱主表
// checkDate 为要检查的日期，格式为 2021-01-01
func getWeeklyRecipeBetweenDate(orgID int64, storeID int64, checkDate string) (dataList []FieldsWeeklyRecipe, err error) {
	//识别日期
	checkDateCarbon := CoreFilter.GetNowTimeCarbon().Parse(fmt.Sprint(checkDate, " 00:00:00"))
	startAt := checkDateCarbon.StartOfDay()
	endAt := checkDateCarbon.EndOfDay()
	//获取数据
	var rawList []FieldsWeeklyRecipe
	_, err = weeklyRecipeDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "dining_date"}).SetPages(CoreSQL2.ArgsPages{
		Page: 1,
		Max:  10,
		Sort: "id",
		Desc: false,
	}).SetDeleteQuery("delete_at", false).SetTimeBetweenQuery("dining_date", startAt.Time, endAt.Time).SetIDQuery("org_id", orgID).SetIDQuery("store_id", storeID).SelectList("").ResultAndCount(&rawList)
	if err != nil {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		dataList = append(dataList, getWeeklyRecipeByID(v.ID))
	}
	//反馈
	return
}

// getWeeklyRecipeByID 通过ID获取WeeklyRecipe数据包
func getWeeklyRecipeByID(id int64) (data FieldsWeeklyRecipe) {
	cacheMark := getWeeklyRecipeCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := weeklyRecipeDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "org_id", "store_id", "submit_org_bind_id", "submit_user_id", "submit_user_name", "audit_at", "audit_status", "audit_org_bind_id", "audit_user_id", "audit_user_name", "name", "remark", "dining_time", "dining_date"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheWeeklyRecipeTime)
	return
}

// 缓冲
func getWeeklyRecipeCacheMark(id int64) string {
	return fmt.Sprint("restaurant:weekly_recipe:id.", id)
}

func deleteWeeklyRecipeCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getWeeklyRecipeCacheMark(id))
}
