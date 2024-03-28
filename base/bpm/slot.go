package BaseBPM

import (
	"errors"
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetSlotList 获取Slot列表参数
type ArgsGetSlotList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//所属主题分类
	ThemeCategoryID int64 `db:"theme_category_id" json:"themeCategoryId" check:"id" empty:"true"`
	//所属主题
	ThemeID int64 `db:"theme_id" json:"themeId" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetSlotList 获取Slot列表
func GetSlotList(args *ArgsGetSlotList) (dataList []FieldsSlot, dataCount int64, err error) {
	dataCount, err = slotDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("theme_category_id", args.ThemeCategoryID).SetIDQuery("theme_id", args.ThemeID).SetSearchQuery([]string{"name"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getSlotByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetSlotByID 获取Slot数据包参数
type ArgsGetSlotByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetSlotByID 获取Slot数
func GetSlotByID(args *ArgsGetSlotByID) (data FieldsSlot, err error) {
	data = getSlotByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreateSlot 创建Slot参数
type ArgsCreateSlot struct {
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//所属主题分类
	ThemeCategoryID int64 `db:"theme_category_id" json:"themeCategoryId" check:"id"`
	//所属主题
	// 插槽可用于的主题域
	ThemeID int64 `db:"theme_id" json:"themeId" check:"id"`
	//值类型
	// 插槽的值类型
	// input 输入框; text 文本域; radio 单选项; checkbox 多选项; select 下拉单选框;
	// date: 日期; time: 时间; datetime: 日期时间;
	// file: 文件ID; files 文件ID列; image: 图片; images 图片列; audio: 音频; video: 视频; videos 视频列; url: URL;
	// email: 邮箱; phone: 电话; id: ID; password: 密码;
	// code: 代码; html: HTML; markdown: Markdown; xml: XML; yaml: YAML;
	ValueType string `db:"value_type" json:"valueType" check:"des" min:"1" max:"3000"`
	//默认值
	// 插槽的默认值
	DefaultValue string `db:"default_value" json:"defaultValue" check:"des" min:"1" max:"3000"`
	//参数
	// 根据组件需求，自定义参数内容
	Params string `db:"params" json:"params"`
}

// CreateSlot 创建Slot
func CreateSlot(args *ArgsCreateSlot) (id int64, err error) {
	//创建数据
	id, err = slotDB.Insert().SetFields([]string{"name", "description", "theme_category_id", "theme_id", "value_type", "default_value", "params"}).Add(map[string]any{
		"name":              args.Name,
		"description":       "",
		"theme_category_id": args.ThemeCategoryID,
		"theme_id":          args.ThemeID,
		"value_type":        args.ValueType,
		"default_value":     args.DefaultValue,
		"params":            args.Params,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateSlot 修改Slot参数
type ArgsUpdateSlot struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//所属主题分类
	ThemeCategoryID int64 `db:"theme_category_id" json:"themeCategoryId" check:"id"`
	//所属主题
	// 插槽可用于的主题域
	ThemeID int64 `db:"theme_id" json:"themeId" check:"id"`
	//值类型
	// 插槽的值类型
	// input 输入框; text 文本域; radio 单选项; checkbox 多选项; select 下拉单选框;
	// date: 日期; time: 时间; datetime: 日期时间;
	// file: 文件ID; files 文件ID列; image: 图片; images 图片列; audio: 音频; video: 视频; videos 视频列; url: URL;
	// email: 邮箱; phone: 电话; id: ID; password: 密码;
	// code: 代码; html: HTML; markdown: Markdown; xml: XML; yaml: YAML;
	ValueType string `db:"value_type" json:"valueType" check:"des" min:"1" max:"3000"`
	//默认值
	// 插槽的默认值
	DefaultValue string `db:"default_value" json:"defaultValue" check:"des" min:"1" max:"3000"`
	//参数
	// 根据组件需求，自定义参数内容
	Params string `db:"params" json:"params"`
}

// UpdateSlot 修改Slot
func UpdateSlot(args *ArgsUpdateSlot) (err error) {
	//更新数据
	err = slotDB.Update().SetFields([]string{"description", "name", "theme_category_id", "theme_id", "value_type", "default_value", "params"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"name":              args.Name,
		"description":       "",
		"theme_category_id": args.ThemeCategoryID,
		"theme_id":          args.ThemeID,
		"value_type":        args.ValueType,
		"default_value":     args.DefaultValue,
		"params":            args.Params,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteSlotCache(args.ID)
	//反馈
	return
}

// ArgsDeleteSlot 删除Slot参数
type ArgsDeleteSlot struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteSlot 删除Slot
func DeleteSlot(args *ArgsDeleteSlot) (err error) {
	//删除数据
	err = slotDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteSlotCache(args.ID)
	//反馈
	return
}

// GetSlotCountByCategoryID 获取分类下的插槽数量
func GetSlotCountByCategoryID(categoryID int64) (count int64) {
	count, _ = slotDB.Select().SetFieldsList([]string{"id"}).SetIDQuery("theme_category_id", categoryID).SetPages(CoreSQL2.ArgsPages{
		Page: 1,
		Max:  1,
		Sort: "id",
		Desc: false,
	}).SelectList("").ResultCount()
	return
}

// GetSlotCountByThemeID 获取主题下的插槽数量
func GetSlotCountByThemeID(themeID int64) (count int64) {
	count, _ = slotDB.Select().SetFieldsList([]string{"id"}).SetIDQuery("theme_id", themeID).SetPages(CoreSQL2.ArgsPages{
		Page: 1,
		Max:  1,
		Sort: "id",
		Desc: false,
	}).SelectList("").ResultCount()
	return
}

// getSlotByID 通过ID获取Slot数据包
func getSlotByID(id int64) (data FieldsSlot) {
	cacheMark := getSlotCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := slotDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "description", "name", "theme_category_id", "theme_id", "value_type", "default_value", "params"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheSlotTime)
	return
}

// 缓冲
func getSlotCacheMark(id int64) string {
	return fmt.Sprint("base:bpm:slot:id.", id)
}

func deleteSlotCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getSlotCacheMark(id))
}
