package TMSUserRunning

import (
	"encoding/json"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	"github.com/lib/pq"
)

// 生成追加日志的部分
func getLogData(des string, desFiles pq.Int64Array) (logData string, err error) {
	newLog := FieldsMissionLogs{
		{
			CreateAt: CoreFilter.GetNowTime(),
			Des:      des,
			DesFiles: desFiles,
		},
	}
	var newLogByte []byte
	newLogByte, err = json.Marshal(newLog)
	if err != nil {
		return
	}
	logData = string(newLogByte)
	return
}
