package FinanceDeposit2

import (
	"errors"
	"fmt"
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// GetOrgSavingByOrgID 获取组织押金
func GetOrgSavingByOrgID(orgID int64) (data FieldsOrgSaving) {
	cacheMark := getOrgSavingCacheMark(orgID)
	err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data)
	if err == nil && data.ID > 0 {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, update_hash, org_id, price FROM finance_deposit2_org_saving WHERE org_id = $1", orgID)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

func GetOrgSavingPriceByOrgID(orgID int64) (price int64) {
	data := GetOrgSavingByOrgID(orgID)
	return data.Price
}

// GetOrgSavingTotalPrice 获取当前储蓄总额
func GetOrgSavingTotalPrice() (price int64) {
	//获取数据
	_ = Router2SystemConfig.MainDB.Get(&price, "SELECT SUM(price) as price FROM finance_deposit2_org_saving")
	//反馈
	return
}

// SetOrgSaving 变更组织押金
func SetOrgSaving(hash string, orgID int64, addPrice int64) (errCode string, err error) {
	orgSavingLock.Lock()
	defer orgSavingLock.Unlock()
	data := GetOrgSavingByOrgID(orgID)
	newHash := CoreFilter.GetSha1Str(fmt.Sprint(orgID, addPrice, CoreFilter.GetRandStr4(6)))
	nowAt := CoreFilter.GetNowTimeCarbon()
	if data.ID < 1 {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO finance_deposit2_org_saving(update_hash, org_id, price) VALUES(:update_hash, :org_id, :price)", map[string]interface{}{
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
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE finance_deposit2_org_saving SET update_at = NOW(), update_hash = :update_hash, price = price + :price WHERE id = :id", map[string]interface{}{
			"update_hash": newHash,
			"price":       addPrice,
			"id":          data.ID,
		})
		if err != nil {
			errCode = "err_update"
			return
		}
	}
	//统计行为
	if addPrice > 0 {
		AnalysisAny2.AppendData("add", "finance_deposit_org_add_price", time.Time{}, orgID, 0, 0, 0, 0, addPrice)
	} else {
		AnalysisAny2.AppendData("reduce", "finance_deposit_org_add_price", time.Time{}, orgID, 0, 0, 0, 0, 0-addPrice)
	}
	//统计总金额
	totalPrice := GetOrgSavingTotalPrice()
	AnalysisAny2.AppendData("re", "finance_deposit_total_org_saving_price", nowAt.Time, 0, 0, 0, 0, 0, totalPrice)
	//删除缓冲
	deleteOrgSavingCache(data.OrgID)
	//反馈
	return
}

// 缓冲
func getOrgSavingCacheMark(orgID int64) string {
	return fmt.Sprint("finance:deposit2:org:saving:org:", orgID)
}

func deleteOrgSavingCache(orgID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getOrgSavingCacheMark(orgID))
}
