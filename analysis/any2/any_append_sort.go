package AnalysisAny2

// AppendDataSort 添加排名类数据
// 废弃模块，直接写入数据，sql自动排序处理即可，不需要专门做写入处理
// desc 是否为倒叙排列
// limit 排名限制，超出的会被裁剪掉
//func AppendDataSort(action, mark string, createAt time.Time, orgID, userID, bindID, param1, param2 int64, data int64, desc bool, limit int) {
//	appendLog := "analysis any2 append data to sort, "
//	//获取配置
//	configData, err := getConfigByMark(mark, false)
//	if err != nil {
//		CoreLog.Warn(appendLog, "get config mark: ", mark, ", ", err)
//		return
//	}
//	//根据配置计算要查询的时间范围
//	nowAt := CoreFilter.GetNowTimeCarbon()
//	if createAt.Unix() > 1000000 {
//		nowAt = nowAt.CreateFromGoTime(createAt)
//	}
//	var findAtMin time.Time
//	var findAtMax time.Time
//	switch configData.Particle {
//	case 0:
//		findAtMin = nowAt.StartOfHour().Time
//		findAtMax = nowAt.EndOfHour().Time
//	case 1:
//		findAtMin = nowAt.StartOfDay().Time
//		findAtMax = nowAt.EndOfDay().Time
//	case 2:
//		findAtMin = nowAt.StartOfWeek().Time
//		findAtMax = nowAt.EndOfWeek().Time
//	case 3:
//		findAtMin = nowAt.StartOfMonth().Time
//		findAtMax = nowAt.EndOfMonth().Time
//	case 4:
//		findAtMin = nowAt.StartOfYear().Time
//		findAtMax = nowAt.EndOfYear().Time
//	default:
//		findAtMin = nowAt.StartOfHour().Time
//		findAtMax = nowAt.EndOfHour().Time
//	}
//	//写入数据
//	AppendData(action, mark, createAt, orgID, userID, bindID, param1, param2, data)
//	//获取排名数据
//	sortList := GetSort(&ArgsGetSort{
//		Mark:    mark,
//		OrgID:   orgID,
//		UserID:  -1,
//		BindID:  -1,
//		Params1: -1,
//		Params2: -1,
//		MinAt:   findAtMin,
//		MaxAt:   findAtMax,
//		Desc:    desc,
//		Limit:   limit,
//	})
//	if len(sortList) > 0 {
//		//如果数量超出限制，则删除超出部分数据
//		//if len(sortList) > limit {
//		//	for k := limit; k < len(sortList); k++ {
//		//		v := sortList[k]
//		//		deleteAnyByID(v.ID)
//		//	}
//		//}
//		//清理缓冲
//		clearAnyCache(configData.ID)
//	}
//}
