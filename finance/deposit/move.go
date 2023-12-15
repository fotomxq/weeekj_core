package FinanceDeposit

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	FinanceDeposit2 "github.com/fotomxq/weeekj_core/v5/finance/deposit2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// runMove 将所有数据同步到新版本中
func runMove() {
	//分批获取数据
	step := 0
	for {
		var dataList []FieldsDepositType
		err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT * FROM finance_deposit WHERE save_price > 0 LIMIT 1000 OFFSET $1", step)
		if err != nil || len(dataList) < 1 {
			break
		}
		for _, v := range dataList {
			syncChangeNewSavings(&v, v.SavePrice, "hash", true)
		}
		//延迟，避免高峰
		step += 1000
	}
}

// 同步存入数据
func syncChangeNewSavings(data *FieldsDepositType, addPrice int64, hash string, isAutoRun bool) {
	//初始化
	var errCode string
	var err error
	//用户和商户储蓄识别
	// createInfo 归属关系
	// fromInfo 存储来源
	var vUserID int64 = 0
	var vOrgID int64 = 0
	var vFromOrgID int64 = 0
	if data.FromInfo.System == "org" && data.FromInfo.ID > 0 {
		vFromOrgID = data.FromInfo.ID
	}
	switch data.CreateInfo.System {
	case "org":
		if data.CreateInfo.ID > 0 {
			vOrgID = data.CreateInfo.ID
		}
	case "user":
		if data.CreateInfo.ID > 0 {
			vUserID = data.CreateInfo.ID
		}
	}
	//存储过程刻意加入hash，可确保移动后不会重复构建数据
	//识别储蓄类型
	switch data.ConfigMark {
	case "user":
		//用户一般不会用到，和savings一样处理
		if vFromOrgID > 0 {
			if vUserID > 0 {
				priceData := FinanceDeposit2.GetUserSavingPriceByOrgID(vFromOrgID, vUserID)
				if isAutoRun && priceData > 0 {
					return
				}
				errCode, err = FinanceDeposit2.SetUserSaving(hash, vFromOrgID, vUserID, addPrice)
			} else {
				if vOrgID > 0 {
					priceData := FinanceDeposit2.GetOrgSavingPriceByOrgID(vOrgID)
					if isAutoRun && priceData > 0 {
						return
					}
					errCode, err = FinanceDeposit2.SetOrgSaving(hash, vOrgID, addPrice)
				}
			}
		} else {
			if vUserID > 0 {
				priceData := FinanceDeposit2.GetUserSavingPriceByOrgID(vFromOrgID, vUserID)
				if isAutoRun && priceData > 0 {
					return
				}
				errCode, err = FinanceDeposit2.SetUserSaving(hash, 0, vUserID, addPrice)
			} else {
				if vOrgID > 0 {
					priceData := FinanceDeposit2.GetOrgSavingPriceByOrgID(vOrgID)
					if isAutoRun && priceData > 0 {
						return
					}
					errCode, err = FinanceDeposit2.SetOrgSaving(hash, vOrgID, addPrice)
				}
			}
		}
	case "saving":
		//用户储蓄
		if vFromOrgID > 0 {
			if vUserID > 0 {
				priceData := FinanceDeposit2.GetUserSavingPriceByOrgID(vFromOrgID, vUserID)
				if isAutoRun && priceData > 0 {
					return
				}
				errCode, err = FinanceDeposit2.SetUserSaving(hash, vFromOrgID, vUserID, addPrice)
			} else {
				if vOrgID > 0 {
					priceData := FinanceDeposit2.GetOrgSavingPriceByOrgID(vOrgID)
					if isAutoRun && priceData > 0 {
						return
					}
					errCode, err = FinanceDeposit2.SetOrgSaving(hash, vOrgID, addPrice)
				}
			}
		} else {
			if vUserID > 0 {
				priceData := FinanceDeposit2.GetUserSavingPriceByOrgID(vFromOrgID, vUserID)
				if isAutoRun && priceData > 0 {
					return
				}
				errCode, err = FinanceDeposit2.SetUserSaving(hash, 0, vUserID, addPrice)
			} else {
				if vOrgID > 0 {
					priceData := FinanceDeposit2.GetOrgSavingPriceByOrgID(vOrgID)
					if isAutoRun && priceData > 0 {
						return
					}
					errCode, err = FinanceDeposit2.SetOrgSaving(hash, vOrgID, addPrice)
				}
			}
		}
	case "savings":
		//用户储蓄
		if vFromOrgID > 0 {
			if vUserID > 0 {
				priceData := FinanceDeposit2.GetUserSavingPriceByOrgID(vFromOrgID, vUserID)
				if isAutoRun && priceData > 0 {
					return
				}
				errCode, err = FinanceDeposit2.SetUserSaving(hash, vFromOrgID, vUserID, addPrice)
			} else {
				if vOrgID > 0 {
					priceData := FinanceDeposit2.GetOrgSavingPriceByOrgID(vOrgID)
					if isAutoRun && priceData > 0 {
						return
					}
					errCode, err = FinanceDeposit2.SetOrgSaving(hash, vOrgID, addPrice)
				}
			}
		} else {
			if vUserID > 0 {
				priceData := FinanceDeposit2.GetUserSavingPriceByOrgID(vFromOrgID, vUserID)
				if isAutoRun && priceData > 0 {
					return
				}
				errCode, err = FinanceDeposit2.SetUserSaving(hash, 0, vUserID, addPrice)
			} else {
				if vOrgID > 0 {
					priceData := FinanceDeposit2.GetOrgSavingPriceByOrgID(vOrgID)
					if isAutoRun && priceData > 0 {
						return
					}
					errCode, err = FinanceDeposit2.SetOrgSaving(hash, vOrgID, addPrice)
				}
			}
		}
	case "deposit":
		//押金处理
		if vFromOrgID > 0 {
			if vUserID > 0 {
				priceData := FinanceDeposit2.GetUserDepositPriceByOrgID(vFromOrgID, vUserID)
				if isAutoRun && priceData > 0 {
					return
				}
				errCode, err = FinanceDeposit2.SetUserDeposit(hash, vFromOrgID, vUserID, addPrice)
			} else {
				if vOrgID > 0 {
					priceData := FinanceDeposit2.GetOrgDepositPriceByOrgID(vOrgID)
					if isAutoRun && priceData > 0 {
						return
					}
					errCode, err = FinanceDeposit2.SetOrgDeposit(hash, vOrgID, addPrice)
				}
			}
		} else {
			if vUserID > 0 {
				priceData := FinanceDeposit2.GetUserDepositPriceByOrgID(vFromOrgID, vUserID)
				if isAutoRun && priceData > 0 {
					return
				}
				errCode, err = FinanceDeposit2.SetUserDeposit(hash, 0, vUserID, addPrice)
			} else {
				if vOrgID > 0 {
					priceData := FinanceDeposit2.GetOrgDepositPriceByOrgID(vOrgID)
					if isAutoRun && priceData > 0 {
						return
					}
					errCode, err = FinanceDeposit2.SetOrgDeposit(hash, vOrgID, addPrice)
				}
			}
		}
	case "free":
		//免费储蓄
		if vFromOrgID > 0 {
			if vUserID > 0 {
				priceData := FinanceDeposit2.GetUserFreePriceByOrgID(vFromOrgID, vUserID)
				if isAutoRun && priceData > 0 {
					return
				}
				errCode, err = FinanceDeposit2.SetUserFree(hash, vFromOrgID, vUserID, addPrice)
			}
		} else {
			if vUserID > 0 {
				priceData := FinanceDeposit2.GetUserFreePriceByOrgID(vFromOrgID, vUserID)
				if isAutoRun && priceData > 0 {
					return
				}
				errCode, err = FinanceDeposit2.SetUserFree(hash, 0, vUserID, addPrice)
			}
		}
	}
	if err != nil {
		CoreLog.Warn("update finance deposit to v2 error, code: ", errCode, ", err:  ", err)
	}
}
