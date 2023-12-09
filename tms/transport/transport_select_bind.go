package TMSTransport

import (
	"fmt"
	ClassConfig "gitee.com/weeekj/weeekj_core/v5/class/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	OrgCoreCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"strings"
)

// 为配送单选择配送员逻辑
func transportSelectBind(isRefund bool, orgID int64, userID int64, bindID int64, goods FieldsTransportGoods, fromAddress, toAddress CoreSQLAddress.FieldsAddress, params CoreSQLConfig.FieldsConfigsType) (newWaitLog []argsAppendLog, newParams CoreSQLConfig.FieldsConfigsType, newBindID int64) {
	//初始化数据
	var err error
	var waitLog []argsAppendLog
	var newLog argsAppendLog
	//检查商品是否在绑定关系中？
	for _, v := range goods {
		if v.System != "mall" {
			continue
		}
		var bindToMallBindID int64
		err = Router2SystemConfig.MainDB.Get(&bindToMallBindID, "SELECT bind_id FROM tms_transport_bind_to_mall WHERE org_id = $1 AND bind_mall_id = $2 AND delete_at < to_timestamp(1000000)", orgID, v.ID)
		if err != nil || bindToMallBindID < 1 {
			continue
		}
		//找到后退出
		if bindToMallBindID > 0 {
			newBindID = bindToMallBindID
			break
		}
	}
	if newBindID > 0 {
		waitLog = append(waitLog, argsAppendLog{
			OrgID:           orgID,
			BindID:          newBindID,
			TransportID:     0,
			TransportBindID: 0,
			Mark:            "bind_to_mall",
			Des:             fmt.Sprint("由于设置了商品和配送员绑定，该配送单强制分配给该配送员，配送人员ID[", newBindID, "]"),
		})
		return
	} else {
		//根据锁定配置，来确定配送员
		newBindID, newLog = transportSelectBindGlobConfig(orgID, userID, goods)
		//如果没有锁定配置，则根据分区来识别处理
		if newBindID < 1 {
			//初始化数据
			var bindData FieldsBind
			var areaID int64
			var errCode string
			//退单和非退单，分配方式存在差异
			if isRefund {
				bindData, areaID, errCode, err = getBindToTransport(&argsGetBindToTransport{
					OrgID:     orgID,
					MapType:   fromAddress.MapType,
					Longitude: fromAddress.Longitude,
					Latitude:  fromAddress.Latitude,
				})
			} else {
				bindData, areaID, errCode, err = getBindToTransport(&argsGetBindToTransport{
					OrgID:     orgID,
					MapType:   toAddress.MapType,
					Longitude: toAddress.Longitude,
					Latitude:  toAddress.Latitude,
				})
			}
			if err != nil {
				bindData = FieldsBind{}
				bindData.ID = 0
				if areaID > 0 {
					waitLog = append(waitLog, argsAppendLog{
						OrgID:           orgID,
						BindID:          bindID,
						TransportID:     0,
						TransportBindID: 0,
						Mark:            "no_bind",
						Des:             fmt.Sprint("无法为分配配送人员，因为所在分区没有配送人员，错误代码[", errCode, "]"),
					})
					params = append(params, CoreSQLConfig.FieldsConfigType{
						Mark: "autoMapAreaID",
						Val:  fmt.Sprint(areaID),
					})
				} else {
					waitLog = append(waitLog, argsAppendLog{
						OrgID:           orgID,
						BindID:          bindID,
						TransportID:     0,
						TransportBindID: 0,
						Mark:            "no_bind",
						Des:             fmt.Sprint("无法为分配配送人员，因为机构没有配送人员，错误代码[", errCode, "]"),
					})
				}
			} else {
				newBindID = bindData.BindID
				waitLog = append(waitLog, argsAppendLog{
					OrgID:           orgID,
					BindID:          bindID,
					TransportID:     0,
					TransportBindID: 0,
					Mark:            "auto_select",
					Des:             fmt.Sprint("根据分区，自动分配配送人员[", newBindID, "]"),
				})
			}
		} else {
			waitLog = append(waitLog, newLog)
		}
	}
	//覆盖值
	newWaitLog = waitLog
	newParams = params
	//反馈
	return
}

// 选择配送员的子模块
// 根据全局一些锁定设置，来找到合适配送员
func transportSelectBindGlobConfig(orgID int64, userID int64, goods FieldsTransportGoods) (newBindID int64, newLog argsAppendLog) {
	//获取商户配置
	tmsTransportLockBind, err := OrgCoreCore.Config.GetConfigValBool(&ClassConfig.ArgsGetConfig{
		BindID:    orgID,
		Mark:      "TMSTransportLockBind",
		VisitType: "admin",
	})
	if err != nil {
		tmsTransportLockBind = false
		err = nil
	}
	//如果锁定配送员，则查询该用户最后一次配送记录
	if tmsTransportLockBind {
		var lastData FieldsTransport
		err = Router2SystemConfig.MainDB.Get(&lastData, "SELECT id, bind_id FROM tms_transport WHERE delete_at < to_timestamp(1000000) AND org_id = $1 AND user_id = $2 ORDER BY id DESC LIMIT 1", orgID, userID)
		if err == nil && lastData.ID > 0 {
			newBindID = lastData.BindID
			newLog = argsAppendLog{
				OrgID:           orgID,
				BindID:          newBindID,
				TransportID:     0,
				TransportBindID: 0,
				Mark:            "lock",
				Des:             fmt.Sprint("由于启动了锁定机制，根据上一次配送人员自动分配配送，配送人员ID[", newBindID, "]"),
			}
		}
	} else {
		//货物超出锁定机制
		var tmsGoodMoreLock string
		tmsGoodMoreLock, err = OrgCoreCore.Config.GetConfigVal(&ClassConfig.ArgsGetConfig{
			BindID:    orgID,
			Mark:      "TMSGoodMoreLock",
			VisitType: "admin",
		})
		if err != nil {
			tmsGoodMoreLock = ""
			err = nil
		}
		if tmsGoodMoreLock != "" {
			//分解配置
			var tmsGoodMoreLocks []string
			tmsGoodMoreLocks = strings.Split(tmsGoodMoreLock, "|")
			if len(tmsGoodMoreLocks) > 0 {
				for _, vMoreLock := range tmsGoodMoreLocks {
					if vMoreLock == "" {
						continue
					}
					//二次拆分
					// 商品系统来源,商品ID,超出购买个数,指定配送成员ID
					var tmsGoodMoreLocks2 []string
					tmsGoodMoreLocks2 = strings.Split(vMoreLock, ",")
					if len(tmsGoodMoreLocks2) != 4 {
						continue
					}
					var tmsGoodMoreLocks2ID int64
					tmsGoodMoreLocks2ID, err = CoreFilter.GetInt64ByString(tmsGoodMoreLocks2[1])
					if err != nil {
						continue
					}
					var tmsGoodMoreLocks2Count int64
					tmsGoodMoreLocks2Count, err = CoreFilter.GetInt64ByString(tmsGoodMoreLocks2[2])
					if err != nil {
						continue
					}
					var tmsGoodMoreLocks2BindID int64
					tmsGoodMoreLocks2BindID, err = CoreFilter.GetInt64ByString(tmsGoodMoreLocks2[3])
					if err != nil {
						continue
					}
					if tmsGoodMoreLocks2ID < 1 || tmsGoodMoreLocks2Count < 1 || tmsGoodMoreLocks2BindID < 1 {
						continue
					}
					//遍历货物
					isFind := false
					for _, vGood := range goods {
						if vGood.System == tmsGoodMoreLocks2[0] && vGood.ID == tmsGoodMoreLocks2ID && int64(vGood.Count) >= tmsGoodMoreLocks2Count {
							isFind = true
							break
						}
					}
					if !isFind {
						continue
					}
					newBindID = tmsGoodMoreLocks2BindID
					newLog = argsAppendLog{
						OrgID:           orgID,
						BindID:          newBindID,
						TransportID:     0,
						TransportBindID: 0,
						Mark:            "more_lock",
						Des:             fmt.Sprint("由于启动了超出货物锁定机制，根据系统配置强制分配配送人员，配送人员ID[", newBindID, "]"),
					}
				}
			}
		}
	}
	//反馈
	return
}
