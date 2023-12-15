package ServiceHousekeeping

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	MarketCore "github.com/fotomxq/weeekj_core/v5/market/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetBindList 获取成员列表参数
type ArgsGetBindList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
}

// GetBindList 获取成员列表
func GetBindList(args *ArgsGetBindList) (dataList []FieldsBind, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	tableName := "service_housekeeping_bind"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, bind_id, map_area_id, all_take_price, all_log_count, all_level, un_finish_count, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetBindID 获取成员ID参数
type ArgsGetBindID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
}

// GetBindID 获取成员ID
func GetBindID(args *ArgsGetBindID) (data FieldsBind, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, bind_id, map_area_id, all_take_price, all_log_count, all_level, un_finish_count, params FROM service_housekeeping_bind WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND ($3 < 1 OR bind_id = $3)", args.ID, args.OrgID, args.BindID)
	return
}

// ArgsGetBindByBind 通过成员获取服务人员信息参数
type ArgsGetBindByBind struct {
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetBindByBind 通过成员获取服务人员信息
func GetBindByBind(args *ArgsGetBindByBind) (data FieldsBind, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, bind_id, map_area_id, all_take_price, all_log_count, all_level, un_finish_count, params FROM service_housekeeping_bind WHERE bind_id = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000)", args.BindID, args.OrgID)
	return
}

// ArgsSetBind 设置成员参数
type ArgsSetBind struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//分区ID
	MapAreaID int64 `db:"map_area_id" json:"mapAreaID" check:"id" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetBind 设置成员
func SetBind(args *ArgsSetBind) (data FieldsBind, err error) {
	data, err = GetBindByBind(&ArgsGetBindByBind{
		OrgID:  args.OrgID,
		BindID: args.BindID,
	})
	if err != nil {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "service_housekeeping_bind", "INSERT INTO service_housekeeping_bind (org_id, bind_id, map_area_id, all_take_price, all_log_count, all_level, params) VALUES (:org_id, :bind_id, :map_area_id, 0, 0, 0, :params)", args, &data)
		return
	} else {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_housekeeping_bind SET delete_at = to_timestamp(0), update_at = NOW(), map_area_id = :map_area_id, params = :params WHERE id = :id", map[string]interface{}{
			"id":          data.ID,
			"map_area_id": args.MapAreaID,
			"params":      args.Params,
		})
		if err == nil {
			data.MapAreaID = args.MapAreaID
			data.Params = args.Params
		}
		return
	}
}

// ArgsDeleteBind 删除成员参数
type ArgsDeleteBind struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteBind 删除成员
func DeleteBind(args *ArgsDeleteBind) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_housekeeping_bind", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// getBindByMarketUserID 通过营销关系，获取用户对应服务人员
func getBindByMarketUserID(orgID, userID int64) (data FieldsBind, err error) {
	var bindData MarketCore.FieldsBind
	bindData, err = MarketCore.GetBindByUserID(&MarketCore.ArgsGetBindByUserID{
		OrgID:      orgID,
		BindUserID: userID,
		BindInfoID: 0,
	})
	if err != nil {
		return
	}
	if bindData.ID < 1 {
		err = errors.New("market core not bind")
		return
	}
	data, err = GetBindByBind(&ArgsGetBindByBind{
		BindID: bindData.BindID,
		OrgID:  orgID,
	})
	if err == nil {
		if data.ID < 1 {
			err = errors.New("no data")
			return
		}
		if data.DeleteAt.Unix() > 100000 {
			err = errors.New("no data")
			return
		}
	}
	return
}
