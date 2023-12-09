package ServiceOrderExcel

import (
	"fmt"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLanguage "gitee.com/weeekj/weeekj_core/v5/core/language"
	Router2Mid "gitee.com/weeekj/weeekj_core/v5/router2/mid"
	ServiceCompany "gitee.com/weeekj/weeekj_core/v5/service/company"
	ServiceOrder "gitee.com/weeekj/weeekj_core/v5/service/order"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
	"strings"
)

// GetDefault 默认表下载
func GetDefault(c any, logErr string, orgID int64, args *ServiceOrder.ArgsGetList) {
	//路由上下文
	ctx := Router2Mid.GetContext(c)
	//时间范围不能超出1个月
	minAt, err := CoreFilter.GetTimeByDefault(args.TimeBetween.MinTime)
	if err != nil {
		Router2Mid.ReportWarnLog(c, logErr+", min at limit, ", err, "err_time")
		return
	}
	maxAt, err := CoreFilter.GetTimeByDefault(args.TimeBetween.MaxTime)
	if err != nil {
		Router2Mid.ReportWarnLog(c, logErr+", max at limit, ", err, "err_time")
		return
	}
	if CoreFilter.GetCarbonByTime(minAt).DiffInDaysWithAbs(CoreFilter.GetCarbonByTime(maxAt)) > 31 {
		Router2Mid.ReportWarnLog(c, logErr+", between time limit, ", err, "err_limit_time")
		return
	}
	//预先加载文件
	fileParams := fmt.Sprint("service_order_default_", orgID, "_args_", args)
	if beforeLoadParamsFile(c, logErr, fileParams) {
		return
	}
	//文件名称
	fileName := fmt.Sprint("订单数据导出.xlsx")
	//获取模板
	replaceExcel, err := getTemplate(fmt.Sprint("service", CoreFile.Sep, "order", CoreFile.Sep, "service_order_default.xlsx"))
	if err != nil {
		Router2Mid.ReportWarnLog(c, logErr+", load excel template failed, ", err, "err_excel_template")
		return
	}
	//获取数量
	_, orderCount, _ := ServiceOrder.GetList(args)
	//预先写入数据
	insertData := map[string]string{
		"A2": fmt.Sprint(CoreFilter.GetTimeToDefaultTime(minAt), " - ", CoreFilter.GetTimeToDefaultTime(maxAt), "          ", "共计", orderCount, "个"),
	}
	//行数
	rowStep := 4
	//读取和写入数据
	var page int64 = 1
	for {
		//获取数据
		args.Pages.Page = page
		orderList, _, err := ServiceOrder.GetList(args)
		if len(orderList) < 1 || err != nil {
			break
		}
		//遍历订单
		for _, vOrder := range orderList {
			//填充数据
			insertData[fmt.Sprint("A", rowStep)] = CoreFilter.GetTimeToDefaultTime(vOrder.CreateAt)
			insertData[fmt.Sprint("B", rowStep)] = fmt.Sprint(vOrder.SerialNumber, " / ", vOrder.SerialNumberDay)
			if vOrder.UserID > 0 {
				if vOrder.CompanyID > 0 {
					insertData[fmt.Sprint("C", rowStep)] = fmt.Sprint(UserCore.GetUserName(vOrder.UserID, false), " / ", ServiceCompany.GetCompanyName(vOrder.CompanyID))
				} else {
					insertData[fmt.Sprint("C", rowStep)] = fmt.Sprint(UserCore.GetUserName(vOrder.UserID, false))
				}
			}
			insertData[fmt.Sprint("D", rowStep)] = fmt.Sprint(CoreLanguage.GetLanguageText(ctx, fmt.Sprint("order_status_", vOrder.Status)), " / ", CoreLanguage.GetLanguageText(ctx, fmt.Sprint("order_refund_status_", vOrder.RefundStatus)))
			insertData[fmt.Sprint("E", rowStep)] = fmt.Sprint("￥", CoreFilter.GetPriceToShowPrice(vOrder.Price), " (", CoreLanguage.GetLanguageText(ctx, fmt.Sprint("order_pay_from_", vOrder.PayFrom)), ")")
			var vGoods []string
			for _, vGood := range vOrder.Goods {
				vGoods = append(vGoods, fmt.Sprint(vGood.From.Name, " x ", vGood.Count))
			}
			insertData[fmt.Sprint("F", rowStep)] = fmt.Sprint(strings.Join(vGoods, "\n\r"))
			insertData[fmt.Sprint("G", rowStep)] = fmt.Sprint(vOrder.AddressTo.Name, " / ", vOrder.AddressTo.Phone, " / ", vOrder.AddressTo.Address)
			insertData[fmt.Sprint("H", rowStep)] = fmt.Sprint(CoreLanguage.GetLanguageText(ctx, fmt.Sprint("order_tms_system_", vOrder.TransportSystem)), " (", vOrder.TransportSN, ")")
			insertData[fmt.Sprint("I", rowStep)] = fmt.Sprint(vOrder.Des)
			//叠加行数
			rowStep += 1
		}
		//下一页
		page += 1
	}
	//写入数据
	quickInsertCol(replaceExcel, "", insertData)
	//设置样式
	quickSetStyle(replaceExcel, "", "A4", fmt.Sprint("A", 4), fmt.Sprint("I", rowStep))
	//保存excel文件到临时文件
	_ = saveTemplateExcel2(c, logErr, fileParams, fileName, replaceExcel)
	//反馈成功
	return
}
