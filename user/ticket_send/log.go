package UserTicketSend

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetLogList 查看发放日志参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//发放ID
	SendID int64 `db:"send_id" json:"sendID" check:"id"`
}

// GetLogList 查看发放日志
func GetLogList(args *ArgsGetLogList) (dataList []FieldsSendLog, dataCount int64, err error) {
	//检查发放ID合法性
	if args.OrgID > 0 {
		var id int64
		err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM user_ticket_send WHERE id = $1 AND org_id = $2", args.SendID, args.OrgID)
		if err != nil || id < 1 {
			err = errors.New("no data")
			return
		}
	}
	//组合数据
	where := ""
	maps := map[string]interface{}{}
	if args.SendID > -1 {
		where = where + "send_id = :send_id"
		maps["send_id"] = args.SendID
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"user_ticket_send_log",
		"id",
		"SELECT id, create_at, send_id, user_id FROM user_ticket_send_log WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}
