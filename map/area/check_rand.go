package MapArea

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

// CheckPointInAreasRand 获取符合条件的分区，但抽取任意一个
func CheckPointInAreasRand(args *ArgsCheckPointInAreas) (data FieldsArea, err error) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint("rand is failed, ", r))
			return
		}
	}()
	var dataList []FieldsArea
	dataList, err = CheckPointInAreas(args)
	if err != nil {
		return
	}
	if len(dataList) == 1 {
		data = dataList[0]
		return
	}
	key := CoreFilter.GetRandNumber(0, len(dataList)-1)
	data = dataList[key]
	return
}
