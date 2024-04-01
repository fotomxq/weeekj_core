package RestaurantWeeklyRecipe

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	"time"
)

// GetWeeklyRecipeMargeByID 获取打包数据
func GetWeeklyRecipeMargeByID(args *ArgsGetWeeklyRecipeByID) (headData FieldsWeeklyRecipe, itemList []FieldsWeeklyRecipeItem, err error) {
	headData, err = GetWeeklyRecipeByID(args)
	if err != nil {
		return
	}
	itemList, _, _ = GetWeeklyRecipeItemList(&ArgsGetWeeklyRecipeItemList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  999,
			Sort: "id",
			Desc: false,
		},
		OrgID:          -1,
		StoreID:        -1,
		WeeklyRecipeID: headData.ID,
		RecipeID:       -1,
		IsRemove:       false,
		Search:         "",
	})
	return
}

// ArgsGetWeeklyRecipeMargeList 获取指定时间段内的菜谱参数
type ArgsGetWeeklyRecipeMargeList struct {
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id"`
	//开始日期
	// 例如: 2021-01-01
	StartDate string `json:"start_date"`
	//结束日期
	// 例如: 2021-01-07
	EndDate string `json:"end_date"`
}

// GetWeeklyRecipeMargeList 获取指定时间段内的菜谱
func GetWeeklyRecipeMargeList(args *ArgsGetWeeklyRecipeMargeList) ([]FieldsWeeklyRecipe, error) {
	//重构时间差
	startAt, err := CoreFilter.GetTimeCarbonByDefault(args.StartDate)
	if err != nil {
		return nil, err
	}
	endAt, err := CoreFilter.GetTimeCarbonByDefault(args.EndDate)
	if err != nil {
		return nil, err
	}
	//初始化数据
	stepAt := startAt.StartOfDay()
	var dataList []FieldsWeeklyRecipe
	//遍历获取数据
	for {
		//获取原始数据
		rawList, err := getWeeklyRecipeBetweenDate(args.OrgID, args.StoreID, CoreFilter.GetTimeToDefaultDate(stepAt.Time))
		if err != nil {
			break
		}
		//遍历数据
		for _, v := range rawList {
			dataList = append(dataList, v)
		}
		//下一天
		stepAt = stepAt.AddDay()
		//检查是否超过结束时间
		if stepAt.DiffInDaysWithAbs(endAt) == 0 {
			break
		}
		//如果超出检查上限，则拒绝
		if stepAt.DiffInDaysWithAbs(endAt) > 99 {
			return nil, errors.New("too many data")
		}
	}
	return dataList, nil
}

// DataWeeklyRecipeMargeMore 打包主数据
type DataWeeklyRecipeMargeMore struct {
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
	//开始日期
	// 例如：2021-01-01
	StartDate string `json:"startDate"`
	//结束日期
	// 例如：2021-01-01
	EndDate string `json:"endDate"`
	//数据
	// 两种方法提供，一种不含该此数据；另外一种会包含此数据
	DataList []DataWeeklyRecipeMargeMoreHeader `json:"dataList"`
}

type DataWeeklyRecipeMargeMoreHeader struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"1000" empty:"true"`
	// 用餐日期
	DiningDate time.Time `db:"dining_date" json:"diningDate"`
	//数据
	DataList []DataWeeklyRecipeMargeMoreHeaderDiningTime `json:"dataList"`
}

type DataWeeklyRecipeMargeMoreHeaderDiningTime struct {
	// 用餐时间
	//0 早餐; 1 午餐; 2 晚餐
	DiningTime int `db:"dining_time" json:"diningTime" check:"intThan0" empty:"true"`
	//数据
	DataList []DataWeeklyRecipeMargeMoreItem `json:"dataList"`
}

type DataWeeklyRecipeMargeMoreItem struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//菜品ID
	RecipeID int64 `db:"recipe_id" json:"recipeID" check:"id"`
	//菜品名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//售价
	Price int64 `db:"price" json:"price" check:"int64Than0"`
}
