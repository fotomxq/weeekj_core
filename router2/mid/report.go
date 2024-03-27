package Router2Mid

import (
	"fmt"
	CoreLanguage "github.com/fotomxq/weeekj_core/v5/core/language"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 内部反馈头上下文
type reportContext struct {
	//总的上下文
	Context any
	//路由上下文
	RouterContext *gin.Context
	//业务上下文
	HeaderContext DataGetContextData
	//反馈头部类型
	HttpCode int
	//反馈编码
	Code string
	//状态
	Status bool
	//错误
	Err error
	//消息
	Msg string
	//覆盖消息
	ReplaceMsg string
	//数量集合
	Count int64
	//数据集合
	Data any
}

// reportDataType 通用数据反馈头
type reportDataType struct {
	//错误信息
	Status bool `json:"status"`
	//错误信息
	Code string `json:"code"`
	//错误描述
	Msg string `json:"msg"`
	//数据个数
	Count int64 `json:"count"`
	//数据集合
	Data any `json:"data"`
}

// ReportErrorBadRequest 反馈错误但不带JSON重新标定
func ReportErrorBadRequest(context any, code string) {
	//实现实例
	ctx := reportGetCtx(context, http.StatusBadRequest, nil, "", false, code, 0, nil)
	//反馈数据
	reportBaseReport(&ctx)
}

func ReportErrorBadRequestLog(context any, message string, err error, code string) {
	//实现实例
	ctx := reportGetCtx(context, http.StatusBadRequest, err, message, false, code, 0, nil)
	//反馈数据
	reportLogWarn(&ctx)
}

func ReportErrorBadRequestLogToParams(context any, message string, err error, code string, newMsg string) {
	//实现实例
	ctx := reportGetCtx(context, http.StatusBadRequest, err, message, false, code, 0, nil)
	ctx.ReplaceMsg = newMsg
	//反馈数据
	reportLogWarn(&ctx)
}

// ReportErrorLog 反馈错误并抛出日志
func ReportErrorLog(context any, message string, err error, code string) {
	//实现实例
	ctx := reportGetCtx(context, 0, err, message, false, code, 0, nil)
	//反馈数据
	reportLogErr(&ctx)
}

// ReportWarnLog 反馈警告并抛出错误
func ReportWarnLog(context any, message string, err error, code string) {
	//实现实例
	ctx := reportGetCtx(context, 0, err, message, false, code, 0, nil)
	//反馈数据
	reportLogWarn(&ctx)
}

func ReportActionCreateNoData(context any, logMsg string, err error, code string) {
	//修正code
	if err != nil && code == "" {
		code = "report_create_failed"
	}
	if err == nil {
		code = ""
	}
	//实现实例
	ctx := reportGetCtx(context, 0, err, logMsg, false, code, 0, nil)
	//判断错误
	ctx.Status = ctx.Err == nil
	//反馈数据
	reportAuto(&ctx)
}

func ReportActionCreate(context any, logMsg string, err error, code string, data interface{}) {
	//修正code
	if err != nil && code == "" {
		code = "report_create_failed"
	}
	if err == nil {
		code = ""
	}
	//实现实例
	ctx := reportGetCtx(context, 0, err, logMsg, false, code, 0, data)
	//判断错误
	ctx.Status = ctx.Err == nil
	//反馈数据
	reportAuto(&ctx)
}

func ReportActionUpdate(context any, logMsg string, err error, code string) {
	//修正code
	if err != nil && code == "" {
		code = "report_update_failed"
	}
	if err == nil {
		code = ""
	}
	//实现实例
	ctx := reportGetCtx(context, 0, err, logMsg, false, code, 0, nil)
	//判断错误
	ctx.Status = ctx.Err == nil
	//反馈数据
	reportAuto(&ctx)
}

func ReportActionDelete(context any, logMsg string, err error, code string) {
	//修正code
	if err != nil && code == "" {
		code = "report_delete_failed"
	}
	if err == nil {
		code = ""
	}
	//实现实例
	ctx := reportGetCtx(context, 0, err, logMsg, false, code, 0, nil)
	//判断错误
	ctx.Status = ctx.Err == nil
	//反馈数据
	reportAuto(&ctx)
}

// ReportData 通用反馈单一数据
func ReportData(context any, errMessage string, err error, code string, data interface{}) {
	//修正code
	if err != nil && code == "" {
		code = "report_data_empty"
	}
	if err == nil {
		code = ""
	}
	//实现实例
	ctx := reportGetCtx(context, 0, err, errMessage, false, code, 0, data)
	//判断错误
	ctx.Status = ctx.Err == nil
	//反馈数据
	reportAuto(&ctx)
}

func ReportDataNoErr(context any, err error, code string, data interface{}) {
	//修正code
	if err != nil && code == "" {
		code = "report_data_empty"
	}
	if err == nil {
		code = ""
	}
	//实现实例
	ctx := reportGetCtx(context, 0, err, "", false, code, 0, data)
	//判断错误
	ctx.Status = true
	//反馈数据
	reportAuto(&ctx)
}

// ReportDataList 通用反馈列表方案
func ReportDataList(context any, errMessage string, err error, code string, dataList interface{}, dataCount int64) {
	//修正code
	if err != nil && code == "" {
		code = "report_data_empty"
	}
	if err == nil {
		code = ""
	}
	//实现实例
	ctx := reportGetCtx(context, 0, err, "", false, code, dataCount, dataList)
	//判断错误
	ctx.Status = ctx.Err == nil
	//反馈数据
	reportAuto(&ctx)
}

// ReportBaseSuccess 反馈成功
func ReportBaseSuccess(context any) {
	//实现实例
	ctx := reportGetCtx(context, 0, nil, "", true, "", 0, nil)
	//反馈数据
	reportAuto(&ctx)
}

// ReportBaseBool 反馈成功或失败
func ReportBaseBool(context any, code string, b bool) {
	if b {
		code = ""
	}
	//实现实例
	ctx := reportGetCtx(context, 0, nil, "", b, code, 0, nil)
	//反馈数据
	reportAuto(&ctx)
}

// BaseData 反馈一般数据
func BaseData(context any, data interface{}) {
	//实现实例
	ctx := reportGetCtx(context, 0, nil, "", true, "", 0, data)
	//反馈数据
	reportAuto(&ctx)
}

// BaseDataCount 反馈数量类数据
func BaseDataCount(context any, count int64) {
	//实现实例
	ctx := reportGetCtx(context, 0, nil, "", true, "", count, nil)
	//反馈数据
	reportAuto(&ctx)
}

// ReportBaseDataList 反馈列队数据
func ReportBaseDataList(context any, count int64, data interface{}) {
	//实现实例
	ctx := reportGetCtx(context, 0, nil, "", true, "", count, data)
	//反馈数据
	reportAuto(&ctx)
}

// ReportBaseError 反馈错误
func ReportBaseError(context any, code string) {
	//修正错误代码
	if code == "" {
		code = "report_error"
	}
	//实现实例
	ctx := reportGetCtx(context, 0, nil, "", false, code, 0, nil)
	//反馈数据
	reportAuto(&ctx)
}

// reportGetLogMsg 日志内部附加处理
func reportGetLogMsg(c *gin.Context, cData DataGetContextData, msg string) string {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("log failed,", r, ",url:", c.Request.RequestURI, ",", cData.LogAppend, ",msg:", msg)
		}
	}()
	appendMsg := fmt.Sprint("url:", c.Request.RequestURI, ".")
	if cData.LogAppend != "" {
		appendMsg = fmt.Sprint(appendMsg, cData.LogAppend, ".")
	}
	userDataIn, b := c.Get("UserData")
	if b {
		userData := userDataIn.(UserCore.DataUserDataType)
		appendMsg = fmt.Sprint("userID:", userData.Info.ID, ",userPhone:", userData.Info.Phone, ",", appendMsg)
	}
	orgDataIn, b := c.Get("OrgData")
	if b {
		orgData := orgDataIn.(OrgCore.FieldsOrg)
		appendMsg = fmt.Sprint("orgID:", orgData.ID, ",", appendMsg)
	}
	orgBindDataIn, b := c.Get("OrgBindData")
	if b {
		bindData := orgBindDataIn.(OrgCore.FieldsBind)
		appendMsg = fmt.Sprint("orgBindID:", bindData.ID, ",", appendMsg)
	}
	if appendMsg != "" {
		msg = fmt.Sprint(appendMsg, msg)
	}
	return msg
}

// reportLogWarn 反馈警告错误
func reportLogWarn(ctx *reportContext) {
	//输出日志
	if ctx.Err != nil {
		CoreLog.Warn(reportGetLogMsg(ctx.RouterContext, ctx.HeaderContext, ctx.Msg), ",", ctx.Err)
	} else {
		CoreLog.Warn(reportGetLogMsg(ctx.RouterContext, ctx.HeaderContext, ctx.Msg))
	}
	//反馈
	reportBaseReport(ctx)
}

// reportLogErr 反馈错误
func reportLogErr(ctx *reportContext) {
	//输出日志
	if ctx.Err != nil {
		CoreLog.Error(reportGetLogMsg(ctx.RouterContext, ctx.HeaderContext, ctx.Msg), ",", ctx.Err)
	} else {
		CoreLog.Error(reportGetLogMsg(ctx.RouterContext, ctx.HeaderContext, ctx.Msg))
	}
	//反馈
	reportBaseReport(ctx)
}

// reportAuto 自动失败是否失败并反馈
func reportAuto(ctx *reportContext) {
	if ctx.Err != nil {
		reportLogWarn(ctx)
		return
	}
	reportBaseReport(ctx)
}

// reportBaseReport 通过反馈结构
func reportBaseReport(ctx *reportContext) {
	//反馈数据
	msg := ""
	if ctx.ReplaceMsg != "" {
		msg = ctx.ReplaceMsg
	} else {
		if ctx.Code != "" {
			msg = CoreLanguage.GetLanguageText(ctx.RouterContext, ctx.Code)
		}
	}
	res := reportDataType{
		Status: ctx.Status,
		Code:   ctx.Code,
		Msg:    msg,
		Count:  ctx.Count,
		Data:   ctx.Data,
	}
	ctx.RouterContext.JSON(ctx.HttpCode, &res)
	ctx.RouterContext.Abort()
}

// reportGin gin直接反馈头
func reportGin(ctxGin *gin.Context, haveLog bool, httpCode int, err error, msg string, status bool, code string, count int64, data interface{}) {
	//实现实例
	ctx := reportGetCtx(&RouterURLPublicC{
		Context:   ctxGin,
		LogAppend: "",
	}, httpCode, err, msg, status, code, count, data)
	//如果需要记录日志
	if haveLog {
		reportAuto(&ctx)
		return
	}
	//反馈数据
	reportBaseReport(&ctx)
}

// reportGetCtx 组装上下文
func reportGetCtx(context any, httpCode int, err error, msg string, status bool, code string, count int64, data interface{}) reportContext {
	//修正http状态
	if httpCode < 1 {
		httpCode = http.StatusOK
	}
	//获取上下文
	ctx, ctxData := GetContextData(context)
	//组装并反馈数据
	return reportContext{
		Context:       context,
		RouterContext: ctx,
		HeaderContext: ctxData,
		HttpCode:      httpCode,
		Code:          code,
		Status:        status,
		Err:           err,
		Msg:           msg,
		Count:         count,
		Data:          data,
	}
}
