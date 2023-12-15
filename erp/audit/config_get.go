package ERPAudit

import (
	"errors"
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"strings"
)

// ArgsGetConfigList 获取配置列表参数
type ArgsGetConfigList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetConfigList 获取配置列表参
func GetConfigList(args *ArgsGetConfigList) (dataList []FieldsConfig, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		args.Search = strings.ReplaceAll(args.Search, " ", "")
		where = where + " AND (REPLACE(name, ' ', '') ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "erp_audit_config"
	var rawList []FieldsConfig
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id "+"FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getConfigByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		vData.StepList = FieldsConfigStepList{}
		dataList = append(dataList, vData)
	}
	return
}

// GetConfig 获取指定的配置
func GetConfig(id int64, orgID int64) (data FieldsConfig) {
	data = getConfigByID(id)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		data = FieldsConfig{}
		return
	}
	return
}

// GetConfigName 获取配置名称
func GetConfigName(id int64) (name string) {
	data := getConfigByID(id)
	return data.Name
}

// 检查配置是否发布
func checkConfigPublish(configID int64, orgID int64) bool {
	data := GetConfig(configID, orgID)
	if data.ID < 1 || !CoreSQL.CheckTimeHaveData(data.PublishAt) {
		return false
	}
	return true
}

// getNewHash 获取新的hash
func getNewHash() string {
	return CoreFilter.GetSha1Str(CoreFilter.GetRandStr4(30))
}

// getConfigByID 获取配置
func getConfigByID(id int64) (data FieldsConfig) {
	cacheMark := getConfigCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, publish_at, hash, org_id, name, step_list, params FROM erp_audit_config WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime3Day)
	return
}

// 检查节点，创建或编辑用
func checkStepList(stepList FieldsConfigStepList) (errCode string, err error) {
	//节点过少
	if len(stepList) < 1 {
		errCode = "err_erp_audit_config_step_less_1"
		err = errors.New("step less 1")
		return
	}
	//节点过多
	if len(stepList) < 1 || len(stepList) > limitConfigStepCount {
		errCode = "err_erp_audit_config_step_too_many"
		err = errors.New("step too many")
		return
	}
	//检查是否存在闭环和是否有重复key
	var checkKeys []string
	var checkComponentKeys []string
	var waitCheckKeys []string
	for _, v := range stepList {
		for _, v2 := range checkKeys {
			if v2 == v.Key {
				errCode = "err_erp_audit_config_step_key_replace"
				err = errors.New("key is replace")
				return
			}
		}
		checkKeys = append(checkKeys, v.Key)
		for _, v2 := range v.ComponentList {
			for _, v3 := range checkComponentKeys {
				if v2.Key == v3 {
					errCode = "err_erp_audit_config_step_component_key_replace"
					err = errors.New("key is replace")
					return
				}
			}
			checkComponentKeys = append(checkComponentKeys, v2.Key)
		}
		for _, v2 := range waitCheckKeys {
			if v.NextStepKey == v2 || v.BanNextStepKey == v2 {
				errCode = "err_erp_audit_config_step_cycle"
				err = errors.New("step have cycle")
				return
			}
		}
		if v.NextStepKey != "" {
			waitCheckKeys = append(waitCheckKeys, v.NextStepKey)
		}
		if v.BanNextStepKey != "" {
			waitCheckKeys = append(waitCheckKeys, v.BanNextStepKey)
		}
	}
	//反馈
	return
}
