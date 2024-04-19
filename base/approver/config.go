package BaseApprover

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetConfigList 获取配置列表参数
type ArgsGetConfigList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//关联的模块标识码
	// erp_project
	ModuleCode string `db:"module_code" json:"moduleCode" check:"des" min:"1" max:"50" empty:"true"`
	//审批分叉标识码
	// 用于识别模块内，不同的审批流程
	ForkCode string `db:"fork_code" json:"forkCode" check:"des" min:"1" max:"50"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetConfigList 获取配置列表
func GetConfigList(args *ArgsGetConfigList) (dataList []DataConfig, dataCount int64, err error) {
	dataCount, err = configItemDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "flow_order"}).SetPages(args.Pages).SetDeleteQuery("delete_at", false).SetIDQuery("org_id", args.OrgID).SetStringQuery("module_code", args.ModuleCode).SetStringQuery("fork_code", args.ForkCode).SetSearchQuery([]string{"name", "description"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData, _ := GetConfigByID(&ArgsGetConfigByID{
			ID:    v.ID,
			OrgID: -1,
		})
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetConfigByID 获取配置参数
type ArgsGetConfigByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetConfigByID 获取配置
func GetConfigByID(args *ArgsGetConfigByID) (data DataConfig, err error) {
	rawData := getConfigByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	if !CoreFilter.EqID2(args.OrgID, rawData.OrgID) {
		err = errors.New("no data")
		return
	}
	var rawItems []FieldsConfigItem
	rawItems, _, err = getConfigItems(&argsGetConfigItems{
		ConfigID:  rawData.ID,
		OrgBindID: -1,
		UserID:    -1,
	})
	var items DataConfigItems
	if err == nil {
		for _, v := range rawItems {
			items = append(items, DataConfigItem{
				FlowOrder: v.FlowOrder,
				OrgBindID: v.OrgBindID,
				UserID:    v.UserID,
			})
		}
	}
	data = DataConfig{
		ID:         rawData.ID,
		CreateAt:   rawData.CreateAt,
		UpdateAt:   rawData.UpdateAt,
		DeleteAt:   rawData.DeleteAt,
		OrgID:      rawData.OrgID,
		ModuleCode: rawData.ModuleCode,
		ForkCode:   rawData.ForkCode,
		Items:      items,
	}
	return
}

// ArgsCreateConfig 创建配置参数
type ArgsCreateConfig struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//关联的模块标识码
	// erp_project
	ModuleCode string `db:"module_code" json:"moduleCode" check:"des" min:"1" max:"50"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Desc string `db:"desc" json:"desc" check:"des" min:"1" max:"300" empty:"true"`
	//审批分叉标识码
	// 用于识别模块内，不同的审批流程
	ForkCode string `db:"fork_code" json:"forkCode" check:"des" min:"1" max:"50"`
	//审批流配置
	Items DataConfigItems `json:"items"`
}

// CreateConfig 创建配置
func CreateConfig(args *ArgsCreateConfig) (configID int64, errCode string, err error) {
	//检查是否已经存在相同的配置？
	checkConfigData := getConfigByModule(args.OrgID, args.ModuleCode, args.ForkCode)
	if checkConfigData.ID > 0 {
		errCode = "err_have_replace"
		return
	}
	//审批流过长
	if len(args.Items) > 999 {
		errCode = "err_approver_flow_order_max"
		return
	}
	//审批流节点必须顺序连续
	checkItemFlow := 0
	for {
		isFind := false
		for _, v := range args.Items {
			if checkItemFlow == v.FlowOrder {
				isFind = true
				break
			}
		}
		if !isFind {
			errCode = "err_approver_flow_order"
			return
		}
		if checkItemFlow >= len(args.Items)-1 {
			break
		}
		checkItemFlow += 1
	}
	//遍历配置行，自动补全信息
	for k, v := range args.Items {
		if v.OrgBindID > 0 && v.UserID < 1 {
			vOrgBindData, _ := OrgCore.GetBind(&OrgCore.ArgsGetBind{
				ID:     v.OrgBindID,
				OrgID:  -1,
				UserID: -1,
			})
			if vOrgBindData.ID > 0 && vOrgBindData.UserID > 0 {
				args.Items[k].UserID = vOrgBindData.UserID
			}
		}
	}
	//创建数据
	configID, err = configDB.Insert().SetFields([]string{"org_id", "module_code", "name", "desc", "fork_code"}).Add(map[string]any{
		"org_id":      args.OrgID,
		"module_code": args.ModuleCode,
		"name":        args.Name,
		"desc":        args.Desc,
		"fork_code":   args.ForkCode,
	}).ExecAndResultID()
	if err != nil {
		errCode = "err_insert"
		return
	}
	//更新配置行
	errCode, err = setConfigItem(&argsSetConfigItem{
		ConfigID: configID,
		Items:    args.Items,
	})
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateConfig 修改配置参数
type ArgsUpdateConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Desc string `db:"desc" json:"desc" check:"des" min:"1" max:"300" empty:"true"`
	//审批流配置
	Items DataConfigItems `json:"items"`
}

// UpdateConfig 修改配置
func UpdateConfig(args *ArgsUpdateConfig) (errCode string, err error) {
	//更新数据
	err = configDB.Update().SetFields([]string{"name", "desc"}).NeedUpdateTime().AddWhereID(args.ID).AddWhereOrgID(args.OrgID).NamedExec(map[string]any{
		"name": args.Name,
		"desc": args.Desc,
	})
	if err != nil {
		errCode = "err_update"
		return
	}
	//删除缓冲
	deleteConfigCache(args.ID)
	//更新配置行
	errCode, err = setConfigItem(&argsSetConfigItem{
		ConfigID: args.ID,
		Items:    args.Items,
	})
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsDeleteConfig 删除Budget参数
type ArgsDeleteConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteConfig 删除Budget
func DeleteConfig(args *ArgsDeleteConfig) (err error) {
	//清理配置行
	err = clearConfigItem(args.ID)
	if err != nil {
		return
	}
	//删除数据
	err = configDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteConfigCache(args.ID)
	//反馈
	return
}

// getConfigByModule 通过模块获取配置
func getConfigByModule(orgID int64, moduleCode string, forkCode string) (data DataConfig) {
	configList, _, _ := GetConfigList(&ArgsGetConfigList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  1,
			Sort: "id",
			Desc: false,
		},
		OrgID:      orgID,
		ModuleCode: moduleCode,
		ForkCode:   forkCode,
		IsRemove:   false,
		Search:     "",
	})
	if len(configList) < 1 {
		return
	}
	var err error
	data, err = GetConfigByID(&ArgsGetConfigByID{
		ID:    configList[0].ID,
		OrgID: -1,
	})
	if err != nil {
		return
	}
	return
}

// getConfigByID 通过ID查询配置
func getConfigByID(id int64) (data FieldsConfig) {
	cacheMark := getConfigCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := configDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "org_id", "module_code", "name", "desc", "fork_code"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheConfigTime)
	return
}

// 缓冲
func getConfigCacheMark(id int64) string {
	return fmt.Sprint("base:approver:config:id.", id)
}

func deleteConfigCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getConfigCacheMark(id))
}
