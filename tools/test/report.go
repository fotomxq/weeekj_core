package ToolsTest

import (
	"testing"
)

//ReportError 反馈错误
func ReportError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

//ReportErrorCode 反馈错误
func ReportErrorCode(t *testing.T, errCode string, err error) {
	if err != nil {
		t.Error(errCode, ", ", err)
	}
}

//ReportData 反馈错误和数据集合
func ReportData(t *testing.T, err error, data interface{}) {
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

//ReportDataList 反馈错误列表
func ReportDataList(t *testing.T, err error, data interface{}, dataCount int64) {
	if err != nil {
		t.Error(err)
	} else {
		if data == nil || dataCount < 1 {
			t.Error("data is empty")
		}
		t.Log(data, dataCount)
	}
}
