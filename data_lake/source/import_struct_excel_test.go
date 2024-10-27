package DataLakeSource

import (
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	"testing"
)

var (
	testImportStructExcelSrc  = "data_source_source_excel1.csv"
	testImportStructExcelSrc2 = "data_source_source_excel1.xlsx"
)

func TestInitImportStructExcel(t *testing.T) {
	TestInit(t)
	baseSrc, _ := CoreFile.BaseWDDir()
	baseSrc = fmt.Sprint(baseSrc, CoreFile.Sep, "..", CoreFile.Sep, "..", CoreFile.Sep, "builds", CoreFile.Sep, "test", CoreFile.Sep, "test_data", CoreFile.Sep)
	testImportStructExcelSrc = fmt.Sprint(baseSrc, testImportStructExcelSrc)
	testImportStructExcelSrc2 = fmt.Sprint(baseSrc, testImportStructExcelSrc2)
	if CoreFile.IsFile(testImportStructExcelSrc) {
		t.Log("testImportStructExcelSrc:", testImportStructExcelSrc)
	} else {
		t.Error("testImportStructExcelSrc not exists: ", testImportStructExcelSrc)
	}
	if CoreFile.IsFile(testImportStructExcelSrc2) {
		t.Log("testImportStructExcelSrc2:", testImportStructExcelSrc2)
	} else {
		t.Error("testImportStructExcelSrc2 not exists: ", testImportStructExcelSrc2)
	}
}

func TestImportStructExcel(t *testing.T) {
	newTableID, errCode, err := ImportStructExcel(&ArgsImportStructExcel{
		TableName:      "data_source_source_excel1_csv",
		TableDesc:      "data_source_source_excel1_csv",
		TipName:        "data_source_source_excel1_csv",
		ChannelName:    "test",
		ChannelTipName: "测试源头",
		Src:            testImportStructExcelSrc,
	})
	if err != nil {
		t.Error(errCode, err)
		return
	}
	_ = ClearFields(newTableID)
	_ = DeleteTable(newTableID)
	newTableID2, errCode2, err2 := ImportStructExcel(&ArgsImportStructExcel{
		TableName:      "data_source_source_excel1_xlsx",
		TableDesc:      "data_source_source_excel1_xlsx",
		TipName:        "data_source_source_excel1_xlsx",
		ChannelName:    "test",
		ChannelTipName: "测试源头",
		Src:            testImportStructExcelSrc,
	})
	if err2 != nil {
		t.Error(errCode2, err2)
		return
	}
	_ = ClearFields(newTableID2)
	_ = DeleteTable(newTableID2)
}

func TestClearImportStructExcel(t *testing.T) {
	TestClear(t)
}
