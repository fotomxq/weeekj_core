package ERPProduct

import (
	"errors"
	BaseBPM "github.com/fotomxq/weeekj_core/v5/base/bpm"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetTemplateList 获取模板列表
type ArgsGetTemplateList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//插槽主题ID
	// BPM模块插槽主题ID，用于关联插槽主题，产品会自动使用相关的插槽用于表单实现
	BPMThemeID int64 `db:"bpm_theme_id" json:"bpmThemeID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetTemplateList 获取品牌列表
func GetTemplateList(args *ArgsGetTemplateList) (dataList []FieldsTemplate, dataCount int64, err error) {
	dataCount, err = templateDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetSearchQuery([]string{"name"}, args.Search).SetIDQuery("bpm_theme_id", args.BPMThemeID).SetIDQuery("org_id", args.OrgID).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getTemplateData(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// GetTemplate 获取模板
func GetTemplate(id int64, orgID int64) (data FieldsTemplate) {
	data = getTemplateData(id)
	if data.ID < 1 {
		return
	}
	if !CoreFilter.EqID2(orgID, data.OrgID) {
		data = FieldsTemplate{}
		return
	}
	return
}

// GetTemplateBPMThemeSlotData 获取模板的插槽数据包
func GetTemplateBPMThemeSlotData(orgID int64, templateID int64) (bpmSlotList []BaseBPM.FieldsSlot, errCode string, err error) {
	//通过绑定关系获取模板数据包
	templateData := GetTemplate(templateID, orgID)
	if templateData.ID < 1 {
		errCode = "err_erp_product_no_exist_template"
		err = errors.New("product bind template not exist")
		return
	}
	//如果模板拿到关联主题
	bmpThemeData, _ := BaseBPM.GetThemeByID(&BaseBPM.ArgsGetThemeByID{
		ID: templateData.BPMThemeID,
	})
	if bmpThemeData.ID < 1 {
		errCode = "err_erp_product_no_exist_bpm_theme"
		err = errors.New("product bind template not exist bpm theme")
		return
	}
	//通过BPM主题，查询关联的插槽
	bpmSlotList, _, err = BaseBPM.GetSlotList(&BaseBPM.ArgsGetSlotList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  999,
			Sort: "id",
			Desc: false,
		},
		ThemeCategoryID: -1,
		ThemeID:         bmpThemeData.ID,
		IsRemove:        false,
		Search:          "",
	})
	if err != nil {
		errCode = "err_erp_product_bpm_slot_data"
		return
	}
	//反馈
	return
}

// ArgsCreateTemplate 创建模板参数
type ArgsCreateTemplate struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//插槽主题ID
	// BPM模块插槽主题ID，用于关联插槽主题，产品会自动使用相关的插槽用于表单实现
	BPMThemeID int64 `db:"bpm_theme_id" json:"bpmThemeID" check:"id" empty:"true"`
}

// CreateTemplate 创建模板
func CreateTemplate(args *ArgsCreateTemplate) (id int64, err error) {
	//创建数据
	id, err = templateDB.Insert().SetFields([]string{"org_id", "name", "bpm_theme_id"}).Add(map[string]any{
		"org_id":       args.OrgID,
		"name":         args.Name,
		"bpm_theme_id": args.BPMThemeID,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateTemplate 更新模板参数
type ArgsUpdateTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//插槽主题ID
	// BPM模块插槽主题ID，用于关联插槽主题，产品会自动使用相关的插槽用于表单实现
	BPMThemeID int64 `db:"bpm_theme_id" json:"bpmThemeID" check:"id" empty:"true"`
}

// UpdateTemplate 更新模板
func UpdateTemplate(args *ArgsUpdateTemplate) (err error) {
	//更新数据
	err = templateDB.Update().SetFields([]string{"name", "bpm_theme_id"}).NeedUpdateTime().AddWhereID(args.ID).AddWhereOrgID(args.OrgID).NamedExec(map[string]interface{}{
		"name":         args.Name,
		"bpm_theme_id": args.BPMThemeID,
	})
	if err != nil {
		return
	}
	deleteTemplateCache(args.ID)
	return
}

// ArgsDeleteTemplate 删除模板参数
type ArgsDeleteTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteTemplate 删除模板
func DeleteTemplate(args *ArgsDeleteTemplate) (err error) {
	err = templateDB.Delete().NeedSoft(true).AddWhereID(args.ID).AddWhereOrgID(args.OrgID).ExecNamed(nil)
	if err != nil {
		return
	}
	deleteTemplateCache(args.ID)
	return
}

// getTemplateData 获取模板数据
func getTemplateData(id int64) (data FieldsTemplate) {
	cacheMark := getTemplateCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := templateDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "org_id", "name", "bpm_theme_id"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheTemplateTime)
	return
}
