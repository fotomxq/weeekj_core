package FinanceDeposit2

import (
	"errors"
	"fmt"
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// GetUserSavingByOrgID 获取组织押金
func GetUserSavingByOrgID(orgID int64, userID int64) (data FieldsUserSaving) {
	cacheMark := getUserSavingCacheMark(orgID, userID)
	err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data)
	if err == nil && data.ID > 0 {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, update_hash, org_id, user_id, price FROM finance_deposit2_user_saving WHERE org_id = $1 AND user_id = $2", orgID, userID)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

func GetUserSavingPriceByOrgID(orgID int64, userID int64) (price int64) {
	data := GetUserSavingByOrgID(orgID, userID)
	return data.Price
}

// GetUserSavingTotalPriceByOrgID 获取当前储蓄总额
func GetUserSavingTotalPriceByOrgID(orgID int64) (price int64) {
	_ = Router2SystemConfig.MainDB.Get(&price, "SELECT SUM(price) as price FROM finance_deposit2_user_saving WHERE ($1 > 0 AND org_id = $1) OR $1 < 1", orgID)
	return
}

// SetUserSaving 变更组织押金
func SetUserSaving(hash string, orgID int64, userID int64, addPrice int64) (errCode string, err error) {
	userSavingLock.Lock()
	defer userSavingLock.Unlock()
	data := GetUserSavingByOrgID(orgID, userID)
	newHash := CoreFilter.GetSha1Str(fmt.Sprint(orgID, userID, addPrice, CoreFilter.GetRandStr4(6)))
	nowAt := CoreFilter.GetNowTimeCarbon()
	if data.ID < 1 {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO finance_deposit2_user_saving(update_hash, org_id, user_id, price) VALUES(:update_hash, :org_id, :user_id, :price)", map[string]interface{}{
			"update_hash": newHash,
			"org_id":      orgID,
			"user_id":     userID,
			"price":       addPrice,
		})
		if err != nil {
			errCode = "err_insert"
			return
		}
	} else {
		if hash != "" && data.UpdateHash != hash {
			errCode = "err_hash"
			err = errors.New("hash error")
			return
		}
		if addPrice == 0 {
			return
		}
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE finance_deposit2_user_saving SET update_at = NOW(), update_hash = :update_hash, price = price + :price WHERE id = :id", map[string]interface{}{
			"update_hash": newHash,
			"price":       addPrice,
			"id":          data.ID,
		})
		if err != nil {
			errCode = "err_update"
			return
		}
	}
	//统计总金额
	totalPrice := GetUserSavingTotalPriceByOrgID(0)
	AnalysisAny2.AppendData("re", "finance_deposit_total_user_saving_price", nowAt.Time, 0, 0, 0, 0, 0, totalPrice)
	if orgID > 0 {
		totalPriceByOrg := GetUserSavingTotalPriceByOrgID(orgID)
		AnalysisAny2.AppendData("re", "finance_deposit_total_user_saving_price", nowAt.Time, orgID, 0, 0, 0, 0, totalPriceByOrg)
	}
	//删除缓冲
	deleteUserSavingCache(data.OrgID, data.UserID)
	//反馈
	return
}

// 缓冲
func getUserSavingCacheMark(orgID int64, userID int64) string {
	return fmt.Sprint("finance:deposit2:user:saving:user:", orgID, ".", userID)
}

func deleteUserSavingCache(orgID int64, userID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getUserSavingCacheMark(orgID, userID))
}
