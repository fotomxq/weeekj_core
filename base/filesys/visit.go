package BaseFileSys

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetVisitList 查看访问记录参数
type ArgsGetVisitList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//指定用户ID
	UserID int64 `json:"userID"`
	//指定文件引用ID
	ClaimID int64 `json:"claimID"`
	//指定文件ID
	FileID int64 `json:"fileID"`
	//IP
	IP string `json:"ip"`
	//最小时间和最大时间
	MinTime time.Time
	MaxTime time.Time
}

// GetVisitList 查看访问记录
func GetVisitList(args *ArgsGetVisitList) (dataList []FieldsFileClaimVisit, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.UserID > 0 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.ClaimID > 0 {
		where = where + " AND claim_id = :claim_id"
		maps["claim_id"] = args.ClaimID
	}
	if args.FileID > 0 {
		where = where + " AND file_id = :file_id"
		maps["file_id"] = args.FileID
	}
	if args.IP != "" {
		where = where + " AND ip = :ip"
		maps["ip"] = args.IP
	}
	if args.MinTime.Unix() > 0 {
		where = where + " AND create_at > :min_time"
		maps["min_time"] = args.MinTime
	}
	if args.MaxTime.Unix() > 0 {
		where = where + " AND create_at > :max_time"
		maps["max_time"] = args.MaxTime
	}
	if where != "" {
		where = " WHERE " + where
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_file_claim_visit",
		"id",
		"SELECT id, create_at, claim_id, file_id, user_id, create_ip FROM core_file_claim_visit"+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "claim_id", "file_id", "user_id", "create_ip"},
	)
	return
}
