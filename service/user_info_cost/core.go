package ServiceUserInfoCost

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	IOTMQTT "github.com/fotomxq/weeekj_core/v5/iot/mqtt"
	"github.com/robfig/cron"
)

//人员成本计算
// 支持对人员信息、房间的成本、消费情况记录

/** 数据来源
- 传感器的数据汇总
- 手动录入成本消耗
- 老人缴费收入
- 手动流入收入
*/

/** 关于传感器
传感器针对不同的类型，将做特别处理。例如在设备扩展中将保留该设备每小时传感数据单位的耗能计数，例如电流检测设备的每小时电流X对应费用Y，计算出每小时的总耗电量设计。
支持的几种模式：
- 电流检测器：根据每秒电流传感量及霍尔计数，计算每小时耗电量，结合本模块的每小时耗电费用，计算每小时总的电量消耗量。
- 水流监测：根据每秒水流监测量，计算每小时耗水量，结合本模块水费，计算每小时总的水费消耗量。
*/

var (
	//定时器
	runTimer *cron.Cron
	runLock  = false
	//OpenSub 是否启动订阅
	OpenSub = false
)

// Init 初始化
func Init() {
	if OpenSub {
		//初始化mqtt订阅
		IOTMQTT.AppendSubFunc(initSub)
	}
}

// 订阅方法集合
func initSub() {
	//请求耗能数据
	if token := IOTMQTT.MQTTClient.Subscribe("service/user/info/cost/list", 0, subUserInfoCostList); token.Wait() && token.Error() != nil {
		CoreLog.Error("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//请求一组房间最新的耗能数据
	if token := IOTMQTT.MQTTClient.Subscribe("service/user/info/cost/last", 0, subUserInfoCostLast); token.Wait() && token.Error() != nil {
		CoreLog.Error("sub mqtt but token is error, " + token.Error().Error())
		return
	}
}
