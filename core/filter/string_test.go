package CoreFilter

import "testing"

func TestCutStringAndEncrypt(t *testing.T) {
	t.Log(CutStringAndEncrypt("1287", 6, 6))
}
