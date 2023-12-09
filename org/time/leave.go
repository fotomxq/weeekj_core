package OrgTime

import (
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetLeaveList 获取请假列表参数
type ArgsGetLeaveList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// 必填
	OrgID int64 `json:"orgID" check:"id"`
	//请假人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//是否审核
	IsAudit bool `json:"isAudit" check:"bool"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetLeaveList 获取请假列表
func GetLeaveList(args *ArgsGetLeaveList) (dataList []FieldsLeave, dataCount int64, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	where = CoreSQL.GetDeleteSQLField(args.IsAudit, where, "audit_at")
	if args.OrgBindID > -1 {
		where = where + " AND org_bind_id = :org_bind_id"
		maps["org_bind_id"] = args.OrgBindID
	}
	if args.Search != "" {
		where = where + " AND (des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	var rawList []FieldsWorkTime
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"org_work_time_leave",
		"id",
		"SELECT id FROM org_work_time_leave WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "delete_at", "audit_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getLeaveByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// CheckLeaveByBindID 检查是否正在请假
func CheckLeaveByBindID(orgBindID int64) bool {
	//获取数据
	var data FieldsLeave
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_work_time_leave WHERE delete_at < to_timestamp(1000000) AND audit_at >= NOW() AND org_bind_id = $1", orgBindID)
	if err != nil || data.ID < 1 {
		return false
	}
	data = getLeaveByID(data.ID)
	if data.ID < 1 {
		return false
	}
	//是否是否正在休假
	isLeave := false
	if !CoreSQL.CheckTimeThanNow(data.StartAt) && CoreSQL.CheckTimeThanNow(data.EndAt) {
		isLeave = true
	}
	return isLeave
}

// ArgsCreateLeave 创建请假参数
type ArgsCreateLeave struct {
	//离开时间
	StartAt string `db:"start_at" json:"startAt" check:"defaultTime"`
	//结束时间
	EndAt string `db:"end_at" json:"endAt" check:"defaultTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//请假人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//请假原因
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
}

// CreateLeave 创建请假
func CreateLeave(args *ArgsCreateLeave) (err error) {
	//获取时间
	var startAt, endAt time.Time
	startAt, err = CoreFilter.GetTimeByDefault(args.StartAt)
	if err != nil {
		return
	}
	endAt, err = CoreFilter.GetTimeByDefault(args.EndAt)
	if err != nil {
		return
	}
	//执行操作
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_work_time_leave (start_at, end_at, org_id, org_bind_id, des, ask_org_bind_id) VALUES (:start_at,:end_at,:org_id,:org_bind_id,:des,0)", map[string]interface{}{
		"start_at":    startAt,
		"end_at":      endAt,
		"org_id":      args.OrgID,
		"org_bind_id": args.OrgBindID,
		"des":         args.Des,
	})
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsAuditLeave 审核请假参数
type ArgsAuditLeave struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//审批人
	AskOrgBindID int64 `db:"ask_org_bind_id" json:"askOrgBindID" check:"id" empty:"true"`
}

// AuditLeave 审核请假参数
func AuditLeave(args *ArgsAuditLeave) (err error) {
	//执行操作
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_work_time_leave SET audit_at = NOW(), ask_org_bind_id = :ask_org_bind_id WHERE id = :id AND org_id = :org_id", args)
	if err != nil {
		return
	}
	//删除缓冲
	deleteLeaveCache(args.ID)
	//获取数据
	data := getLeaveByID(args.ID)
	if data.ID > 0 {
		pushLeaveStatus(&data)
	}
	//反馈
	return
}

// ArgsDeleteLeave 删除请假参数
type ArgsDeleteLeave struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteLeave 删除请假
func DeleteLeave(args *ArgsDeleteLeave) (err error) {
	//执行操作
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_work_time_leave", "id = :id AND org_id = :org_id", args)
	if err != nil {
		return
	}
	//删除缓冲
	deleteLeaveCache(args.ID)
	//获取数据
	data := getLeaveByID(args.ID)
	if data.ID > 0 {
		pushLeaveStatus(&data)
	}
	//反馈
	return
}

// 获取数据
func getLeaveByID(id int64) (data FieldsLeave) {
	cacheMark := getLeaveCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, delete_at, audit_at, start_at, end_at, org_id, org_bind_id, des, ask_org_bind_id FROM org_work_time_leave WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// 缓冲
func getLeaveCacheMark(id int64) string {
	return fmt.Sprint("org:time:leave:id:", id)
}

func deleteLeaveCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getLeaveCacheMark(id))
}

// 检查成员是否休假并推送通知
func pushLeaveStatus(data *FieldsLeave) {
	//通知下班
	pushNatsWork(data.OrgID, []int64{data.OrgBindID}, []int64{}, CheckIsWorkByOrgBindID(data.OrgBindID))
}
