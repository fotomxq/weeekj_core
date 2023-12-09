package UserChat

import (
	"database/sql"
	"errors"
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	FinanceDeposit "gitee.com/weeekj/weeekj_core/v5/finance/deposit"
	FinancePay "gitee.com/weeekj/weeekj_core/v5/finance/pay"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	UserTicket "gitee.com/weeekj/weeekj_core/v5/user/ticket"
	"github.com/jmoiron/sqlx"
	"time"
)

// ArgsGetMessageList 获取消息列表参数
type ArgsGetMessageList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//聊天室ID
	GroupID int64 `db:"group_id" json:"groupID" check:"id"`
	//获取该数据的用户ID
	// -1 跳过
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

type DataGetMessageList struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//参与聊天室
	GroupID int64 `db:"group_id" json:"groupID"`
	//发起用户
	UserID int64 `db:"user_id" json:"userID"`
	//消息类型
	// 0 普通消息; 1 红包领取; 2 优惠券; 3 定位坐标; 4 语音消息; 5 提示信息（发起了视频、语音通话）
	MessageType int `db:"message_type" json:"messageType"`
	//消息内容
	/**
	0 普通消息，存储消息内容
	4 语音消息，存储语音转化为base64文本数据包
	*/
	Message string `db:"message" json:"message"`
	//扩展参数
	/**
	3 定位坐标中，此处将存储坐标系统的address_xy: xy位置、address: 地址详情
	4 语音消息，此处将存储语音转译后的文字信息message_text: 语音消息文本
	*/
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
	//红包消息结构
	MoneyData FieldsMessageMoney `json:"moneyData"`
	//票据消息结构
	TicketData FieldsMessageTicket `json:"ticketData"`
}

// GetMessageList 获取消息列表
func GetMessageList(args *ArgsGetMessageList) (dataList []DataGetMessageList, dataCount int64, err error) {
	//检查用户
	if args.UserID > 0 {
		if !checkChatUser(args.GroupID, args.UserID) {
			err = errors.New("no permission")
			return
		}
	}
	//组合数据
	where := "group_id = :group_id"
	maps := map[string]interface{}{
		"group_id": args.GroupID,
	}
	tableName := "user_chat_message"
	var messageList []FieldsMessage
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&messageList,
		tableName,
		"id",
		"SELECT id, create_at, group_id, user_id, message_type, message, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	//检查结构体，补充数据
	if err == nil && len(messageList) > 0 {
		for k := 0; k < len(messageList); k++ {
			v := messageList[k]
			appendData := DataGetMessageList{
				ID:          v.ID,
				CreateAt:    v.CreateAt,
				GroupID:     v.GroupID,
				UserID:      v.UserID,
				MessageType: v.MessageType,
				Message:     v.Message,
				Params:      v.Params,
				MoneyData:   FieldsMessageMoney{},
				TicketData:  FieldsMessageTicket{},
			}
			switch v.MessageType {
			case 1:
				_ = Router2SystemConfig.MainDB.Get(&appendData.MoneyData, "SELECT id, create_at, user_id, group_id, message_id, config_mark, price, take_type, count_limit, take_list FROM user_chat_message_money WHERE message_id = $1", v.ID)
			case 2:
				_ = Router2SystemConfig.MainDB.Get(&appendData.TicketData, "SELECT id, create_at, user_id, group_id, message_id, config_id, use_count, take_type, count_limit, take_list FROM user_chat_message_ticket WHERE message_id = $1", v.ID)
			}
			dataList = append(dataList, appendData)
		}
	}
	//反馈
	return
}

// ArgsCreateMessage 添加新的消息参数
type ArgsCreateMessage struct {
	//参与聊天室
	GroupID int64 `db:"group_id" json:"groupID" check:"id"`
	//发起用户
	UserID int64 `db:"user_id" json:"userID"`
	//消息类型
	// 0 普通消息; 1 红包领取; 2 优惠券; 3 定位坐标; 4 语音消息; 5 提示信息（发起了视频、语音通话）
	MessageType int `db:"message_type" json:"messageType"`
	//消息内容
	/**
	0 普通消息，存储消息内容
	4 语音消息，存储语音转化为base64文本数据包
	*/
	Message string `db:"message" json:"message" check:"des" min:"1" max:"600"`
	//扩展参数
	/**
	3 定位坐标中，此处将存储坐标系统的address_xy: xy位置、address: 地址详情
	4 语音消息，此处将存储语音转译后的文字信息message_text: 语音消息文本
	*/
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
	//红包数据集合
	MoneyData ArgsCreateMessageMoney `json:"moneyData"`
	//票据消息结构
	TicketData ArgsCreateMessageTicket `json:"ticketData"`
}
type ArgsCreateMessageMoney struct {
	//储蓄配置
	ConfigMark string `db:"config_mark" json:"configMark" check:"mark" empty:"true"`
	//金额
	Price int64 `db:"price" json:"price"`
	//发放办法
	// 0 全部发放给领取的第一个人; 1 随机发放(0.01 - 总金额-尚未领取人数x0.01)
	TakeType int `db:"take_type" json:"takeType"`
	//领取人数限制
	CountLimit int `db:"count_limit" json:"countLimit"`
}
type ArgsCreateMessageTicket struct {
	//票据配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//发放张数
	UseCount int64 `db:"use_count" json:"useCount"`
	//发放办法
	// 0 全部发放给领取的第一个人; 1 随机发放(1张 - 总张数-尚未领取人数x1张)
	TakeType int `db:"take_type" json:"takeType"`
	//领取人数限制
	CountLimit int `db:"count_limit" json:"countLimit"`
}

// CreateMessage 添加新的消息
func CreateMessage(args *ArgsCreateMessage) (errCode string, err error) {
	//检查用户必须在聊天室
	if !checkChatUser(args.GroupID, args.UserID) {
		errCode = "no_permission"
		err = errors.New("no permission")
		return
	}
	//获取房间的商户信息
	var orgID int64
	err = Router2SystemConfig.MainDB.Get(&orgID, "SELECT org_id FROM user_chat_group WHERE id = $1", args.GroupID)
	if err != nil {
		errCode = "group_not_exist"
		return
	}
	//创建数据
	defer func() {
		if e := recover(); e != nil {
			errCode = "system"
			err = errors.New(fmt.Sprint(e))
			return
		}
	}()
	tx := Router2SystemConfig.MainDB.MustBegin()
	var stmt *sqlx.NamedStmt
	stmt, err = tx.PrepareNamed("INSERT INTO user_chat_message (group_id, user_id, message_type, message, params) VALUES (:group_id,:user_id,:message_type,:message,:params) RETURNING id;")
	if err != nil {
		errCode = "insert_message"
		if err2 := tx.Rollback(); err2 != nil {
			err = errors.New(fmt.Sprint("create message failed, message create, ", err, ", rollback failed, ", err2))
			return
		}
		return
	}
	var messageID int64
	messageID, err = CoreSQL.LastRowsAffectedCreate(tx, stmt, args, err)
	if err != nil {
		errCode = "insert_message_id"
		if err2 := tx.Rollback(); err2 != nil {
			err = errors.New(fmt.Sprint("create message failed, message id lost, ", err, ", rollback failed, ", err2))
			return
		}
		return
	}
	//检查消息类型
	switch args.MessageType {
	case 1:
		//红包
		//储蓄存储来源
		fromInfo := CoreSQLFrom.FieldsFrom{}
		var userChatMoneyFromOrg bool
		userChatMoneyFromOrg, err = BaseConfig.GetDataBool("UserChatMoneyFromOrg")
		if err != nil {
			userChatMoneyFromOrg = false
		}
		if userChatMoneyFromOrg {
			fromInfo = CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     orgID,
				Mark:   "",
				Name:   "",
			}
		}
		//检查用户是否具有金额
		depositPrice := FinanceDeposit.GetPriceByFrom(&FinanceDeposit.ArgsGetByFrom{
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     args.UserID,
				Mark:   "",
				Name:   "",
			},
			FromInfo:   fromInfo,
			ConfigMark: args.MoneyData.ConfigMark,
		})
		if depositPrice < args.MoneyData.Price {
			errCode = "no_more_money"
			err = errors.New(fmt.Sprint("no more money, user have config: ", args.MoneyData.ConfigMark, ", price: ", depositPrice, ", need deposit price: ", args.MoneyData.Price))
			if err2 := tx.Rollback(); err2 != nil {
				//errCode = "insert"
				err = errors.New(fmt.Sprint("no more money by create message failed, ", err, ", rollback failed, ", err2))
				return
			}
			return
		}
		//发起红包处理
		var resultMoney sql.Result
		resultMoney, err = tx.NamedExec("INSERT INTO user_chat_message_money (user_id, group_id, message_id, config_mark, price, take_type, count_limit, take_list) VALUES (:user_id,:group_id,:message_id,:config_mark,:price,:take_type,:count_limit,:take_list)", map[string]interface{}{
			"user_id":     args.UserID,
			"group_id":    args.GroupID,
			"message_id":  messageID,
			"config_mark": args.MoneyData.ConfigMark,
			"price":       args.MoneyData.Price,
			"take_type":   args.MoneyData.TakeType,
			"count_limit": args.MoneyData.CountLimit,
			"take_list":   FieldsMessageMoneyTakeList{},
		})
		if err != nil {
			errCode = "insert_money"
			err = errors.New("no more money")
			if err2 := tx.Rollback(); err2 != nil {
				err = errors.New(fmt.Sprint("create message money failed, message money id lost, ", err, ", rollback failed, ", err2))
				return
			}
			return
		}
		err = CoreSQL.LastRowsAffected(tx, resultMoney, err)
		if err != nil {
			return
		}
	case 2:
		//优惠券
		//检查用户票据是否足够？
		//TODO: 票据接近过期后，可能会发起转让，出现过期时间被重置的问题
		ticketCount, _ := UserTicket.GetTicketCount(&UserTicket.ArgsGetTicketCount{
			ConfigID: args.TicketData.ConfigID,
			UserID:   args.UserID,
		})
		if ticketCount < args.TicketData.UseCount {
			errCode = "no_more_ticket"
			err = errors.New("no more ticket")
			if err2 := tx.Rollback(); err2 != nil {
				//errCode = "insert"
				err = errors.New(fmt.Sprint("no more ticket by create message failed, ", err, ", rollback failed, ", err2))
				return
			}
			return
		}
		//发起优惠券处理
		var resultTicket sql.Result
		resultTicket, err = tx.NamedExec("INSERT INTO user_chat_message_ticket (user_id, group_id, message_id, config_id, use_count, take_type, count_limit, take_list) VALUES (:user_id,:group_id,:message_id,:config_id,:use_count,:take_type,:count_limit,:take_list)", map[string]interface{}{
			"user_id":     args.UserID,
			"group_id":    args.GroupID,
			"message_id":  messageID,
			"config_id":   args.TicketData.ConfigID,
			"use_count":   args.TicketData.UseCount,
			"take_type":   args.TicketData.TakeType,
			"count_limit": args.TicketData.CountLimit,
			"take_list":   FieldsMessageTicketTakeList{},
		})
		if err != nil {
			errCode = "insert_ticket"
			err = errors.New("no more ticket")
			if err2 := tx.Rollback(); err2 != nil {
				err = errors.New(fmt.Sprint("create message ticket failed, message ticket id lost, ", err, ", rollback failed, ", err2))
				return
			}
			return
		}
		err = CoreSQL.LastRowsAffected(tx, resultTicket, err)
		if err != nil {
			return
		}
	}
	//执行事务关系
	err = tx.Commit()
	//更新访问时间
	if err == nil {
		updateLastAtByChat(args.GroupID, args.UserID)
		pushUpdateByGroup(args.GroupID, 2)
	}
	//反馈
	return
}

// ArgsTakeMessageMoneyOrTicket 领取消息红包或优惠券参数
type ArgsTakeMessageMoneyOrTicket struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//参与聊天室
	GroupID int64 `db:"group_id" json:"groupID" check:"id"`
	//消息ID
	MessageID int64 `db:"message_id" json:"messageID" check:"id"`
}

// TakeMessageMoneyOrTicket 领取消息红包或优惠券
func TakeMessageMoneyOrTicket(args *ArgsTakeMessageMoneyOrTicket) (priceOrCount int64, errCode string, err error) {
	//锁定机制
	takeMoneyOrTicketLock.Lock()
	defer takeMoneyOrTicketLock.Unlock()
	//检查用户必须在聊天室
	if !checkChatUser(args.GroupID, args.UserID) {
		errCode = "no_permission"
		err = errors.New("no permission")
		return
	}
	//获取消息
	var messageData FieldsMessage
	err = Router2SystemConfig.MainDB.Get(&messageData, "SELECT id, user_id, message_type FROM user_chat_message WHERE id = $1 AND group_id = $2", args.MessageID, args.GroupID)
	if err != nil || messageData.ID < 1 {
		errCode = "no_message"
		err = errors.New("no message")
		return
	}
	//获取房间的商户信息
	var orgID int64
	err = Router2SystemConfig.MainDB.Get(&orgID, "SELECT org_id FROM user_chat_group WHERE id = $1", args.GroupID)
	if err != nil {
		errCode = "group_not_exist"
		return
	}
	//检查类型
	switch messageData.MessageType {
	case 1:
		//红包
		//检查是否领取过？
		// 同时获取红包发送渠道
		var moneyData FieldsMessageMoney
		err = Router2SystemConfig.MainDB.Get(&moneyData, "SELECT id, config_mark, price, take_type, count_limit, take_list FROM user_chat_message_money WHERE message_id = $1", args.MessageID)
		if err != nil || messageData.ID < 1 {
			errCode = "no_money_message"
			err = errors.New("no message money")
			return
		}
		var havePrice int64
		for _, v := range moneyData.TakeList {
			if v.UserID == args.UserID {
				errCode = "have_take"
				err = errors.New("user have take money")
				return
			}
			havePrice += v.Price
		}
		if len(moneyData.TakeList) >= moneyData.CountLimit {
			errCode = "count_limit"
			err = errors.New("take count more than limit")
			return
		}
		var price int64
		switch moneyData.TakeType {
		case 0:
			price = moneyData.Price
		case 1:
			price = CoreFilter.GetMaxRand(moneyData.Price, havePrice, int64(len(moneyData.TakeList)), int64(moneyData.CountLimit))
		}
		if price < 1 {
			errCode = "price_0"
			err = errors.New("price less 1")
			return
		}
		//储蓄存储来源
		fromInfo := CoreSQLFrom.FieldsFrom{}
		var userChatMoneyFromOrg bool
		userChatMoneyFromOrg, err = BaseConfig.GetDataBool("UserChatMoneyFromOrg")
		if err != nil {
			userChatMoneyFromOrg = false
		}
		if userChatMoneyFromOrg {
			fromInfo = CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     orgID,
				Mark:   "",
				Name:   "",
			}
		}
		//根据设定抽取金额
		//发起从消息发起人，给目标用户打款
		var payData FinancePay.FieldsPayType
		payData, _, err = FinancePay.CreateQuickPay(&FinancePay.ArgsCreate{
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: "user_chat",
				ID:     messageData.ID,
				Mark:   "",
				Name:   "聊天红包",
			},
			PaymentCreate: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     messageData.UserID,
				Mark:   "",
				Name:   "",
			},
			PaymentChannel: CoreSQLFrom.FieldsFrom{
				System: "deposit",
				ID:     0,
				Mark:   moneyData.ConfigMark,
				Name:   "",
			},
			PaymentFrom: fromInfo,
			TakeCreate: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     args.UserID,
				Mark:   "",
				Name:   "",
			},
			TakeChannel: CoreSQLFrom.FieldsFrom{
				System: "deposit",
				ID:     0,
				Mark:   moneyData.ConfigMark,
				Name:   "",
			},
			TakeFrom: fromInfo,
			Des:      "消息红包",
			ExpireAt: CoreFilter.GetNowTimeCarbon().AddMinutes(30).Time,
			Currency: 86,
			Price:    price,
			Params:   nil,
		})
		if err != nil {
			errCode = "create_pay"
			return
		}
		//添加到已经付款名单
		moneyData.TakeList = append(moneyData.TakeList, FieldsMessageMoneyTake{
			CreateAt: CoreFilter.GetNowTime(),
			UserID:   args.UserID,
			PayID:    payData.ID,
			Price:    price,
		})
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_chat_message_money SET take_list = :take_list WHERE id = :id", map[string]interface{}{
			"id":        moneyData.ID,
			"take_list": moneyData.TakeList,
		})
		if err != nil {
			errCode = "update"
			return
		}
		priceOrCount = price
		pushUpdateByGroup(args.GroupID, 2)
		return
	case 2:
		//票据
		//检查是否领取过？
		// 同时获取发送渠道
		var ticketData FieldsMessageTicket
		err = Router2SystemConfig.MainDB.Get(&ticketData, "SELECT id, config_id, use_count, take_type, count_limit, take_list FROM user_chat_message_ticket WHERE message_id = $1", args.MessageID)
		if err != nil || messageData.ID < 1 {
			errCode = "no_ticket_message"
			err = errors.New("no message ticket")
			return
		}
		var haveCount int64
		for _, v := range ticketData.TakeList {
			if v.UserID == args.UserID {
				errCode = "have_take"
				err = errors.New("user have take ticket")
				return
			}
			haveCount += v.GetCount
		}
		if len(ticketData.TakeList) >= ticketData.CountLimit {
			errCode = "count_limit"
			err = errors.New("take count more than limit")
			return
		}
		var takeCount int64
		switch ticketData.TakeType {
		case 0:
			takeCount = ticketData.UseCount
		case 1:
			takeCount = CoreFilter.GetMaxRand(int64(ticketData.UseCount), haveCount, int64(len(ticketData.TakeList)), int64(ticketData.CountLimit))
		}
		if takeCount < 1 {
			errCode = "ticket_count_0"
			err = errors.New("ticket count less 1")
			return
		}
		//扣除发送方票据
		err = UserTicket.UseTicket(&UserTicket.ArgsUseTicket{
			ID:          0,
			OrgID:       orgID,
			ConfigID:    ticketData.ConfigID,
			UserID:      messageData.UserID,
			Count:       takeCount,
			UseFromName: "消息赠送前扣除",
		})
		if err != nil {
			errCode = "ues_ticket"
			err = errors.New(fmt.Sprint("use user ticket failed, ", err))
			return
		}
		//给与用户票据
		_, _, err = UserTicket.AddTickets(&UserTicket.ArgsAddTickets{
			OrgID:  orgID,
			UserID: args.UserID,
			Data: []UserTicket.ArgsAddTicketsChild{
				{
					ConfigID:    ticketData.ConfigID,
					Count:       takeCount,
					UseFromName: "消息赠送",
				},
			},
			CanRefundConfigIDs: []int64{},
		})
		if err != nil {
			errCode = "add_ticket"
			err = errors.New(fmt.Sprint("add user ticket failed, ", err))
			return
		}
		//添加到已经付款名单
		ticketData.TakeList = append(ticketData.TakeList, FieldsMessageTicketTake{
			CreateAt: CoreFilter.GetNowTime(),
			UserID:   args.UserID,
			GetCount: takeCount,
		})
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_chat_message_ticket SET take_list = :take_list WHERE id = :id", map[string]interface{}{
			"id":        ticketData.ID,
			"take_list": ticketData.TakeList,
		})
		if err != nil {
			errCode = "update"
			return
		}
		priceOrCount = takeCount
		pushUpdateByGroup(args.GroupID, 2)
		return
	default:
		errCode = "not_support"
		err = errors.New(fmt.Sprint("message type not support, message id: ", messageData.ID, ", type: ", messageData.MessageType))
		return
	}
}
