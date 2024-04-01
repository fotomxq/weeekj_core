package BaseBPM

import (
	"errors"
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetEventList 获取Event列表参数
type ArgsGetEventList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//所属主题分类
	ThemeCategoryID int64 `db:"theme_category_id" json:"themeCategoryId" check:"id" empty:"true"`
	//所属主题
	ThemeID int64 `db:"theme_id" json:"themeId" check:"id" empty:"true"`
	//事件编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"300" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetEventList 获取Event列表
func GetEventList(args *ArgsGetEventList) (dataList []FieldsEvent, dataCount int64, err error) {
	dataCount, err = eventDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("theme_category_id", args.ThemeCategoryID).SetIDQuery("theme_id", args.ThemeID).SetStringQuery("code", args.Code).SetSearchQuery([]string{"name", "description"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getEventByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetEventByID 获取Event数据包参数
type ArgsGetEventByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetEventByID 获取Event数
func GetEventByID(args *ArgsGetEventByID) (data FieldsEvent, err error) {
	data = getEventByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreateEvent 创建Event参数
type ArgsCreateEvent struct {
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
	//所属主题分类
	ThemeCategoryID int64 `db:"theme_category_id" json:"themeCategoryId" check:"id"`
	//所属主题
	// 插槽可用于的主题域
	ThemeID int64 `db:"theme_id" json:"themeId" check:"id"`
	//事件编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"300"`
	//事件类型
	// nats - NATS事件
	EventType string `db:"event_type" json:"eventType" check:"intThan0"`
	//事件地址
	// nats - 触发的地址
	EventURL string `db:"event_url" json:"eventURL" check:"des" min:"1" max:"600"`
	//事件固定参数
	// nats - 事件附带的固定参数，如果为空则根据流程阶段事件触发填入
	EventParams string `db:"event_params" json:"eventParams" check:"des" min:"1" max:"1000" empty:"true"`
}

// CreateEvent 创建Event
func CreateEvent(args *ArgsCreateEvent) (id int64, err error) {
	//创建数据
	id, err = eventDB.Insert().SetFields([]string{"name", "description", "theme_category_id", "theme_id", "code", "event_type", "event_url", "event_params"}).Add(map[string]any{
		"name":              args.Name,
		"description":       args.Description,
		"theme_category_id": args.ThemeCategoryID,
		"theme_id":          args.ThemeID,
		"code":              args.Code,
		"event_type":        args.EventType,
		"event_url":         args.EventURL,
		"event_params":      args.EventParams,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateEvent 修改Event参数
type ArgsUpdateEvent struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
	//所属主题分类
	ThemeCategoryID int64 `db:"theme_category_id" json:"themeCategoryId" check:"id"`
	//所属主题
	// 插槽可用于的主题域
	ThemeID int64 `db:"theme_id" json:"themeId" check:"id"`
	//事件编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"300"`
	//事件类型
	// nats - NATS事件
	EventType string `db:"event_type" json:"eventType" check:"intThan0"`
	//事件地址
	// nats - 触发的地址
	EventURL string `db:"event_url" json:"eventURL" check:"des" min:"1" max:"600"`
	//事件固定参数
	// nats - 事件附带的固定参数，如果为空则根据流程阶段事件触发填入
	EventParams string `db:"event_params" json:"eventParams" check:"des" min:"1" max:"1000" empty:"true"`
}

// UpdateEvent 修改Event
func UpdateEvent(args *ArgsUpdateEvent) (err error) {
	//更新数据
	err = eventDB.Update().SetFields([]string{"name", "description", "theme_category_id", "theme_id", "code", "event_type", "event_url", "event_params"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"name":              args.Name,
		"description":       args.Description,
		"theme_category_id": args.ThemeCategoryID,
		"theme_id":          args.ThemeID,
		"code":              args.Code,
		"event_type":        args.EventType,
		"event_url":         args.EventURL,
		"event_params":      args.EventParams,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteEventCache(args.ID)
	//反馈
	return
}

// ArgsDeleteEvent 删除Event参数
type ArgsDeleteEvent struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteEvent 删除Event
func DeleteEvent(args *ArgsDeleteEvent) (err error) {
	//删除数据
	err = eventDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteEventCache(args.ID)
	//反馈
	return
}

// GetEventCountByCategoryID 获取分类下的事件数量
func GetEventCountByCategoryID(categoryID int64) (count int64) {
	count, _ = eventDB.Select().SetFieldsList([]string{"id"}).SetDeleteQuery("delete_at", false).SetIDQuery("theme_category_id", categoryID).SetPages(CoreSQL2.ArgsPages{
		Page: 1,
		Max:  1,
		Sort: "id",
		Desc: false,
	}).SelectList("").ResultCount()
	return
}

// GetEventCountByThemeID 获取主题下的事件数量
func GetEventCountByThemeID(themeID int64) (count int64) {
	count, _ = eventDB.Select().SetFieldsList([]string{"id"}).SetDeleteQuery("delete_at", false).SetIDQuery("theme_id", themeID).SetPages(CoreSQL2.ArgsPages{
		Page: 1,
		Max:  1,
		Sort: "id",
		Desc: false,
	}).SelectList("").ResultCount()
	return
}

// getEventByID 通过ID获取Event数据包
func getEventByID(id int64) (data FieldsEvent) {
	cacheMark := getEventCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := eventDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "name", "description", "theme_category_id", "theme_id", "code", "event_type", "event_url", "event_params"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheEventTime)
	return
}

// 缓冲
func getEventCacheMark(id int64) string {
	return fmt.Sprint("base:bpm:event:id.", id)
}

func deleteEventCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getEventCacheMark(id))
}
