package UserReport

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetReportList 获取反馈列表参数
type ArgsGetReportList struct {
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

// GetReportList 获取反馈列表
func GetReportList(args *ArgsGetReportList) (dataList []FieldsReport, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		where = where + " AND (userName ILIKE '%' || :search || '%' OR ip ILIKE '%' || :search || '%' OR phone ILIKE '%' || :search || '%' OR content ILIKE '%' || :search || '%' OR report_content ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "user_report"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, delete_at, report_at, org_id, user_id, from_info, ip, user_name, nation_code, phone, email, files, report_user_id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "report_at"},
	)
	return
}

// ArgsGetReport 获取指定反馈参数
type ArgsGetReport struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetReport 获取指定反馈
func GetReport(args *ArgsGetReport) (data FieldsReport, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, delete_at, report_at, org_id, user_id, from_info, ip, user_name, nation_code, phone, email, files, content, report_user_id, report_content FROM user_report WHERE id = $1 AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	return
}

// ArgsCreateReport 创建新的反馈参数
type ArgsCreateReport struct {
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//举报内容来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//IP
	IP string `db:"ip" json:"ip"`
	//用户昵称
	UserName string `db:"user_name" json:"userName" check:"name" empty:"true"`
	//绑定手机号的国家代码
	NationCode string `db:"nation_code" json:"nationCode" check:"nationCode" empty:"true"`
	//手机号码，绑定后的手机
	Phone string `db:"phone" json:"phone" check:"phone" empty:"true"`
	//邮箱，如果不存在手机则必须存在的
	Email string `db:"email" json:"email" check:"email" empty:"true"`
	//截图
	Files pq.Int64Array `db:"files" json:"files" check:"ids" empty:"true"`
	//举报内容描述
	Content string `db:"content" json:"content" check:"des" min:"1" max:"3000"`
}

// CreateReport 创建新的反馈
func CreateReport(args *ArgsCreateReport) (err error) {
	//检查用户或IP是否重复提交反馈，已有反馈必须处理完成后才能继续提交
	if args.UserID > 0 {
		var id int64
		err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM user_report WHERE user_id = $1 AND report_at < to_timestamp(1000000) AND delete_at < to_timestamp(1000000)", args.UserID)
		if err == nil && id > 0 {
			err = errors.New("user have report")
			return
		}
	} else {
		var id int64
		err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM user_report WHERE ip = $1 AND report_at < to_timestamp(1000000) AND delete_at < to_timestamp(1000000)", args.IP)
		if err == nil && id > 0 {
			err = errors.New("ip have report")
			return
		}
	}
	//创建反馈
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_report (org_id, user_id, from_info, ip, user_name, nation_code, phone, email, files, content, report_user_id, report_content) VALUES (:org_id,:user_id,:from_info,:ip,:user_name,:nation_code,:phone,:email,:files,:content,0,'')", args)
	return
}

// ArgsReportData 反馈问题参数
type ArgsReportData struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//反馈人
	ReportUserID int64 `db:"report_user_id" json:"reportUserID" check:"id"`
	//反馈内容
	ReportContent string `db:"report_content" json:"reportContent" check:"des" min:"1" max:"3000"`
}

// ReportData 反馈问题
func ReportData(args *ArgsReportData) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_report SET report_at = NOW(), report_user_id = :report_user_id, report_content = :report_content WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteReport 删除反馈参数
type ArgsDeleteReport struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteReport 删除反馈
func DeleteReport(args *ArgsDeleteReport) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "user_report", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
