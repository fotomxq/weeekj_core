package BaseDistribution

import (
	"time"
)

type FieldsDistribution struct {
	//唯一标记
	Mark string `bson:"Mark" json:"mark"`
	//创建时间
	CreateAt time.Time `bson:"CreateAt" json:"createAt"`
	//更新时间
	UpdateAt time.Time `bson:"UpdateAt" json:"updateAt"`
	//名称
	Name string `bson:"Name" json:"name"`
	//默认过期间隔
	ExpireInterval int64 `bson:"ExpireInterval" json:"expireInterval"`
	//默认心跳相应的触发方法
	// 注意，不发送和接收任何参数的方法，只做自动心跳和测试响应时间
	DefaultAction string `bson:"DefaultAction" json:"defaultAction"`
	DefaultFunc   string `bson:"DefaultFunc" json:"defaultFunc"`
}

type FieldsDistributionChild struct {
	//唯一标记
	Mark string `bson:"Mark" json:"mark"`
	//创建时间
	CreateAt time.Time `bson:"CreateAt" json:"createAt"`
	//更新时间
	UpdateAt time.Time `bson:"UpdateAt" json:"updateAt"`
	//服务器名称
	ServerName string `bson:"ServerName" json:"serverName"`
	ServerIP   string `bson:"ServerIP" json:"serverIP"`
	ServerPort string `bson:"ServerPort" json:"serverPort"`
	//过期时间，unix时间戳
	ExpireTime int64 `bson:"ExpireTime" json:"expireTime"`
	//累计执行次数，服务启动后的总累计次数
	RunCount int `bson:"RunCount" json:"runCount"`
	//统计的小时时间，对应time.getHour()的数据
	RunHour int `bson:"RunHour" json:"runHour"`
	//小时累计次数，统计更新时间对应的本小时内的累计使用次数
	RunHourCount int `bson:"RunHourCount" json:"runHourCount"`
	//服务器的响应时间，统计最近10次的平均值
	// 根据心跳包自动识别，单位：毫秒
	LastCountAvg int64 `bson:"LastCountAvg" json:"lastCountAvg"`
	//最近10次的列队
	LastCount []int64 `bson:"LastCount" json:"lastCount"`
}

//服务内部的子run方法
// 所有数据关联child记录
type FieldsDistributionChildRun struct {
	//唯一标记
	Mark string `bson:"Mark" json:"mark"`
	//创建时间
	CreateAt time.Time `bson:"CreateAt" json:"createAt"`
	//更新时间
	UpdateAt time.Time `bson:"UpdateAt" json:"updateAt"`
	//run名称
	RunMark string `bson:"RunMark" json:"runMark"`
	//服务器信息
	ServerIP   string `bson:"ServerIP" json:"serverIP"`
	ServerPort string `bson:"ServerPort" json:"serverPort"`
	//过期时间，unix时间戳
	// 该信息不影响负载均衡处理，只是做记录
	ExpireTime int64 `bson:"ExpireTime" json:"expireTime"`
	//累计执行次数，服务启动后的总累计次数
	RunCount int `bson:"RunCount" json:"runCount"`
	//统计的小时时间，对应time.getHour()的数据
	RunHour int `bson:"RunHour" json:"runHour"`
	//小时累计次数，统计更新时间对应的本小时内的累计使用次数
	RunHourCount int `bson:"RunHourCount" json:"runHourCount"`
}
