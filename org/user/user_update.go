package OrgUser

import (
	"errors"
	"fmt"
	BaseExpireTip "github.com/fotomxq/weeekj_core/v5/base/expire_tip"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime2 "github.com/fotomxq/weeekj_core/v5/core/sql/time2"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ServiceOrder "github.com/fotomxq/weeekj_core/v5/service/order"
	UserAddress "github.com/fotomxq/weeekj_core/v5/user/address"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	UserIntegral "github.com/fotomxq/weeekj_core/v5/user/integral"
	UserSubscription "github.com/fotomxq/weeekj_core/v5/user/subscription"
	UserTicket "github.com/fotomxq/weeekj_core/v5/user/ticket"
)

// ArgsUpdateUserData 强制更新指定用户的数据参数
type ArgsUpdateUserData struct {
	//组织ID
	OrgID int64 `json:"orgID" check:"id"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//用户地址ID
	UserAddressID int64 `json:"userAddressID" check:"id" empty:"true"`
}

// UpdateUserData 强制更新指定用户的数据
func UpdateUserData(args *ArgsUpdateUserData) (data FieldsOrgUser, err error) {
	//如果用户ID存在，则优先走用户ID，否则走用户地址ID
	if args.UserID < 1 {
		var addressData UserAddress.FieldsAddress
		addressData, err = UserAddress.GetID(&UserAddress.ArgsGetID{
			ID:       args.UserAddressID,
			UserID:   -1,
			IsRemove: false,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("not find user address, ", err))
			return
		}
		if addressData.UserID < 1 {
			err = errors.New("user address user id less 1")
			return
		}
		args.UserID = addressData.UserID
	}
	//更新数据
	err = updateByUserID(args.OrgID, args.UserID)
	if err != nil {
		return
	}
	//获取数据
	data = getUserData(args.OrgID, args.UserID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//反馈
	return
}

// 更新用户的聚合数据
func updateByUserID(orgID int64, userID int64) (err error) {
	//查询该用户数据包
	var userData UserCore.FieldsUserType
	userData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
		ID:    userID,
		OrgID: -1,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("user data not exist, ", err))
		return
	}
	//查询用户10条地址
	var addressList []UserAddress.FieldsAddress
	addressList, _ = UserAddress.GetAddressByUserID(userData.ID, 10)
	var newAddressList FieldsOrgUserAddress
	for _, v := range addressList {
		newAddressList = append(newAddressList, FieldsAddress{
			ID:         v.ID,
			UpdateAt:   v.UpdateAt,
			Country:    v.Country,
			Province:   v.Province,
			City:       v.City,
			Address:    v.Address,
			MapType:    v.MapType,
			Longitude:  v.Longitude,
			Latitude:   v.Latitude,
			Name:       v.Name,
			NationCode: v.NationCode,
			Phone:      v.Phone,
		})
	}
	//查询用户积分
	var userIntegral int64
	userIntegral = UserIntegral.GetUserCount(orgID, userData.ID)
	//查询用户订阅
	var userSubs []UserSubscription.FieldsSub
	userSubs, _, _ = UserSubscription.GetSubList(&UserSubscription.ArgsGetSubList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  100,
			Sort: "expire_at",
			Desc: true,
		},
		OrgID:        orgID,
		ConfigID:     -1,
		UserID:       userData.ID,
		NeedIsExpire: true,
		IsExpire:     false,
		IsRemove:     false,
		Search:       "",
	})
	newUserSubs := FieldsOrgUserSubs{}
	for _, v := range userSubs {
		if v.ConfigID < 1 {
			continue
		}
		if v.ExpireAt.Unix() < CoreFilter.GetNowTime().Unix() {
			continue
		}
		newUserSubs = append(newUserSubs, FieldsOrgUserSub{
			ConfigID: v.ConfigID,
			ExpireAt: v.ExpireAt,
		})
	}
	//查询用户票据
	var userTickets []UserTicket.FieldsTicket
	userTickets, _, _ = UserTicket.GetTicketList(&UserTicket.ArgsGetTicketList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  100,
			Sort: "expire_at",
			Desc: true,
		},
		OrgID:        orgID,
		ConfigID:     -1,
		UserID:       userData.ID,
		NeedIsExpire: true,
		IsExpire:     false,
		IsRemove:     false,
		Search:       "",
	})
	newUserTickets := FieldsOrgUserTickets{}
	for _, v := range userTickets {
		if v.ConfigID < 1 {
			continue
		}
		if v.ExpireAt.Unix() < CoreFilter.GetNowTime().Unix() {
			continue
		}
		newUserTickets = append(newUserTickets, FieldsOrgUserTicket{
			ConfigID: v.ConfigID,
			Count:    v.Count,
			ExpireAt: v.ExpireAt,
		})
	}
	//获取储蓄账户数据
	var depositList []FinanceDeposit.FieldsDepositType
	depositList, _, _ = FinanceDeposit.GetList(&FinanceDeposit.ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  100,
			Sort: "id",
			Desc: false,
		},
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     userData.ID,
			Mark:   "",
			Name:   "",
		},
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     orgID,
			Mark:   "",
			Name:   "",
		},
		ConfigMark: "",
		MinPrice:   0,
		MaxPrice:   0,
	})
	newDepositData := FieldsOrgUserDeposits{}
	for _, v := range depositList {
		newDepositData = append(newDepositData, FieldsOrgUserDeposit{
			Mark:  v.ConfigMark,
			Price: v.SavePrice,
		})
	}
	//获取最后一次订单
	var lastOrderList []ServiceOrder.FieldsOrder
	lastOrderList, _, _ = ServiceOrder.GetList(&ServiceOrder.ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  1,
			Sort: "id",
			Desc: true,
		},
		SystemMark:             "",
		OrgID:                  orgID,
		UserID:                 userData.ID,
		CompanyID:              -1,
		CreateFrom:             -1,
		Status:                 []int{3, 4},
		RefundStatus:           []int{},
		TransportID:            -1,
		NeedTransportAllowAuto: false,
		TransportAllowAuto:     false,
		PayStatus:              []int{},
		PayID:                  -1,
		PayFrom:                "",
		GoodFrom:               CoreSQLFrom.FieldsFrom{},
		TimeBetween:            CoreSQLTime2.DataCoreTime{},
		IsRemove:               false,
		IsHistory:              false,
		Search:                 "",
	})
	var lastOrder FieldsOrgUserOrder
	if len(lastOrderList) > 0 {
		lastOrder = FieldsOrgUserOrder{
			ID:                 lastOrderList[0].ID,
			CreateAt:           lastOrderList[0].CreateAt,
			UpdateAt:           lastOrderList[0].UpdateAt,
			DeleteAt:           lastOrderList[0].DeleteAt,
			ExpireAt:           lastOrderList[0].ExpireAt,
			SystemMark:         lastOrderList[0].SystemMark,
			OrgID:              lastOrderList[0].OrgID,
			UserID:             lastOrderList[0].UserID,
			CreateFrom:         lastOrderList[0].CreateFrom,
			SerialNumber:       lastOrderList[0].SerialNumber,
			SerialNumberDay:    lastOrderList[0].SerialNumberDay,
			Status:             lastOrderList[0].Status,
			RefundStatus:       lastOrderList[0].RefundStatus,
			AddressFrom:        lastOrderList[0].AddressFrom,
			AddressTo:          lastOrderList[0].AddressTo,
			Goods:              lastOrderList[0].Goods,
			Exemptions:         lastOrderList[0].Exemptions,
			AllowAutoAudit:     lastOrderList[0].AllowAutoAudit,
			TransportID:        lastOrderList[0].TransportID,
			TransportAllowAuto: lastOrderList[0].TransportAllowAuto,
			TransportTaskAt:    lastOrderList[0].TransportTaskAt,
			TransportPayAfter:  lastOrderList[0].TransportPayAfter,
			TransportIDs:       lastOrderList[0].TransportIDs,
			PriceList:          lastOrderList[0].PriceList,
			PricePay:           lastOrderList[0].PricePay,
			Currency:           lastOrderList[0].Currency,
			Price:              lastOrderList[0].Price,
			PriceTotal:         lastOrderList[0].PriceTotal,
			PayStatus:          lastOrderList[0].PayStatus,
			PayID:              lastOrderList[0].PayID,
			PayList:            lastOrderList[0].PayList,
			Des:                lastOrderList[0].Des,
			Logs:               lastOrderList[0].Logs,
			Params:             lastOrderList[0].Params,
		}
	}
	//找到数据集合
	dataID := getUserDataIDByUserID(orgID, userID)
	if dataID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_user_data SET update_at = NOW(), name = :name, phone = :phone, address_list = :address_list, user_integral = :user_integral, user_subs = :user_subs, user_tickets = :user_tickets, deposit_data = :deposit_data, last_order = :last_order, params = :params WHERE id = :id", map[string]interface{}{
			"id":            dataID,
			"name":          userData.Name,
			"phone":         userData.Phone,
			"address_list":  newAddressList,
			"user_integral": userIntegral,
			"user_subs":     newUserSubs,
			"user_tickets":  newUserTickets,
			"deposit_data":  newDepositData,
			"last_order":    lastOrder,
			"params":        CoreSQLConfig.FieldsConfigsType{},
		})
		if err != nil {
			return
		}
	} else {
		//创建新的数据
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_user_data (org_id, user_id, name, phone, address_list, user_integral, user_subs, user_tickets, deposit_data, last_order, params) VALUES (:org_id,:user_id,:name,:phone,:address_list,:user_integral,:user_subs,:user_tickets,:deposit_data,:last_order,:params)", map[string]interface{}{
			"org_id":        orgID,
			"user_id":       userData.ID,
			"name":          userData.Name,
			"phone":         userData.Phone,
			"address_list":  newAddressList,
			"user_integral": userIntegral,
			"user_subs":     newUserSubs,
			"user_tickets":  newUserTickets,
			"deposit_data":  newDepositData,
			"last_order":    lastOrder,
			"params":        CoreSQLConfig.FieldsConfigsType{},
		})
		if err != nil {
			return
		}
	}
	//清理缓冲
	deleteUserCache(orgID, userID)
	//请求过期删除
	// 删除超过60天数据
	_ = BaseExpireTip.AppendTip(&BaseExpireTip.ArgsAppendTip{
		OrgID:      orgID,
		UserID:     userID,
		SystemMark: "org_user_data",
		BindID:     0,
		Hash:       "",
		ExpireAt:   CoreFilter.GetNowTimeCarbon().AddDays(60).Time,
	})
	//反馈
	return
}
