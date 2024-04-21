package ERPProduct

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetModelTypeList 获取规格列表参数
type ArgsGetModelTypeList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//规格ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetModelTypeList 获取规格列表
func GetModelTypeList(args *ArgsGetModelTypeList) (dataList []FieldsModelType, dataCount int64, err error) {
	dataCount, err = modelTypeDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("org_id", args.OrgID).SetIDQuery("brand_id", args.BrandID).SetSearchQuery([]string{"name"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getModelType(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// GetModelType 获取规格
func GetModelType(id int64, orgID int64) (data FieldsModelType) {
	data = getModelType(id)
	if data.ID < 1 {
		return
	}
	if !CoreFilter.EqID2(orgID, data.OrgID) {
		data = FieldsModelType{}
		return
	}
	return
}

// GetModelTypeByCode 通过编码获取规格
func GetModelTypeByCode(code string, orgID int64) (data FieldsModelType) {
	_ = modelTypeDB.Get().GetByCodeAndOrgID(code, orgID).Result(&data)
	if data.ID < 1 {
		return
	}
	data = getModelType(data.ID)
	if data.ID < 1 {
		return
	}
	return
}

// ArgsCreateModelType 创建规格参数
type ArgsCreateModelType struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//规格编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"300"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//规格ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id"`
}

// CreateModelType 创建规格
func CreateModelType(args *ArgsCreateModelType) (id int64, err error) {
	//检查编码是否重复
	findCodeData := GetModelTypeByCode(args.Code, args.OrgID)
	if findCodeData.ID > 0 {
		err = errors.New("code is repeat")
		return
	}
	//创建数据
	id, err = modelTypeDB.Insert().SetFields([]string{"org_id", "code", "name", "brand_id"}).Add(map[string]any{
		"org_id":   args.OrgID,
		"code":     args.Code,
		"name":     args.Name,
		"brand_id": args.BrandID,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateModelType 修改规格参数
type ArgsUpdateModelType struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}

// UpdateModelType 修改规格
func UpdateModelType(args *ArgsUpdateModelType) (err error) {
	//更新数据
	err = modelTypeDB.Update().SetFields([]string{"name"}).NeedUpdateTime().AddWhereID(args.ID).AddWhereOrgID(args.OrgID).NamedExec(map[string]interface{}{
		"name": args.Name,
	})
	if err != nil {
		return
	}
	deleteModelTypeCache(args.ID)
	return
}

// ArgsDeleteModelType 删除规格参数
type ArgsDeleteModelType struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteModelType 删除规格
func DeleteModelType(args *ArgsDeleteModelType) (err error) {
	err = modelTypeDB.Delete().NeedSoft(true).AddWhereID(args.ID).AddWhereOrgID(args.OrgID).ExecNamed(nil)
	if err != nil {
		return
	}
	deleteModelTypeCache(args.ID)
	return
}

// 获取规格数据
func getModelType(id int64) (data FieldsModelType) {
	cacheMark := getModelTypeCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := modelTypeDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "org_id", "code", "name", "brand_id"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheModelTypeTime)
	return
}

// 缓冲
func getModelTypeCacheMark(id int64) string {
	return fmt.Sprint("erp:product:model:type:id.", id)
}

func deleteModelTypeCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getModelTypeCacheMark(id))
}
