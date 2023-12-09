package CoreFilter

import "testing"

func TestGetMaxRand(t *testing.T) {
	result := GetMaxRand(100, 0, 0, 3)
	t.Log("result: ", result)
	result2 := GetMaxRand(100, result, 1, 3)
	t.Log("result2: ", result2)
	result3 := GetMaxRand(100, result+result2, 2, 3)
	t.Log("result3: ", result3)
	result4 := GetMaxRand(100, result+result2+result3, 3, 3)
	t.Log("result4: ", result4)
}
