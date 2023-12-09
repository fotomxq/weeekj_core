package MapUserArea

import "github.com/robfig/cron"

//用户辖区判断处理模块
// 用于判断用户GPS和电子围栏，是否在范围内，或是否超出区域
// 1\ 超出区域后，可设置联动到组织任务
// 2\ 判断用户是否在范围内

var(
	//用户专用电子围栏系统专用mark
	fenceMark = "user_area"
	//定时器
	runTimer *cron.Cron
	runAreaLock = false
)
