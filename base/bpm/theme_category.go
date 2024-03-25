package BaseBPM

import (
	"errors"
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetThemeCategoryList 获取ThemeCategory列表参数
type ArgsGetThemeCategoryList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetThemeCategoryList 获取ThemeCategory列表
func GetThemeCategoryList(args *ArgsGetThemeCategoryList) (dataList []FieldsThemeCategory, dataCount int64, err error) {
	dataCount, err = themeCategoryDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetSearchQuery([]string{"name", "description"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getThemeCategoryByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetThemeCategoryByID 获取ThemeCategory数据包参数
type ArgsGetThemeCategoryByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetThemeCategoryByID 获取ThemeCategory数
func GetThemeCategoryByID(args *ArgsGetThemeCategoryByID) (data FieldsThemeCategory, err error) {
	data = getThemeCategoryByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreateThemeCategory 创建ThemeCategory参数
type ArgsCreateThemeCategory struct {
	//主题名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//主题描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
}

// CreateThemeCategory 创建ThemeCategory
func CreateThemeCategory(args *ArgsCreateThemeCategory) (id int64, err error) {
	//创建数据
	id, err = themeCategoryDB.Insert().SetFields([]string{"name", "description"}).Add(map[string]any{
		"name":        args.Name,
		"description": args.Description,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateThemeCategory 修改ThemeCategory参数
type ArgsUpdateThemeCategory struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//主题名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//主题描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
}

// UpdateThemeCategory 修改ThemeCategory
func UpdateThemeCategory(args *ArgsUpdateThemeCategory) (err error) {
	//更新数据
	err = themeCategoryDB.Update().SetFields([]string{"name", "description"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"name":        args.Name,
		"description": args.Description,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteThemeCategoryCache(args.ID)
	//反馈
	return
}

// ArgsDeleteThemeCategory 删除ThemeCategory参数
type ArgsDeleteThemeCategory struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteThemeCategory 删除ThemeCategory
func DeleteThemeCategory(args *ArgsDeleteThemeCategory) (err error) {
	//删除数据
	err = themeCategoryDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteThemeCategoryCache(args.ID)
	//反馈
	return
}

// getThemeCategoryByID 通过ID获取ThemeCategory数据包
func getThemeCategoryByID(id int64) (data FieldsThemeCategory) {
	cacheMark := getThemeCategoryCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := themeCategoryDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "name", "description"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheThemeCategoryTime)
	return
}

// 缓冲
func getThemeCategoryCacheMark(id int64) string {
	return fmt.Sprint("base:bpm:theme:category:id.", id)
}

func deleteThemeCategoryCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getThemeCategoryCacheMark(id))
}
