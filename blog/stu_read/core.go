package BlogStuRead

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

//查看并阅读完成学习模块
/**
1. 用于用户完成指定学习任务
2. 统计用户是否完成了指定的学习内容
3. 内容根据博客核心复刻，支持性质一样
*/

var (
	//OpenSub 启动订阅
	OpenSub = false
)

func Init() {
	if OpenSub {
		_ = BaseService.SetService(&BaseService.ArgsSetService{
			ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
			Name:         "博客阅读学习日志",
			Description:  "",
			EventSubType: "all",
			Code:         "blog_stu_read_log",
			EventType:    "nats",
			EventURL:     "/blog/stu_read/log",
			//TODO:待补充
			EventParams: "",
		})
	}
}
