package BaseFileSys2

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
)

func DeleteClaim(claimId int64, orgID int64, userID int64) (err error) {
	data := getClaimByID(claimId)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) || !CoreFilter.EqID2(userID, data.UserID) {
		err = errors.New("no data")
		return
	}
	err = claimDB.Delete().AddWhereID(data.ID).Exec("id = $1", data.ID)
	if err != nil {
		return
	}
	deleteClaimCache(data.ID)
	count := claimDB.Analysis().DataNoDelete().Count("file_id = $1", data.FileID)
	if count < 1 {
		err = DeleteCore(data.FileID, -1, -1)
		if err != nil {
			return
		}
	}
	return
}
