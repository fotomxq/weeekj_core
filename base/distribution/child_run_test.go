package BaseDistribution

import (
	"testing"
	"time"
)

func TestInit3(t *testing.T) {
	TestInit(t)
}

func TestSetChildRun(t *testing.T) {
	if err := SetChildRun(&ArgsSetChildRun{
		Mark: "test", IP: "1.1.1.1", Port: "9000", RunMark: "clear", ExpireAddTime: 150,
	}); err == nil {
		t.Error("set child run failed")
	}
	TestSetChild(t)
	if err := SetChildRun(&ArgsSetChildRun{
		Mark: "test", IP: "1.1.1.1", Port: "9000", RunMark: "clear", ExpireAddTime: 150,
	}); err != nil {
		t.Error(err)
	}
}

func TestGetChildRun(t *testing.T) {
	data, err := GetChildRun(&ArgsGetChildRun{
		Mark: "test", IP: "1.1.1.1", Port: "9000",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestRun(t *testing.T) {
	//启动服务
	go Run()
	//开始封闭测试
	//该部分内容无法进行完整测试，相关测试放到框架的test内部
	//这里仅对非负载处理进行测试，如删除多余数据的巡逻方案
	if err := DeleteService(&ArgsDeleteService{
		Mark: "test",
	}); err != nil {
		t.Error(err)
	}
	//等待3秒
	time.Sleep(time.Second * 3)
	//检查是否还存在子服务数据？
	data, err := GetChildRun(&ArgsGetChildRun{
		Mark: "test", IP: "1.1.1.1", Port: "9000",
	})
	if err == nil && len(data) > 0 {
		t.Error("run data is exist, ", data)
	}
}

//清理测试数据
func TestDeleteAll(t *testing.T) {
	if err := DeleteChild(&ArgsDeleteChild{
		Mark: "test", IP: "1.1.1.1", Port: "9000",
	}); err != nil {
		t.Error(err)
	}
	if err := DeleteService(&ArgsDeleteService{
		Mark: "test",
	}); err != nil {
		t.Error(err)
	}
}
