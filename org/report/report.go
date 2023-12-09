package OrgReport

import (
	"errors"
	"fmt"
	BaseFileSys2 "gitee.com/weeekj/weeekj_core/v5/base/filesys2"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsGetReportList 获取反馈列表参数
type ArgsGetReportList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//是否需要已经回复参数
	NeedIsReport bool `json:"needIsReport" check:"bool"`
	IsReport     bool `json:"isReport" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

type DataGetReportList struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt string `db:"create_at" json:"createAt"`
	//投诉来源
	// 可以全部为0或空，则代表匿名
	FromSystem string `db:"from_system" json:"fromSystem"`
	FromID     int64  `db:"from_id" json:"fromID"`
	FromName   string `db:"from_name" json:"fromName"`
	//建议内容
	Des string `db:"des" json:"des"`
	//投诉目标
	// 可以不含投诉目标，具体看投诉人意愿和业务逻辑需要
	TargetSystem string `db:"target_system" json:"targetSystem"`
	TargetID     int64  `db:"target_id" json:"targetID"`
	TargetName   string `db:"target_name" json:"targetName"`
	//投诉图片
	DesFiles    pq.Int64Array `db:"des_files" json:"desFiles"`
	DesFileURLs []string      `json:"desFileURLs"`
	//反馈内容
	ReportAt  string `db:"report_at" json:"reportAt"`
	ReportDes string `db:"report_des" json:"reportDes"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// GetReportList 获取反馈列表
func GetReportList(args *ArgsGetReportList) (dataList []DataGetReportList, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.NeedIsReport {
		where = CoreSQL.GetDeleteSQLField(args.IsReport, where, "report_at")
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(des ILIKE '%' || :search || '%' OR report_des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "org_report"
	var rawList []FieldsReport
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "report_at"},
	)
	for _, v := range rawList {
		vData := getReportByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, DataGetReportList{
			ID:           vData.ID,
			CreateAt:     CoreFilter.GetTimeToDefaultTime(vData.CreateAt),
			FromSystem:   vData.FromSystem,
			FromID:       vData.FromID,
			FromName:     vData.FromName,
			Des:          CoreFilter.SubStrQuick(vData.Des, 10),
			TargetSystem: vData.TargetSystem,
			TargetID:     vData.TargetID,
			TargetName:   vData.TargetName,
			DesFiles:     vData.DesFiles,
			DesFileURLs:  BaseFileSys2.GetPublicURLsByClaimIDs(vData.DesFiles),
			ReportAt:     CoreFilter.GetTimeToDefaultTime(vData.ReportAt),
			ReportDes:    CoreFilter.SubStrQuick(vData.ReportDes, 10),
			Params:       vData.Params,
		})
	}
	return
}

// ArgsGetReport 获取反馈数
type ArgsGetReport struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}
type DataGetReport struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt string `db:"create_at" json:"createAt"`
	//投诉来源
	// 可以全部为0或空，则代表匿名
	FromSystem string `db:"from_system" json:"fromSystem"`
	FromID     int64  `db:"from_id" json:"fromID"`
	FromName   string `db:"from_name" json:"fromName"`
	//建议内容
	Des string `db:"des" json:"des"`
	//建议附图
	DesFiles []string `db:"des_files" json:"desFiles"`
	//投诉目标
	// 可以不含投诉目标，具体看投诉人意愿和业务逻辑需要
	TargetSystem string `db:"target_system" json:"targetSystem"`
	TargetID     int64  `db:"target_id" json:"targetID"`
	TargetName   string `db:"target_name" json:"targetName"`
	//反馈内容
	ReportAt    string   `db:"report_at" json:"reportAt"`
	ReportDes   string   `db:"report_des" json:"reportDes"`
	ReportFiles []string `db:"report_files" json:"reportFiles"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// GetReport 获取反馈数据包
func GetReport(args *ArgsGetReport) (data DataGetReport, err error) {
	rawData := getReportByID(args.ID)
	if rawData.ID < 1 || !CoreFilter.EqID2(args.OrgID, rawData.OrgID) {
		err = errors.New("no data")
		return
	}
	data = DataGetReport{
		ID:           rawData.ID,
		CreateAt:     CoreFilter.GetTimeToDefaultTime(rawData.CreateAt),
		FromSystem:   rawData.FromSystem,
		FromID:       rawData.FromID,
		FromName:     rawData.FromName,
		Des:          rawData.Des,
		DesFiles:     BaseFileSys2.GetPublicURLsByClaimIDs(rawData.DesFiles),
		TargetSystem: rawData.TargetSystem,
		TargetID:     rawData.TargetID,
		TargetName:   rawData.TargetName,
		ReportAt:     CoreFilter.GetTimeToDefaultTime(rawData.ReportAt),
		ReportDes:    rawData.ReportDes,
		ReportFiles:  BaseFileSys2.GetPublicURLsByClaimIDs(rawData.ReportFiles),
		Params:       rawData.Params,
	}
	return
}

// ArgsAddReport 添加新的反馈参数
type ArgsAddReport struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//投诉来源
	// 可以全部为0或空，则代表匿名
	FromSystem string `db:"from_system" json:"fromSystem" check:"mark"`
	FromID     int64  `db:"from_id" json:"fromID" check:"id" empty:"true"`
	FromName   string `db:"from_name" json:"fromName" check:"name"`
	//投诉目标
	// 可以不含投诉目标，具体看投诉人意愿和业务逻辑需要
	TargetSystem string `db:"target_system" json:"targetSystem"`
	TargetID     int64  `db:"target_id" json:"targetID"`
	TargetName   string `db:"target_name" json:"targetName"`
	//建议内容
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000"`
	//建议附图
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// AddReport 添加新的反馈
func AddReport(args *ArgsAddReport) (err error) {
	if len(args.DesFiles) > 10 {
		err = errors.New("too many files")
		return
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_report (org_id, from_system, from_id, from_name, des, des_files, target_system, target_id, target_name, params) VALUES (:org_id,:from_system,:from_id,:from_name,:des,:des_files,:target_system,:target_id,:target_name,:params)", args)
	if err != nil {
		return
	}
	return
}

// ArgsReReport 处理反馈参数
type ArgsReReport struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//反馈内容
	ReportAt    string        `db:"report_at" json:"reportAt"`
	ReportDes   string        `db:"report_des" json:"reportDes"`
	ReportFiles pq.Int64Array `db:"report_files" json:"reportFiles"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// ReReport 处理反馈
func ReReport(args *ArgsReReport) (errCode string, err error) {
	//获取时间
	var reportAt time.Time
	reportAt, err = CoreFilter.GetTimeByDefault(args.ReportAt)
	if err != nil {
		errCode = "err_time"
		return
	}
	//获取数据
	rawData := getReportByID(args.ID)
	if rawData.ID < 1 || !CoreFilter.EqID2(args.OrgID, rawData.OrgID) {
		errCode = "err_no_data"
		err = errors.New("no data")
		return
	}
	//检查是否已经回复
	if CoreSQL.CheckTimeHaveData(rawData.ReportAt) {
		errCode = "err_have_report"
		err = errors.New("have report")
		return
	}
	//修改反馈
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_report SET report_at = :report_at, report_des = :report_des, report_files = :report_files, params = :params WHERE id = :id", map[string]interface{}{
		"id":           args.ID,
		"report_at":    reportAt,
		"report_des":   args.ReportDes,
		"report_files": args.ReportFiles,
		"params":       args.Params,
	})
	if err != nil {
		errCode = "err_insert"
		return
	}
	//删除缓冲
	deleteReportCache(args.ID)
	//反馈
	return
}

// 获取ID
func getReportByID(id int64) (data FieldsReport) {
	cacheMark := getReportCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, from_system, from_id, from_name, des, des_files, target_system, target_id, target_name, report_at, report_des, report_files, params FROM org_report WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 10800)
	return
}

type ArgsDeleteReport struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// 删除反馈
func DeleteReport(args *ArgsDeleteReport) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_report", "id = :id AND org_id = :org_id", args)
	if err != nil {
		return
	}
	deleteReportCache(args.ID)
	return
}

// 缓冲
func getReportCacheMark(id int64) string {
	return fmt.Sprint("org:report:id:", id)
}

func deleteReportCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getReportCacheMark(id))
}
