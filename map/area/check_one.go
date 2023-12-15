package MapArea

import (
	CoreSQLGPS "github.com/fotomxq/weeekj_core/v5/core/sql/gps"
	MapMathArea "github.com/fotomxq/weeekj_core/v5/map/math/area"
	MapMathArgs "github.com/fotomxq/weeekj_core/v5/map/math/args"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsCheckPointInArea 检查某点在哪组分区参数
type ArgsCheckPointInArea struct {
	//地图制式
	// 0 / 1 / 2 / 3
	// WGS-84 / GCJ-02 / BD-09 / 2000-china
	MapType int `db:"map_type" json:"mapType"`
	//要检查的点
	Point CoreSQLGPS.FieldsPoint `db:"point" json:"point"`
	//地图ID
	AreaID int64 `db:"area_id" json:"areaID" check:"id"`
}

// CheckPointInArea 检查某点在哪组分区
func CheckPointInArea(args *ArgsCheckPointInArea) bool {
	//mapType
	mapType := getMapType(args.MapType)
	if mapType == "" {
		return false
	}
	var areaData FieldsArea
	if err := Router2SystemConfig.MainDB.Get(&areaData, "SELECT id, map_type, points FROM map_area WHERE id = $1 AND delete_at < to_timestamp(1000000)", args.AreaID); err != nil || areaData.ID < 1 {
		return false
	}
	var points []MapMathArea.ParamsAreaPoint
	for _, vPoint := range areaData.Points {
		points = append(points, MapMathArea.ParamsAreaPoint{
			Longitude: vPoint.Longitude,
			Latitude:  vPoint.Latitude,
		})
	}
	area := MapMathArea.ParamsArea{
		ID:        areaData.ID,
		PointType: getMapType(areaData.MapType),
		Points:    points,
	}
	if area.PointType == "" {
		return false
	}
	if MapMathArea.CheckXYInArea(&MapMathArea.ArgsCheckXYInArea{
		Point: MapMathArgs.ParamsPoint{
			PointType: mapType,
			Longitude: args.Point.Longitude,
			Latitude:  args.Point.Latitude,
		},
		Area: area,
	}) {
		return true
	}
	return false
}
