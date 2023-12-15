package FinancePay

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
)

// ArgsUpdateStatusRemove 销毁支付参数
type ArgsUpdateStatusRemove struct {
	//操作人
	CreateInfo CoreSQLFrom.FieldsFrom
	//ID
	ID int64
	//key
	Key string
	//补充扩展
	Params []CoreSQLConfig.FieldsConfigType
}

// UpdateStatusRemove 销毁支付
func UpdateStatusRemove(args *ArgsUpdateStatusRemove) (errCode string, err error) {
	errCode, err = updateStatus(&argsUpdateStatus{
		CreateInfo: args.CreateInfo,
		ID:         args.ID,
		Key:        args.Key,
		PrevStatus: []int{0, 1},
		Status:     4,
		SetQuery:   "",
		SetMaps:    nil,
		Params:     args.Params,
	})
	return
}

// ArgsUpdateStatusFailed 交易失败处理参数
type ArgsUpdateStatusFailed struct {
	//操作人
	CreateInfo CoreSQLFrom.FieldsFrom
	//ID
	ID int64
	//key
	Key string
	//失败原因
	FailedCode    string
	FailedMessage string
	//补充扩展
	Params []CoreSQLConfig.FieldsConfigType
}

// UpdateStatusFailed 交易失败处理
func UpdateStatusFailed(args *ArgsUpdateStatusFailed) (errCode string, err error) {
	errCode, err = updateStatus(&argsUpdateStatus{
		CreateInfo: args.CreateInfo,
		ID:         args.ID,
		Key:        args.Key,
		PrevStatus: []int{0, 1, 7},
		Status:     2,
		SetQuery:   ", failed_code = :failed_code, failed_message = :failed_message",
		SetMaps: map[string]interface{}{
			"failed_code":    args.FailedCode,
			"failed_message": args.FailedMessage,
		},
		Params: args.Params,
	})
	return
}

// ArgsUpdateStatusExpire 交易过期参数
type ArgsUpdateStatusExpire struct {
	//操作人
	CreateInfo CoreSQLFrom.FieldsFrom
	//ID
	ID int64
	//key
	Key string
	//补充扩展
	Params []CoreSQLConfig.FieldsConfigType
}

// UpdateStatusExpire 交易过期
func UpdateStatusExpire(args *ArgsUpdateStatusExpire) (errCode string, err error) {
	errCode, err = updateStatus(&argsUpdateStatus{
		CreateInfo: args.CreateInfo,
		ID:         args.ID,
		Key:        args.Key,
		PrevStatus: []int{0, 1},
		Status:     5,
		SetQuery:   "",
		SetMaps:    nil,
		Params:     args.Params,
	})
	return
}
