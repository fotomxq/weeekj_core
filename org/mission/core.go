package OrgMission

import (
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	ClassTag "github.com/fotomxq/weeekj_core/v5/class/tag"
	"github.com/robfig/cron"
)

//组织任务处理模块

var (
	//定时器
	runTimer    *cron.Cron
	runAutoLock = false
	//Sort 任务分类
	Sort = ClassSort.Sort{
		SortTableName: "org_mission_sort",
	}
	//Tag 任务标签
	Tag = ClassTag.Tag{
		TagTableName: "org_mission_tags",
	}
)
