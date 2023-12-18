package BaseIPAddr

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

var (
	isRun = false
)

func TestInit(t *testing.T) {
	if isRun {
		return
	}
	isRun = true
	ToolsTest.Init(t)
	if err := Init(); err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestSetOpenBan(t *testing.T) {
	SetOpenBan(true)
}

func TestSetOpenWhite(t *testing.T) {
	SetOpenWhite(true)
}

func TestSetBan(t *testing.T) {
	if err := SetBan(&ArgsSetBan{
		IP: "0.0.0.1", B: true, ExpireTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTime().Add(time.Second * 30)),
	}); err != nil {
		t.Error(err)
	}
	if !CheckIsBan("0.0.0.1") {
		t.Error("check is ban is false")
	}
	SetOpenBan(false)
	if CheckIsBan("0.0.0.1") {
		t.Error("change open ban is false, but ip is ban")
	}
	SetOpenBan(true)
}

func TestSetWhite(t *testing.T) {
	if err := SetWhite(&ArgsSetWhite{
		IP: "0.0.0.2", B: true, ExpireTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTime().Add(time.Second * 30)),
	}); err != nil {
		t.Error(err)
	}
	if !CheckIsWhite("0.0.0.2") {
		t.Error("check is white is false")
	}
	SetOpenWhite(false)
	if !CheckIsWhite("0.0.0.2") {
		t.Error("change open white is false, but ip is not white")
	}
	SetOpenWhite(true)
	if err := SetWhite(&ArgsSetWhite{
		IP: "192.168.1.252", B: true, ExpireTime: "2021-01-25T16:00:00.000Z",
	}); err != nil {
		t.Error(err)
	}
	if err := SetWhite(&ArgsSetWhite{
		IP: "192.168.1.252", B: true, ExpireTime: "2021-03-25T16:00:00.000Z",
	}); err != nil {
		t.Error(err)
	}
}

func TestCheckAuto(t *testing.T) {
	if CheckAuto("0.0.0.1") {
		t.Error("check auto 0.0.0.1 is true")
	}
	if !CheckAuto("0.0.0.2") {
		t.Error("check auto 0.0.0.2 is false")
	}
}

// 检查一个不存在的IP
func TestCheckIsBanNotHave(t *testing.T) {
	if CheckAuto("123.222.312.442") {
		t.Error("check auto 0.0.0.1 is true")
	}
}

func TestGetList(t *testing.T) {
	dataList, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		Search: "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(dataList, dataCount)
	}
}

func TestCheckIP(t *testing.T) {
	if err := checkIP("0.0.0.1"); err != nil {
		t.Error(err)
	}
	if err := checkIP("7126378612783"); err == nil {
		t.Error("ip is error, but check is ok, ip: 7126378612783")
	}
}

func TestCheckIsBan(t *testing.T) {
	if !CheckIsBan("0.0.0.1") {
		t.Error("check is ban is false")
	}
}

func TestCheckIsWhite(t *testing.T) {
	if !CheckIsWhite("0.0.0.2") {
		t.Error("check is white is false")
	}
}

func TestSetIP(t *testing.T) {
	if err := SetIP(&ArgsSetIP{
		IP: "*", IsMatch: true, IsBan: true, IsWhite: true, ExpireTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTime().Add(time.Second * 30)),
	}); err != nil {
		t.Error(err)
	}
}

func TestGetAddressByIP(t *testing.T) {
	address, err := GetAddressByIP(&ArgsGetAddressByIP{
		IP: "0.0.0.0",
	})
	if err != nil {
		t.Error(err)
	}
	if address != "" {
		t.Error("address is empty")
	}
}

func TestClearIP(t *testing.T) {
	if err := ClearIP(&ArgsClearIP{
		IP: "0.0.0.1",
	}); err != nil {
		t.Error(err)
	}
	if CheckIsBan("0.0.0.1") {
		t.Error("check is white is true")
	}
}

func TestClearAll(t *testing.T) {
	if err := ClearAll(); err != nil {
		t.Error(err)
	}
	if CheckIsWhite("0.0.0.2") {
		t.Error("check is white is true")
	}
}
