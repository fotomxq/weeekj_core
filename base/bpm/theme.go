package BaseBPM

import (
	"errors"
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetThemeList 获取Theme列表参数
type ArgsGetThemeList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//所属主题分类
	CategoryID int64 `db:"category_id" json:"categoryId" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetThemeList 获取Theme列表
func GetThemeList(args *ArgsGetThemeList) (dataList []FieldsTheme, dataCount int64, err error) {
	dataCount, err = themeDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("category_id", args.CategoryID).SetSearchQuery([]string{"name", "description"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getThemeByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetThemeByID 获取Theme数据包参数
type ArgsGetThemeByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetThemeByID 获取Theme数
func GetThemeByID(args *ArgsGetThemeByID) (data FieldsTheme, err error) {
	data = getThemeByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreateTheme 创建Theme参数
type ArgsCreateTheme struct {
	//所属主题分类
	CategoryID int64 `db:"category_id" json:"categoryId" check:"id"`
	//主题名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//主题描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
}

// CreateTheme 创建Theme
func CreateTheme(args *ArgsCreateTheme) (id int64, err error) {
	//创建数据
	id, err = themeDB.Insert().SetFields([]string{"category_id", "name", "description"}).Add(map[string]any{
		"category_id": args.CategoryID,
		"name":        args.Name,
		"description": args.Description,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateTheme 修改Theme参数
type ArgsUpdateTheme struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//所属主题分类
	CategoryID int64 `db:"category_id" json:"categoryId" check:"id"`
	//主题名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//主题描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
}

// UpdateTheme 修改Theme
func UpdateTheme(args *ArgsUpdateTheme) (err error) {
	//更新数据
	err = themeDB.Update().SetFields([]string{"category_id", "name", "description"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"category_id": args.CategoryID,
		"name":        args.Name,
		"description": args.Description,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteThemeCache(args.ID)
	//反馈
	return
}

// ArgsDeleteTheme 删除Theme参数
type ArgsDeleteTheme struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteTheme 删除Theme
func DeleteTheme(args *ArgsDeleteTheme) (err error) {
	//删除数据
	err = themeDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteThemeCache(args.ID)
	//反馈
	return
}

// getThemeByID 通过ID获取Theme数据包
func getThemeByID(id int64) (data FieldsTheme) {
	cacheMark := getThemeCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := themeDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "category_id", "name", "description"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheThemeTime)
	return
}

// 缓冲
func getThemeCacheMark(id int64) string {
	return fmt.Sprint("base:bpm:theme:id.", id)
}

func deleteThemeCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getThemeCacheMark(id))
}
