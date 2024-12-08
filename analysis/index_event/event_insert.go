package AnalysisIndexEvent

import "errors"

// InsertEvent 插入新的预警事件
func InsertEvent(args *FieldsEvent) (err error) {
	//检查指标是否已触发过风险
	var findData FieldsEvent
	findData, err = getEventBySystem(args.FromSystem, args.FromID, args.FromType)
	if err == nil && findData.ID > 0 {
		err = nil
		return
	}
	//如果存在数据，则更新
	if err == nil && findData.ID > 0 {
		//如果存在数据，且日期不同，则软删除数据
		if findData.YearMD != args.YearMD {
			err = eventDB.GetDelete().DeleteByID(findData.ID)
			if err != nil {
				err = errors.New("delete event error: " + err.Error())
				return
			}
			findData.ID = 0
		} else {
			args.ID = findData.ID
			err = eventDB.GetUpdate().UpdateByID(args)
			if err != nil {
				return
			}
			return
		}
	}
	if err != nil || findData.ID < 1 {
		//写入数据
		_, err = eventDB.GetInsert().InsertRow(args)
		if err != nil {
			return
		}
	}
	//反馈
	return
}

// InsertEventList 插入一组新的预警事件
func InsertEventList(args []FieldsEvent) (err error) {
	for _, v := range args {
		err = InsertEvent(&v)
		if err != nil {
			return
		}
	}
	return
}
