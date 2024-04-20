package ERPProduct

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetTemplateBindList 获取绑定关系列表参数
type ArgsGetTemplateBindList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//模板ID
	TemplateID int64 `db:"template_id" json:"templateID" check:"id" empty:"true"`
	//分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id" empty:"true"`
	//品牌ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id" empty:"true"`
	//规格型号
	// 三选一，分类ID、品牌ID、规格型号ID
	ModelTypeID int64 `db:"model_type_id" json:"modelTypeID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
}

// GetTemplateBindList 获取绑定关系列表
func GetTemplateBindList(args *ArgsGetTemplateBindList) (dataList []FieldsTemplateBind, dataCount int64, err error) {
	dataCount, err = templateBindDB.Select().SetFieldsList([]string{"id", "org_id", "template_id", "category_id", "brand_id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at"}).SetPages(args.Pages).SelectList("((delete_at < to_timestamp(1000000) AND $1 = false) OR (delete_at >= to_timestamp(1000000) AND $1 = true)) AND (org_id = $2 OR $2 < 0) AND (template_id = $3 OR $3 < 0) AND (category_id = $4 OR $4 < 0) AND (brand_id = $5 OR $5 < 0) AND (model_type_id = $6 OR $6 < 0)", args.IsRemove, args.OrgID, args.TemplateID, args.CategoryID, args.BrandID, args.ModelTypeID).ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := GetTemplateBindData(&ArgsGetTemplateBindData{
			OrgID:       v.OrgID,
			TemplateID:  v.TemplateID,
			CategoryID:  v.CategoryID,
			BrandID:     v.BrandID,
			ModelTypeID: v.ModelTypeID,
		})
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// getTemplateBindRecursionByCategoryID 查询分类对应的绑定关系
// 给与产品最初的分类ID，从最底层到最高层追溯到绑定模板关系
func getTemplateBindRecursionByCategoryID(orgID int64, categoryID int64) (data FieldsTemplateBind) {
	var dataList []FieldsTemplateBind
	_ = templateBindDB.Select().SetFieldsList([]string{"id", "create_at", "update_at", "delete_at", "org_id", "template_id", "brand_id"}).SetFieldsSort([]string{"id"}).SetPages(CoreSQL2.ArgsPages{
		Page: 1,
		Max:  1,
		Sort: "id",
		Desc: false,
	}).SetDeleteQuery("delete_at", false).SetIDQuery("org_id", orgID).SetIDQuery("category_id", categoryID).Result(&dataList)
	if len(dataList) < 1 {
		categoryData := Sort.GetByIDNoErr(categoryID, orgID)
		if categoryData.ID < 1 {
			return
		} else {
			if categoryData.ParentID < 1 {
				return
			} else {
				return getTemplateBindRecursionByCategoryID(orgID, categoryData.ParentID)
			}
		}
	}
	data = dataList[0]
	return
}

// ArgsGetTemplateBindData 获取品牌绑定关系参数
type ArgsGetTemplateBindData struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//模板ID
	TemplateID int64 `db:"template_id" json:"templateID" check:"id"`
	//分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id" empty:"true"`
	//品牌ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id" empty:"true"`
	//规格型号
	// 三选一，分类ID、品牌ID、规格型号ID
	ModelTypeID int64 `db:"model_type_id" json:"modelTypeID" check:"id" empty:"true"`
}

// GetTemplateBindData 获取品牌绑定关系
func GetTemplateBindData(args *ArgsGetTemplateBindData) (data FieldsTemplateBind) {
	//必须三选一
	if args.CategoryID < 1 && args.BrandID < 1 && args.ModelTypeID < 1 {
		return
	}
	//获取缓冲
	cacheMark := getTemplateBindCacheMark(args.OrgID, args.TemplateID, args.CategoryID, args.BrandID, args.ModelTypeID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	//获取数据
	err := templateBindDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "org_id", "template_id", "category_id", "brand_id", "model_type_id"}).AppendWhere("(org_id = $1 OR $1 < 0) AND template_id = $2 AND (category_id = $3 OR $3 < 0) AND (brand_id = $4 OR $4 < 0) AND (model_type_id = $5 OR $5 < 0)", args.OrgID, args.TemplateID, args.CategoryID, args.BrandID, args.ModelTypeID).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	//保存缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheTemplateBindTime)
	//反馈
	return
}

// ArgsCheckTemplateBind 检查品牌绑定关系参数
type ArgsCheckTemplateBind struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//模板ID
	TemplateID int64 `db:"template_id" json:"templateID" check:"id"`
	//分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id" empty:"true"`
	//品牌ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id" empty:"true"`
	//规格型号
	// 三选一，分类ID、品牌ID、规格型号ID
	ModelTypeID int64 `db:"model_type_id" json:"modelTypeID" check:"id" empty:"true"`
}

// CheckTemplateBind 检查品牌绑定关系
func CheckTemplateBind(args *ArgsCheckTemplateBind) (b bool) {
	//获取数据
	data := GetTemplateBindData(&ArgsGetTemplateBindData{
		OrgID:       args.OrgID,
		TemplateID:  args.TemplateID,
		CategoryID:  args.CategoryID,
		BrandID:     args.BrandID,
		ModelTypeID: args.ModelTypeID,
	})
	if data.ID < 1 {
		return
	}
	b = true
	//反馈
	return
}

// ArgsCreateTemplateBind 添加新品牌绑定关系参数
type ArgsCreateTemplateBind struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//模板ID
	TemplateID int64 `db:"template_id" json:"templateID" check:"id"`
	//分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id" empty:"true"`
	//品牌ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id" empty:"true"`
	//规格型号
	// 三选一，分类ID、品牌ID、规格型号ID
	ModelTypeID int64 `db:"model_type_id" json:"modelTypeID" check:"id" empty:"true"`
}

// CreateTemplateBind 添加新品牌绑定关系
func CreateTemplateBind(args *ArgsCreateTemplateBind) (id int64, err error) {
	//检查数据
	data := GetTemplateBindData(&ArgsGetTemplateBindData{
		OrgID:       args.OrgID,
		TemplateID:  args.TemplateID,
		CategoryID:  args.CategoryID,
		BrandID:     args.BrandID,
		ModelTypeID: args.ModelTypeID,
	})
	if data.ID > 0 {
		if CoreFilter.CheckHaveTime(data.DeleteAt) {
			id = data.ID
			err = templateBindDB.Update().NeedSoft(false).NeedUpdateTime().AddWhereID(data.ID).SetFields([]string{"delete_at"}).NamedExec(map[string]any{
				"delete_at": time.Time{},
			})
			return
		} else {
			if args.TemplateID == data.TemplateID && args.CategoryID == data.CategoryID && args.BrandID == data.BrandID && args.ModelTypeID == data.ModelTypeID {
				return
			} else {
				err = errors.New("have replace")
				return
			}
		}
	} else {
		//创建数据
		id, err = templateBindDB.Insert().SetFields([]string{"org_id", "template_id", "category_id", "brand_id", "model_type_id"}).Add(map[string]any{
			"org_id":        args.OrgID,
			"template_id":   args.TemplateID,
			"category_id":   args.CategoryID,
			"brand_id":      args.BrandID,
			"model_type_id": args.ModelTypeID,
		}).ExecAndResultID()
		if err != nil {
			return
		}
	}
	//反馈
	return
}

// ArgsDeleteTemplateBind 删除产品绑定关系参数
type ArgsDeleteTemplateBind struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//模板ID
	TemplateID int64 `db:"template_id" json:"templateID" check:"id"`
	//分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id" empty:"true"`
	//品牌ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id" empty:"true"`
	//规格型号
	// 三选一，分类ID、品牌ID、规格型号ID
	ModelTypeID int64 `db:"model_type_id" json:"modelTypeID" check:"id" empty:"true"`
}

// DeleteTemplateBind 删除产品绑定关系
func DeleteTemplateBind(args *ArgsDeleteTemplateBind) (err error) {
	data := GetTemplateBindData(&ArgsGetTemplateBindData{
		OrgID:       args.OrgID,
		TemplateID:  args.TemplateID,
		CategoryID:  args.CategoryID,
		BrandID:     args.BrandID,
		ModelTypeID: args.ModelTypeID,
	})
	if data.ID < 1 {
		return
	}
	err = templateBindDB.Delete().NeedSoft(true).AddWhereID(data.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	deleteTemplateBindCache(args.OrgID, args.TemplateID, args.CategoryID, args.BrandID, args.ModelTypeID)
	return
}

// 缓冲
func getTemplateBindCacheMark(orgID, templateID, categoryID, brandID, modelTypeID int64) string {
	return fmt.Sprint("erp:product:template:bind:org.", orgID, ".template.", templateID, ".category.", categoryID, ".brand.", brandID, ".modelType.", modelTypeID)
}

func deleteTemplateBindCache(orgID, templateID, categoryID, brandID, modelTypeID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getTemplateBindCacheMark(orgID, templateID, categoryID, brandID, modelTypeID))
}
