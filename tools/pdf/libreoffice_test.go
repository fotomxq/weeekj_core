package ToolsPDF

import (
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	"testing"
)

func TestConvertToPDF(t *testing.T) {
	fileDir, _ := CoreFile.BaseWDDir()
	srcFileSrc := fileDir + CoreFile.Sep + "test_out.xlsx"
	outFileSrc := fileDir + CoreFile.Sep + "test_out.pdf"
	if CoreFile.IsFile(srcFileSrc) {
		t.Log("is file: ", srcFileSrc)
	}
	ConvertToPDF(srcFileSrc, outFileSrc)
}
