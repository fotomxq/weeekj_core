package ServiceInfoExchange

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetInfoList 获取列表参数
type ArgsGetInfoList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//信息类型
	// none 普通类型; recruitment 招聘信息; rent 租房信息; thing 物品交易
	InfoType string `db:"info_type" json:"infoType" check:"mark" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//是否需要审核参数
	NeedIsAudit bool `json:"needIsAudit" check:"bool"`
	//是否已经审核
	IsAudit bool `json:"isAudit" check:"bool"`
	//是否需要发布参数
	NeedIsPublish bool `json:"needIsPublish" check:"bool"`
	//是否已经发布
	IsPublish bool `json:"isPublish" check:"bool"`
	//价格区间
	PriceMin int64 `db:"price_min" json:"priceMin" check:"price" empty:"true"`
	PriceMax int64 `db:"price_max" json:"priceMax" check:"price" empty:"true"`
	//是否需要过期参数
	NeedIsExpire bool `json:"needIsExpire" check:"bool"`
	IsExpire     bool `json:"isExpire" check:"bool"`
	//关联的订单
	OrderID     int64 `db:"order_id" json:"orderID" check:"id" empty:"true"`
	WaitOrderID int64 `db:"wait_order_id" json:"waitOrderID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetInfoList 获取列表
func GetInfoList(args *ArgsGetInfoList) (dataList []FieldsInfo, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.InfoType != "" {
		where = where + " AND info_type = :info_type"
		maps["info_type"] = args.InfoType
	}
	if args.SortID > -1 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if len(args.Tags) > 0 {
		where = where + " AND tags @> :tags"
		maps["tags"] = args.Tags
	}
	if args.NeedIsAudit {
		if args.IsAudit {
			where = where + " AND audit_at >= to_timestamp(1000000)"
		} else {
			where = where + " AND audit_at < to_timestamp(1000000)"
		}
	}
	if args.NeedIsPublish {
		if args.IsPublish {
			where = where + " AND publish_at >= to_timestamp(1000000)"
		} else {
			where = where + " AND publish_at < to_timestamp(1000000)"
		}
	}
	if args.PriceMin > 0 {
		where = where + " AND price >= :price_min"
		maps["price_min"] = args.PriceMin
	}
	if args.PriceMax > 0 {
		where = where + " AND price <= :price_max"
		maps["price_max"] = args.PriceMax
	}
	if args.NeedIsExpire {
		if args.IsExpire {
			where = where + " AND expire_at < NOW()"
		} else {
			where = where + " AND expire_at >= NOW()"
		}
	}
	if args.OrderID > 0 {
		where = where + " AND order_id = :order_id"
		maps["order_id"] = args.OrderID
	}
	if args.WaitOrderID > 0 {
		where = where + " AND wait_order_id = :wait_order_id"
		maps["wait_order_id"] = args.WaitOrderID
	}
	if args.Search != "" {
		where = where + " AND (title ILIKE '%' || :search || '%' OR title_des ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "service_info_exchange"
	var rawList []FieldsInfo
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "publish_at", "audit_at", "expire_at"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getInfoByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		vData.Des = ""
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetInfoID 获取指定信息参数
type ArgsGetInfoID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// GetInfoID 获取指定信息
func GetInfoID(args *ArgsGetInfoID) (data FieldsInfo, err error) {
	data = getInfoByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	if !CoreFilter.EqID2(args.OrgID, data.OrgID) || !CoreFilter.EqID2(args.UserID, data.UserID) {
		err = errors.New("no data")
		return
	}
	return
}

type ArgsGetInfoMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

func GetInfoMore(args *ArgsGetInfoMore) (dataList []FieldsInfo, err error) {
	for _, v := range args.IDs {
		vData := getInfoByID(v)
		if vData.ID < 1 || (!args.HaveRemove && CoreSQL.CheckTimeHaveData(vData.DeleteAt)) {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetInfoPublishID 获取公开信息
func GetInfoPublishID(args *ArgsGetInfoID) (data FieldsInfo, err error) {
	data = getInfoByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	if !CoreFilter.EqID2(args.OrgID, data.OrgID) || !CoreFilter.EqID2(args.UserID, data.UserID) || CoreSQL.CheckTimeHaveData(data.DeleteAt) || !CoreSQL.CheckTimeHaveData(data.AuditAt) || !CoreSQL.CheckTimeHaveData(data.PublishAt) {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetInfoCountByUser 获取用户发布数量参数
type ArgsGetInfoCountByUser struct {
	//信息类型
	// none 普通类型; recruitment 招聘信息; rent 租房信息; thing 物品交易
	InfoType string `db:"info_type" json:"infoType" check:"mark" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// GetInfoCountByUser 获取用户发布数量
func GetInfoCountByUser(args *ArgsGetInfoCountByUser) (count int64) {
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM service_info_exchange WHERE user_id = $1 AND delete_at < to_timestamp(1000000) AND ($2 = '' OR info_type = $2)", args.UserID, args.InfoType)
	return
}

// 获取信息交互数据
func getInfoByID(id int64) (data FieldsInfo) {
	cacheMark := getInfoCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, publish_at, audit_at, audit_des, info_type, org_id, user_id, sort_id, tags, title, title_des, des, cover_file_ids, currency, price, order_id, wait_order_id, order_finish, address, params, expire_at, limit_count FROM service_info_exchange WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 86400)
	return
}

// 根据等待订单获取信息
func getInfoByWaitOrderID(waitOrderID int64) (data FieldsInfo) {
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM service_info_exchange WHERE wait_order_id = $1", waitOrderID)
	if data.ID < 1 {
		return
	}
	data = getInfoByID(data.ID)
	return
}

// 根据订单获取信息
func getInfoByOrderID(orderID int64) (data FieldsInfo) {
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM service_info_exchange WHERE order_id = $1", orderID)
	if data.ID < 1 {
		return
	}
	data = getInfoByID(data.ID)
	return
}
