package UserChat

import (
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	UserTicket "github.com/fotomxq/weeekj_core/v5/user/ticket"
	"testing"
)

var (
	newMessageData DataGetMessageList
)

func TestInitMessage(t *testing.T) {
	TestInitChat(t)
	TestInviteUser(t)
}

func TestCreateMessage(t *testing.T) {
	errCode, err := CreateMessage(&ArgsCreateMessage{
		GroupID:     newGroupData.ID,
		UserID:      newUserData.ID,
		MessageType: 0,
		Message:     "普通消息",
		Params:      nil,
		MoneyData:   ArgsCreateMessageMoney{},
		TicketData:  ArgsCreateMessageTicket{},
	})
	ToolsTest.ReportErrorCode(t, errCode, err)
}

func TestGetMessageList(t *testing.T) {
	dataList, dataCount, err := GetMessageList(&ArgsGetMessageList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		GroupID: newGroupData.ID,
		UserID:  -1,
		Search:  "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil && len(dataList) > 0 {
		newMessageData = dataList[0]
	}
}

func TestTakeMessageMoneyOrTicket_Money(t *testing.T) {
	_, errCode, err := FinanceDeposit.SetByFrom(&FinanceDeposit.ArgsSetByFrom{
		UpdateHash: "",
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     newUserData.ID,
			Mark:   "",
			Name:   "",
		},
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     newGroupData.OrgID,
			Mark:   "",
			Name:   "",
		},
		ConfigMark:      "deposit",
		AppendSavePrice: 10,
	})
	ToolsTest.ReportErrorCode(t, errCode, err)
	errCode, err = CreateMessage(&ArgsCreateMessage{
		GroupID:     newGroupData.ID,
		UserID:      newUserData.ID,
		MessageType: 1,
		Message:     "红包消息",
		Params:      nil,
		MoneyData: ArgsCreateMessageMoney{
			ConfigMark: "deposit",
			Price:      3,
			TakeType:   0,
			CountLimit: 1,
		},
		TicketData: ArgsCreateMessageTicket{},
	})
	ToolsTest.ReportErrorCode(t, errCode, err)
	TestGetMessageList(t)
	priceOrCount, errCode, err := TakeMessageMoneyOrTicket(&ArgsTakeMessageMoneyOrTicket{
		UserID:    newGroupData.UserID,
		GroupID:   newGroupData.ID,
		MessageID: newMessageData.ID,
	})
	ToolsTest.ReportErrorCode(t, errCode, err)
	if err == nil {
		depositPrice := FinanceDeposit.GetPriceByFrom(&FinanceDeposit.ArgsGetByFrom{
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     newGroupData.UserID,
				Mark:   "",
				Name:   "",
			},
			FromInfo: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     newGroupData.OrgID,
				Mark:   "",
				Name:   "",
			},
			ConfigMark: "deposit",
		})
		t.Log("priceOrCount: ", priceOrCount, ", depositPrice: ", depositPrice)
		if depositPrice != 3 {
			t.Error("depositPrice not 3")
		}
	}
}

func TestTakeMessageMoneyOrTicket_Ticket(t *testing.T) {
	//设置用户票据
	ticketConfig, err := UserTicket.CreateConfig(&UserTicket.ArgsCreateConfig{
		DefaultExpireTime: 6000,
		OrgID:             newGroupData.OrgID,
		Mark:              "",
		Title:             "测试票据",
		Des:               "",
		CoverFileID:       0,
		DesFiles:          []int64{},
		ExemptionPrice:    10,
		ExemptionDiscount: 0,
		ExemptionMinPrice: 0,
		UseOrder:          false,
		LimitTimeType:     0,
		LimitCount:        0,
		StyleID:           0,
		Params:            nil,
	})
	ToolsTest.ReportData(t, err, ticketConfig)
	err = UserTicket.AddTicket(&UserTicket.ArgsAddTicket{
		OrgID:       newGroupData.OrgID,
		ConfigID:    ticketConfig.ID,
		UserID:      newGroupData.UserID,
		Count:       10,
		UseFromName: "测试赠送",
	})
	ToolsTest.ReportError(t, err)
	errCode, err := CreateMessage(&ArgsCreateMessage{
		GroupID:     newGroupData.ID,
		UserID:      newGroupData.UserID,
		MessageType: 2,
		Message:     "票据消息",
		Params:      nil,
		MoneyData:   ArgsCreateMessageMoney{},
		TicketData: ArgsCreateMessageTicket{
			ConfigID:   ticketConfig.ID,
			UseCount:   3,
			TakeType:   1,
			CountLimit: 1,
		},
	})
	ToolsTest.ReportErrorCode(t, errCode, err)
	TestGetMessageList(t)
	priceOrCount, errCode, err := TakeMessageMoneyOrTicket(&ArgsTakeMessageMoneyOrTicket{
		UserID:    newUserData.ID,
		GroupID:   newGroupData.ID,
		MessageID: newMessageData.ID,
	})
	ToolsTest.ReportErrorCode(t, errCode, err)
	if err == nil {
		t.Log("priceOrCount: ", priceOrCount)
		newUserTicketCount, _ := UserTicket.GetTicketCount(&UserTicket.ArgsGetTicketCount{
			ConfigID: ticketConfig.ID,
			UserID:   newUserData.ID,
		})
		if newUserTicketCount != 3 {
			t.Error("newUserTicketCount not 3, is: ", newUserTicketCount)
		}
	}
}

func TestClearMessage(t *testing.T) {
	TestOutChat(t)
	TestClearChat(t)
}
