package MallCore

import (
	"testing"

	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	UserSubscription "github.com/fotomxq/weeekj_core/v5/user/subscription"
	UserTicket "github.com/fotomxq/weeekj_core/v5/user/ticket"
)

var (
	newProductData   FieldsCore
	newTicketConfig  UserTicket.FieldsConfig
	newUserSubConfig UserSubscription.FieldsConfig
)

func TestInitProduct(t *testing.T) {
	TestInit(t)
	TestCreateTransport(t)
	var err error
	//创建票据数据包
	newTicketConfig, err = UserTicket.CreateConfig(&UserTicket.ArgsCreateConfig{
		DefaultExpireTime: 300,
		OrgID:             TestOrg.OrgData.ID,
		Mark:              "test",
		Title:             "测试票据",
		Des:               "测试票据描述",
		CoverFileID:       0,
		DesFiles:          []int64{},
		ExemptionPrice:    100,
		ExemptionDiscount: 0,
		ExemptionMinPrice: 0,
		UseOrder:          false,
		LimitTimeType:     0,
		LimitCount:        0,
		StyleID:           0,
		Params:            nil,
	})
	ToolsTest.ReportData(t, err, newTicketConfig)
	err = UserTicket.AddTicket(&UserTicket.ArgsAddTicket{
		OrgID:       newTicketConfig.OrgID,
		ConfigID:    newTicketConfig.ID,
		UserID:      TestOrg.UserInfo.ID,
		Count:       3,
		UseFromName: "test",
	})
	ToolsTest.ReportError(t, err)
	//创建会员数据包
	newUserSubConfig, err = UserSubscription.CreateConfig(&UserSubscription.ArgsCreateConfig{
		OrgID:             TestOrg.OrgData.ID,
		Mark:              "test_sub",
		TimeType:          1,
		TimeN:             1,
		Currency:          86,
		Price:             1,
		PriceOld:          2,
		Title:             "测试会员",
		Des:               "测试会员描述",
		CoverFileID:       0,
		DesFiles:          []int64{},
		UserGroups:        []int64{},
		ExemptionPrice:    0,
		ExemptionDiscount: 100,
		ExemptionMinPrice: 0,
		Limits:            nil,
		ExemptionTime:     nil,
		StyleID:           0,
		Params:            nil,
	})
	ToolsTest.ReportData(t, err, newUserSubConfig)
	err = UserSubscription.SetSub(&UserSubscription.ArgsSetSub{
		OrgID:       TestOrg.OrgData.ID,
		ConfigID:    newUserSubConfig.ID,
		UserID:      TestOrg.UserInfo.ID,
		ExpireAt:    CoreFilter.GetNowTimeCarbon().AddDay().Time,
		HaveExpire:  false,
		UseFrom:     "test",
		UseFromName: "测试",
	})
	ToolsTest.ReportError(t, err)
}

func TestCreateProduct(t *testing.T) {
	data, errCode, err := CreateProduct(&ArgsCreateProduct{
		OrgID:                 TestOrg.OrgData.ID,
		ProductType:           0,
		IsVirtual:             false,
		SortID:                0,
		Tags:                  []int64{},
		Title:                 "测试商品",
		TitleDes:              "测试商品小标题",
		Des:                   "测试商品描述",
		CoverFileIDs:          []int64{},
		DesFiles:              []int64{},
		Currency:              86,
		PriceReal:             10,
		PriceExpireAt:         CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().AddDay().Time),
		Price:                 30,
		Integral:              40,
		IntegralPrice:         0,
		IntegralTransportFree: false,
		UserSubPrice: []FieldsUserSubPrice{
			{
				ID:    newUserSubConfig.ID,
				Price: 5,
			},
		},
		UserTicket:  []int64{newTicketConfig.ID},
		TransportID: newTransportData.ID,
		Address: CoreSQLAddress.FieldsAddress{
			Country:    86,
			Province:   0,
			City:       0,
			Address:    "测试地址",
			MapType:    0,
			Longitude:  11,
			Latitude:   11,
			Name:       "测试地址",
			NationCode: "86",
			Phone:      "17777777777",
		},
		WarehouseProductID: 0,
		Weight:             500,
		Count:              10,
		OtherOptions:       DataOtherOptions{},
		GivingTickets:      nil,
		Params:             nil,
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		newProductData = data
	}
}

func TestUpdateProduct(t *testing.T) {
	errCode, err := UpdateProduct(&ArgsUpdateProduct{
		ID:                    newProductData.ID,
		OrgID:                 newProductData.OrgID,
		IsVirtual:             newProductData.IsVirtual,
		Sort:                  newProductData.Sort,
		SortID:                newProductData.SortID,
		Tags:                  newProductData.Tags,
		Title:                 newProductData.Title,
		TitleDes:              newProductData.TitleDes,
		Des:                   newProductData.Des,
		CoverFileIDs:          newProductData.CoverFileIDs,
		DesFiles:              newProductData.DesFiles,
		Currency:              newProductData.Currency,
		PriceReal:             newProductData.PriceReal,
		PriceExpireAt:         CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().AddDay().Time),
		Price:                 newProductData.Price,
		Integral:              newProductData.Integral,
		IntegralPrice:         newProductData.IntegralPrice,
		IntegralTransportFree: newProductData.IntegralTransportFree,
		UserSubPrice:          newProductData.UserSubPrice,
		UserTicket:            newProductData.UserTicket,
		TransportID:           newProductData.TransportID,
		Address:               newProductData.Address,
		WarehouseProductID:    newProductData.WarehouseProductID,
		Weight:                newProductData.Weight,
		Count:                 newProductData.Count,
		OtherOptions:          DataOtherOptions{},
		GivingTickets:         newProductData.GivingTickets,
		Params:                newProductData.Params,
	})
	if err != nil {
		t.Error(errCode, err)
	}
}

func TestGetProduct(t *testing.T) {
	var err error
	newProductData, err = GetProduct(&ArgsGetProduct{
		ID:    newProductData.ID,
		OrgID: newProductData.OrgID,
	})
	ToolsTest.ReportData(t, err, newProductData)
}

func TestGetProductList(t *testing.T) {
	dataList, dataCount, err := GetProductList(&ArgsGetProductList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:            -1,
		SortID:           0,
		Tags:             nil,
		NeedIsPublish:    false,
		IsPublish:        false,
		PriceMin:         0,
		PriceMax:         0,
		NeedHaveVIP:      false,
		HaveVIP:          false,
		Tickets:          []int64{},
		NeedHaveCount:    false,
		HaveCount:        false,
		NeedHaveIntegral: false,
		HaveIntegral:     false,
		ParentID:         0,
		TransportID:      0,
		IsRemove:         false,
		Search:           "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetProducts(t *testing.T) {
	data, err := GetProducts(&ArgsGetProducts{
		IDs:        []int64{newProductData.ID},
		HaveRemove: false,
		OrgID:      0,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetProductsName(t *testing.T) {
	data, err := GetProductsName(&ArgsGetProducts{
		IDs:        []int64{newProductData.ID},
		HaveRemove: false,
		OrgID:      newProductData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetProductsByTickets(t *testing.T) {
	data, err := GetProductsByTickets(&ArgsGetProductsByTickets{
		TicketIDs: []int64{2},
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetProductTop(t *testing.T) {
	data, err := GetProductTop(&ArgsGetProduct{
		ID:    newProductData.ID,
		OrgID: newProductData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetProductAddress(t *testing.T) {
	data, err := getProductAddress(newProductData.ID)
	ToolsTest.ReportData(t, err, data)
}

func TestCheckProductSub(t *testing.T) {
	errCode, err := CheckProductSub(&ArgsCheckProductSub{
		ID:         newProductData.ID,
		OrgID:      newProductData.OrgID,
		SortID:     newProductData.SortID,
		Tags:       []int64{},
		UserSubs:   []int64{},
		UserTicket: []int64{},
		BuyCount:   1,
	})
	if err != nil {
		t.Error(errCode, err)
	}
}

func TestCheckProductTicket(t *testing.T) {
	errCode, err := CheckProductTicket(&ArgsCheckProductTicket{
		ID:         newProductData.ID,
		OrgID:      newProductData.OrgID,
		SortID:     newProductData.SortID,
		Tags:       []int64{},
		UserSubs:   []int64{},
		UserTicket: []int64{},
		BuyCount:   1,
	})
	if err != nil {
		t.Error(errCode, err)
	}
}

func TestDeleteProduct(t *testing.T) {
	err := DeleteProduct(&ArgsDeleteProduct{
		ID:    newProductData.ID,
		OrgID: newProductData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestReturnProduct(t *testing.T) {
	err := ReturnProduct(&ArgsReturnProduct{
		ID:    newProductData.ID,
		OrgID: newProductData.OrgID,
	})
	ToolsTest.ReportError(t, err)
	TestDeleteProduct(t)
}

func TestClearProduct(t *testing.T) {
	err := UserTicket.DeleteConfig(&UserTicket.ArgsDeleteConfig{
		ID:    newTicketConfig.ID,
		OrgID: newTicketConfig.OrgID,
	})
	ToolsTest.ReportError(t, err)
	err = UserSubscription.DeleteConfig(&UserSubscription.ArgsDeleteConfig{
		ID:    newUserSubConfig.ID,
		OrgID: newUserSubConfig.OrgID,
	})
	ToolsTest.ReportError(t, err)
	TestDeleteTransport(t)
	TestClearTransport(t)
}
