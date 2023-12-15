package MapArea

import (
	"errors"
	"fmt"
	CoreSQLGPS "github.com/fotomxq/weeekj_core/v5/core/sql/gps"
	MapMathArea "github.com/fotomxq/weeekj_core/v5/map/math/area"
	MapMathArgs "github.com/fotomxq/weeekj_core/v5/map/math/args"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsCheckPointInAreas 检查某点在哪组分区参数
type ArgsCheckPointInAreas struct {
	//地图制式
	// 0 / 1 / 2 / 3
	// WGS-84 / GCJ-02 / BD-09 / 2000-china
	MapType int `db:"map_type" json:"mapType"`
	//要检查的点
	Point CoreSQLGPS.FieldsPoint `db:"point" json:"point"`
	//组织ID
	// 可选
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否为上级ID
	// 上级将强制约束为0，否则必须>0
	IsParent bool `db:"is_parent" json:"isParent" check:"bool" empty:"true"`
	//是否启用优先级机制
	NeedLevel bool `json:"needLevel" check:"bool"`
	//标识码
	// 可选
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
}

// CheckPointInAreas 检查某点在哪组分区
func CheckPointInAreas(args *ArgsCheckPointInAreas) (dataList []FieldsArea, err error) {
	//mapType
	mapType := getMapType(args.MapType)
	if mapType == "" {
		err = errors.New("map type is error")
		return
	}
	//找出所有符合条件的分区
	var page int64 = 1
	for {
		var findAreaList []FieldsArea
		err = Router2SystemConfig.MainDB.Select(&findAreaList, "SELECT id, map_type, points, params FROM map_area WHERE ($1 < 1 OR org_id = $1) AND (($2 = TRUE AND parent_id = 0) OR ($2 = FALSE AND parent_id > 0)) AND ($3 = '' OR mark = $3) AND delete_at < to_timestamp(1000000) LIMIT 10 OFFSET "+fmt.Sprint((page-1)*10), args.OrgID, args.IsParent, args.Mark)
		if err != nil {
			err = nil
			break
		}
		if len(findAreaList) < 1 {
			break
		}
		//遍历数据，检查在范围内的点
		for _, vArea := range findAreaList {
			var points []MapMathArea.ParamsAreaPoint
			for _, vPoint := range vArea.Points {
				points = append(points, MapMathArea.ParamsAreaPoint{
					Longitude: vPoint.Longitude,
					Latitude:  vPoint.Latitude,
				})
			}
			area := MapMathArea.ParamsArea{
				ID:        vArea.ID,
				PointType: getMapType(vArea.MapType),
				Points:    points,
			}
			if area.PointType == "" {
				continue
			}
			if MapMathArea.CheckXYInArea(&MapMathArea.ArgsCheckXYInArea{
				Point: MapMathArgs.ParamsPoint{
					PointType: mapType,
					Longitude: args.Point.Longitude,
					Latitude:  args.Point.Latitude,
				},
				Area: area,
			}) {
				dataList = append(dataList, vArea)
			}
		}
		page += 1
	}
	if len(dataList) < 1 {
		err = errors.New("no area find")
		return
	}
	//如果启动优先级
	if args.NeedLevel && len(dataList) > 1 {
		// 当前最靠前的优先级级别
		var nowLevel int64 = -1
		//找出序列中最小值
		for _, v := range dataList {
			level, b := v.Params.GetValInt64("level")
			if !b {
				continue
			}
			if level > 0 && level < nowLevel {
				nowLevel = level
			}
		}
		if nowLevel == -1 {
			return
		}
		var newDataList []FieldsArea
		for _, v := range dataList {
			level, b := v.Params.GetValInt64("level")
			if !b {
				continue
			}
			if level == nowLevel {
				newDataList = append(newDataList, v)
			}
		}
		dataList = newDataList
		if len(dataList) < 1 {
			err = errors.New("no area level")
			return
		}
	}
	return
}

// 获取mapType对应的值
func getMapType(mapType int) string {
	switch mapType {
	case 0:
		return "WGS-84"
	case 1:
		return "GCJ-02"
	case 2:
		return "BD-09"
	default:
		return ""
	}
}
