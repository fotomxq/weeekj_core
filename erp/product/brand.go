package ERPProduct

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetBrandList 获取品牌列表参数
type ArgsGetBrandList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetBrandList 获取品牌列表
func GetBrandList(args *ArgsGetBrandList) (dataList []FieldsBrand, dataCount int64, err error) {
	dataCount, err = brandDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SelectList("((delete_at < to_timestamp(1000000) AND $1 = false) OR (delete_at >= to_timestamp(1000000) AND $1 = true)) AND (name LIKE $2 OR $2 = '')", args.IsRemove, "%"+args.Search+"%").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getBrand(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// GetBrand 获取品牌
func GetBrand(id int64, orgID int64) (data FieldsBrand) {
	data = getBrand(id)
	if data.ID < 1 {
		return
	}
	if !CoreFilter.EqID2(orgID, data.OrgID) {
		data = FieldsBrand{}
		return
	}
	return
}

// GetBrandByCode 通过编码获取品牌
func GetBrandByCode(code string, orgID int64) (data FieldsBrand) {
	_ = brandDB.Get().GetByCodeAndOrgID(code, orgID).Result(&data)
	if data.ID < 1 {
		return
	}
	data = getBrand(data.ID)
	if data.ID < 1 {
		return
	}
	return
}

// ArgsCreateBrand 创建品牌参数
type ArgsCreateBrand struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//品牌编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"300"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}

// CreateBrand 创建品牌
func CreateBrand(args *ArgsCreateBrand) (id int64, err error) {
	//检查编码是否重复
	findCodeData := GetBrandByCode(args.Code, args.OrgID)
	if findCodeData.ID > 0 {
		err = errors.New("code is repeat")
		return
	}
	//创建数据
	id, err = brandDB.Insert().SetFields([]string{"org_id", "code", "name"}).Add(map[string]any{
		"org_id": args.OrgID,
		"code":   args.Code,
		"name":   args.Name,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateBrand 修改品牌参数
type ArgsUpdateBrand struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}

// UpdateBrand 修改品牌
func UpdateBrand(args *ArgsUpdateBrand) (err error) {
	//更新数据
	err = brandDB.Update().SetFields([]string{"name"}).NeedUpdateTime().AddWhereID(args.ID).AddWhereOrgID(args.OrgID).NamedExec(map[string]interface{}{
		"name": args.Name,
	})
	if err != nil {
		return
	}
	deleteBrandCache(args.ID)
	return
}

// ArgsDeleteBrand 删除品牌参数
type ArgsDeleteBrand struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteBrand 删除品牌
func DeleteBrand(args *ArgsDeleteBrand) (err error) {
	err = brandDB.Delete().NeedSoft(true).AddWhereID(args.ID).AddWhereOrgID(args.OrgID).ExecNamed(nil)
	if err != nil {
		return
	}
	deleteBrandCache(args.ID)
	return
}

// 获取品牌数据
func getBrand(id int64) (data FieldsBrand) {
	cacheMark := getBrandCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := brandDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "org_id", "code", "name"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheBrandTime)
	return
}
