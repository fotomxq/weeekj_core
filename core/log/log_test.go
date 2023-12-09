package CoreLog

import (
	"testing"
)

func TestInit(t *testing.T) {
	Init(true, "weeekj", false)
	go Run()
	Info("test info")
	Error("test error")
	Debug("test debug")
}
