package OrgMission

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsGetMissionList 获取任务列表参数
type ArgsGetMissionList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//任意一种形式包含此人
	OperateBindID int64 `db:"operate_bind_id" json:"operateBindID" check:"id" empty:"true"`
	//状态
	// 0 未完成; 1 已完成; 2 放弃; 3 删除或取消
	Status pq.Int32Array `db:"status" json:"status"`
	//创建人
	CreateBindID int64 `json:"createBindID" check:"id" empty:"true"`
	//执行人
	BindID int64 `json:"bindID" check:"id" empty:"true"`
	//其他执行人
	OtherBindID int64 `json:"otherBindID" check:"id" empty:"true"`
	//上级任务
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//级别
	Level int `db:"level" json:"level" check:"intThan0" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签列
	Tags int64 `db:"tags" json:"tags" check:"id" empty:"true"`
	//开始时间范围
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetMissionList 获取任务列表
func GetMissionList(args *ArgsGetMissionList) (dataList []FieldsMission, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.CreateBindID > -1 {
		where = where + " AND create_bind_id = :create_bind_id"
		maps["create_bind_id"] = args.CreateBindID
	}
	if args.BindID > -1 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.OtherBindID > -1 {
		where = where + " AND :other_bind_id = ANY(other_bind_ids)"
		maps["other_bind_id"] = args.OtherBindID
	}
	if args.OperateBindID > -1 {
		where = where + " AND (create_bind_id = :operate_bind_id OR bind_id = :operate_bind_id OR :operate_bind_id = ANY(other_bind_ids))"
		maps["operate_bind_id"] = args.OperateBindID
	}
	if len(args.Status) > 0 {
		where = where + " AND status = ANY(status)"
		maps["status"] = args.Status
	}
	if args.ParentID > -1 {
		where = where + " AND parent_id = :parent_id"
		maps["parent_id"] = args.ParentID
	}
	if args.Level > -1 {
		where = where + " AND level = :level"
		maps["level"] = args.Level
	}
	if args.SortID > -1 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if args.Tags > -1 {
		where = where + " AND :tags = ANY(tags)"
		maps["tags"] = args.Tags
	}
	if args.Search != "" {
		where = where + " AND (title ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "org_mission"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, auto_id, status, org_id, create_bind_id, bind_id, other_bind_ids, title, des, des_files, start_at, end_at, tip_id, level, sort_id, tags, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "start_at", "end_at", "level"},
	)
	return
}

// ArgsGetMissionID 查看任务详情参数
type ArgsGetMissionID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//修改操作人
	// 用于验证
	OperateBindID int64 `db:"operate_bind_id" json:"operateBindID" check:"id"`
}

// GetMissionID 查看任务详情
func GetMissionID(args *ArgsGetMissionID) (data FieldsMission, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, auto_id, status, org_id, create_bind_id, bind_id, other_bind_ids, title, des, des_files, start_at, end_at, tip_id, level, sort_id, tags, params FROM org_mission WHERE id = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000) AND ($3 < 1 OR create_bind_id = $3 OR bind_id = $3 OR $3 = ANY(other_bind_ids))", args.ID, args.OrgID, args.OperateBindID)
	return
}

// ArgsCheckMissionOperate 检查任务是否属于该操作人参数
type ArgsCheckMissionOperate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//修改操作人
	// 用于验证
	OperateBindID int64 `db:"operate_bind_id" json:"operateBindID" check:"id"`
}

// CheckMissionOperate 检查任务是否属于该操作人
func CheckMissionOperate(args *ArgsCheckMissionOperate) (b bool) {
	var data FieldsMission
	if err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_mission WHERE id = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000) AND ($3 < 1 OR create_bind_id = $3 OR bind_id = $3 OR $3 = ANY(other_bind_ids))", args.ID, args.OrgID, args.OperateBindID); err != nil || data.ID < 1 {
		return
	}
	return true
}

// ArgsCreateMission 创建新的任务参数
type ArgsCreateMission struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//创建人
	CreateBindID int64 `db:"create_bind_id" json:"createBindID" check:"id"`
	//执行人
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//其他执行人
	OtherBindIDs pq.Int64Array `db:"other_bind_ids" json:"otherBindIDs" check:"ids" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//执行时间
	StartAt time.Time `db:"start_at" json:"startAt" check:"isoTime" empty:"true"`
	//结束时间
	EndAt time.Time `db:"end_at" json:"endAt" check:"isoTime" empty:"true"`
	//是否需提醒
	// -1 不需要; 0 需要等待提醒中; >0 已经触发提醒的ID
	TipID int64 `db:"tip_id" json:"tipID" check:"id" empty:"true"`
	//上级任务
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//级别
	Level int `db:"level" json:"level" check:"intThan0" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateMission 新的任务参数
func CreateMission(args *ArgsCreateMission) (data FieldsMission, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_mission", "INSERT INTO org_mission (auto_id, status, org_id, create_bind_id, bind_id, other_bind_ids, title, des, des_files, start_at, end_at, tip_id, parent_id, level, sort_id, tags, params) VALUES (0,0,:org_id,:create_bind_id,:bind_id,:other_bind_ids,:title,:des,:des_files,:start_at,:end_at,:tip_id,:parent_id,:level,:sort_id,:tags,:params)", args, &data)
	if err == nil {
		_ = CreateLog(&ArgsCreateLog{
			OrgID:       args.OrgID,
			BindID:      args.CreateBindID,
			MissionID:   data.ID,
			ContentMark: "update",
			Content:     "创建新的任务",
		})
	}
	return
}

// ArgsUpdateMission 修改任务参数
type ArgsUpdateMission struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//修改操作人
	// 用于验证
	OperateBindID int64 `db:"operate_bind_id" json:"operateBindID" check:"id"`
	//状态
	// 0 未完成; 1 已完成; 2 放弃; 3 删除或取消
	Status int `db:"status" json:"status"`
	//执行人
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//其他执行人
	OtherBindIDs pq.Int64Array `db:"other_bind_ids" json:"otherBindIDs" check:"ids" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//执行时间
	StartAt time.Time `db:"start_at" json:"startAt" check:"isoTime" empty:"true"`
	//结束时间
	EndAt time.Time `db:"end_at" json:"endAt" check:"isoTime" empty:"true"`
	//是否需提醒
	// -1 不需要; 0 需要等待提醒中; >0 已经触发提醒的ID
	TipID int64 `db:"tip_id" json:"tipID" check:"id" empty:"true"`
	//上级任务
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//级别
	Level int `db:"level" json:"level" check:"intThan0" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateMission 修改任务
func UpdateMission(args *ArgsUpdateMission) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_mission SET update_at = NOW(), status = :status, bind_id = :bind_id, other_bind_ids = :other_bind_ids, title = :title, des = :des, des_files = :des_files, start_at = :start_at, end_at = :end_at, tip_id = :tip_id, parent_id = :parent_id, level = :level, sort_id = :sort_id, tags = :tags, params = :params WHERE id = :id AND org_id = :org_id AND (:operate_bind_id < 1 OR create_bind_id = :operate_bind_id OR bind_id = :operate_bind_id)", args)
	if err == nil {
		switch args.Status {
		case 1:
			_ = CreateLog(&ArgsCreateLog{
				OrgID:       args.OrgID,
				BindID:      args.OperateBindID,
				MissionID:   args.ID,
				ContentMark: "finish",
				Content:     "完成任务",
			})
		case 2:
			_ = CreateLog(&ArgsCreateLog{
				OrgID:       args.OrgID,
				BindID:      args.OperateBindID,
				MissionID:   args.ID,
				ContentMark: "cancel",
				Content:     "放弃任务",
			})
		case 3:
			_ = CreateLog(&ArgsCreateLog{
				OrgID:       args.OrgID,
				BindID:      args.OperateBindID,
				MissionID:   args.ID,
				ContentMark: "cancel",
				Content:     "取消任务",
			})
		default:
			_ = CreateLog(&ArgsCreateLog{
				OrgID:       args.OrgID,
				BindID:      args.OperateBindID,
				MissionID:   args.ID,
				ContentMark: "update",
				Content:     "更新任务",
			})
		}
	}
	return
}

// ArgsUpdateMissionStatus 修改任务状态参数
type ArgsUpdateMissionStatus struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//修改操作人
	// 用于验证
	OperateBindID int64 `db:"operate_bind_id" json:"operateBindID" check:"id"`
	//状态
	// 0 未完成; 1 已完成; 2 放弃; 3 删除或取消
	Status int `db:"status" json:"status"`
}

// UpdateMissionStatus 修改任务状态参数
func UpdateMissionStatus(args *ArgsUpdateMissionStatus) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_mission SET update_at = NOW(), status = :status WHERE id = :id AND org_id = :org_id AND (:operate_bind_id < 1 OR create_bind_id = :operate_bind_id OR bind_id = :operate_bind_id)", args)
	if err == nil {
		switch args.Status {
		case 1:
			_ = CreateLog(&ArgsCreateLog{
				OrgID:       args.OrgID,
				BindID:      args.OperateBindID,
				MissionID:   args.ID,
				ContentMark: "finish",
				Content:     "完成任务",
			})
		case 2:
			_ = CreateLog(&ArgsCreateLog{
				OrgID:       args.OrgID,
				BindID:      args.OperateBindID,
				MissionID:   args.ID,
				ContentMark: "cancel",
				Content:     "放弃任务",
			})
		case 3:
			_ = CreateLog(&ArgsCreateLog{
				OrgID:       args.OrgID,
				BindID:      args.OperateBindID,
				MissionID:   args.ID,
				ContentMark: "cancel",
				Content:     "取消任务",
			})
		default:
			_ = CreateLog(&ArgsCreateLog{
				OrgID:       args.OrgID,
				BindID:      args.OperateBindID,
				MissionID:   args.ID,
				ContentMark: "update",
				Content:     "更新任务",
			})
		}
	}
	return
}

// ArgsDeleteMission 删除任务参数
type ArgsDeleteMission struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//任意一种形式包含此人
	OperateBindID int64 `db:"operate_bind_id" json:"operateBindID" check:"id" empty:"true"`
}

// DeleteMission 删除任务
func DeleteMission(args *ArgsDeleteMission) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_mission", "id = :id AND org_id = :org_id AND (:operate_bind_id < 1 OR create_bind_id = :operate_bind_id OR bind_id = :operate_bind_id)", args)
	if err == nil {
		_ = CreateLog(&ArgsCreateLog{
			OrgID:       args.OrgID,
			BindID:      args.OperateBindID,
			MissionID:   args.ID,
			ContentMark: "delete",
			Content:     "删除任务",
		})
	}
	return
}
