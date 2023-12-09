package IOTDevice

import (
	ClassSort "gitee.com/weeekj/weeekj_core/v5/class/sort"
	ClassTag "gitee.com/weeekj/weeekj_core/v5/class/tag"
	CoreHighf "gitee.com/weeekj/weeekj_core/v5/core/highf"
)

var (
	//Sort 设备组织下分类
	Sort = ClassSort.Sort{
		SortTableName: "iot_core_device_sort",
	}
	//Tag 设备组织下标签
	Tag = ClassTag.Tag{
		TagTableName: "iot_core_device_tags",
	}
	//OpenSub 是否启动订阅
	OpenSub = false
	//删除旧的日志拦截器
	autoLogDeleteBlocker CoreHighf.BlockerWait
)

func Init() {
	autoLogDeleteBlocker.Init(1800)
	if OpenSub {
		//订阅消息列队
		subNats()
	}
}
