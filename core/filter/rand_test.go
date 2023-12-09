package CoreFilter

import "testing"

func TestRandomWeightedValue(t *testing.T) {
	r := RandomWeightedValue([]int{10, 30, 60})
	t.Log(r)
}
