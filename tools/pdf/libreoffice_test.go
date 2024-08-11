package ToolsPDF

import (
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	"testing"
)

func TestConvertToPDF(t *testing.T) {
	fileDir, _ := CoreFile.BaseWDDir()
	srcFileSrc := fileDir + CoreFile.Sep + "test_out.xlsx"
	newFileDirSrc := CoreFile.GetDir(srcFileSrc)
	outFileSrc := newFileDirSrc
	if CoreFile.IsFile(srcFileSrc) {
		t.Log("is file: ", srcFileSrc)
	}
	if !ConvertToPDF(srcFileSrc, outFileSrc) {
		t.Error("convert to pdf fail")
		return
	}
}
