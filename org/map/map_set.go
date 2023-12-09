package OrgMap

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	UserSystemTip "gitee.com/weeekj/weeekj_core/v5/user/system_tip"
	"github.com/lib/pq"
)

// ArgsCreateMap 创建商户的位置信息参数
type ArgsCreateMap struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//上级ID
	// 用于叠加展示
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//展示小图标
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	// 轮播图片组
	CoverFileIDs pq.Int64Array `db:"cover_file_ids" json:"coverFileIDs" check:"ids" empty:"true"`
	//展示信息
	Name string `db:"name" json:"name" check:"name"`
	//展示介绍信息
	Des string `db:"des" json:"des" check:"des" min:"0" max:"1000" empty:"true"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country" empty:"true"`
	//省份 编码
	// eg: 710000
	Province int `db:"province" json:"province" check:"province" empty:"true"`
	//所属城市
	City int `db:"city" json:"city" check:"city" empty:"true"`
	//街道详细信息
	Address string `db:"address" json:"address" check:"address" empty:"true"`
	//地图制式
	// 0 WGS-84 / 1 GCJ-02 / 2 BD-09
	MapType int `db:"map_type" json:"mapType" check:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	//广告可用次数
	AdCountLimit int64 `db:"ad_count_limit" json:"adCountLimit"`
	//查看最短时间长度
	ViewTimeLimit int64 `db:"view_time_limit" json:"viewTimeLimit"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateMap 创建商户的位置信息
func CreateMap(args *ArgsCreateMap) (data FieldsMap, err error) {
	//写入新的数据
	if args.ViewTimeLimit < 0 {
		args.ViewTimeLimit = 0
	}
	var mapID int64
	mapID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO org_map (org_id, user_id, parent_id, cover_file_id, cover_file_ids, name, des, country, province, city, address, map_type, longitude, latitude, ad_count, ad_count_limit, view_time_limit, params) VALUES (:org_id,:user_id,:parent_id,:cover_file_id,:cover_file_ids,:name,:des,:country,:province,:city,:address,:map_type,:longitude,:latitude,0,:ad_count_limit,:view_time_limit,:params)", args)
	if err != nil {
		return
	}
	//获取数据
	data = getMapByID(mapID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//反馈
	return
}

// ArgsUpdateMap 修改组织地图
type ArgsUpdateMap struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//上级ID
	// 用于叠加展示
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//展示小图标
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	// 轮播图片组
	CoverFileIDs pq.Int64Array `db:"cover_file_ids" json:"coverFileIDs" check:"ids" empty:"true"`
	//展示信息
	Name string `db:"name" json:"name" check:"name"`
	//展示介绍信息
	Des string `db:"des" json:"des" check:"des" min:"0" max:"1000" empty:"true"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country" empty:"true"`
	//省份 编码
	// eg: 710000
	Province int `db:"province" json:"province" check:"province" empty:"true"`
	//所属城市
	City int `db:"city" json:"city" check:"city" empty:"true"`
	//街道详细信息
	Address string `db:"address" json:"address" check:"address" empty:"true"`
	//地图制式
	// 0 WGS-84 / 1 GCJ-02 / 2 BD-09
	MapType int `db:"map_type" json:"mapType" check:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	//广告可用次数
	AdCountLimit int64 `db:"ad_count_limit" json:"adCountLimit"`
	//查看最短时间长度
	ViewTimeLimit int64 `db:"view_time_limit" json:"viewTimeLimit"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateMap 修改组织地图参数
func UpdateMap(args *ArgsUpdateMap) (err error) {
	//获取数据
	mapData := getMapByID(args.ID)
	//判断所有权
	if (mapData.OrgID != args.OrgID && args.OrgID > 0) && (mapData.UserID != args.UserID && args.UserID > 0) {
		err = errors.New("map not this user or org")
		return
	}
	//更新数据
	if args.ViewTimeLimit < 0 {
		args.ViewTimeLimit = mapData.ViewTimeLimit
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_map SET update_at = NOW(), audit_at = to_timestamp(0), parent_id = :parent_id, cover_file_id = :cover_file_id, cover_file_ids = :cover_file_ids, name = :name, des = :des, country = :country, province = :province, city = :city, address = :address, map_type = :map_type, longitude = :longitude, latitude = :latitude, ad_count_limit = :ad_count_limit, view_time_limit = :view_time_limit, params = :params WHERE id = :id", map[string]interface{}{
		"id":              args.ID,
		"parent_id":       args.ParentID,
		"cover_file_id":   args.CoverFileID,
		"cover_file_ids":  args.CoverFileIDs,
		"name":            args.Name,
		"des":             args.Des,
		"country":         args.Country,
		"province":        args.Province,
		"city":            args.City,
		"address":         args.Address,
		"map_type":        args.MapType,
		"longitude":       args.Longitude,
		"latitude":        args.Latitude,
		"ad_count_limit":  args.AdCountLimit,
		"view_time_limit": args.ViewTimeLimit,
		"params":          args.Params,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteMapCache(args.ID)
	//反馈
	return
}

// ArgsSetMap 设置商户的位置信息参数
type ArgsSetMap struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//上级ID
	// 用于叠加展示
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//展示小图标
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	// 轮播图片组
	CoverFileIDs pq.Int64Array `db:"cover_file_ids" json:"coverFileIDs" check:"ids" empty:"true"`
	//展示信息
	Name string `db:"name" json:"name" check:"name"`
	//展示介绍信息
	Des string `db:"des" json:"des" check:"des" min:"0" max:"1000" empty:"true"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country" empty:"true"`
	//省份 编码
	// eg: 710000
	Province int `db:"province" json:"province" check:"province" empty:"true"`
	//所属城市
	City int `db:"city" json:"city" check:"city" empty:"true"`
	//街道详细信息
	Address string `db:"address" json:"address" check:"address" empty:"true"`
	//地图制式
	// 0 WGS-84 / 1 GCJ-02 / 2 BD-09
	MapType int `db:"map_type" json:"mapType" check:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	//广告可用次数
	AdCountLimit int64 `db:"ad_count_limit" json:"adCountLimit"`
	//查看最短时间长度
	ViewTimeLimit int64 `db:"view_time_limit" json:"viewTimeLimit"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetMap 设置商户的位置信息
func SetMap(args *ArgsSetMap) (data FieldsMap, err error) {
	//查询该组织的位置信息
	var id int64
	if args.OrgID > 0 {
		_ = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM org_map WHERE org_id = $1 AND delete_at < to_timestamp(1000000)", args.OrgID)
	} else {
		if args.UserID > 0 {
			_ = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM org_map WHERE user_id = $1 AND delete_at < to_timestamp(1000000)", args.UserID)
		} else {
			_ = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM org_map WHERE org_id = 0 AND map_type = $1 AND longitude = $2 AND latitude = $3 AND delete_at < to_timestamp(1000000)", args.MapType, args.Longitude, args.Latitude)
		}
	}
	if id > 0 {
		//重新获取数据
		data = getMapByID(data.ID)
		//更新数据
		if args.ViewTimeLimit < 0 {
			args.ViewTimeLimit = data.ViewTimeLimit
		}
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_map SET update_at = NOW(), audit_at = to_timestamp(0), parent_id = :parent_id, cover_file_id = :cover_file_id, cover_file_ids = :cover_file_ids, name = :name, des = :des, country = :country, province = :province, city = :city, address = :address, map_type = :map_type, longitude = :longitude, latitude = :latitude, ad_count_limit = :ad_count_limit, view_time_limit = :view_time_limit, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:user_id < 1 OR user_id = :user_id)", map[string]interface{}{
			"id":              id,
			"org_id":          args.OrgID,
			"user_id":         args.UserID,
			"parent_id":       args.ParentID,
			"cover_file_id":   args.CoverFileID,
			"cover_file_ids":  args.CoverFileIDs,
			"name":            args.Name,
			"des":             args.Des,
			"country":         args.Country,
			"province":        args.Province,
			"city":            args.City,
			"address":         args.Address,
			"map_type":        args.MapType,
			"longitude":       args.Longitude,
			"latitude":        args.Latitude,
			"ad_count_limit":  args.AdCountLimit,
			"view_time_limit": args.ViewTimeLimit,
			"params":          args.Params,
		})
		if err != nil {
			return
		}
		deleteMapCache(id)
	} else {
		//写入新的数据
		if args.ViewTimeLimit < 0 {
			args.ViewTimeLimit = 0
		}
		id, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO org_map (org_id, user_id, parent_id, cover_file_id, cover_file_ids, name, des, country, province, city, address, map_type, longitude, latitude, ad_count, ad_count_limit, view_time_limit, params) VALUES (:org_id,:user_id,:parent_id,:cover_file_id,:name,:des,:country,:province,:city,:address,:map_type,:longitude,:latitude,0,:ad_count_limit,:view_time_limit,:params)", args)
	}
	if err != nil {
		return
	}
	data = getMapByID(id)
	return
}

// ArgsAuditMap 审核地址信息参数
type ArgsAuditMap struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// AuditMap 审核地址信息
func AuditMap(args *ArgsAuditMap) (err error) {
	//修改审核状态
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_map SET audit_at = NOW() WHERE id = :id", args)
	if err != nil {
		return
	}
	//删除缓冲
	deleteMapCache(args.ID)
	//推送nats
	pushNatsMapAudit(args.ID)
	//获取数据
	data := getMapByID(args.ID)
	//审核通知
	if data.UserID > 0 {
		UserSystemTip.SendSuccess(data.UserID, "商户地图", data.ID, data.Name)
	}
	//反馈
	return
}
