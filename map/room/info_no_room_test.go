package MapRoom

import "testing"

// 测试序列是否正常
func TestGetInfoNoInRoomInList(t *testing.T) {
	var rawList, dataList []DataInfoNoInRoom
	makeKey := 0
	for {
		if makeKey > 23 {
			break
		}
		rawList = append(rawList, DataInfoNoInRoom{
			InfoID: int64(makeKey + 1),
			Status: 0,
		})
		makeKey += 1
	}
	page := 3
	limit := 10
	step := (page - 1) * limit
	if step < 1 {
		step = 0
	}
	max := step + limit
	if len(rawList) <= step {
		t.Log("no data")
		return
	}
	if len(rawList) >= max {
		dataList = rawList[step:max]
	} else {
		dataList = rawList[step:]
	}
	t.Log("dataList len: ", len(dataList), ", start key: ", dataList[0].InfoID, ", end key: ", dataList[len(dataList)-1].InfoID)
}
