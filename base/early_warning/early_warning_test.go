package BaseEarlyWarning

import (
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	templateData FieldsTemplateType
	toData       FieldsToType
	bindData     FieldsBindType
	sendData     FieldsWaitType

	//测试template mark 名称
	testTemplateMark = "test-template"
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
	//只是测试是否可用，会直接跳出
	go Run()
}
