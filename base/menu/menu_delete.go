package BaseMenu

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteMenu 删除目录参数
type ArgsDeleteMenu struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteMenu 删除目录
func DeleteMenu(args *ArgsDeleteMenu) (errCode string, err error) {
	//检查是否存在下一级
	count := GetMenuCountByParentID(args.ID)
	if count > 0 {
		errCode = "err_base_menu_have_child"
		err = errors.New("have child")
		return
	}
	//删除目录
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "core_menu", "id = :id AND org_id = :org_id", args)
	if err != nil {
		errCode = "err_delete"
		return
	}
	//删除缓冲
	deleteMenuCache(args.ID)
	//反馈
	return
}
