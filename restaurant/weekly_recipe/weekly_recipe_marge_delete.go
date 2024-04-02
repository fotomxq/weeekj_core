package RestaurantWeeklyRecipeMarge

// ArgsDeleteWeeklyRecipe 删除WeeklyRecipe参数
type ArgsDeleteWeeklyRecipe struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true"`
}

// DeleteWeeklyRecipe 删除WeeklyRecipe
func DeleteWeeklyRecipe(args *ArgsDeleteWeeklyRecipe) (err error) {
	//删除数据
	err = weeklyRecipeMargeDB.Delete().NeedSoft(true).AddWhereID(args.ID).AddWhereOrgID(args.OrgID).SetWhereOrThan("store_id", args.StoreID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteWeeklyRecipeCache(args.ID)
	//反馈
	return
}
