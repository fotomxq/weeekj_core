package ServiceAD

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLGPS "github.com/fotomxq/weeekj_core/v5/core/sql/gps"
	MapArea "github.com/fotomxq/weeekj_core/v5/map/area"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"github.com/robfig/cron"
)

//广告模块
// 广告和分区直接绑定，为客户提供广告投放的服务
// 广告不绑定分区，也可以进行投放，具体需要根据业务逻辑区分

var (
	//定时器
	runTimer       *cron.Cron
	runEndLock     = false
	runHistoryLock = false
)

// ArgsPutAD 根据分区ID，自动投放指定mark的广告参数
type ArgsPutAD struct {
	//组织ID
	// 如果为0，则表示为平台方
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分区ID
	// 分区如果没有给与，则直接抽取mark
	AreaID int64 `db:"area_id" json:"areaID" check:"id" empty:"true"`
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
}

// DataAD 根据分区ID，自动投放指定mark的广告数据
type DataAD struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//广告标识码，识别位置信息
	Mark string `db:"mark" json:"mark"`
	//名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//描述组图
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// PutAD 根据分区ID，自动投放指定mark的广告
// 本方法将依赖于绑定关系及权重计算后投放
func PutAD(args *ArgsPutAD) (data DataAD, errCode string, err error) {
	defer func() {
		//覆盖位置信息
		data.Mark = args.Mark
	}()
	//如果分区不存在，则根据mark找到合适的分区
	if args.AreaID < 1 {
		data, err = getAdByAuto(args.OrgID, args.Mark)
		if err != nil {
			errCode = "glob_not_exist"
			return
		}
		return
	}
	//获取所有绑定关系
	var bindList []FieldsBind
	err = Router2SystemConfig.MainDB.Select(&bindList, "SELECT bind.factor as factor, bind.ad_id as ad_id FROM service_ad_bind as bind, service_ad as ad WHERE bind.org_id = $1 AND bind.area_id = $2 AND bind.delete_at < to_timestamp(1000000) AND ad.id = bind.ad_id AND ad.delete_at < to_timestamp(1000000) AND ad.mark = $3 AND bind.start_at <= NOW() AND bind.end_at >= NOW()", args.OrgID, args.AreaID, args.Mark)
	if err != nil || len(bindList) < 1 {
		data, err = getAdByAuto(args.OrgID, args.Mark)
		if err != nil {
			errCode = "glob_not_exist"
			return
		}
		return
	}
	//根据绑定关系，计算权重
	countFactor := 0
	for _, v := range bindList {
		countFactor += v.Factor
	}
	result := CoreFilter.GetRandNumber(0, countFactor)
	if result < 1 {
		result = 1
	}
	//第二次遍历，按照偏移值找出step + min~max之间的值
	key := 0
	nextFactor := 0
	for k, v := range bindList {
		if result >= nextFactor && result <= nextFactor+v.Factor {
			key = k
			break
		}
		nextFactor += v.Factor
	}
	var findData FieldsBind
	var isFind bool
	for k, v := range bindList {
		if k != key {
			continue
		}
		findData = v
		isFind = true
		break
	}
	if !isFind {
		findData = bindList[0]
	}
	//找到符合条件的绑定ID
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT name, des, cover_file_id, des_files, params FROM service_ad WHERE id = $1 AND delete_at < to_timestamp(1000000)", findData.AdID)
	if err == nil {
		err = appendAnalysisData(&argsAppendAnalysisData{
			OrgID:  args.OrgID,
			AreaID: args.AreaID,
			AdID:   findData.AdID,
			Count:  1,
		})
		if err != nil {
			errCode = "analysis_insert"
			err = errors.New("append analysis, " + err.Error())
			return
		}
	} else {
		data, err = getAdByAuto(args.OrgID, args.Mark)
		if err != nil {
			errCode = "glob_not_exist"
			return
		}
		return
	}
	//如果没有数据，则获取全局广告
	if data.ID < 1 {
		data, err = getAdByAuto(args.OrgID, args.Mark)
		if err != nil {
			errCode = "glob_not_exist"
			return
		}
		return
	}
	//反馈数据
	return
}

// ArgsPutADByGPS 根据定位识别分区后，再投放广告参数
type ArgsPutADByGPS struct {
	//地图制式
	// 0 / 1 / 2
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType"`
	//要检查的点
	Point CoreSQLGPS.FieldsPoint `db:"point" json:"point"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
}

// PutADByGPS 根据定位识别分区后，再投放广告
func PutADByGPS(args *ArgsPutADByGPS) (data DataAD, errCode string, err error) {
	defer func() {
		//覆盖位置信息
		data.Mark = args.Mark
	}()
	//获取符合条件的分区
	var areaData MapArea.FieldsArea
	areaData, err = MapArea.CheckPointInAreasRand(&MapArea.ArgsCheckPointInAreas{
		MapType:   args.MapType,
		Point:     args.Point,
		OrgID:     args.OrgID,
		IsParent:  false,
		NeedLevel: true,
		Mark:      "ad",
	})
	if err != nil || areaData.ID < 1 {
		//不存在分区，则获取全局广告
		data, err = getAdByAuto(args.OrgID, args.Mark)
		if err != nil {
			err = errors.New(fmt.Sprint("get ad by auto, ", err))
			errCode = "mark_not_exist"
			return
		}
		return
	}
	//如果存在数据，则投放数据
	return PutAD(&ArgsPutAD{
		OrgID:  args.OrgID,
		AreaID: areaData.ID,
		Mark:   args.Mark,
	})
}

// 自动获取广告设计
// 先从商户获取，否则从平台获取
func getAdByAuto(orgID int64, mark string) (data DataAD, err error) {
	var dataList []DataAD
	if orgID > 0 {
		dataList, err = putADByOrgMark(&argsPutADByOrgMark{
			OrgID: orgID,
			Mark:  mark,
		})
	}
	if orgID < 1 || err != nil {
		dataList, err = putADByGlobMark(&argsPutADByGlobMark{
			Mark: mark,
		})
	}
	//抽取一个符合条件的数据
	if len(dataList) == 1 {
		data = dataList[0]
	} else if len(dataList) > 1 {
		key := CoreFilter.GetRandNumber(0, len(dataList)-1)
		for k, v := range dataList {
			if k == key {
				data = v
				break
			}
		}
		if data.ID < 1 {
			data = dataList[0]
		}
	}
	//如果存在数据，则反馈
	if err == nil {
		if orgID < 1 {
			orgID = 0
		}
		_ = appendAnalysisData(&argsAppendAnalysisData{
			OrgID:  orgID,
			AreaID: 0,
			AdID:   data.ID,
			Count:  1,
		})
		data.Mark = mark
	}
	return
}

// 获取全局mark广告
// 当没有分区、没有
type argsPutADByOrgMark struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
}

func putADByOrgMark(args *argsPutADByOrgMark) (dataList []DataAD, err error) {
	//获取数据
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, mark, name, des, cover_file_id, des_files, params FROM service_ad WHERE mark = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000) LIMIT 1", args.Mark, args.OrgID)
	if err != nil {
		err = errors.New("glob not exist, " + err.Error())
		return
	}
	//覆盖位置信息
	for k, _ := range dataList {
		dataList[k].Mark = args.Mark
	}
	//反馈数据
	return
}

// 获取全局mark广告
// 当没有分区、没有
type argsPutADByGlobMark struct {
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
}

func putADByGlobMark(args *argsPutADByGlobMark) (dataList []DataAD, err error) {
	//获取数据
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, mark, name, des, cover_file_id, des_files, params FROM service_ad WHERE mark = $1 AND org_id = 0 AND delete_at < to_timestamp(1000000) LIMIT 1", args.Mark)
	if err != nil {
		err = errors.New("glob not exist, " + err.Error())
		return
	}
	//覆盖位置信息
	for k, _ := range dataList {
		dataList[k].Mark = args.Mark
	}
	//反馈数据
	return
}
