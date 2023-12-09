package CoreFilter

import "testing"

func TestRoundToTwoDecimalPlaces(t *testing.T) {
	data := RoundToTwoDecimalPlaces(55.800000000000004)
	t.Logf("data: %v", data)
}
