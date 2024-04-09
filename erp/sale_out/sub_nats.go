package ERPSaleOut

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	ERPProduct "github.com/fotomxq/weeekj_core/v5/erp/product"
	MallCoreMod "github.com/fotomxq/weeekj_core/v5/mall/core/mod"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	ServiceOrderMod "github.com/fotomxq/weeekj_core/v5/service/order/mod"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//订单完成
	CoreNats.SubDataByteNoErr("service_order_status", "/service/order/status", subNatsOrderFinish)
	//订单退货
	CoreNats.SubDataByteNoErr("service_order_update", "/service/order/update", subNatsOrderRefund)
	//创建虚拟订单
	CoreNats.SubDataByteNoErr("service_order_create_wait_virtual_finish", "/service/order/create_wait_virtual_finish", subNatsCreateWaitVirtualFinish)
}

// 订单完成
func subNatsOrderFinish(_ *nats.Msg, action string, orderID int64, _ string, _ []byte) {
	//日志
	appendLog := "erp sale out sub nats order finish, "
	//必须是订单完成
	if action != "finish" {
		return
	}
	//获取订单
	orderData := ServiceOrderMod.GetByIDNoErr(orderID)
	if orderData.ID < 1 {
		return
	}
	//如果存在虚拟订单标记，则跳出
	if orderData.Params.GetValNoBool("virtual_sync") == "true" {
		return
	}
	//获取开关
	if orderData.OrgID < 1 {
		return
	}
	erpSaleOutSyncServiceOrder := OrgCore.Config.GetConfigValBoolNoErr(orderData.OrgID, "ERPSaleOutSyncServiceOrder")
	if !erpSaleOutSyncServiceOrder {
		return
	}
	//分析订单商品并构建出货单
	for _, vGood := range orderData.Goods {
		//跳过非正常商品
		if vGood.From.System != "mall" {
			continue
		}
		if vGood.From.Mark != "" {
			continue
		}
		//找到商品数据
		vMallProduct := MallCoreMod.GetProductNoErr(vGood.From.ID, -1)
		if vMallProduct.ID < 1 || vMallProduct.WarehouseProductID < 1 {
			continue
		}
		//修正价格设计
		if vGood.Price > 0 {
			vMallProduct.Price = vGood.Price
		}
		//找到商品ERP产品
		vErpProductData := ERPProduct.GetProductByIDNoErr(vMallProduct.WarehouseProductID)
		if vErpProductData.ID < 1 {
			continue
		}
		//找到产品的供货商
		var erpFromCompanyData ERPProduct.FieldsProductCompany
		if vErpProductData.CompanyID > 0 {
			erpFromCompanyData = ERPProduct.GetProductCompany(vErpProductData.OrgID, vErpProductData.ID, vErpProductData.CompanyID)
		} else {
			//根据最低成本价随机抽取一个供货商
			erpFromCompanyData = ERPProduct.GetProductCompanyRand(vErpProductData.ID, 1)
		}
		//构建出货单
		_, errCode, err := CreateLog(&ArgsCreateLog{
			OrgID:               vErpProductData.OrgID,
			CreateOrgBindID:     0,
			FromOrgBindID:       0,
			SellOrgBindID:       0,
			ConfirmOrgBindID:    0,
			SaleType:            0,
			SyncSystem:          "",
			SyncHash:            "",
			OrderID:             orderData.ID,
			MallProductID:       vMallProduct.ID,
			ERPProductID:        vErpProductData.ID,
			FromCompanyID:       erpFromCompanyData.CompanyID,
			BuyUserID:           orderData.UserID,
			BuyCompanyID:        orderData.CompanyID,
			ToAddress:           orderData.AddressTo,
			SellCount:           vGood.Count,
			SellRawPrice:        vGood.Count * vMallProduct.Price,
			SellRawTaxPrice:     vGood.Count * vMallProduct.Price,
			SellFinalPrice:      vGood.Count * vGood.Price,
			SellFinalTaxPrice:   vGood.Count * vGood.Price,
			SellExPrice:         vGood.Count * (vGood.Price - vMallProduct.Price),
			PurchasePrice:       vGood.Count * erpFromCompanyData.CostPrice,
			PurchaseTaxPrice:    vGood.Count * erpFromCompanyData.TaxCostPrice,
			PurchaseFixPrice:    vGood.Count * erpFromCompanyData.CostPrice,
			PurchaseFixTaxPrice: vGood.Count * erpFromCompanyData.TaxCostPrice,
			ProfitPrice:         vGood.Count*vGood.Price - vGood.Count*erpFromCompanyData.CostPrice,
			ProfitTaxPrice:      vGood.Count*vGood.Price - vGood.Count*erpFromCompanyData.TaxCostPrice,
		})
		if err != nil {
			CoreLog.Error(appendLog, "create log failed, ", errCode, ", err: ", err)
		}
	}
}

// 订单退货
func subNatsOrderRefund(_ *nats.Msg, action string, orderID int64, _ string, _ []byte) {
	//日志
	appendLog := "erp sale out sub nats order refund, "
	//必须是订单完成
	if action != "refund" {
		return
	}
	//获取订单
	orderData := ServiceOrderMod.GetByIDNoErr(orderID)
	if orderData.ID < 1 {
		return
	}
	//获取开关
	if orderData.OrgID < 1 {
		return
	}
	erpSaleOutSyncServiceOrder := OrgCore.Config.GetConfigValBoolNoErr(orderData.OrgID, "ERPSaleOutSyncServiceOrder")
	if !erpSaleOutSyncServiceOrder {
		return
	}
	//分析订单商品并构建出货单
	for _, vGood := range orderData.Goods {
		//跳过非正常商品
		if vGood.From.System != "mall" {
			continue
		}
		if vGood.From.Mark != "" {
			continue
		}
		//找到商品数据
		vMallProduct := MallCoreMod.GetProductNoErr(vGood.From.ID, -1)
		if vMallProduct.ID < 1 || vMallProduct.WarehouseProductID < 1 {
			continue
		}
		//找到商品ERP产品
		vErpProductData := ERPProduct.GetProductByIDNoErr(vMallProduct.WarehouseProductID)
		if vErpProductData.ID < 1 {
			continue
		}
		//找到产品的供货商
		var erpFromCompanyData ERPProduct.FieldsProductCompany
		if vErpProductData.CompanyID > 0 {
			erpFromCompanyData = ERPProduct.GetProductCompany(vErpProductData.OrgID, vErpProductData.ID, vErpProductData.CompanyID)
		} else {
			//根据最低成本价随机抽取一个供货商
			erpFromCompanyData = ERPProduct.GetProductCompanyRand(vErpProductData.ID, 1)
		}
		//构建出货单
		_, errCode, err := CreateLog(&ArgsCreateLog{
			OrgID:               vErpProductData.OrgID,
			CreateOrgBindID:     0,
			FromOrgBindID:       0,
			SellOrgBindID:       0,
			ConfirmOrgBindID:    0,
			SaleType:            1,
			SyncSystem:          "",
			SyncHash:            "",
			OrderID:             orderData.ID,
			MallProductID:       vMallProduct.ID,
			ERPProductID:        vErpProductData.ID,
			FromCompanyID:       erpFromCompanyData.CompanyID,
			BuyUserID:           orderData.UserID,
			BuyCompanyID:        orderData.CompanyID,
			ToAddress:           orderData.AddressTo,
			SellCount:           0 - vGood.Count,
			SellRawPrice:        0 - vGood.Count*vMallProduct.Price,
			SellRawTaxPrice:     0 - vGood.Count*vMallProduct.Price,
			SellFinalPrice:      0 - vGood.Count*vGood.Price,
			SellFinalTaxPrice:   0 - vGood.Count*vGood.Price,
			SellExPrice:         0 - vGood.Count*(vGood.Price-vMallProduct.Price),
			PurchasePrice:       0 - vGood.Count*erpFromCompanyData.CostPrice,
			PurchaseTaxPrice:    0 - vGood.Count*erpFromCompanyData.TaxCostPrice,
			PurchaseFixPrice:    0 - vGood.Count*erpFromCompanyData.CostPrice,
			PurchaseFixTaxPrice: 0 - vGood.Count*erpFromCompanyData.TaxCostPrice,
			ProfitPrice:         0 - (vGood.Count*vGood.Price - vGood.Count*erpFromCompanyData.CostPrice),
			ProfitTaxPrice:      0 - (vGood.Count*vGood.Price - vGood.Count*erpFromCompanyData.TaxCostPrice),
		})
		if err != nil {
			CoreLog.Error(appendLog, "create log failed, ", errCode, ", err: ", err)
		}
	}
}

func subNatsCreateWaitVirtualFinish(_ *nats.Msg, action string, orderID int64, _ string, params []byte) {
	//日志
	appendLog := "erp sale out sub nats order finish, "
	//必须是订单完成
	if action != "finish" {
		return
	}
	//获取订单
	orderData := ServiceOrderMod.GetByIDNoErr(orderID)
	if orderData.ID < 1 {
		return
	}
	//解析参数
	type paramsType struct {
		Products []struct {
			//商品ID
			ID int64 `db:"id" json:"id" check:"id"`
			//选项Key
			// 如果给空，则该商品必须也不包含选项
			OptionKey string `db:"option_key" json:"optionKey" check:"mark" empty:"true"`
			//购买数量
			// 如果为0，则只判断单价的价格
			BuyCount int `db:"buy_count" json:"buyCount" check:"int64Than0"`
			//进货价
			PriceIn int64 `db:"price_in" json:"priceIn" check:"price" empty:"true"`
			//销售价格
			PriceOut int64 `db:"price_out" json:"priceOut" check:"price" empty:"true"`
		}
	}
	var paramsRaw paramsType
	if err := CoreNats.ReflectDataByte(params, &paramsRaw); err != nil {
		CoreLog.Error(appendLog, "reflect data byte, ", err)
		return
	}
	//获取开关
	if orderData.OrgID < 1 {
		return
	}
	erpSaleOutSyncServiceOrder := OrgCore.Config.GetConfigValBoolNoErr(orderData.OrgID, "ERPSaleOutSyncServiceOrder")
	if !erpSaleOutSyncServiceOrder {
		return
	}
	//分析订单商品并构建出货单
	for _, vGood := range orderData.Goods {
		//跳过非正常商品
		if vGood.From.System != "mall" {
			continue
		}
		if vGood.From.Mark != "" {
			continue
		}
		//找到商品数据
		vMallProduct := MallCoreMod.GetProductNoErr(vGood.From.ID, -1)
		if vMallProduct.ID < 1 || vMallProduct.WarehouseProductID < 1 {
			continue
		}
		//找到商品ERP产品
		vErpProductData := ERPProduct.GetProductByIDNoErr(vMallProduct.WarehouseProductID)
		if vErpProductData.ID < 1 {
			continue
		}
		//找到产品的供货商
		var erpFromCompanyData ERPProduct.FieldsProductCompany
		if vErpProductData.CompanyID > 0 {
			erpFromCompanyData = ERPProduct.GetProductCompany(vErpProductData.OrgID, vErpProductData.ID, vErpProductData.CompanyID)
		} else {
			//根据最低成本价随机抽取一个供货商
			erpFromCompanyData = ERPProduct.GetProductCompanyRand(vErpProductData.ID, 1)
		}
		//修正进出价格
		for _, v2 := range paramsRaw.Products {
			if vGood.From.ID == v2.ID && vGood.From.Mark == v2.OptionKey {
				vMallProduct.Price = v2.PriceOut
				erpFromCompanyData.CostPrice = v2.PriceIn
				erpFromCompanyData.TaxCostPrice = v2.PriceIn
				break
			}
		}
		//构建出货单
		_, errCode, err := CreateLog(&ArgsCreateLog{
			OrgID:               vErpProductData.OrgID,
			CreateOrgBindID:     0,
			FromOrgBindID:       0,
			SellOrgBindID:       0,
			ConfirmOrgBindID:    0,
			SaleType:            0,
			SyncSystem:          "",
			SyncHash:            "",
			OrderID:             orderData.ID,
			MallProductID:       vMallProduct.ID,
			ERPProductID:        vErpProductData.ID,
			FromCompanyID:       erpFromCompanyData.CompanyID,
			BuyUserID:           orderData.UserID,
			BuyCompanyID:        orderData.CompanyID,
			ToAddress:           orderData.AddressTo,
			SellCount:           vGood.Count,
			SellRawPrice:        vGood.Count * vMallProduct.Price,
			SellRawTaxPrice:     vGood.Count * vMallProduct.Price,
			SellFinalPrice:      vGood.Count * vGood.Price,
			SellFinalTaxPrice:   vGood.Count * vGood.Price,
			SellExPrice:         vGood.Count * (vGood.Price - vMallProduct.Price),
			PurchasePrice:       vGood.Count * erpFromCompanyData.CostPrice,
			PurchaseTaxPrice:    vGood.Count * erpFromCompanyData.TaxCostPrice,
			PurchaseFixPrice:    vGood.Count * erpFromCompanyData.CostPrice,
			PurchaseFixTaxPrice: vGood.Count * erpFromCompanyData.TaxCostPrice,
			ProfitPrice:         vGood.Count*vGood.Price - vGood.Count*erpFromCompanyData.CostPrice,
			ProfitTaxPrice:      vGood.Count*vGood.Price - vGood.Count*erpFromCompanyData.TaxCostPrice,
		})
		if err != nil {
			CoreLog.Error(appendLog, "create log failed, ", errCode, ", err: ", err)
		}
	}
}
