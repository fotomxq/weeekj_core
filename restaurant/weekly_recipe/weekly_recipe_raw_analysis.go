package RestaurantWeeklyRecipeMarge

import CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"

// ArgsGetRawAnalysis 周分化单统计数据参数
type ArgsGetRawAnalysis struct {
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true" index:"true"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true" index:"true"`
	//菜谱类型ID
	RecipeTypeID int64 `db:"recipe_type_id" json:"recipeTypeID" check:"id" empty:"true" index:"true"`
	//时间范围
	BetweenAt CoreSQL2.ArgsTimeBetween `json:"betweenAt"`
}

// DataGetRawAnalysis 周分化单统计数据
type DataGetRawAnalysis struct {
	//原材料ID
	MaterialID int64 `db:"material_id" json:"materialID" check:"id" empty:"true" index:"true"`
	//原材料名称
	MaterialName string `db:"material_name" json:"materialName" check:"des" min:"1" max:"300" empty:"true"`
	//用量
	UseCount float64 `db:"use_count" json:"useCount" check:"intThan0"`
	//使用次数
	CountUse int `db:"count_use" json:"countUse" check:"intThan0"`
	//单价均价
	PriceAvg float64 `db:"price_avg" json:"priceAvg"`
	//合计总价
	PriceTotal float64 `db:"price_total" json:"priceTotal"`
}

// GetRawAnalysis 周分化单统计数据
func GetRawAnalysis(args *ArgsGetRawAnalysis) (dataList []DataGetRawAnalysis, err error) {
	//时间周期
	var betweenAt CoreSQL2.FieldsTimeBetween
	betweenAt, err = args.BetweenAt.GetFields()
	if err != nil {
		return
	}
	//获取数据
	err = weeklyRecipeRawDB.DB.GetPostgresql().Select(&dataList, "select material_id, max(material_name) as material_name, sum(use_count) as use_count, count(id) as count_use, avg(price) as price_avg, sum(total_price) as price_total from restaurant_weekly_recipe_raw where dining_date >= $1 and dining_date <= $2 and ($3 < 0 or org_id = $3) and ($4 < 0 or store_id = $4) and ($5 < 0 or recipe_type_id = $5) group by material_id order by use_count desc;", betweenAt.MinTime.Format("20060102"), betweenAt.MaxTime.Format("20060102"), args.OrgID, args.StoreID, args.RecipeTypeID)
	if err != nil {
		return
	}
	//反馈
	return
}
