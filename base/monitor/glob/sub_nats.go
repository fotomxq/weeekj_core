package BaseMonitorGlob

import (
	"fmt"
	BaseSystemMission "github.com/fotomxq/weeekj_core/v5/base/system_mission"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"strings"
)

func subNats() {
	//任务调度，监控系统
	BaseSystemMission.ReginSub(&runSysM, subNatsRun)
}

func subNatsRun() {
	//日志
	logAppend := "base monitor glob run, "
	//捕捉异常
	defer func() {
		//进度控制
		runSysM.Bind.UpdateNextAtFutureSec(runSec)
		//跳出处理
		if r := recover(); r != nil {
			runSysM.Update(fmt.Sprint("发生错误: ", r), "run.error", 0)
			CoreLog.Error(logAppend, r)
		}
	}()
	//进度监听
	runSysM.Start("开始分析", "start", 5)
	//初始化数据集合
	var rawData DataGlob
	//获取机器信息
	hostInfo, _ := host.Info()
	rawData.HostName = hostInfo.Hostname
	//获取CPU性能
	rawData.CPUCountsCoreTotal, _ = cpu.Counts(false)
	rawData.CPUCountsTotal, _ = cpu.Counts(true)
	cpuTimeInfo, _ := cpu.Times(false)
	for _, v := range cpuTimeInfo {
		if v.CPU == "cpu-total" {
			rawData.CPUTimeIdle = v.Idle
			rawData.CPUTimeUser = v.User
			rawData.CPUTimeSystem = v.System
			rawData.CPUPercent = rawData.CPUTimeUser / (rawData.CPUTimeIdle + rawData.CPUTimeUser + rawData.CPUTimeSystem)
			break
		}
	}
	//进程
	runSysM.Update("分析完成CPU", "cpu", 1)
	//获取内存性能
	memInfo, _ := mem.VirtualMemory()
	rawData.MemTotal = memInfo.Total
	rawData.MemFree = memInfo.Free
	rawData.MemUse = memInfo.Used
	rawData.MemPercent = memInfo.UsedPercent
	memSwapInfo, _ := mem.SwapMemory()
	rawData.MemSwapTotal = memSwapInfo.Total
	rawData.MemSwapFree = memSwapInfo.Free
	rawData.MemSwapUse = memSwapInfo.Used
	rawData.MemSwapPercent = memSwapInfo.UsedPercent
	//进程
	runSysM.Update("分析完成内存", "mem", 1)
	//磁盘
	diskInfo, _ := disk.Partitions(true)
	for _, v := range diskInfo {
		//获取磁盘的具体信息
		vDiskInfo, _ := disk.Usage(v.Device)
		rawData.Disk = append(rawData.Disk, DataGlobDisk{
			Name:        v.Device,
			Total:       vDiskInfo.Total,
			Free:        vDiskInfo.Free,
			Used:        vDiskInfo.Used,
			UsedPercent: vDiskInfo.UsedPercent,
		})
	}
	diskIOInfo, _ := disk.IOCounters()
	for k, v := range diskIOInfo {
		for k2, v2 := range rawData.Disk {
			if k != v2.Name {
				continue
			}
			//找到对应的数据开始处理
			v2.ReadCount = v.ReadCount
			v2.WriteCount = v.WriteCount
			v2.ReadBytes = v.ReadBytes
			v2.WriteBytes = v.WriteBytes
			v2.ReadTime = v.ReadTime
			v2.WriteTime = v.WriteTime
			rawData.Disk[k2] = v2
		}
	}
	//进程
	runSysM.Update("分析完成磁盘", "disk", 1)
	//网络
	netInfo, _ := net.IOCounters(false)
	for _, v := range netInfo {
		if v.Name != "all" {
			continue
		}
		rawData.NetSendBytes = v.BytesSent
		rawData.NetReadBytes = v.BytesRecv
		rawData.NetSendPackCount = v.PacketsSent
		rawData.NetReadPackCount = v.PacketsRecv
	}
	//进程
	runSysM.Update("分析完成网络", "net", 1)
	//进程
	subNatsRunPID(&rawData)
	//进程
	runSysM.Update("分析完成进程", "pid", 1)
	//整理进程数据
	rawData.ServiceStep = fmt.Sprint(rawData.NowPID.PID)
	//写入集合
	Router2SystemConfig.MainCache.SetStruct(fmt.Sprint(cacheDataKey, ".", rawData.ServiceStep), rawData, 3600)
	//进程完结
	runSysM.Finish()
}

func subNatsRunPID(rawData *DataGlob) {
	//日志
	logAppend := "base monitor glob run, pid, "
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error(logAppend, r)
		}
	}()
	rawData.OtherPIDs = []DataGlobPID{}
	pidInfo, _ := process.Processes()
	for _, v := range pidInfo {
		//组装数据
		appendData := DataGlobPID{
			Name:          "",
			PID:           v.Pid,
			CreateTime:    0,
			CPUPercent:    0,
			MemoryPercent: 0,
			IOReadCount:   0,
			IOWriteCount:  0,
			IOReadBytes:   0,
			IOWriteBytes:  0,
		}
		appendData.Name, _ = v.Name()
		appendData.CreateTime, _ = v.CreateTime()
		appendData.CPUPercent, _ = v.CPUPercent()
		appendData.MemoryPercent, _ = v.MemoryPercent()
		vIO, _ := v.IOCounters()
		if vIO != nil {
			appendData.IOReadCount = vIO.ReadCount
			appendData.IOReadBytes = vIO.ReadBytes
			appendData.IOWriteCount = vIO.WriteCount
			appendData.IOWriteBytes = vIO.WriteBytes
		}
		//写入其他进程序列
		rawData.OtherPIDs = append(rawData.OtherPIDs, appendData)
		//检查是否为当前pid
		if Router2SystemConfig.ServicePIDName == "" {
			//检查是否存在pid名称？
			if strings.Contains(appendData.Name, Router2SystemConfig.ServicePIDName) {
				rawData.NowPID = appendData
			}
		}
	}
}
