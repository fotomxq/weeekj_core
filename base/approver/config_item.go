package BaseApprover

import (
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// argsGetConfigItems 获取配置行列表参数
type argsGetConfigItems struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//审批人用户ID
	// 用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// getConfigItems 获取配置行列表
func getConfigItems(args *argsGetConfigItems) (dataList []FieldsConfigItem, dataCount int64, err error) {
	dataCount, err = configItemDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "flow_order"}).SetPages(CoreSQL2.ArgsPages{
		Page: 1,
		Max:  999,
		Sort: "flow_order",
		Desc: false,
	}).SetDeleteQuery("delete_at", false).SetIDQuery("config_id", args.ConfigID).SetIDQuery("org_bind_id", args.OrgBindID).SetIDQuery("user_id", args.UserID).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getConfigItemByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// argsSetConfigItem 批量设置配置行参数
type argsSetConfigItem struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//审批流配置
	Items DataConfigItems `json:"items"`
}

// setConfigItem 批量设置配置行
func setConfigItem(args *argsSetConfigItem) (errCode string, err error) {
	//不能超过99个审批
	if len(args.Items) > 99 {
		errCode = "err_approver_flow_order_max"
		err = fmt.Errorf("replace flow order：%d", len(args.Items))
		return

	}
	//遍历数据，不能存在重复排序
	for _, v := range args.Items {
		haveCount := 0
		for _, v2 := range args.Items {
			if v.FlowOrder == v2.FlowOrder {
				haveCount += 1
			}
		}
		if haveCount > 1 {
			errCode = "err_approver_flow_order_repeat"
			err = fmt.Errorf("replace flow order：%d", v.FlowOrder)
			return
		}
	}
	//获取配置行列表
	dataList, _, _ := getConfigItems(&argsGetConfigItems{
		ConfigID:  args.ConfigID,
		OrgBindID: -1,
		UserID:    -1,
	})
	//找到现有数据，如果存在则更新/删除
	for _, v := range dataList {
		vIsFind := false
		for _, v2 := range args.Items {
			if v.FlowOrder == v2.FlowOrder {
				vIsFind = true
				break
			}
		}
		if vIsFind {
			err = configItemDB.Update().SetFields([]string{"flow_order", "org_bind_id", "user_id"}).NeedUpdateTime().AddWhereID(v.ID).NamedExec(map[string]any{
				"flow_order":  v.FlowOrder,
				"org_bind_id": v.OrgBindID,
				"user_id":     v.UserID,
			})
			if err != nil {
				errCode = "err_update"
				return
			}
			deleteConfigItemCache(v.ID)
		} else {
			err = configItemDB.Delete().NeedSoft(true).AddWhereID(v.ID).ExecNamed(nil)
			if err != nil {
				errCode = "err_delete"
				return
			}
			deleteConfigItemCache(v.ID)
		}
	}
	//遍历新数据，如果不存在则创建
	for _, v := range args.Items {
		vIsFind := false
		for _, v2 := range dataList {
			if v.FlowOrder == v2.FlowOrder {
				vIsFind = true
				break
			}
		}
		if !vIsFind {
			_, err = configItemDB.Insert().SetFields([]string{"config_id", "flow_order", "org_bind_id", "user_id"}).Add(map[string]any{
				"config_id":   args.ConfigID,
				"flow_order":  v.FlowOrder,
				"org_bind_id": v.OrgBindID,
				"user_id":     v.UserID,
			}).ExecAndResultID()
			if err != nil {
				return
			}
		}
	}
	//反馈
	return
}

// clearConfigItem 清理配置行
func clearConfigItem(configID int64) (err error) {
	//获取配置列表
	dataList, _, _ := getConfigItems(&argsGetConfigItems{
		ConfigID:  configID,
		OrgBindID: -1,
		UserID:    -1,
	})
	if len(dataList) < 1 {
		return
	}
	//删除数据
	err = configItemDB.Delete().NeedSoft(true).SetWhereAnd("config_id", configID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	for _, v := range dataList {
		deleteConfigItemCache(v.ID)
	}
	//反馈
	return
}

// getConfigItemByID 通过ID查询配置行
func getConfigItemByID(id int64) (data FieldsConfigItem) {
	cacheMark := getConfigItemCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := configItemDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "config_id", "flow_order", "org_bind_id", "user_id"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheConfigItemTime)
	return
}

// 缓冲
func getConfigItemCacheMark(id int64) string {
	return fmt.Sprint("base:approver:config:item:id.", id)
}

func deleteConfigItemCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getConfigItemCacheMark(id))
}
