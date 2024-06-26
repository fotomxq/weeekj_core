package BlogExam

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

//在线考试模块
/**
1. 支持题库，可以设置题库、题目
2. 题目支持单选题、判断题、多选题、填空题、问答题
3. 除问答题外，其他题目可以设置标准答案，其中填空题如果不设置标准答案则不做判断处理
4. 可以抽取不同试卷题库的题目组成单独的考试
5. 考试可以设置开始和到期时间，用户开始问答后反馈数据集，问答结束后收集数据并反馈考试结果
6. 考试可以设置答案是否在前端直接判断，如果直接判断则同时会反馈正确答案方便前端处理；否则必须交给后台统一处理
7. 系统会记录考生每次考试的时间、统计次数、考试用时、得分，暂不记录具体答题细节
*/
var (
	//OpenSub 启动订阅
	OpenSub = false
)

func Init() {
	if OpenSub {
		_ = BaseService.SetService(&BaseService.ArgsSetService{
			ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
			Name:         "博客考试日志",
			Description:  "",
			EventSubType: "all",
			Code:         "blog_exam_log",
			EventType:    "nats",
			EventURL:     "/blog/exam/log",
			//TODO:待补充
			EventParams: "",
		})
	}
}
