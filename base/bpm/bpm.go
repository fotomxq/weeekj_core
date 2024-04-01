package BaseBPM

import (
	"errors"
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetBPMList 获取BPM列表参数
type ArgsGetBPMList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//所属主题
	ThemeID int64 `db:"theme_id" json:"themeId" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetBPMList 获取BPM列表
func GetBPMList(args *ArgsGetBPMList) (dataList []FieldsBPM, dataCount int64, err error) {
	dataCount, err = bpmDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("theme_id", args.ThemeID).SetSearchQuery([]string{"name", "description"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getBPMByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetBPMByID 获取BPM数据包参数
type ArgsGetBPMByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetBPMByID 获取BPM数
func GetBPMByID(args *ArgsGetBPMByID) (data FieldsBPM, err error) {
	data = getBPMByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreateBPM 创建BPM参数
type ArgsCreateBPM struct {
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
	//所属主题
	ThemeID int64 `db:"theme_id" json:"themeId" check:"id"`
	//流程节点数量
	NodeCount int `db:"node_count" json:"nodeCount" check:"int64Than0"`
	//流程json文件内容
	JSONNode string `db:"json_node" json:"jsonNode"`
}

// CreateBPM 创建BPM
func CreateBPM(args *ArgsCreateBPM) (id int64, err error) {
	//创建数据
	id, err = bpmDB.Insert().SetFields([]string{"name", "description", "theme_id", "node_count", "json_node"}).Add(map[string]any{
		"name":        args.Name,
		"description": args.Description,
		"theme_id":    args.ThemeID,
		"node_count":  args.NodeCount,
		"json_node":   args.JSONNode,
	}).ExecAndResultID()
	if err != nil {
		fmt.Println(err)
		return
	}
	//反馈
	return
}

// ArgsUpdateBPM 修改BPM参数
type ArgsUpdateBPM struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
	//所属主题
	ThemeID int64 `db:"theme_id" json:"themeId" check:"id"`
	//流程节点数量
	NodeCount int `db:"node_count" json:"nodeCount" check:"int64Than0"`
	//流程json文件内容
	JSONNode string `db:"json_node" json:"jsonNode"`
}

// UpdateBPM 修改BPM
func UpdateBPM(args *ArgsUpdateBPM) (err error) {
	//更新数据
	err = bpmDB.Update().SetFields([]string{"name", "description", "theme_id", "node_count", "json_node"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]interface{}{
		"name":        args.Name,
		"description": args.Description,
		"theme_id":    args.ThemeID,
		"node_count":  args.NodeCount,
		"json_node":   args.JSONNode,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteBPMCache(args.ID)
	//反馈
	return
}

// ArgsDeleteBPM 删除BPM参数
type ArgsDeleteBPM struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteBPM 删除BPM
func DeleteBPM(args *ArgsDeleteBPM) (err error) {
	//删除数据
	err = bpmDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteBPMCache(args.ID)
	//反馈
	return
}

// GetBPMCountByThemeID 获取主题下的bpm数量
func GetBPMCountByThemeID(themeID int64) (count int64) {
	count, _ = bpmDB.Select().SetFieldsList([]string{"id"}).SetDeleteQuery("delete_at", false).SetIDQuery("theme_id", themeID).SetPages(CoreSQL2.ArgsPages{
		Page: 1,
		Max:  1,
		Sort: "id",
		Desc: false,
	}).SelectList("").ResultCount()
	return
}

// getBPMByID 通过ID获取bpm数据包
func getBPMByID(id int64) (data FieldsBPM) {
	cacheMark := getBPMCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := bpmDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "name", "description", "theme_id", "node_count", "json_node"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheBPMTime)
	return
}

// 缓冲
func getBPMCacheMark(id int64) string {
	return fmt.Sprint("base:bpm:id.", id)
}

func deleteBPMCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getBPMCacheMark(id))
}
