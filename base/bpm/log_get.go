package BaseBPM

import (
	"errors"
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetLogList 获取Log列表参数
type ArgsGetLogList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//管理单元
	UnitID int64 `db:"unit_id" json:"unitID" check:"id" empty:"true"`
	//操作用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//操作组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//BPM ID
	BPMID int64 `db:"bpm_id" json:"bpmId" check:"id" check:"id" empty:"true"`
}

// GetLogList 获取Log列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	dataCount, err = logDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at"}).SetPages(args.Pages).SetIDQuery("org_id", args.OrgID).SetIDQuery("unit_id", args.UnitID).SetIDQuery("user_id", args.UserID).SetIDQuery("org_bind_id", args.OrgBindID).SetIDQuery("bpm_id", args.BPMID).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getLogByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetLogByID 获取Log数据包参数
type ArgsGetLogByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetLogByID 获取Log数
func GetLogByID(args *ArgsGetLogByID) (data FieldsLog, err error) {
	data = getLogByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// getLogByID 通过ID获取Log数据包
func getLogByID(id int64) (data FieldsLog) {
	cacheMark := getLogCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := logDB.Get().SetFieldsOne([]string{"id", "create_at", "org_id", "unit_id", "user_id", "org_bind_id", "bpm_id", "node_id", "node_number", "node_content"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheLogTime)
	return
}

// 缓冲
func getLogCacheMark(id int64) string {
	return fmt.Sprint("base:bpm:log:id.", id)
}

func deleteLogCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getLogCacheMark(id))
}
