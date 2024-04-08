package BaseService

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetServiceList 获取服务列表参数
type ArgsGetServiceList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//事件编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"300" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetServiceList 获取服务列表
func GetServiceList(args *ArgsGetServiceList) (dataList []FieldsService, dataCount int64, err error) {
	dataCount, err = serviceDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "expire_at"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetStringQuery("code", args.Code).SetSearchQuery([]string{"name", "description"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getServiceByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetServiceByID 获取Service数据包参数
type ArgsGetServiceByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetServiceByID 获取Service数
func GetServiceByID(args *ArgsGetServiceByID) (data FieldsService, err error) {
	data = getServiceByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// getServiceByCode 通过编码查询服务
func getServiceByCode(code string) (data FieldsService) {
	err := serviceDB.Get().SetFieldsOne([]string{"id"}).SetStringQuery("code", code).SetDeleteQuery("delete_at", false).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	data = getServiceByID(data.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return

}

// argsSetService 设置Service参数
type argsSetService struct {
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
	//事件订阅方式
	// server 服务器订阅; client 客户端订阅; all 服务器和客户端都订阅
	EventSubType string `db:"event_sub_type" json:"eventSubType" check:"intThan0"`
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

// setService 设置Service
func setService(args *argsSetService) (err error) {
	//检查订阅方式
	switch args.EventSubType {
	case "server":
	case "client":
	case "all":
	default:
		err = errors.New("event sub type error")
		return
	}
	//检查事件类型
	switch args.EventType {
	case "nats":
	default:
		err = errors.New("event type error")
		return
	}
	//检查过期时间
	if args.ExpireAt.Unix() == 0 {
		args.ExpireAt = CoreFilter.GetNowTimeCarbon().AddDay().Time
	}
	//尝试获取code服务
	data := getServiceByCode(args.Code)
	if data.ID > 0 {
		err = updateService(&argsUpdateService{
			ID:           data.ID,
			ExpireAt:     args.ExpireAt,
			Name:         args.Name,
			Description:  args.Description,
			EventSubType: args.EventSubType,
			Code:         args.Code,
			EventType:    args.EventType,
			EventURL:     args.EventURL,
			EventParams:  args.EventParams,
		})
		if err != nil {
			return
		}
	} else {
		_, err = createService(&argsCreateService{
			ExpireAt:     args.ExpireAt,
			Name:         args.Name,
			Description:  args.Description,
			EventSubType: args.EventSubType,
			Code:         args.Code,
			EventType:    args.EventType,
			EventURL:     args.EventURL,
			EventParams:  args.EventParams,
		})
		if err != nil {
			return
		}
	}
	return
}

// argsCreateService 创建Service参数
type argsCreateService struct {
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
	//事件订阅方式
	// server 服务器订阅; client 客户端订阅; all 服务器和客户端都订阅
	EventSubType string `db:"event_sub_type" json:"eventSubType" check:"intThan0"`
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

// createService 创建Service
func createService(args *argsCreateService) (id int64, err error) {
	//创建数据
	id, err = serviceDB.Insert().SetFields([]string{"expire_at", "name", "description", "event_sub_type", "code", "event_type", "event_url", "event_params"}).Add(map[string]any{
		"expire_at":      args.ExpireAt,
		"name":           args.Name,
		"description":    args.Description,
		"event_sub_type": args.EventSubType,
		"code":           args.Code,
		"event_type":     args.EventType,
		"event_url":      args.EventURL,
		"event_params":   args.EventParams,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// argsUpdateService 修改Service参数
type argsUpdateService struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
	//事件订阅方式
	// server 服务器订阅; client 客户端订阅; all 服务器和客户端都订阅
	EventSubType string `db:"event_sub_type" json:"eventSubType" check:"intThan0"`
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

// updateService 修改Service
func updateService(args *argsUpdateService) (err error) {
	//更新数据
	err = serviceDB.Update().SetFields([]string{"expire_at", "name", "description", "event_sub_type", "code", "event_type", "event_url", "event_params"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"expire_at":      args.ExpireAt,
		"name":           args.Name,
		"description":    args.Description,
		"event_sub_type": args.EventSubType,
		"code":           args.Code,
		"event_type":     args.EventType,
		"event_url":      args.EventURL,
		"event_params":   args.EventParams,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteServiceCache(args.ID)
	//反馈
	return
}

// argsDeleteService 删除Service参数
type argsDeleteService struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// deleteService 删除Service
func deleteService(args *argsDeleteService) (err error) {
	//删除数据
	err = serviceDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteServiceCache(args.ID)
	//反馈
	return
}

// getServiceByID 通过ID获取Service数据包
func getServiceByID(id int64) (data FieldsService) {
	cacheMark := getServiceCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := serviceDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "expire_at", "name", "description", "event_sub_type", "code", "event_type", "event_url", "event_params"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheServiceTime)
	return
}

// 缓冲
func getServiceCacheMark(id int64) string {
	return fmt.Sprint("base:service:id.", id)
}

func deleteServiceCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getServiceCacheMark(id))
}
