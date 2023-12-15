package TMSUserRunning

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	MapMathArgs "github.com/fotomxq/weeekj_core/v5/map/math/args"
	MapMathConversion "github.com/fotomxq/weeekj_core/v5/map/math/conversion"
	MapMathPoint "github.com/fotomxq/weeekj_core/v5/map/math/point"
	MapTMap "github.com/fotomxq/weeekj_core/v5/map/tmap"
)

// ArgsGetRunPrice 自动计算跑腿费用参数
type ArgsGetRunPrice struct {
	//期望上门时间
	WaitAt string `db:"wait_at" json:"waitAt" check:"isoTime"`
	//物品类型
	// order 订单类; 其他类型需前端约定
	GoodType string `db:"good_type" json:"goodType" check:"mark"`
	//物品重量，单位克
	GoodWidget int `db:"good_widget" json:"goodWidget" check:"intThan0" empty:"true"`
	//发货地址
	FromAddress CoreSQLAddress.FieldsAddress `db:"from_address" json:"fromAddress" check:"address_data" empty:"true"`
	//送货地址
	ToAddress CoreSQLAddress.FieldsAddress `db:"to_address" json:"toAddress" check:"address_data" empty:"true"`
}

type DataGetRunPrice struct {
	//米数
	BetweenM int64 `json:"betweenM"`
	//公里数
	BetweenKM int64 `json:"betweenKM"`
	//距离收费
	BetweenPrice int64 `json:"betweenPrice"`
	//重量增收费用
	WidgetPrice int64 `json:"widgetPrice"`
	//特殊时间段追加费用
	WaitPrice int64 `json:"waitPrice"`
	//总费用
	TotalPrice int64 `json:"totalPrice"`
}

// GetRunPrice 自动计算跑腿费用
func GetRunPrice(args *ArgsGetRunPrice) (data DataGetRunPrice) {
	/** 旧的计算形式，考虑项目快速上线，改为规则写死的设计逻辑
	//获取配置
	// 0-5,1000|6-10,1500
	// 0-5,1000|6-10,1500代表0-5公里10元, 6-10公里15元
	tmsRunningPriceConfig := BaseConfig.GetDataStringNoErr("TMSRunningPriceConfig")
	tmsRunningPriceConfigs := strings.Split(tmsRunningPriceConfig, "|")
	//计算米数距离
	addressBetween, err := MapMathPoint.GetDistanceM(&MapMathPoint.ArgsGetDistance{
		StartPoint: MapMathArgs.ParamsPoint{
			PointType: MapMathConversion.ConversionMapType(args.FromAddress.MapType),
			Longitude: args.FromAddress.Longitude,
			Latitude:  args.FromAddress.Latitude,
		},
		EndPoint: MapMathArgs.ParamsPoint{
			PointType: MapMathConversion.ConversionMapType(args.ToAddress.MapType),
			Longitude: args.ToAddress.Longitude,
			Latitude:  args.ToAddress.Latitude,
		},
	})
	betweenM = CoreFilter.GetInt64ByFloat64(addressBetween)
	//失败后，按照最高规格走
	if err != nil || addressBetween < 1 {
		for _, v := range tmsRunningPriceConfigs {
			vStrArr := strings.Split(v, ",")
			if len(vStrArr) != 2 {
				continue
			}
			vPrice, _ := CoreFilter.GetInt64ByString(vStrArr[1])
			if price < vPrice {
				price = vPrice
			}
		}
		return
	}
	//计算公里数
	miles := addressBetween / 1000
	//是否比所有范围都大
	tooMaxAll := true
	//遍历找到负荷条件的
	for _, v := range tmsRunningPriceConfigs {
		vStrArr := strings.Split(v, ",")
		if len(vStrArr) != 2 {
			continue
		}
		vStr1Arr := strings.Split(vStrArr[0], "-")
		if len(vStr1Arr) != 2 {
			continue
		}
		v1, _ := CoreFilter.GetFloat64ByString(vStr1Arr[0])
		v2, _ := CoreFilter.GetFloat64ByString(vStr1Arr[1])
		vPrice, _ := CoreFilter.GetInt64ByString(vStrArr[1])
		//检查条件
		if v1 <= miles && v2 >= miles {
			price = vPrice
			tooMaxAll = false
			break
		}
	}
	//如果没有符合条件的，按照最大的处理
	if tooMaxAll {
		for _, v := range tmsRunningPriceConfigs {
			vStrArr := strings.Split(v, ",")
			if len(vStrArr) != 2 {
				continue
			}
			vPrice, _ := CoreFilter.GetInt64ByString(vStrArr[1])
			if price < vPrice {
				price = vPrice
			}
		}
	}
	*/
	//以下采用新计算形式，注意如果未来需要引入参数配置，需在这个模式基础上做额外的调整
	// 调整后，请删除上述旧的方案设计
	//增加基础费用
	data.TotalPrice = 0
	//根据时间段，计算费用
	nowAt, err := CoreFilter.GetTimeCarbonByDefault(args.WaitAt)
	if err != nil {
		err = nil
		nowAt = CoreFilter.GetNowTimeCarbon()
	}
	//0-7点，增加6元
	if nowAt.Hour() >= 0 && nowAt.Hour() < 7 {
		data.WaitPrice += 600
	}
	//22-24点，增加3元
	if nowAt.Hour() >= 22 && nowAt.Hour() < 24 {
		data.WaitPrice += 300
	}
	//通过天地图，计算导航距离
	tMapData, err := MapTMap.GetDrive(CoreSQL2.ArgsGPS{
		Longitude: args.FromAddress.Longitude,
		Latitude:  args.FromAddress.Latitude,
	}, CoreSQL2.ArgsGPS{
		Longitude: args.ToAddress.Longitude,
		Latitude:  args.ToAddress.Latitude,
	}, "0")
	if err != nil {
		CoreLog.Warn("tms running, sum price, tmap get drive error: ", err)
		err = nil
		//通过本地化距离计算模块，计算实际距离
		//计算米数距离
		addressBetween, err := MapMathPoint.GetDistanceM(&MapMathPoint.ArgsGetDistance{
			StartPoint: MapMathArgs.ParamsPoint{
				PointType: MapMathConversion.ConversionMapType(args.FromAddress.MapType),
				Longitude: args.FromAddress.Longitude,
				Latitude:  args.FromAddress.Latitude,
			},
			EndPoint: MapMathArgs.ParamsPoint{
				PointType: MapMathConversion.ConversionMapType(args.ToAddress.MapType),
				Longitude: args.ToAddress.Longitude,
				Latitude:  args.ToAddress.Latitude,
			},
		})
		if err == nil {
			//计算出最终的米数
			data.BetweenM = CoreFilter.GetInt64ByFloat64(addressBetween)
		} else {
			CoreLog.Warn("tms running, sum price, math point get distance error: ", err)
			err = nil
		}
	} else {
		//使用天地图，计算距离
		data.BetweenM = CoreFilter.GetInt64ByFloat64(tMapData.Distance * 1000)
	}
	//折合米为公里数
	data.BetweenKM = data.BetweenM / 1000
	//如果在0-3公里范围，则每公里增加1元
	if data.BetweenKM > 0 && data.BetweenKM < 3 {
		data.BetweenPrice = CoreFilter.GetRoundToInt64(float64(data.BetweenKM)) * 100
	} else {
		//如果大于等于3公里，起步价4元，则每公里增加2元
		if data.BetweenKM >= 3 {
			data.BetweenPrice = 400 + CoreFilter.GetRoundToInt64(float64(data.BetweenKM))*200
		}
	}
	//计算重量附加费用
	data.WidgetPrice = 0
	//5-10公斤，2元起，每公斤增加2元
	if args.GoodWidget >= 5000 && args.GoodWidget < 10000 {
		data.WidgetPrice = 200 + int64(args.GoodWidget/1000)*200
	} else {
		//10公斤以上，每公斤增加10元
		if args.GoodWidget >= 10000 {
			data.WidgetPrice = int64(args.GoodWidget/1000) * 1000
		}
	}
	//汇总费用
	data.TotalPrice = data.WaitPrice + data.BetweenPrice + data.WidgetPrice
	//如果少于10元，按照10元计费
	if data.TotalPrice < 1000 {
		data.TotalPrice = 1000
	}
	//反馈
	return
}
