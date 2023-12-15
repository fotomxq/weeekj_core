package BaseSMS

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"strconv"
	"time"
)

// ArgsCreateSMS 创建新的短信请求参数
type ArgsCreateSMS struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `json:"configID"`
	//会话
	Token int64 `json:"token"`
	//电话
	NationCode string `json:"nationCode"`
	Phone      string `json:"phone"`
	//短信内容
	Params CoreSQLConfig.FieldsConfigsType `json:"params"`
	//创建来源
	FromInfo CoreSQLFrom.FieldsFrom `json:"fromInfo"`
}

// CreateSMS 创建新的短信请求
func CreateSMS(args *ArgsCreateSMS) (errCode string, err error) {
	//获取基础配置
	var configData FieldsConfig
	configData, err = getConfigByDefault(args.OrgID, args.ConfigID)
	if err != nil {
		errCode = "config_not_exist"
		err = errors.New("cannot get config, " + err.Error())
		return
	}
	//验证短信内容合规
	isCheck, b := configData.Params.GetValBool("check")
	if !b {
		isCheck = false
	}
	if !isCheck {
		for _, v := range configData.TemplateParams {
			isFind := false
			for _, v2 := range args.Params {
				if v.Mark == v2.Mark {
					isFind = true
					break
				}
			}
			if !isFind {
				errCode = "template_lost"
				err = errors.New("sms template error")
				return
			}
		}
	}
	//计算时间间隔
	if args.Token > 0 {
		//获取该来源是否存在其他验证码？且没有超时？
		waitTime := CoreFilter.GetNowTimeCarbon().SubSeconds(int(configData.TimeSpacing)).Time
		var oldData FieldsSMS
		err = Router2SystemConfig.MainDB.Get(&oldData, "SELECT id FROM core_sms WHERE config_id = $1 AND token = $2 AND expire_at >= NOW() AND create_at >= $3", args.ConfigID, args.Token, waitTime)
		if err == nil {
			errCode = "have_other_sms"
			err = errors.New("token have sms wait check")
			return
		}
	}
	//生成过期时间
	var expireAt time.Time
	expireAt, err = CoreFilter.GetTimeByAdd(configData.DefaultExpire)
	if err != nil {
		errCode = "expire_failed"
		err = errors.New("cannot get expire time, " + err.Error())
		return
	}
	//记录数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO core_sms(org_id, expire_at, send_at, failed_msg, is_check, config_id, token, nation_code, phone, params, from_info) VALUES (:org_id, :expire_at, to_timestamp(0), '', false, :config_id, :token, :nation_code, :phone, :params, :from_info)", map[string]interface{}{
		"org_id":      configData.OrgID,
		"expire_at":   expireAt,
		"config_id":   configData.ID,
		"token":       args.Token,
		"nation_code": args.NationCode,
		"phone":       args.Phone,
		"params":      args.Params,
		"from_info":   args.FromInfo,
	})
	if err != nil {
		errCode = "insert_data"
		err = errors.New("cannot create sms data, " + err.Error())
	}
	return
}

// ArgsCreateSMSCheck 创建验证用的短信验证码参数
type ArgsCreateSMSCheck struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `json:"configID"`
	//会话
	Token int64 `json:"token"`
	//电话
	NationCode string `json:"nationCode"`
	Phone      string `json:"phone"`
	//创建来源
	FromInfo CoreSQLFrom.FieldsFrom `json:"fromInfo"`
}

// CreateSMSCheck 创建验证用的短信验证码
func CreateSMSCheck(args *ArgsCreateSMSCheck) (errCode string, err error) {
	//获取基础配置
	var configData FieldsConfig
	configData, err = getConfigByDefault(args.OrgID, args.ConfigID)
	if err != nil {
		errCode = "config_not_exist"
		err = errors.New("cannot get config, " + err.Error())
		return
	}
	isCheck, b := configData.Params.GetValBool("check")
	if !b || !isCheck {
		err = errors.New("sms config not check")
		return
	}
	//计算过期时间
	var expireTime int64
	expireTime, err = CoreFilter.GetTimeBetweenAdd(configData.DefaultExpire)
	if err != nil {
		err = errors.New("cannot get between time, " + err.Error())
		return
	}
	//验证码
	randNumber := CoreFilter.GetRandNumber(1000, 9999)
	randStr := CoreFilter.GetStringByInt(randNumber)
	var params CoreSQLConfig.FieldsConfigsType
	valMark, b := configData.TemplateParams.GetVal("val")
	if !b {
		err = errors.New("config error")
		return
	}
	if valMark != "__skip" {
		params = append(params, CoreSQLConfig.FieldsConfigType{
			Mark: valMark,
			Val:  randStr,
		})
	}
	timeMark, b := configData.TemplateParams.GetVal("time")
	if !b {
		err = errors.New("config error")
		return
	}
	switch timeMark {
	case "__skip":
		break
	case "second":
		params = append(params, CoreSQLConfig.FieldsConfigType{
			Mark: "time",
			Val:  strconv.FormatInt(expireTime, 10),
		})
	case "minutes":
		params = append(params, CoreSQLConfig.FieldsConfigType{
			Mark: "time",
			Val:  strconv.FormatInt(expireTime/60, 10),
		})
	default:
		params = append(params, CoreSQLConfig.FieldsConfigType{
			Mark: "time",
			Val:  strconv.FormatInt(expireTime, 10),
		})
		break
	}
	return CreateSMS(&ArgsCreateSMS{
		OrgID:      args.OrgID,
		ConfigID:   args.ConfigID,
		Token:      args.Token,
		NationCode: args.NationCode,
		Phone:      args.Phone,
		Params:     params,
		FromInfo:   CoreSQLFrom.FieldsFrom{},
	})
}
