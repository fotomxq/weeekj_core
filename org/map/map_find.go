package OrgMap

import (
	"errors"
	MapMathConversion "gitee.com/weeekj/weeekj_core/v5/map/math/conversion"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsFindMapByArea 查询GPS附件的商户信息列参数
type ArgsFindMapByArea struct {
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country"`
	//省份 编码
	// eg: 710000
	Province int `db:"province" json:"province"`
	//所属城市
	City int `db:"city" json:"city"`
	//地图制式
	// 0 WGS-84 / 1 GCJ-02 / 2 BD-09
	MapType int `db:"map_type" json:"mapType"`
	//中心点
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	//半径
	Radius float64 `db:"radius" json:"radius"`
	//是否包含不可点击的广告数据
	// 0 不包含 / 1 包含
	IncludeDisable int `db:"include_disable" json:"includeDisable"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// FindMapByArea 查询GPS附件的商户信息列
func FindMapByArea(args *ArgsFindMapByArea) (dataList []FieldsMap, err error) {
	for k := 0; k < 3; k++ {
		var result []MapMathConversion.ArgsConversionGPS
		result, err = MapMathConversion.ConversionMapTypeInt(&MapMathConversion.ArgsConversionMapTypeInt{
			SrcType:  args.MapType,
			DestType: k,
			Data: []MapMathConversion.ArgsConversionGPS{
				{
					Longitude: args.Longitude,
					Latitude:  args.Latitude,
				},
			},
		})
		if err != nil {
			return
		}
		longitudeMin := result[0].Longitude - args.Radius
		longitudeMax := result[0].Longitude + args.Radius
		latitudeMin := result[0].Latitude - args.Radius
		latitudeMax := result[0].Latitude + args.Radius
		var resultDataList []FieldsMap
		err = Router2SystemConfig.MainDB.Select(&resultDataList, "SELECT id FROM org_map WHERE delete_at < to_timestamp(1000000) AND audit_at > to_timestamp(1000000) AND country = $1 AND province = $2 AND city = $3 AND map_type = $4 AND longitude >= $5 AND longitude <= $6 AND latitude >= $7 AND latitude <= $8 AND (($9 = 1 AND ad_count < ad_count_limit) OR $9 = 0) AND parent_id = 0 AND (name ILIKE '%' || $10 || '%' OR des ILIKE '%' || $10 || '%') LIMIT 100", args.Country, args.Province, args.City, k, longitudeMin, longitudeMax, latitudeMin, latitudeMax, args.IncludeDisable, args.Search)
		if err != nil {
			err = nil
			continue
		}
		if len(resultDataList) < 1 {
			continue
		}
		for _, v := range resultDataList {
			vData := getMapByID(v.ID)
			if vData.ID < 1 {
				continue
			}
			dataList = append(dataList, vData)
		}
	}
	if len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// FindMapByAreaV2 FindMapByArea 查询GPS附件的商户信息列
func FindMapByAreaV2(args *ArgsFindMapByArea) (dataList []FieldsMap, err error) {
	for k := 0; k < 3; k++ {
		var result []MapMathConversion.ArgsConversionGPS
		result, err = MapMathConversion.ConversionMapTypeInt(&MapMathConversion.ArgsConversionMapTypeInt{
			SrcType:  args.MapType,
			DestType: k,
			Data: []MapMathConversion.ArgsConversionGPS{
				{
					Longitude: args.Longitude,
					Latitude:  args.Latitude,
				},
			},
		})
		if err != nil {
			return
		}
		longitudeMin := result[0].Longitude - args.Radius
		longitudeMax := result[0].Longitude + args.Radius
		latitudeMin := result[0].Latitude - args.Radius
		latitudeMax := result[0].Latitude + args.Radius
		var resultDataList []FieldsMap
		err = Router2SystemConfig.MainDB.Select(&resultDataList, "SELECT id FROM org_map WHERE delete_at < to_timestamp(1000000) AND audit_at > to_timestamp(1000000) AND country = $1 AND province = $2 AND city = $3 AND map_type = $4 AND longitude >= $5 AND longitude <= $6 AND latitude >= $7 AND latitude <= $8 AND parent_id = 0 AND (name ILIKE '%' || $9 || '%' OR des ILIKE '%' || $9 || '%') LIMIT 100", args.Country, args.Province, args.City, k, longitudeMin, longitudeMax, latitudeMin, latitudeMax, args.Search)
		if err != nil {
			err = nil
			continue
		}
		if len(resultDataList) < 1 {
			continue
		}
		for _, v := range resultDataList {
			vData := getMapByID(v.ID)
			if vData.ID < 1 {
				continue
			}
			dataList = append(dataList, vData)
		}
	}
	if len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	return
}
