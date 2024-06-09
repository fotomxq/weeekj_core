package RestaurantRawMaterials

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetRawList 获取Raw列表参数
type ArgsGetRawList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetRawList 获取Raw列表
func GetRawList(args *ArgsGetRawList) (dataList []FieldsRaw, dataCount int64, err error) {
	dataCount, err = rawDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("org_id", args.OrgID).SetSearchQuery([]string{"name"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getRawByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetRawByID 获取Raw数据包参数
type ArgsGetRawByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetRawByID 获取Raw数
func GetRawByID(args *ArgsGetRawByID) (data FieldsRaw, err error) {
	data = getRawByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = errors.New("no data")
		return
	}
	return
}

// GetRawNameByID 获取菜品名称
func GetRawNameByID(id int64) (name string) {
	data := getRawByID(id)
	if data.ID < 1 {
		return
	}
	return data.Name
}

// ArgsCreateRaw 创建Raw参数
type ArgsCreateRaw struct {
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}

// CreateRaw 创建Raw
func CreateRaw(args *ArgsCreateRaw) (id int64, err error) {
	//创建数据
	id, err = rawDB.Insert().SetFields([]string{"org_id", "name"}).Add(map[string]any{
		"org_id": args.OrgID,
		"name":   args.Name,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateRaw 修改Raw参数
type ArgsUpdateRaw struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}

// UpdateRaw 修改Raw
func UpdateRaw(args *ArgsUpdateRaw) (err error) {
	//更新数据
	err = rawDB.Update().SetFields([]string{"org_id", "name"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"org_id": args.OrgID,
		"name":   args.Name,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteRawCache(args.ID)
	//反馈
	return
}

// ArgsDeleteRaw 删除Raw参数
type ArgsDeleteRaw struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteRaw 删除Raw
func DeleteRaw(args *ArgsDeleteRaw) (err error) {
	//删除数据
	err = rawDB.Delete().NeedSoft(true).AddWhereID(args.ID).AddWhereOrgID(args.OrgID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteRawCache(args.ID)
	//反馈
	return
}

// getRawByID 通过ID获取Raw数据包
func getRawByID(id int64) (data FieldsRaw) {
	cacheMark := getRawCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := rawDB.Get().SetDefaultFields().GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheRawTime)
	return
}

// 缓冲
func getRawCacheMark(id int64) string {
	return fmt.Sprint("restaurant:raw_materials:id.", id)
}

func deleteRawCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getRawCacheMark(id))
}
