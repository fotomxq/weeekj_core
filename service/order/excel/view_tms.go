package ServiceOrderExcel

import (
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	ERPProduct "github.com/fotomxq/weeekj_core/v5/erp/product"
	MallCore "github.com/fotomxq/weeekj_core/v5/mall/core"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2Mid "github.com/fotomxq/weeekj_core/v5/router2/mid"
	ServiceOrder "github.com/fotomxq/weeekj_core/v5/service/order"
	TMSTransport "github.com/fotomxq/weeekj_core/v5/tms/transport"
)

// GetViewTMS 查看详情和配送服务打印表
func GetViewTMS(c any, logErr string, orgID int64, orderID int64) {
	//路由上下文
	//ctx := Router2Mid.GetContext(c)
	//获取订单
	orderData, err := ServiceOrder.GetByID(&ServiceOrder.ArgsGetByID{
		ID:     orderID,
		OrgID:  orgID,
		UserID: -1,
	})
	if err != nil {
		Router2Mid.ReportWarnLog(c, logErr+", get order data, ", err, "err_no_data")
		return
	}
	//预先加载文件
	fileParams := fmt.Sprint("service_order_view_tms_", orgID, "_id_", orderID)
	if beforeLoadParamsFile(c, logErr, fileParams) {
		return
	}
	//文件名称
	fileName := fmt.Sprint("订单", orderData.SerialNumber, "详情和配送打印表.xlsx")
	//获取模板
	replaceExcel, err := getTemplate(fmt.Sprint("service", CoreFile.Sep, "order", CoreFile.Sep, "service_order_view_tms.xlsx"))
	if err != nil {
		Router2Mid.ReportWarnLog(c, logErr+", load excel template failed, ", err, "err_excel_template")
		return
	}
	//预先写入数据
	insertData := map[string]string{
		"H2":  fmt.Sprint("单号NO：", orderData.SerialNumber),
		"C4":  fmt.Sprint(orderData.AddressTo.Name),
		"G4":  fmt.Sprint(orderData.AddressTo.Phone),
		"C5":  fmt.Sprint(orderData.AddressTo.Address),
		"C6":  fmt.Sprint(orderData.CreateAt.Format("2006-01-02 15:04:05")),
		"G6":  fmt.Sprint(orderData.TransportTaskAt.Format("2006-01-02 15:04:05")),
		"C23": orderData.Des,
	}
	//写入货物信息
	var goodAllPrice float64 = 0
	vRowKey := 8
	for _, vGood := range orderData.Goods {
		if vRowKey > 21 {
			break
		}
		if vGood.From.System != "mall" {
			insertData[fmt.Sprint("C", vRowKey)] = vGood.From.Name
			insertData[fmt.Sprint("E", vRowKey)] = "件"
			insertData[fmt.Sprint("F", vRowKey)] = fmt.Sprint(vGood.Count)
			insertData[fmt.Sprint("G", vRowKey)] = fmt.Sprint(vGood.Price)
			continue
		}
		vMallProductData, _ := MallCore.GetProduct(&MallCore.ArgsGetProduct{
			ID:    vGood.From.ID,
			OrgID: -1,
		})
		if vMallProductData.ID < 1 {
			continue
		}
		vERPProductData := ERPProduct.GetProductBySN(vMallProductData.OrgID, vMallProductData.Code)
		insertData[fmt.Sprint("B", vRowKey)] = vERPProductData.Code
		vTitleOption := ""
		if vGood.OptionKey != "" {
			vTitleOption = MallCore.GetProductOtherOptionsName(vMallProductData.ID, vGood.OptionKey)
			if vTitleOption != "" {
				vTitleOption = fmt.Sprint("#", vTitleOption)
			}
		}
		insertData[fmt.Sprint("C", vRowKey)] = fmt.Sprint(vMallProductData.Title, vTitleOption)
		insertData[fmt.Sprint("E", vRowKey)] = "件"
		insertData[fmt.Sprint("F", vRowKey)] = fmt.Sprint(vGood.Count)
		vGoodPrice := CoreFilter.RoundToTwoDecimalPlaces(float64(MallCore.GetProductLastPrice(vMallProductData.ID, vGood.OptionKey)) / 100)
		insertData[fmt.Sprint("G", vRowKey)] = fmt.Sprint(vGoodPrice)
		insertData[fmt.Sprint("H", vRowKey)] = fmt.Sprint(CoreFilter.RoundToTwoDecimalPlaces(vGoodPrice * float64(vGood.Count)))
		goodAllPrice += vGoodPrice * float64(vGood.Count)
		//叠加
		vRowKey += 1
	}
	insertData["H22"] = fmt.Sprint(goodAllPrice)
	//写入配送信息
	if orderData.TransportID > 0 {
		tmsData, _ := TMSTransport.GetTransport(&TMSTransport.ArgsGetTransport{
			ID:     orderData.TransportID,
			OrgID:  -1,
			InfoID: -1,
			UserID: -1,
		})
		if tmsData.ID > 0 {
			bindData, _ := OrgCore.GetBind(&OrgCore.ArgsGetBind{
				ID:     tmsData.BindID,
				OrgID:  -1,
				UserID: -1,
			})
			insertData["F26"] = bindData.Name
			insertData["I26"] = bindData.Phone
		}
	}
	//写入数据
	quickInsertCol(replaceExcel, "", insertData)
	//设置样式
	//quickSetStyle(replaceExcel, "", "A4", fmt.Sprint("A", 4), fmt.Sprint("I", rowStep))
	//保存excel文件到临时文件
	_ = saveTemplateExcel2(c, logErr, fileParams, fileName, replaceExcel)
	//反馈成功
	return
}
