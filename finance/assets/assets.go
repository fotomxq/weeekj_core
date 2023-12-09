package FinanceAssets

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

//资产数据

// ArgsGetAssetsList 获取资产列表参数
type ArgsGetAssetsList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 如果留空，则说明该资产被转移给组织自身
	UserID int64 `db:"user_id" json:"userID"`
	//资产产品
	ProductID int64 `db:"product_id" json:"productID"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove"`
}

// GetAssetsList 获取资产列表
func GetAssetsList(args *ArgsGetAssetsList) (dataList []FieldsAssets, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := CoreSQL.GetDeleteSQL(args.IsRemove, "")
	if args.OrgID > 0 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > 0 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.ProductID > 0 {
		where = where + " AND product_id = :product_id"
		maps["product_id"] = args.ProductID
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"finance_assets",
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, user_id, product_id, count FROM finance_assets WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "count"},
	)
	return
}

// ArgsGetAssetsByID 获取指定资产数据参数
type ArgsGetAssetsByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 用于检查
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetAssetsByID 获取指定资产数据
func GetAssetsByID(args *ArgsGetAssetsByID) (data FieldsAssets, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, user_id, product_id, count FROM finance_assets WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND "+CoreSQL.GetDeleteSQL(false, ""), args.ID, args.OrgID)
	return
}

// ArgsSetAssets 修改资产数据参数
type ArgsSetAssets struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//实际操作人，组织绑定成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//用户ID
	// 如果留空，则说明该资产被转移给组织自身
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//资产产品
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//增减数量
	Count int64 `db:"count" json:"count"`
	//变动原因
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
}

// SetAssets 修改资产数据
func SetAssets(args *ArgsSetAssets) (data FieldsAssets, err error) {
	//获取商品
	err = checkProductAndOrg(args.ProductID, args.OrgID)
	if err != nil {
		err = errors.New("not find product, " + err.Error())
		return
	}
	//获取存储数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, user_id, product_id, count FROM finance_assets WHERE org_id = $1 AND user_id = $2 AND product_id = $3 AND "+CoreSQL.GetDeleteSQL(false, ""), args.OrgID, args.UserID, args.ProductID)
	if err == nil {
		data.Count += args.Count
		if data.Count < 0 {
			err = errors.New("count less 1")
			return
		}
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE finance_assets SET update_at = NOW(), count = count + :count WHERE id = :id", map[string]interface{}{
			"id":    data.ID,
			"count": args.Count,
		})
	} else {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "finance_assets", "INSERT INTO finance_assets (org_id, user_id, product_id, count) VALUES (:org_id, :user_id, :product_id, :count)", args, &data)
	}
	if err == nil {
		//修改产品总数
		err = updateProductCountInc(&argsUpdateProductCountInc{
			ID:    data.ProductID,
			Count: args.Count,
		})
		if err != nil {
			return
		}
		//新增日志
		err = appendLog(&argsAppendLog{
			OrgID:     data.OrgID,
			BindID:    args.BindID,
			UserID:    data.UserID,
			ProductID: data.ProductID,
			Count:     args.Count,
			Des:       args.Des,
		})
		if err != nil {
			return
		}
	}
	return
}

// ArgsClearAssetsByProductID 清理指定产品的所有资产参数
type ArgsClearAssetsByProductID struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//资产产品
	ProductID int64 `db:"product_id" json:"productID"`
}

// ClearAssetsByProductID 清理指定产品的所有资产
func ClearAssetsByProductID(args *ArgsClearAssetsByProductID) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "finance_assets", "org_id = :org_id AND product_id = :product_id", args)
	return
}

// ArgsClearAssetsByOrgID 清理组织的所有资产参数
type ArgsClearAssetsByOrgID struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// ClearAssetsByOrgID 清理组织的所有资产
func ClearAssetsByOrgID(args *ArgsClearAssetsByOrgID) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "finance_assets", "org_id = :org_id", args)
	return
}
