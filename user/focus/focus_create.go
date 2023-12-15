package UserFocus

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserFocus2 "github.com/fotomxq/weeekj_core/v5/user/focus2"
)

// ArgsCreate 创建关注参数
type ArgsCreate struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//绑定组织
	OrgID int64 `db:"org_id" json:"orgID"`
	//关注类型
	Mark string `db:"mark" json:"mark"`
	//关注内容来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
}

// Create 创建关注
// Deprecated
func Create(args *ArgsCreate) (data FieldsFocus, err error) {
	var fromInfo string
	fromInfo, err = args.FromInfo.GetRaw()
	if err != nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, delete_at, user_id, org_id, mark, from_info FROM user_focus WHERE user_id = $1 AND org_id = $2 AND mark = $3 AND from_info @> $4", args.UserID, args.OrgID, args.Mark, fromInfo)
	if err == nil {
		if data.DeleteAt.Unix() > 0 {
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_focus SET delete_at = to_timestamp(0) WHERE id = :id", map[string]interface{}{
				"id": data.ID,
			})
			if err == nil {
				err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, delete_at, user_id, org_id, mark, from_info FROM user_focus WHERE id = $1", data.ID)
			}
		}
	} else {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "user_focus", "INSERT INTO user_focus (user_id, org_id, mark, from_info) VALUES (:user_id, :org_id, :mark, :from_info)", args, &data)
	}
	//废弃处理，由于模块即将处理，所以此处为固定处理方案
	newMark := args.Mark
	switch args.Mark {
	case "like":
	case "focus":
	default:
		newMark = "like"
	}
	newSystem := ""
	switch args.FromInfo.System {
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
			UserID:  args.UserID,
			Mark:    newMark,
			System:  newSystem,
			BindID:  args.FromInfo.ID,
			IsFocus: true,
		})
	}
	return
}
