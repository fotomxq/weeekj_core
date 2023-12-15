package ServiceAD

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	MapArea "github.com/fotomxq/weeekj_core/v5/map/area"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetBindList 获取绑定关系参数
type ArgsGetBindList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//广告ID
	// -1 跳过
	AdID int64 `db:"ad_id" json:"adID" check:"id" empty:"true"`
	//分区ID
	// -1 跳过
	AreaID int64 `db:"area_id" json:"areaID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
}

// GetBindList 获取绑定关系
func GetBindList(args *ArgsGetBindList) (dataList []FieldsBind, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.AdID > -1 {
		where = where + " AND ad_id = :ad_id"
		maps["ad_id"] = args.AdID
	}
	if args.AreaID > -1 {
		where = where + " AND area_id = :area_id"
		maps["area_id"] = args.AreaID
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"service_ad_bind",
		"id",
		"SELECT id, create_at, update_at, delete_at, start_at, end_at, org_id, area_id, ad_id, factor, params FROM service_ad_bind WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "start_at", "end_at", "factor"},
	)
	return
}

// ArgsSetBind 设置绑定关系参数
type ArgsSetBind struct {
	//启动时间
	StartAt time.Time `db:"start_at" json:"startAt" check:"isoTime"`
	//结束时间
	EndAt time.Time `db:"end_at" json:"endAt" check:"isoTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分区ID
	AreaID int64 `db:"area_id" json:"areaID" check:"id"`
	//广告ID
	AdID int64 `db:"ad_id" json:"adID" check:"id"`
	//权重因子
	// 当一个分区存在多个广告绑定时，权重将作为随机的重要计算参数，影响随机的倾向性
	// 同一个分区下将所有广告进行合计，然后取随机数，在因子+偏移值范围的广告将作为最终投放标定
	// 注意，广告也可能会被用户数据分析模块影响，该模块将注入新的因子作为投放参数，作为因子2存在。
	// 投放概率 = 广告因子 + 用户数据分析因子 / (因子总数)
	Factor int `db:"factor" json:"factor" check:"intThan0"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetBind 设置绑定关系
func SetBind(args *ArgsSetBind) (data FieldsBind, err error) {
	//检查分区和组织关联性
	var areaData MapArea.FieldsArea
	areaData, err = MapArea.GetByID(&MapArea.ArgsGetByID{
		ID:    args.AreaID,
		OrgID: args.OrgID,
	})
	if err != nil {
		err = errors.New("area not exist, " + err.Error())
		return
	}
	//检查分区标识码是否为ad
	if areaData.Mark != "ad" {
		err = errors.New("area mark not ad")
		return
	}
	//获取绑定关系
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, area_id, ad_id, factor, params FROM service_ad_bind WHERE org_id = $1 AND area_id = $2 AND ad_id = $3 AND delete_at < to_timestamp(1000000)")
	if err == nil && data.ID > 0 {
		//设置绑定关系
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_ad_bind SET start_at = :start_at, end_at = :end_at, factor = :factor, params = :params WHERE id = :id", map[string]interface{}{
			"id":       data.ID,
			"start_at": args.StartAt,
			"end_at":   args.EndAt,
			"factor":   args.Factor,
			"params":   args.Params,
		})
		if err == nil {
			data.Factor = args.Factor
			data.Params = args.Params
		}
		return
	}
	//创建新的绑定关系
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "service_ad_bind", "INSERT INTO service_ad_bind(start_at, end_at, org_id, area_id, ad_id, factor, params) VALUES(:start_at, :end_at, :org_id, :area_id, :ad_id, :factor, :params)", args, &data)
	return
}

// ArgsDeleteBind 删除绑定关系参数
type ArgsDeleteBind struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 验证用
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteBind 删除绑定关系
func DeleteBind(args *ArgsDeleteBind) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_ad_bind", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
