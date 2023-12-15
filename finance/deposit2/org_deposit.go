package FinanceDeposit2

import (
	"errors"
	"fmt"
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// GetOrgDepositByOrgID 获取组织押金
func GetOrgDepositByOrgID(orgID int64) (data FieldsOrgDeposit) {
	cacheMark := getOrgDepositCacheMark(orgID)
	err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data)
	if err == nil && data.ID > 0 {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, update_hash, org_id, price FROM finance_deposit2_org_deposit WHERE org_id = $1", orgID)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

func GetOrgDepositPriceByOrgID(orgID int64) (price int64) {
	data := GetOrgDepositByOrgID(orgID)
	return data.Price
}

// GetOrgDepositTotalPrice 获取当前储蓄总额
func GetOrgDepositTotalPrice() (price int64) {
	_ = Router2SystemConfig.MainDB.Get(&price, "SELECT SUM(price) as price FROM finance_deposit2_org_deposit")
	return
}

// SetOrgDeposit 变更组织押金
func SetOrgDeposit(hash string, orgID int64, addPrice int64) (errCode string, err error) {
	orgDepositLock.Lock()
	defer orgDepositLock.Unlock()
	data := GetOrgDepositByOrgID(orgID)
	newHash := CoreFilter.GetSha1Str(fmt.Sprint(orgID, addPrice, CoreFilter.GetRandStr4(6)))
	nowAt := CoreFilter.GetNowTimeCarbon()
	if data.ID < 1 {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO finance_deposit2_org_deposit(update_hash, org_id, price) VALUES(:update_hash, :org_id, :price)", map[string]interface{}{
			"update_hash": newHash,
			"org_id":      orgID,
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
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE finance_deposit2_org_deposit SET update_at = NOW(), update_hash = :update_hash, price = price + :price WHERE id = :id", map[string]interface{}{
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
	totalPrice := GetOrgDepositTotalPrice()
	AnalysisAny2.AppendData("re", "finance_deposit_total_org_deposit_price", nowAt.Time, 0, 0, 0, 0, 0, totalPrice)
	//删除缓冲
	deleteOrgDepositCache(data.OrgID)
	//反馈
	return
}

// 缓冲
func getOrgDepositCacheMark(orgID int64) string {
	return fmt.Sprint("finance:deposit2:org:deposit:org:", orgID)
}

func deleteOrgDepositCache(orgID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getOrgDepositCacheMark(orgID))
}
