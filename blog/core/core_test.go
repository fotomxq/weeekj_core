package BlogCore

import (
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"github.com/mozillazg/go-pinyin"
	"testing"
)

var (
	isInit = false
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
	}
	isInit = true
}

func TestPinyin(t *testing.T) {
	newKeys := pinyin.LazyPinyin("d f", pinyin.NewArgs())
	t.Log(newKeys)
}
