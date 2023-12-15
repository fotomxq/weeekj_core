package ERPAudit

import (
	"errors"
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetStepList 获取审批列表参数
type ArgsGetStepList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//流程配置
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//是否审批完成
	NeedIsFinish bool `json:"needIsFinish" check:"bool"`
	IsFinish     bool `json:"isFinish" check:"bool"`
	//最终状态
	// -1 跳过; 0 无状态; 1 审批通过; 2 拒绝审批
	FinishStatus int `db:"finish_status" json:"finishStatus"`
	//可能存在访问或编辑能力的组织成员
	OrgBindID int64 `json:"orgBindID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetStepList 获取审批列表
func GetStepList(args *ArgsGetStepList) (dataList []FieldsStep, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ConfigID > -1 {
		//检查是否发布
		if !checkConfigPublish(args.ConfigID, args.OrgID) {
			err = errors.New("config no publish")
			return
		}
		//组合条件
		where = where + " AND config_id = :config_id"
		maps["config_id"] = args.ConfigID
	}
	if args.NeedIsFinish {
		where = CoreSQL.GetDeleteSQLField(args.IsFinish, where, "finish_at")
	}
	if args.FinishStatus > -1 {
		where = where + " AND finish_status = :finish_status"
		maps["finish_status"] = args.FinishStatus
	}
	if args.OrgBindID > -1 {
		where = where + " AND (:org_bind_id = ANY(can_view_org_bind_ids) OR :org_bind_id = ANY(can_edit_org_bind_ids))"
		maps["org_bind_id"] = args.OrgBindID
	}
	if args.Search != "" {
		where = where + " AND (sn ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "erp_audit_step"
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
		vData := getStepByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetStep 获取指定审批数据
func GetStep(id int64, orgID int64, orgBindID int64) (data FieldsStep) {
	data = getStepByID(id)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) || CoreSQL.CheckTimeHaveData(data.DeleteAt) || (!CoreFilter.EqHaveID2(orgBindID, data.CanViewOrgBindIDs) && !CoreFilter.EqHaveID2(orgBindID, data.CanEditOrgBindIDs)) {
		data = FieldsStep{}
		return
	}
	return
}

// getStepByID 获取审批
func getStepByID(id int64) (data FieldsStep) {
	cacheMark := getStepCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, finish_at, finish_status, org_id, config_id, sn, name, create_org_bind_id, can_view_org_bind_ids, can_edit_org_bind_ids, have_org_bind_ids, now_step_child_key, params FROM erp_audit_step WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime2Day)
	return
}
