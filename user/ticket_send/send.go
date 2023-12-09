package UserTicketSend

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	UserTicket "gitee.com/weeekj/weeekj_core/v5/user/ticket"
)

// ArgsGetSendList 获取批量给予列表参数
type ArgsGetSendList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//发放的票据配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//是否自动发放，如果不是，则需绑定广告
	NeedAuto bool `db:"need_auto" json:"needAuto" check:"bool" empty:"true"`
	IsAuto   bool `db:"is_auto" json:"isAuto" check:"bool" empty:"true"`
	//是否已经领取完成
	NeedIsFinish bool `json:"needIsFinish" check:"bool" empty:"true"`
	IsFinish     bool `json:"isFinish" check:"bool" empty:"true"`
}

// GetSendList 获取批量给予列表
func GetSendList(args *ArgsGetSendList) (dataList []FieldsSend, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ConfigID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "config_id = :config_id"
		maps["config_id"] = args.ConfigID
	}
	if args.NeedAuto {
		if where != "" {
			where = where + " AND "
		}
		if args.IsAuto {
			where = where + " need_auto = true"
		} else {
			where = where + "need_auto = false"
		}
	}
	if args.NeedIsFinish {
		if where != "" {
			where = where + " AND "
		}
		if args.IsFinish {
			where = where + " finish_at >= to_timestamp(1000000)"
		} else {
			where = where + "finish_at < to_timestamp(1000000)"
		}
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"user_ticket_send",
		"id",
		"SELECT id, create_at, finish_at, org_id, send_count, need_user_sub_config_id, need_auto, config_id, per_count FROM user_ticket_send WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "finish_at"},
	)
	return
}

// ArgsCreateSend 创建新增赠与参数
type ArgsCreateSend struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否必须是会员配置ID
	NeedUserSubConfigID int64 `db:"need_user_sub_config_id" json:"needUserSubConfigID" check:"id" empty:"true"`
	//是否自动发放，如果不是，则需绑定广告
	NeedAuto bool `db:"need_auto" json:"needAuto" check:"bool"`
	//发放的票据配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//每个用户发放几张
	PerCount int64 `db:"per_count" json:"perCount" check:"int64Than0"`
}

// CreateSend 创建新增赠与
func CreateSend(args *ArgsCreateSend) (data FieldsSend, err error) {
	//如果存在商户
	if args.OrgID > 0 {
		//检查票据是否属于组织？
		_, err = UserTicket.GetConfigByID(&UserTicket.ArgsGetConfigByID{
			ID:    args.ConfigID,
			OrgID: args.OrgID,
		})
		if err != nil {
			err = errors.New("config not org self")
			return
		}
		/**
		//商户不能创建超过10条数据
		var count int64
		err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM user_ticket_send WHERE org_id = $1 AND finish_at < to_timestamp(1000000)", args.OrgID)
		if err == nil && count > 10 {
			err = errors.New("too many")
			return
		}
		*/
	}
	//创建数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "user_ticket_send", "INSERT INTO user_ticket_send (org_id, need_user_sub_config_id, need_auto, config_id, per_count) VALUES (:org_id, :need_user_sub_config_id, :need_auto, :config_id, :per_count)", args, &data)
	return
}

// ArgsDeleteSend 删除赠与参数
type ArgsDeleteSend struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteSend 删除赠与
func DeleteSend(args *ArgsDeleteSend) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_ticket_send", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err == nil {
		_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_ticket_send_log", "send_id = :send_id", map[string]interface{}{
			"send_id": args.ID,
		})
	}
	return
}
