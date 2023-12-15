package TMSTransport

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	MapMathArgs "github.com/fotomxq/weeekj_core/v5/map/math/args"
	MapMathPoint "github.com/fotomxq/weeekj_core/v5/map/math/point"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 更新配送员统计数据
func updateTransportBindAnalysis(data FieldsTransport, newBindID int64) (err error) {
	if newBindID < 1 {
		return
	}
	//生成统计数据
	var km int64
	km, err = MapMathPoint.GetDistance(&MapMathPoint.ArgsGetDistance{
		StartPoint: MapMathArgs.ParamsPoint{
			PointType: data.FromAddress.GetMapType(),
			Longitude: data.FromAddress.Longitude,
			Latitude:  data.FromAddress.Latitude,
		},
		EndPoint: MapMathArgs.ParamsPoint{
			PointType: data.ToAddress.GetMapType(),
			Longitude: data.ToAddress.Longitude,
			Latitude:  data.ToAddress.Latitude,
		},
	})
	if err != nil {
		err = nil
		km = 0
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport_bind SET km_1_day = km_1_day + :km_1_day, count_1_day = count_1_day + 1 WHERE bind_id = :bind_id", map[string]interface{}{
		"bind_id":  newBindID,
		"km_1_day": km,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("update bind km and count, bind id: ", newBindID, ", ", err))
		return
	}
	//反馈
	return
}
