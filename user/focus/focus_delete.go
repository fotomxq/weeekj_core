package UserFocus

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserFocus2 "github.com/fotomxq/weeekj_core/v5/user/focus2"
)

// ArgsDelete 删除关注参数
type ArgsDelete struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// Delete 删除关注
// Deprecated
func Delete(args *ArgsDelete) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "user_focus", "id = :id AND (:user_id < 1 OR user_id = :user_id)", args)
	if err != nil {
		return
	}
	var data FieldsFocus
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT * FROM user_focus WHERE id = $1", args.ID)
	if err == nil {
		//废弃处理，由于模块即将处理，所以此处为固定处理方案
		newMark := data.Mark
		switch data.Mark {
		case "like":
		case "focus":
		default:
			newMark = "like"
		}
		newSystem := ""
		switch data.FromInfo.System {
		case "blog":
			newSystem = "blog_content"
		case "user":
			newSystem = "user_core"
		case "info_exchange":
			newSystem = "info_exchange"
		case "mall":
			newSystem = "mall_product"
		default:
		}
		if newSystem != "" {
			_ = UserFocus2.SetFocus(UserFocus2.ArgsSetFocus{
				UserID:  data.UserID,
				Mark:    newMark,
				System:  newSystem,
				BindID:  data.FromInfo.ID,
				IsFocus: false,
			})
		}
	}
	return
}

// ArgsDeleteByUserFrom 根据来源删除关注参数
type ArgsDeleteByUserFrom struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//关注类型
	Mark string `db:"mark" json:"mark"`
	//关注内容来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
}

// DeleteByUserFrom 根据来源删除关注
func DeleteByUserFrom(args *ArgsDeleteByUserFrom) (err error) {
	where := "user_id = :user_id AND mark = :mark"
	maps := map[string]interface{}{
		"mark":    args.Mark,
		"user_id": args.UserID,
	}
	where, maps, err = args.FromInfo.GetListAnd("from_info", "from_info", where, maps)
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "user_focus", where, args)
	return
}

// ArgsDeleteByOrg 删除组织的所有关注参数
type ArgsDeleteByOrg struct {
	//绑定组织
	OrgID int64 `db:"org_id" json:"orgID"`
}

// DeleteByOrg 删除组织的所有关注
func DeleteByOrg(args *ArgsDeleteByOrg) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "user_focus", "org_id = :org_id", args)
	return
}

// ArgsDeleteByFrom 删除某个来源的所有关注参数
type ArgsDeleteByFrom struct {
	//关注内容来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
}

// DeleteByFrom 删除某个来源的所有关注
func DeleteByFrom(args *ArgsDeleteByFrom) (err error) {
	var fromInfo string
	fromInfo, err = args.FromInfo.GetRaw()
	if err != nil {
		return
	}
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "user_focus", "from_info @> :from_info", map[string]interface{}{
		"from_info": fromInfo,
	})
	return
}
