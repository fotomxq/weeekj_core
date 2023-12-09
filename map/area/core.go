package MapArea

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLGPS "gitee.com/weeekj/weeekj_core/v5/core/sql/gps"
	CoreSQLIDs "gitee.com/weeekj/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	MapMathArea "gitee.com/weeekj/weeekj_core/v5/map/math/area"
	MapMathArgs "gitee.com/weeekj/weeekj_core/v5/map/math/args"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetList 获取列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分区标识码
	// 后端不做任何限制，该信息作为前端抽取数据类型使用
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//归属关系
	// -1 跳过
	// 可以作为行政分区和下级配送分区关系的设置，只有平台方能设置没有上级的分区
	// 其他分区必须指定行政分区作为上级，否则无法建立分区
	// 上级分区必须同属一个城市，且所有点不能超越范围
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//所属国家 国家代码
	// -1 跳过
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country" empty:"true"`
	//所属城市
	// -1 跳过
	City int `db:"city" json:"city" check:"city" empty:"true"`
	//地图制式
	// -1 跳过
	// 0 / 1 / 2
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType" check:"intThan0" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 获取列表
func GetList(args *ArgsGetList) (dataList []FieldsArea, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Mark != "" {
		where = where + " AND mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.ParentID > -1 {
		where = where + " AND parent_id = :parent_id"
		maps["parent_id"] = args.ParentID
	}
	if args.Country > -1 {
		where = where + " AND country = :country"
		maps["country"] = args.Country
	}
	if args.City > -1 {
		where = where + " AND city = :city"
		maps["city"] = args.City
	}
	if args.MapType > -1 {
		where = where + " AND map_type = :map_type"
		maps["map_type"] = args.MapType
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"map_area",
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, mark, parent_id, name, des, country, city, map_type, points, params FROM map_area WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsSearchName 搜索专用方法参数
type ArgsSearchName struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分区标识码
	// 后端不做任何限制，该信息作为前端抽取数据类型使用
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//归属关系
	// -1 跳过
	// 可以作为行政分区和下级配送分区关系的设置，只有平台方能设置没有上级的分区
	// 其他分区必须指定行政分区作为上级，否则无法建立分区
	// 上级分区必须同属一个城市，且所有点不能超越范围
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//所属国家 国家代码
	// -1 跳过
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country" empty:"true"`
	//所属城市
	// -1 跳过
	City int `db:"city" json:"city" check:"city" empty:"true"`
	//地图制式
	// -1 跳过
	// 0 / 1 / 2
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType" check:"intThan0" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// SearchName 搜索专用方法
func SearchName(args *ArgsSearchName) (dataMaps map[int64]string, dataCount int64, err error) {
	var dataList []FieldsArea
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Mark != "" {
		where = where + " AND mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.ParentID > -1 {
		where = where + " AND parent_id = :parent_id"
		maps["parent_id"] = args.ParentID
	}
	if args.Country > -1 {
		where = where + " AND country = :country"
		maps["country"] = args.Country
	}
	if args.City > -1 {
		where = where + " AND city = :city"
		maps["city"] = args.City
	}
	if args.MapType > -1 {
		where = where + " AND map_type = :map_type"
		maps["map_type"] = args.MapType
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"map_area",
		"id",
		"SELECT id, name FROM map_area WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err == nil {
		dataMaps = map[int64]string{}
		for _, v := range dataList {
			dataMaps[v.ID] = v.Name
		}
	}
	return
}

// ArgsGetByID 获取ID参数
type ArgsGetByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetByID 获取ID
func GetByID(args *ArgsGetByID) (data FieldsArea, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, mark, parent_id, name, des, country, city, map_type, points, params FROM map_area WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000)", args.ID, args.OrgID)
	return
}

// ArgsGetMore 批量获取参数
type ArgsGetMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
}

// GetMore 批量获取
func GetMore(args *ArgsGetMore) (dataList []FieldsArea, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "map_area", "id, create_at, update_at, delete_at, org_id, mark, parent_id, name, des, country, city, map_type, points, params", args.IDs, args.HaveRemove)
	return
}

func GetMoreMap(args *ArgsGetMore) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsNameAndDelete("map_area", args.IDs, args.HaveRemove)
	return
}

// ArgsCreate 创建分区参数
type ArgsCreate struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分区标识码
	// 后端不做任何限制，该信息作为前端抽取数据类型使用
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//归属关系
	// 可以作为行政分区和下级配送分区关系的设置，只有平台方能设置没有上级的分区
	// 其他分区必须指定行政分区作为上级，否则无法建立分区
	// 上级分区必须同属一个城市，且所有点不能超越范围
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country"`
	//所属城市
	City int `db:"city" json:"city" check:"city"`
	//地图制式
	// 0 / 1 / 2
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标系
	Points CoreSQLGPS.FieldsPoints `db:"points" json:"points"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// Create 创建分区
func Create(args *ArgsCreate) (data FieldsArea, errCode string, err error) {
	//检查上级
	if err = checkParent(args.OrgID, args.ParentID); err != nil {
		errCode = "check_parent"
		return
	}
	//如果存在上级，则检查坐标系在上级范围内
	if errCode, err = checkParentPoint(args.MapType, args.Points, args.ParentID); err != nil {
		return
	}
	//创建数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "map_area", "INSERT INTO map_area (org_id, mark, parent_id, name, des, country, city, map_type, points, params) VALUES (:org_id, :mark, :parent_id, :name, :des, :country, :city, :map_type, :points, :params)", args, &data)
	if err != nil {
		errCode = "insert"
		return
	}
	return
}

// ArgsUpdate 修改分区参数
type ArgsUpdate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID"`
	//分区标识码
	// 后端不做任何限制，该信息作为前端抽取数据类型使用
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//归属关系
	// 可以作为行政分区和下级配送分区关系的设置，只有平台方能设置没有上级的分区
	// 其他分区必须指定行政分区作为上级，否则无法建立分区
	// 上级分区必须同属一个城市，且所有点不能超越范围
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country"`
	//所属城市
	City int `db:"city" json:"city" check:"city"`
	//地图制式
	// 0 / 1 / 2
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标系
	Points CoreSQLGPS.FieldsPoints `db:"points" json:"points"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// Update 修改分区
func Update(args *ArgsUpdate) (errCode string, err error) {
	//检查上级
	if err = checkParent(args.OrgID, args.ParentID); err != nil {
		errCode = "check_parent"
		return
	}
	//如果存在上级，则检查坐标系在上级范围内
	if errCode, err = checkParentPoint(args.MapType, args.Points, args.ParentID); err != nil {
		return
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE map_area SET update_at = NOW(), mark = :mark, parent_id = :parent_id, name = :name, des = :des, country = :country, city = :city, map_type = :map_type, points = :points, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		errCode = "update"
		return
	}
	return
}

// ArgsDelete 删除分区参数
type ArgsDelete struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// Delete 删除分区
// 如果删除上级分区，将递归删除下属所有分区
func Delete(args *ArgsDelete) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "map_area", "(id = :id OR parent_id = :id) AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// checkParent 检查上级
func checkParent(orgID int64, parentID int64) (err error) {
	if parentID == 0 {
		return
	}
	var data FieldsArea
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, parent_id FROM map_area WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000)", parentID, orgID)
	if err != nil {
		err = errors.New("parent not exist, " + err.Error())
		return
	}
	if data.ParentID > 0 {
		err = errors.New("parent have parent id")
		return
	}
	return
}

// checkParentPoint 检查点是否在上级分区范围内？
func checkParentPoint(pointMapType int, points CoreSQLGPS.FieldsPoints, parentID int64) (errCode string, err error) {
	if parentID == 0 {
		return
	}
	var data FieldsArea
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, map_type, points FROM map_area WHERE id = $1 AND delete_at < to_timestamp(1000000)", parentID)
	if err != nil {
		errCode = "parent_not_exist"
		err = errors.New("parent not exist, " + err.Error())
		return
	}
	var parentPoints []MapMathArea.ParamsAreaPoint
	for _, v := range data.Points {
		parentPoints = append(parentPoints, MapMathArea.ParamsAreaPoint{
			Longitude: v.Longitude,
			Latitude:  v.Latitude,
		})
	}
	parentArea := MapMathArea.ParamsArea{
		ID:        data.ID,
		PointType: getMapType(data.MapType),
		Points:    parentPoints,
	}
	for _, v := range points {
		if !MapMathArea.CheckXYInArea(&MapMathArea.ArgsCheckXYInArea{
			Point: MapMathArgs.ParamsPoint{
				PointType: getMapType(pointMapType),
				Longitude: v.Longitude,
				Latitude:  v.Latitude,
			},
			Area: parentArea,
		}) {
			errCode = fmt.Sprint("point: ", v.Latitude, ",", v.Longitude)
			err = errors.New(fmt.Sprint("point not in parent area, point: ", v.Latitude, ",", v.Longitude))
			return
		}
	}
	return
}
