package BaseMonitorPostgresql

import "time"

type FieldsData struct {
	//动作类型
	// select/get/insert/update/delete/analysis
	Action string
	//消息内容
	Msg string
	//是否为事务关系
	IsBegin bool
	//开始时间
	StartAt time.Time
	//结束时间
	EndAt time.Time
	//执行时间
	RunSec int64
	//反馈尺寸 单位字节
	ResultSize int64
	//是否存在报错
	Err error
}

type FieldsAnalysis struct {
	//连接池当前总数
	ConnectCount int64
	//最大连接数量
	ConnectMaxCount int64
	//执行次数
	AllCount int64
	//细分执行次数
	SelectCount   int64
	GetCount      int64
	InsertCount   int64
	UpdateCount   int64
	DeleteCount   int64
	AnalysisCount int64
	//总错误次数
	AllErrCount int64
	//细分错误次数
	SelectErrCount   int64
	GetErrCount      int64
	InsertErrCount   int64
	UpdateErrCount   int64
	DeleteErrCount   int64
	AnalysisErrCount int64
	//事物关系次数
	AllBeginCount      int64
	SelectBeginCount   int64
	GetBeginCount      int64
	InsertBeginCount   int64
	UpdateBeginCount   int64
	DeleteBeginCount   int64
	AnalysisBeginCount int64
}
