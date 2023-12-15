package MallCore

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetTransportList 获取配送模版参数
type ArgsGetTransportList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetTransportList 获取配送模版
func GetTransportList(args *ArgsGetTransportList) (dataList []FieldsTransport, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "mall_core_transport"
	var rawList []FieldsTransport
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getTransportID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetTransportID 获取指定的模版ID参数
type ArgsGetTransportID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetTransportID 获取指定的模版ID
func GetTransportID(args *ArgsGetTransportID) (data FieldsTransport, err error) {
	data = getTransportID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		data = FieldsTransport{}
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetTransports 获取多个模版参数
type ArgsGetTransports struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// GetTransports 获取多个模版
func GetTransports(args *ArgsGetTransports) (dataList []FieldsTransport, err error) {
	for _, v := range args.IDs {
		vData := getTransportID(v)
		if vData.ID < 1 || !CoreFilter.EqID2(args.OrgID, vData.OrgID) {
			continue
		}
		if !args.HaveRemove && CoreSQL.CheckTimeHaveData(vData.DeleteAt) {
			continue
		}
		dataList = append(dataList, vData)
	}
	if len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	return
}

func ArgsGetTransportsName(args *ArgsGetTransports) (dataList map[int64]string, err error) {
	dataList, err = CoreSQLIDs.GetIDsNameAndDelete("mall_core_transport", args.IDs, args.HaveRemove)
	return
}

// ArgsCreateTransport 创建新的模版参数
type ArgsCreateTransport struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//模版名称
	Name string `db:"name" json:"name" check:"title" min:"1" max:"300"`
	//计费规则
	// 0 无配送费；1 按件计算；2 按重量计算；3 按公里数计算
	// 后续所有计费标准必须统一，否则系统将拒绝创建或修改
	Rules int `db:"rules" json:"rules" check:"intThan0" empty:"true"`
	//首费标准
	// 0 留空；1 N件开始收费；2 N重量开始收费；3 N公里开始收费
	RulesUnit int `db:"rules_unit" json:"rulesUnit" check:"intThan0" empty:"true"`
	//首费金额
	RulesPrice int64 `db:"rules_price" json:"rulesPrice" check:"price" empty:"true"`
	//增费标准
	// 0 无增量；1 每N件增加费用；2 每N重量增加费用；3 每N公里增加费用
	AddType int `db:"add_type" json:"addType" check:"intThan0" empty:"true"`
	//增费单位
	AddUnit int `db:"add_unit" json:"addUnit" check:"intThan0" empty:"true"`
	//增费金额
	// 单位增加的费用
	AddPrice int64 `db:"add_price" json:"addPrice" check:"price" empty:"true"`
	//免邮条件
	// 0 无免费；1 按件免费；2 按重量免费; 3 公里数内免费
	FreeType int `db:"free_type" json:"freeType" check:"intThan0" empty:"true"`
	//免邮数量
	FreeUnit int `db:"free_unit" json:"freeUnit" check:"intThan0" empty:"true"`
}

// CreateTransport 创建新的模版
func CreateTransport(args *ArgsCreateTransport) (data FieldsTransport, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "mall_core_transport", "INSERT INTO mall_core_transport (org_id, name, rules, rules_unit, rules_price, add_type, add_unit, add_price, free_type, free_unit) VALUES (:org_id,:name,:rules,:rules_unit,:rules_price,:add_type,:add_unit,:add_price,:free_type,:free_unit)", args, &data)
	if err != nil {
		return
	}
	return
}

// ArgsUpdateTransport 修改模版参数
type ArgsUpdateTransport struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，检查项
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//模版名称
	Name string `db:"name" json:"name" check:"title" min:"1" max:"300"`
	//计费规则
	// 0 无配送费；1 按件计算；2 按重量计算；3 按公里数计算
	// 后续所有计费标准必须统一，否则系统将拒绝创建或修改
	Rules int `db:"rules" json:"rules" check:"intThan0" empty:"true"`
	//首费标准
	// 0 留空；1 N件开始收费；2 N重量开始收费；3 N公里开始收费
	RulesUnit int `db:"rules_unit" json:"rulesUnit" check:"intThan0" empty:"true"`
	//首费金额
	RulesPrice int64 `db:"rules_price" json:"rulesPrice" check:"price" empty:"true"`
	//增费标准
	// 0 无增量；1 每N件增加费用；2 每N重量增加费用；3 每N公里增加费用
	AddType int `db:"add_type" json:"addType" check:"intThan0" empty:"true"`
	//增费单位
	AddUnit int `db:"add_unit" json:"addUnit" check:"intThan0" empty:"true"`
	//增费金额
	// 单位增加的费用
	AddPrice int64 `db:"add_price" json:"addPrice" check:"price" empty:"true"`
	//免邮条件
	// 0 无免费；1 按件免费；2 按重量免费; 3 公里数内免费
	FreeType int `db:"free_type" json:"freeType" check:"intThan0" empty:"true"`
	//免邮数量
	FreeUnit int `db:"free_unit" json:"freeUnit" check:"intThan0" empty:"true"`
}

// UpdateTransport 修改模版
func UpdateTransport(args *ArgsUpdateTransport) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE mall_core_transport SET update_at = NOW(), name = :name, rules = :rules, rules_unit = :rules_unit, rules_price = :rules_price, add_type = :add_type, add_unit = :add_unit, add_price = :add_price, free_type = :free_type, free_unit = :free_unit WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteTransportCache(args.ID)
	//反馈
	return
}

// ArgsDeleteTransport 删除模版参数
type ArgsDeleteTransport struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，检查项
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteTransport 删除模版
// 必须删除关联后才能删除
func DeleteTransport(args *ArgsDeleteTransport) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "mall_core_transport", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteTransportCache(args.ID)
	//反馈
	return
}

func getTransportID(id int64) (data FieldsTransport) {
	cacheMark := getTransportCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, name, rules, rules_unit, rules_price, add_type, add_unit, add_price, free_type, free_unit FROM mall_core_transport WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 86400)
	return
}
