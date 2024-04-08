package BaseService

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisList 获取分析列表参数
type ArgsGetAnalysisList struct {
	//服务ID
	ServiceID int64 `json:"serviceID" check:"id"`
	//时间段
	BetweenAt CoreSQL2.ArgsTimeBetween `json:"betweenAt"`
}

// GetAnalysisList 获取分析列表
func GetAnalysisList(args *ArgsGetAnalysisList) (dataList []FieldsAnalysis, dataCount int64, err error) {
	dataCount, err = analysisDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"create_at"}).SetPages(CoreSQL2.ArgsPages{
		Page: 1,
		Max:  9999,
		Sort: "create_at",
		Desc: false,
	}).SetIDQuery("service_id", args.ServiceID).SetTimeBetweenByArgQuery("create_at", args.BetweenAt).SelectList("").ResultAndCount(&dataList)
	if err != nil {
		return
	}
	for k, v := range dataList {
		vData := getAnalysisByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// argsAppendAnalysisData 追加分析数据参数
type argsAppendAnalysisData struct {
	//服务ID
	ServiceID int64 `db:"service_id" json:"serviceID" check:"id"`
	//增加服务端发送消息次数
	SendCount int64 `db:"send_count" json:"sendCount" check:"intThan0"`
	//增加服务端接收次数
	ReceiveCount int64 `db:"receive_count" json:"receiveCount" check:"intThan0"`
}

// appendAnalysisData 追加分析数据
func appendAnalysisData(args *argsAppendAnalysisData) (err error) {
	nowAt := CoreFilter.GetNowTimeCarbon()
	minAt := nowAt.StartOfHour()
	maxAt := nowAt.EndOfHour()
	var dataList []FieldsAnalysis
	_, err = analysisDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"create_at"}).SetPages(CoreSQL2.ArgsPages{
		Page: 1,
		Max:  1,
		Sort: "create_at",
		Desc: false,
	}).SetIDQuery("service_id", args.ServiceID).SetTimeBetweenQuery("create_at", minAt.Time, maxAt.Time).SelectList("").ResultAndCount(&dataList)
	if err == nil && len(dataList) > 0 {
		data := dataList[0]
		data.SendCount += args.SendCount
		data.ReceiveCount += args.ReceiveCount
		err = analysisDB.Update().AddWhereID(data.ID).SetFields([]string{"send_count", "receive_count"}).NeedSoft(false).NamedExec(map[string]any{
			"send_count":    data.SendCount,
			"receive_count": data.ReceiveCount,
		})
		if err != nil {
			return
		}
	} else {
		err = analysisDB.Insert().SetFields([]string{"service_id", "send_count", "receive_count"}).Add(map[string]any{
			"service_id":    args.ServiceID,
			"send_count":    args.SendCount,
			"receive_count": args.ReceiveCount,
		}).ExecAndCheckID()
		if err != nil {
			return
		}
	}
	return
}

// getAnalysisByID 通过ID获取分析数据包
func getAnalysisByID(id int64) (data FieldsAnalysis) {
	cacheMark := getAnalysisCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := analysisDB.Get().SetFieldsOne([]string{"id", "create_at", "service_id", "send_count", "receive_count"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheAnalysisTime)
	return
}

// 缓冲
func getAnalysisCacheMark(id int64) string {
	return fmt.Sprint("base:service:analysis:id.", id)
}

func deleteAnalysisCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getAnalysisCacheMark(id))
}
