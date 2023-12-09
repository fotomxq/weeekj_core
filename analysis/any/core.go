package AnalysisAny

import (
	"github.com/robfig/cron"
)

//任意统计数据模块支持
/**
主要解决任意模块写入统计数据，并通过各类接口反馈数据和同步数据
1. 任意模块、组织或个人的统计数据，并提供缓存能力
2. 发生差异时，自动触发MQTT推送
3. API接口可随时调取相关统计数据
4. 数据会可以产生历史记录，同时会在超出N条后自动归档处理
5. 没有UI交互处理，所有配置全部采用配置注入形式实现

注意本模块无法解决的问题:
1. 统计数据混合运算
2. 单一标签只能推送一条数据
3. 只支持int64/string类型统计数字类数据的记录，其中string主要用于存储json数据集合
*/

var (
	//定时器
	runTimer    *cron.Cron
	runFileLock = false
	//NoMqtt 是否不推送MQTT
	// 用于外部服务更新数据
	NoMqtt = false
)

func Init() {
	//初始化变量
	NoMqtt = false
}
