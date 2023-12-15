package BaseSMS

import (
	"errors"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	ClassConfig "github.com/fotomxq/weeekj_core/v5/class/config"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//短信模块
// 提供多家渠道的短信发送、消息模版处理
// 警告，本模块不支持并发，请在并发服务前将该模块摘出

// ArgsGetSMSList 获取短信列表参数
type ArgsGetSMSList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//短信配置
	ConfigID int64 `json:"configID" check:"id" empty:"true"`
	//手机号
	NationCode string `json:"nationCode" check:"nationCode" empty:"true"`
	Phone      string `json:"phone" check:"phone" empty:"true"`
	//是否过期
	NeedIsExpire bool `json:"needIsExpire" check:"bool"`
	IsExpire     bool `json:"isExpire" check:"bool"`
	//是否已经发送
	NeedIsSend bool `json:"needIsSend" check:"bool"`
	IsSend     bool `json:"isSend" check:"bool"`
	//是否发送失败
	NeedIsFailed bool `json:"needIsFailed" check:"bool"`
	IsFailed     bool `json:"isFailed" check:"bool"`
}

// GetSMSList 获取短信列表
func GetSMSList(args *ArgsGetSMSList) (dataList []FieldsSMS, dataCount int64, err error) {
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
	if args.NationCode != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "nation_code = :nation_code"
		maps["nation_code"] = args.NationCode
	}
	if args.Phone != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "phone = :phone"
		maps["phone"] = args.Phone
	}
	if args.NeedIsExpire {
		if where != "" {
			where = where + " AND "
		}
		if args.IsExpire {
			where = where + "expire_at > to_timestamp(1000000)"
		} else {
			where = where + "expire_at <= to_timestamp(1000000)"
		}
	}
	if args.NeedIsSend {
		if where != "" {
			where = where + " AND "
		}
		if args.IsSend {
			where = where + "send_at > to_timestamp(1000000)"
		} else {
			where = where + "send_at <= to_timestamp(1000000)"
		}
	}
	if args.NeedIsFailed {
		if where != "" {
			where = where + " AND "
		}
		if args.IsFailed {
			where = where + "failed_msg != ''"
		} else {
			where = where + "failed_msg = ''"
		}
	}
	if where == "" {
		where = "true"
	}
	tableName := "core_sms"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, expire_at, send_at, failed_msg, is_check, config_id, token, nation_code, phone, params, from_info "+"FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "expire_at", "send_at"},
	)
	return
}

// 获取短信配置
func getConfigByDefault(orgID, configID int64) (configData FieldsConfig, err error) {
	//采用默认短信配置发送验证码
	if configID < 1 {
		if orgID > 0 {
			configID, err = OrgCore.Config.GetConfigValInt64(&ClassConfig.ArgsGetConfig{
				BindID:    orgID,
				Mark:      "VerificationCodeSMSDefault",
				VisitType: "admin",
			})
			if err != nil {
				//err = errors.New("get config by VerificationCodeSMSDefault")
				//return
			}
		}
		if configID < 1 {
			orgID = 0
			configID, err = BaseConfig.GetDataInt64("VerificationCodeSMSDefault")
			if err != nil {
				err = errors.New("get config by VerificationCodeSMSDefault")
				return
			}
		}
		if configID < 1 {
			err = errors.New("no config")
			return
		}
	}
	//获取配置
	configData, err = GetConfigByID(&ArgsGetConfigByID{
		ID:    configID,
		OrgID: orgID,
	})
	if err != nil {
		err = errors.New("cannot get config, " + err.Error())
		return
	}
	return
}
