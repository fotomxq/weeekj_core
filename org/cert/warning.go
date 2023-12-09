package OrgCert

import (
	"errors"
	"fmt"
	BaseSMS "gitee.com/weeekj/weeekj_core/v5/base/sms"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	OrgCoreCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
	UserMessage "gitee.com/weeekj/weeekj_core/v5/user/message"
	"github.com/lib/pq"
	"time"
)

// ArgsGetWarningList 获取异常列表参数
type ArgsGetWarningList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子组织ID
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//是否反馈
	NeedIsFinish bool `db:"need_is_finish" json:"needIsFinish" check:"bool"`
	IsFinish     bool `db:"is_finish" json:"isFinish" check:"bool"`
	//证件标识码
	ConfigMarks pq.StringArray `db:"config_marks" json:"configMarks" check:"marks" empty:"true"`
	//证件配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//时间段
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetWarningList 获取异常列表参数
func GetWarningList(args *ArgsGetWarningList) (dataList []FieldsWarning, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > 0 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ChildOrgID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "child_org_id = :child_org_id"
		maps["child_org_id"] = args.ChildOrgID
	}
	where = CoreSQL.GetNeedChange(where, "finish_at", args.NeedIsFinish, args.IsFinish)
	if len(args.ConfigMarks) > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "config_mark = ANY(:config_marks)"
		maps["config_marks"] = args.ConfigMarks
	}
	if args.ConfigID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "config_id = :config_id"
		maps["config_id"] = args.ConfigID
	}
	if args.TimeBetween.MinTime != "" && args.TimeBetween.MaxTime != "" {
		var timeBetween CoreSQLTime.FieldsCoreTime
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
		var newWhere string
		newWhere, maps = CoreSQLTime.GetBetweenByTime("create_at", timeBetween, maps)
		if newWhere != "" {
			if where != "" {
				where = where + " AND "
			}
			where = where + newWhere
		}
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(msg ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "org_cert_warning"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, finish_at, org_id, child_org_id, cert_id, config_id, config_mark, msg FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// argsCreateWarning 创建新的记录参数
type argsCreateWarning struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//子商户
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//异常证件ID
	CertID int64 `db:"cert_id" json:"certID" check:"id"`
	//消息
	Msg string `db:"msg" json:"msg" check:"des" min:"1" max:"600"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
	//配置标题
	ConfigName string `json:"configName"`
	//绑定来源
	// user 用户 / org 商户 / org_bind 商户成员 / finance_assets 财务资产 /
	BindFrom string `db:"bind_from" json:"bindFrom"`
	//绑定ID
	// 根据配置决定，可能是用户ID\组织ID\或其他任意主体
	BindID int64 `db:"bind_id" json:"bindID"`
}

// createWarning 创建新的记录
func createWarning(args *argsCreateWarning) (err error) {
	//获取证件信息
	var certData FieldsCert
	certData, err = GetCert(&ArgsGetCert{
		ID:         args.CertID,
		OrgID:      args.OrgID,
		ChildOrgID: args.ChildOrgID,
	})
	if err != nil {
		return
	}
	//获取配置
	var configData FieldsConfig
	configData, err = GetConfigByID(&ArgsGetConfigByID{
		ID:    certData.ConfigID,
		OrgID: -1,
	})
	if err != nil {
		return
	}
	//检查是否存在未处理
	var data FieldsWarning
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at FROM org_cert_warning WHERE org_id = $1 AND child_org_id = $2 AND cert_id = $3 AND finish_at < to_timestamp(1000000)", args.OrgID, args.ChildOrgID, args.CertID)
	if err == nil && data.ID > 0 {
		//检查是否存在每年提醒开关？
		needTipEveryYear, b := configData.Params.GetValBool("needTipEveryYear")
		if !b {
			needTipEveryYear = false
		}
		if !needTipEveryYear {
			err = errors.New("have warning")
			return
		} else {
			//检查是否达到一年
			if data.CreateAt.Unix() >= CoreFilter.GetNowTimeCarbon().SubMonths(11).Time.Unix() {
				err = errors.New("have warning")
				return
			}
		}
	}
	//创建记录
	var warningData FieldsWarning
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_cert_warning", "INSERT INTO org_cert_warning (org_id, child_org_id, cert_id, config_id, config_mark, msg, send_msg_at, send_sms_at, params) VALUES (:org_id, :child_org_id, :cert_id, :config_id, :config_mark, :msg, to_timestamp(0), to_timestamp(0), :params)", map[string]interface{}{
		"org_id":       args.OrgID,
		"child_org_id": args.ChildOrgID,
		"cert_id":      args.CertID,
		"config_id":    configData.ID,
		"config_mark":  configData.Mark,
		"msg":          args.Msg,
		"params":       args.Params,
	}, &warningData)
	if err == nil {
		//根据需求，发送用户消息和短信消息
		needUserMessage, b := args.Params.GetValBool("needUserMessage")
		if b && needUserMessage {
			switch args.BindFrom {
			case "org":
				var orgData OrgCoreCore.FieldsOrg
				orgData, err = OrgCoreCore.GetOrg(&OrgCoreCore.ArgsGetOrg{
					ID: args.BindID,
				})
				if err != nil {
					break
				}
				var params CoreSQLConfig.FieldsConfigsType
				params = append(params, CoreSQLConfig.FieldsConfigType{
					Mark: "orgCertWarningID",
					Val:  fmt.Sprint(warningData.ID),
				})
				_, err = UserMessage.Create(&UserMessage.ArgsCreate{
					WaitSendAt:    time.Now(),
					SendUserID:    0,
					ReceiveUserID: orgData.UserID,
					Title:         args.Msg,
					Content:       args.Msg,
					Files:         []int64{},
					Params:        params,
				})
				if err != nil {
					break
				}
			case "user":
				var params CoreSQLConfig.FieldsConfigsType
				params = append(params, CoreSQLConfig.FieldsConfigType{
					Mark: "orgCertWarningID",
					Val:  fmt.Sprint(warningData.ID),
				})
				_, err = UserMessage.Create(&UserMessage.ArgsCreate{
					WaitSendAt:    time.Now(),
					SendUserID:    0,
					ReceiveUserID: args.BindID,
					Title:         args.Msg,
					Content:       args.Msg,
					Files:         []int64{},
					Params:        params,
				})
				if err != nil {
					break
				}
			}
		}
		needSMS, b := args.Params.GetValBool("needSMS")
		if b && needSMS {
			smsConfigID, b := args.Params.GetValInt64("smsConfigID")
			if b && smsConfigID > 0 {
				switch args.BindFrom {
				case "org":
					var orgData OrgCoreCore.FieldsOrg
					orgData, err = OrgCoreCore.GetOrg(&OrgCoreCore.ArgsGetOrg{
						ID: args.BindID,
					})
					if err != nil {
						break
					}
					var userData UserCore.FieldsUserType
					userData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
						ID:    orgData.UserID,
						OrgID: -1,
					})
					if err != nil {
						break
					}
					if userData.NationCode == "" || userData.Phone == "" {
						break
					}
					var params CoreSQLConfig.FieldsConfigsType
					params = append(params, CoreSQLConfig.FieldsConfigType{
						Mark: "$1",
						Val:  args.ConfigName,
					})
					_, err = BaseSMS.CreateSMS(&BaseSMS.ArgsCreateSMS{
						OrgID:      orgData.ID,
						ConfigID:   smsConfigID,
						Token:      0,
						NationCode: userData.NationCode,
						Phone:      userData.Phone,
						Params:     params,
						FromInfo: CoreSQLFrom.FieldsFrom{
							System: "",
							ID:     0,
							Mark:   "",
							Name:   "组织证件",
						},
					})
					if err != nil {
						break
					}
				case "user":
					var userData UserCore.FieldsUserType
					userData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
						ID:    args.BindID,
						OrgID: -1,
					})
					if err != nil {
						break
					}
					if userData.NationCode == "" || userData.Phone == "" {
						break
					}
					var params CoreSQLConfig.FieldsConfigsType
					params = append(params, CoreSQLConfig.FieldsConfigType{
						Mark: "$1",
						Val:  args.ConfigName,
					})
					_, err = BaseSMS.CreateSMS(&BaseSMS.ArgsCreateSMS{
						OrgID:      userData.OrgID,
						ConfigID:   smsConfigID,
						Token:      0,
						NationCode: userData.NationCode,
						Phone:      userData.Phone,
						Params:     params,
						FromInfo: CoreSQLFrom.FieldsFrom{
							System: "",
							ID:     0,
							Mark:   "",
							Name:   "组织证件",
						},
					})
					if err != nil {
						break
					}
				}
			}
		}
	}
	return
}

// ArgsUpdateWarningFinish 标记已经处理参数
type ArgsUpdateWarningFinish struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子商户
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
}

// UpdateWarningFinish 标记已经处理
func UpdateWarningFinish(args *ArgsUpdateWarningFinish) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_cert_warning SET finish_at = NOW() WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:child_org_id < 1 OR child_org_id = :child_org_id)", args)
	return
}
