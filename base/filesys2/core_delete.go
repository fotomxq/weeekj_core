package BaseFileSys2

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
)

func DeleteCore(id int64, orgID int64, userID int64) (err error) {
	data := getCoreByID(id)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) || !CoreFilter.EqID2(userID, data.UserID) {
		err = errors.New("no data")
		return
	}
	err = coreDB.Delete().AddWhereID(data.ID).Exec("id = $1", data.ID)
	if err != nil {
		return
	}
	deleteCoreCache(data.ID)
	return
}
