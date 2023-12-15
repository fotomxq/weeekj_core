package TMSTransport

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetBindToMallList 获取绑定列表参数
type ArgsGetBindToMallList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//商品ID
	BindMallID int64 `db:"bind_mall_id" json:"bindMallID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
}

// GetBindToMallList 获取绑定列表
func GetBindToMallList(args *ArgsGetBindToMallList) (dataList []FieldsBindToMall, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.BindID > -1 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.BindMallID > -1 {
		where = where + " AND bind_mall_id = :bind_mall_id"
		maps["bind_mall_id"] = args.BindMallID
	}
	tableName := "tms_transport_bind_to_mall"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, delete_at, org_id, bind_id, bind_mall_id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "delete_at"},
	)
	return
}

// ArgsSetBindToMall 设置绑定关系参数
type ArgsSetBindToMall struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//绑定商品
	BindMallID int64 `db:"bind_mall_id" json:"bindMallID" check:"id"`
}

// SetBindToMall 设置绑定关系
func SetBindToMall(args *ArgsSetBindToMall) (err error) {
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM tms_transport_bind_to_mall WHERE org_id = $1 AND bind_id = $2 AND bind_mall_id = $3", args.OrgID, args.BindID, args.BindMallID)
	if err == nil && id > 0 {
		return
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO tms_transport_bind_to_mall (org_id, bind_id, bind_mall_id) VALUES (:org_id,:bind_id,:bind_mall_id)", args)
	return
}

// ArgsDeleteBindToMall 删除绑定关系参数
type ArgsDeleteBindToMall struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteBindToMall 删除绑定关系
func DeleteBindToMall(args *ArgsDeleteBindToMall) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "tms_transport_bind_to_mall", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
