package BaseSystemMission

import (
	"errors"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
)

// ArgsGetMissionList 获取服务列表参数
type ArgsGetMissionList struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetMissionList 获取服务列表
func GetMissionList(args *ArgsGetMissionList) (dataList []FieldsMission, dataCount int64, err error) {
	var rawList []FieldsMission
	dataCount, err = missionDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "update_at"}).SetPages(args.Pages).SelectList("($1 < 0 OR org_id = $1) AND (name ILIKE '%' || $2 || '%')", args.OrgID, args.Search).ResultAndCount(&rawList)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getMission(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetMissionByMark 获取指定服务
func GetMissionByMark(orgID int64, mark string) (data FieldsMission) {
	data = getMissionByMark(mark, orgID)
	return
}
