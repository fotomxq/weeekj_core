package FinanceSafe

import (
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime2 "github.com/fotomxq/weeekj_core/v5/core/sql/time2"
	FinanceLog "github.com/fotomxq/weeekj_core/v5/finance/log"
	FinancePay "github.com/fotomxq/weeekj_core/v5/finance/pay"
	"time"
)

func runSafe() {
	//获取配置，对最近N天数据进行反复审计处理
	// FinanceSafeCheckBeforeTime
	financeSafeCheckBeforeTimeStr, err := BaseConfig.GetDataString("FinanceSafeCheckBeforeTime")
	if err != nil {
		CoreLog.Error("get config by financeSafeCheckBeforeTime, ", err)
	}
	financeSafeCheckBeforeTime, err := CoreFilter.GetTimeByAdd(financeSafeCheckBeforeTimeStr)
	if err != nil {
		runLog("get financeSafeCheckBeforeTime time, ", err)
		financeSafeCheckBeforeTime = CoreFilter.GetNowTime().AddDate(0, 0, 3)
	}
	//最大金额浮动上限
	financeSafeMaxPrice, err := BaseConfig.GetDataInt64("FinanceSafeMaxPrice")
	if err != nil {
		runLog("get financeSafeMaxPrice time, ", err)
		financeSafeMaxPrice = 1000.00
	}
	//获取负载均衡配置项
	financeSafeFrequencyOneFrom, err := BaseConfig.GetDataInt64("FinanceSafeFrequencyOneFrom")
	if err != nil {
		runLog("get financeSafeFrequencyOneFrom config, ", err)
		financeSafeFrequencyOneFrom = 100
	}
	financeSafeFrequencyOneTo, err := BaseConfig.GetDataInt64("FinanceSafeFrequencyOneTo")
	if err != nil {
		runLog("get financeSafeFrequencyOneTo config, ", err)
		financeSafeFrequencyOneTo = 100
	}
	financeSafeFrequencyAll, err := BaseConfig.GetDataInt64("FinanceSafeFrequencyAll")
	if err != nil {
		runLog("get financeSafeFrequencyAll config, ", err)
		financeSafeFrequencyOneTo = 200
	}
	financeSafeFrequencyTime, err := BaseConfig.GetDataInt64("FinanceSafeFrequencyTime")
	if err != nil {
		runLog("get financeSafeFrequencyTime config, ", err)
		financeSafeFrequencyTime = 60
	}
	//获取该时间节点之后的所有数据，分批获取
	var page int64 = 1
	for {
		dataList, _, err := FinanceLog.GetList(&FinanceLog.ArgsGetList{
			Pages: CoreSQLPages.ArgsDataList{
				Page: page,
				Max:  100,
				Sort: "_id",
				Desc: true,
			},
			Key:            "",
			Status:         []int{},
			PaymentCreate:  CoreSQLFrom.FieldsFrom{},
			PaymentChannel: CoreSQLFrom.FieldsFrom{},
			PaymentFrom:    CoreSQLFrom.FieldsFrom{},
			TakeCreate:     CoreSQLFrom.FieldsFrom{},
			TakeChannel:    CoreSQLFrom.FieldsFrom{},
			TakeFrom:       CoreSQLFrom.FieldsFrom{},
			CreateInfo:     CoreSQLFrom.FieldsFrom{},
			TimeBetween:    CoreSQLTime2.FieldsCoreTime{},
			IsHistory:      false,
			Search:         "",
		})
		if err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		needBreak := false
		for _, v := range dataList {
			if v.CreateAt.Unix() < financeSafeCheckBeforeTime.Unix() {
				needBreak = true
				break
			}
			//是否已经存在异常安全
			if b := checkHaveRecord(CoreSQLFrom.FieldsFrom{System: "log", ID: v.ID}, "LogHash"); b {
				continue
			}
			//开始检查该数据节点
			// 检查本节点的hash是否匹配
			hashStr := fmt.Sprint(v.PayID, v.Key, v.Currency, v.Price, v.PaymentCreate, v.PaymentFrom, v.PaymentChannel, v.TakeCreate, v.TakeFrom, v.TakeChannel, FinanceLog.GetHash())
			hash, err := CoreFilter.GetSha256Str(hashStr)
			if err != nil {
				runLog("get hash, ", err)
			} else {
				if hash != v.Hash {
					if err := createRecord(&argsCreateRecord{
						CreateInfo: CoreSQLFrom.FieldsFrom{
							System: "log",
							ID:     v.ID,
							Mark:   "",
							Name:   "",
						},
						PaymentCreate: v.PaymentCreate,
						PaymentFrom:   v.PaymentFrom,
						TakeCreate:    v.TakeCreate,
						TakeFrom:      v.TakeFrom,
						PayID:         v.PayID,
						PayLogID:      v.ID,
						Message:       "hash错误",
						Code:          "LogHash",
						NeedEW:        false,
					}); err != nil {
						runLog("check log hash, ", err)
					}
				}
			}
			//检查交易数据是否匹配
			if v.PayID < 1 {
				payData, err := FinancePay.GetOne(&FinancePay.ArgsGetOne{
					ID:  v.PayID,
					Key: "",
				})
				if err != nil {
					if b := checkHaveRecord(CoreSQLFrom.FieldsFrom{System: "log", ID: v.ID}, "PayLost"); b {
						continue
					}
					if err := createRecord(&argsCreateRecord{
						CreateInfo: CoreSQLFrom.FieldsFrom{
							System: "log",
							ID:     v.ID,
							Mark:   "",
							Name:   "",
						},
						PaymentCreate: v.PaymentCreate,
						PaymentFrom:   v.PaymentFrom,
						TakeCreate:    v.TakeCreate,
						TakeFrom:      v.TakeFrom,
						PayID:         v.PayID,
						PayLogID:      v.ID,
						Message:       "系统异常，日志存在交易丢失",
						Code:          "PayLost",
						NeedEW:        false,
					}); err != nil {
						runLog("pay not exist, ", err)
					}
				}
				//检查金额是否一致
				if payData.Price != v.Price {
					if b := checkHaveRecord(CoreSQLFrom.FieldsFrom{System: "pay", ID: payData.ID}, "PayPrice"); b {
						continue
					}
					if err := createRecord(&argsCreateRecord{
						CreateInfo: CoreSQLFrom.FieldsFrom{
							System: "log",
							ID:     v.ID,
							Mark:   "",
							Name:   "",
						},
						PaymentCreate: v.PaymentCreate,
						PaymentFrom:   v.PaymentFrom,
						TakeCreate:    v.TakeCreate,
						TakeFrom:      v.TakeFrom,
						PayID:         v.PayID,
						PayLogID:      v.ID,
						Message:       "金额不一致",
						Code:          "PayPrice",
						NeedEW:        false,
					}); err != nil {
						runLog("check pay price, ", err)
					}
				}
			}
			//金额为负数，预警
			if v.Price <= 0 {
				if b := checkHaveRecord(CoreSQLFrom.FieldsFrom{System: "pay", ID: v.PayID}, "PayLimit0"); b {
					continue
				}
				if err := createRecord(&argsCreateRecord{
					CreateInfo: CoreSQLFrom.FieldsFrom{
						System: "log",
						ID:     v.ID,
						Mark:   "",
						Name:   "",
					},
					PaymentCreate: v.PaymentCreate,
					PaymentFrom:   v.PaymentFrom,
					TakeCreate:    v.TakeCreate,
					TakeFrom:      v.TakeFrom,
					PayID:         v.PayID,
					PayLogID:      v.ID,
					Message:       "金额为负数",
					Code:          "PayLimit0",
					NeedEW:        false,
				}); err != nil {
					runLog("pay limit less 0, ", err)
				}
			}
			if financeSafeMaxPrice > 0 {
				if v.Price > financeSafeMaxPrice {
					if b := checkHaveRecord(CoreSQLFrom.FieldsFrom{System: "pay", ID: v.PayID}, "PayLimitMax"); b {
						continue
					}
					if err := createRecord(&argsCreateRecord{
						CreateInfo: CoreSQLFrom.FieldsFrom{
							System: "log",
							ID:     v.ID,
							Mark:   "",
							Name:   "",
						},
						PaymentCreate: v.PaymentCreate,
						PaymentFrom:   v.PaymentFrom,
						TakeCreate:    v.TakeCreate,
						TakeFrom:      v.TakeFrom,
						PayID:         v.PayID,
						PayLogID:      v.ID,
						Message:       "金额超出最大预警额",
						Code:          "PayLimitMax",
						NeedEW:        false,
					}); err != nil {
						runLog("pay limit max, ", err)
					}
				}
			}
			if safeCode := runSafeFrequency(financeSafeFrequencyOneFrom, financeSafeFrequencyOneTo, financeSafeFrequencyAll, financeSafeFrequencyTime, &v); safeCode != "" {
				if b := checkHaveRecord(CoreSQLFrom.FieldsFrom{System: "pay", ID: v.PayID}, safeCode); b {
					continue
				}
				if err := createRecord(&argsCreateRecord{
					CreateInfo: CoreSQLFrom.FieldsFrom{
						System: "log",
						ID:     v.ID,
						Mark:   "",
						Name:   "",
					},
					PaymentCreate: v.PaymentCreate,
					PaymentFrom:   v.PaymentFrom,
					TakeCreate:    v.TakeCreate,
					TakeFrom:      v.TakeFrom,
					PayID:         v.PayID,
					PayLogID:      v.ID,
					Message:       "高频交易",
					Code:          safeCode,
					NeedEW:        false,
				}); err != nil {
					runLog("frequency, ", err)
				}
			}
		}
		if needBreak {
			break
		}
		page += 1
		//强制延迟10毫秒，消峰处理
		time.Sleep(time.Millisecond * 10)
	}
	//全部收尾后，处理预警服务发送问题
	runEW()
}

// 检查高频交易
func runSafeFrequency(financeSafeFrequencyOneFrom, financeSafeFrequencyOneTo, financeSafeFrequencyAll, financeSafeFrequencyTime int64, logData *FinanceLog.FieldsLogType) string {
	//平台总的计数
	var allCount int64
	//检查付款方
	if logData.PaymentCreate.System != "" {
		//获取数据列表，检查付款方的频率计数
		dataList, _, err := FinanceLog.GetList(&FinanceLog.ArgsGetList{
			Pages: CoreSQLPages.ArgsDataList{
				Page: 1,
				Max:  financeSafeFrequencyOneFrom,
				Sort: "_id",
				Desc: true,
			},
			Key:            "",
			Status:         nil,
			PaymentCreate:  logData.PaymentCreate,
			PaymentChannel: CoreSQLFrom.FieldsFrom{},
			PaymentFrom:    CoreSQLFrom.FieldsFrom{},
			TakeCreate:     CoreSQLFrom.FieldsFrom{},
			TakeChannel:    CoreSQLFrom.FieldsFrom{},
			TakeFrom:       CoreSQLFrom.FieldsFrom{},
			CreateInfo:     CoreSQLFrom.FieldsFrom{},
			TimeBetween:    CoreSQLTime2.FieldsCoreTime{},
			IsHistory:      false,
			Search:         "",
		})
		if err != nil {
			return ""
		}
		if len(dataList) < 1 {
			return ""
		}
		//计算第一个和最后一个数据的相差时间，如果少于financeSafeFrequencyTime，则按照高频处理
		if dataList[0].CreateAt.Unix()+financeSafeFrequencyTime <= dataList[len(dataList)-1].CreateAt.Unix() {
			return "PayFrequencyOneFrom"
		}
		//获取平台总的来源个数
		allCount += int64(len(dataList))
	}
	//检查收款方
	if logData.TakeCreate.System != "" {
		//获取数据列表
		dataList, _, err := FinanceLog.GetList(&FinanceLog.ArgsGetList{
			Pages: CoreSQLPages.ArgsDataList{
				Page: 1,
				Max:  financeSafeFrequencyOneTo,
				Sort: "_id",
				Desc: true,
			},
			Key:            "",
			Status:         nil,
			PaymentCreate:  CoreSQLFrom.FieldsFrom{},
			PaymentChannel: CoreSQLFrom.FieldsFrom{},
			PaymentFrom:    CoreSQLFrom.FieldsFrom{},
			TakeCreate:     logData.TakeCreate,
			TakeChannel:    CoreSQLFrom.FieldsFrom{},
			TakeFrom:       CoreSQLFrom.FieldsFrom{},
			CreateInfo:     CoreSQLFrom.FieldsFrom{},
			TimeBetween:    CoreSQLTime2.FieldsCoreTime{},
			IsHistory:      false,
			Search:         "",
		})
		if err != nil {
			return ""
		}
		if len(dataList) < 1 {
			return ""
		}
		//计算第一个和最后一个数据的相差时间，如果少于financeSafeFrequencyTime，则按照高频处理
		if dataList[0].CreateAt.Unix()+financeSafeFrequencyTime <= dataList[len(dataList)-1].CreateAt.Unix() {
			return "PayFrequencyOneTo"
		}
		//获取平台总的来源个数
		allCount += int64(len(dataList))
	}
	//如果总量超出平台最大限制
	if allCount > financeSafeFrequencyAll {
		return "PayFrequencyAll"
	}
	//反馈
	return ""
}

// run的日志
func runLog(content string, err error) {
	if err != nil {
		CoreLog.Error("finance safe run error, ", content, err)
	} else {
		CoreLog.Error("finance safe run error, ", content)
	}
}
