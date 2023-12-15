package ERPAudit

import (
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetStepChildList 获取节点列表参数
type ArgsGetStepChildList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//所属流程
	StepID int64 `db:"step_id" json:"stepID" check:"id"`
}

// GetStepChildList 获取节点列表
func GetStepChildList(args *ArgsGetStepChildList) (dataList []FieldsStepChild, dataCount int64, err error) {
	where := "step_id = :step_id"
	maps := map[string]interface{}{
		"step_id": args.StepID,
	}
	tableName := "erp_audit_step_child"
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
		[]string{"id", "create_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getStepChildByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetStepChild 获取指定节点
func GetStepChild(id int64, orgID int64, orgBindID int64) (data FieldsStepChild) {
	data = getStepChildByID(id)
	if data.ID < 1 {
		data = FieldsStepChild{}
		return
	}
	stepData := GetStep(data.StepID, orgID, orgBindID)
	if stepData.ID < 1 {
		data = FieldsStepChild{}
		return
	}
	return
}

// 通过bind和key获取节点
func getStepChildByKey(stepID int64, key string) (data FieldsStepChild) {
	cacheMark := getStepChildKeyCacheMark(stepID, key)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, step_id, audit_at, ban_at, expire_at, key, name, audit_mode, audit_org_bind_group, audit_org_role_ids, audit_org_bind_ids, wait_audit_org_bind_ids, finish_audit_org_binds, next_step_key, ban_next_step_key, params FROM erp_audit_step_child WHERE step_id = $1 AND key = $2", stepID, key)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Day)
	return
}

// getStepChildByID 获取节点
func getStepChildByID(id int64) (data FieldsStepChild) {
	cacheMark := getStepChildCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, step_id, audit_at, ban_at, expire_at, key, name, audit_mode, audit_org_bind_group, audit_org_role_ids, audit_org_bind_ids, wait_audit_org_bind_ids, finish_audit_org_binds, next_step_key, ban_next_step_key, params FROM erp_audit_step_child WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Day)
	return
}
