package ServiceInfoExchange

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsTakeInfo 参与报名参数
type ArgsTakeInfo struct {
	//参与用户
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//信息ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id"`
	//备注信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"500" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// TakeInfo 参与报名
func TakeInfo(args *ArgsTakeInfo) (errCode string, err error) {
	//查询信息报名信息
	infoData, _ := GetInfoPublishID(&ArgsGetInfoID{
		ID:     args.InfoID,
		OrgID:  -1,
		UserID: -1,
	})
	if !CoreSQL.CheckTimeThanNow(infoData.ExpireAt) {
		errCode = "err_expire"
		err = errors.New("no data")
		return
	}
	//检查报名人数
	if infoData.LimitCount == 0 {
		errCode = "err_limit"
		err = errors.New("no data")
		return
	}
	if infoData.LimitCount > 0 {
		var count int64
		_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM service_info_exchange_take WHERE info_id = $1", args.InfoID)
		if count > infoData.LimitCount {
			errCode = "err_too_many"
			err = errors.New("to many")
			return
		}
	}
	//检查是否报名过
	var id int64
	_ = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM service_info_exchange_take WHERE info_id = $1 AND user_id = $2", args.InfoID, args.UserID)
	if id > 0 {
		errCode = "err_have_replace"
		err = errors.New("have data")
		return
	}
	//创建数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_info_exchange_take (user_id, info_id, des, params) VALUES (:user_id,:info_id,:des,:params)", map[string]interface{}{
		"info_id": args.InfoID,
		"user_id": args.UserID,
		"des":     args.Des,
		"params":  args.Params,
	})
	if err != nil {
		errCode = "err_insert"
	}
	return
}
