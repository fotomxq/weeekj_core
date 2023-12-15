package BaseLog

import (
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

var (
	isInit  = false
	newData FieldsLog
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
	}
	isInit = true
	CoreFile.BaseSrc = "../../builds/test"
	Router2SystemConfig.RootDir = CoreFile.BaseSrc
}

func TestRun(t *testing.T) {
	//构建虚拟数据
	fileDir := CoreFile.BaseSrc + CoreFile.Sep + "log" + CoreFile.Sep + CoreFilter.GetNowTimeCarbon().Format("20060102")
	if err := CoreFile.CreateFolder(fileDir); err != nil {
		t.Error(err)
		return
	}
	fileSrc := fileDir + CoreFile.Sep + "gin." + CoreFilter.GetNowTimeCarbon().Format("2006010215") + ".log"
	if err := CoreFile.WriteFile(fileSrc, []byte("12312381273_hour")); err != nil {
		t.Error(err)
		return
	}
	fileSrc = fileDir + CoreFile.Sep + "gin." + CoreFilter.GetNowTimeCarbon().Format("20060102") + ".log"
	if err := CoreFile.WriteFile(fileSrc, []byte("124apple_apple_day")); err != nil {
		t.Error(err)
		return
	}
	//执行一次
	runSave()
}

func TestGetList(t *testing.T) {
	dataList, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		Mark:     "",
		IP:       "",
		LogType:  "",
		TimeType: "",
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil {
		if len(dataList) < 1 {
			t.Error("no data list")
		} else {
			newData = dataList[0]
		}
	}
}

func TestGetByID(t *testing.T) {
	if newData.ID < 1 {
		return
	}
	data, err := GetByID(&ArgsGetByID{
		ID: newData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestDownload(t *testing.T) {
	minTime := CoreFilter.GetNowTimeCarbon().SubHour().Format(time.RFC3339)
	maxTime := CoreFilter.GetNowTimeCarbon().Format(time.RFC3339)
	fileName, dataByte, err := Download(&ArgsDownload{
		TimeBetween: CoreSQLTime.DataCoreTime{
			MinTime: minTime,
			MaxTime: maxTime,
		},
	})
	ToolsTest.ReportData(t, err, dataByte)
	//必定没有数据，激活巡逻程序，再次获取
	runDownload()
	fileName, dataByte, err = Download(&ArgsDownload{
		TimeBetween: CoreSQLTime.DataCoreTime{
			MinTime: minTime,
			MaxTime: maxTime,
		},
	})
	ToolsTest.ReportData(t, err, dataByte)
	if err != nil {
		t.Error(err)
	} else {
		if len(dataByte) < 1 {
			t.Error("no data byte < 1")
		} else {
			t.Log("fileName: ", fileName, ", download byte: ", dataByte)
		}
	}
}
