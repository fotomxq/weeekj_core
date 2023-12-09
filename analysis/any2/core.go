package AnalysisAny2

import (
	"github.com/robfig/cron/v3"
	"time"
)

/**
第二代混合统计支持，主要改进有：
1. 不需要配置，直接指定mark即可记录数据
2. 提供各种图表所需的数据投放形式，方便前端整合数据
3. 自动归档，无需监管处理。但也保留支持归档处理方案
4. 更好的redis缓冲设计，减少不必要的全局删除方法
*/

var (
	//定时器
	runTimer    *cron.Cron
	runFileLock = false
	//缓冲时间
	cacheExpire = 604800
)

type waitAppendDataType struct {
	//写入模式
	// add 递增数据; re 覆盖数据
	Action string
	//标识码
	Mark string
	//创建时间
	CreateAt time.Time
	//组织ID
	// 可留空
	OrgID int64
	//用户ID
	// 可留空
	UserID int64
	//绑定ID
	BindID int64
	//扩展参数1
	// 可选扩展参数，默认给0无视
	Param1 int64
	//扩展参数2
	// 可选扩展参数，默认给0无视
	Param2 int64
	//数据
	Data int64
}
