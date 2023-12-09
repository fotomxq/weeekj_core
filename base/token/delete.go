package BaseToken

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteByID 删除参数
type ArgsDeleteByID struct {
	//ID
	ID int64 `db:"id" json:"id"`
}

// DeleteByID 删除
func DeleteByID(args *ArgsDeleteByID) (err error) {
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_token", "id", args)
	if err != nil {
		return
	}
	return
}

// ArgsDeleteByIDFrom 删除指定范围的ID参数
type ArgsDeleteByIDFrom struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//来源
	FromInfo CoreSQLFrom.FieldsFrom `json:"fromInfo"`
}

// DeleteByIDFrom 删除指定范围的ID
func DeleteByIDFrom(args *ArgsDeleteByIDFrom) (err error) {
	where := "id = :id"
	maps := map[string]interface{}{
		"id": args.ID,
	}
	where, maps, err = args.FromInfo.GetListAnd("from_info", "from_info", where, maps)
	if err != nil {
		return
	}
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_token", where, maps)
	if err != nil {
		return
	}
	return
}

// ArgsDeleteByFrom 删除某个From所有的数据参数
type ArgsDeleteByFrom struct {
	//来源
	FromInfo CoreSQLFrom.FieldsFrom `json:"fromInfo"`
	//登陆渠道
	LoginInfo CoreSQLFrom.FieldsFrom `json:"loginInfo"`
}

// DeleteByFrom 删除某个From所有的数据
func DeleteByFrom(args *ArgsDeleteByFrom) (err error) {
	var where string
	var maps map[string]interface{}
	where, maps, err = args.FromInfo.GetListAnd("from_info", "from_info", where, maps)
	if err != nil {
		return
	}
	where, maps, err = args.LoginInfo.GetListAnd("login_info", "login_info", where, maps)
	if err != nil {
		return
	}
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_token", where, maps)
	if err != nil {
		return
	}
	return
}

// DeleteByExpire 将所有过期人群踢下线
func DeleteByExpire() {
	_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_token", "expire_at < NOW()", nil)
	return
}

// DeleteAll 将所有人踢下线
func DeleteAll() {
	_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_token", "true", nil)
	return
}

// ArgsClearAndCreate 登陆成功后的处理规则参数
type ArgsClearAndCreate struct {
	//要删除的tokenID
	OldTokenID int64
	//登陆来源
	FromInfo CoreSQLFrom.FieldsFrom
	//登陆渠道
	LoginInfo CoreSQLFrom.FieldsFrom
	//key长度限制
	LimitKeyLen int
	//IP地址
	IP string
	//过期时间
	// 请使用RFC3339Nano结构时间，eg: 2020-11-03T08:31:13.314Z
	// JS中使用new Date().toISOString()
	ExpireAt string
	//是否记住我
	IsRemember bool `db:"is_remember" json:"isRemember"`
}

// ClearAndCreate 登陆成功后的处理规则
func ClearAndCreate(args *ArgsClearAndCreate) (data FieldsTokenType, errCode string, err error) {
	//如果存在旧的ID，则删除
	if args.OldTokenID > 0 {
		err = DeleteByID(&ArgsDeleteByID{
			ID: args.OldTokenID,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("cannot remove old token, token id: ", args.OldTokenID, ", err: ", err))
			err = nil
		}
	}
	//清理该用户同来源和渠道的旧的数据
	err = DeleteByFrom(&ArgsDeleteByFrom{
		FromInfo: args.FromInfo, LoginInfo: args.LoginInfo,
	})
	if err != nil {
		//err = errors.New(fmt.Sprint("delete by from all token, from info: ", args.FromInfo, ", login info: ", args.LoginInfo, ", err: ", err))
		//return
		//不记录错误信息，避免无登陆的异常行为
		err = nil
	}
	//创建新的token
	data, errCode, err = Create(&ArgsCreate{
		FromInfo:    args.FromInfo,
		LoginInfo:   args.LoginInfo,
		LimitKeyLen: args.LimitKeyLen,
		IP:          args.IP,
		ExpireAt:    args.ExpireAt,
		IsRemember:  args.IsRemember,
	})
	if err != nil {
		err = errors.New("create token failed, " + err.Error())
		return
	}
	//反馈
	return
}
