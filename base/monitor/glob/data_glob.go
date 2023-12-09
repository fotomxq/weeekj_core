package BaseMonitorGlob

type DataGlob struct {
	//服务标识码
	ServiceMark string `json:"serviceMark"`
	//分布式分支，采用进程PID
	ServiceStep string `json:"serviceStep"`
	//host
	HostName string `json:"hostName"`
	//CPU信息
	CPUCountsCoreTotal int `json:"cpuCountsCoreTotal"`
	CPUCountsTotal     int `json:"cpuCountsTotal"`
	//CPU占用时间
	CPUTimeIdle   float64 `json:"cpuTimeIdle"`
	CPUTimeUser   float64 `json:"cpuTimeUser"`
	CPUTimeSystem float64 `json:"cpuTimeSystem"`
	//CPU占用率
	CPUPercent float64 `json:"cpuPercent"`
	//内存
	MemTotal   uint64  `json:"memTotal"`
	MemUse     uint64  `json:"memUse"`
	MemFree    uint64  `json:"memFree"`
	MemPercent float64 `json:"memPercent"`
	//交换区内存
	MemSwapTotal   uint64  `json:"memSwapTotal"`
	MemSwapUse     uint64  `json:"memSwapUse"`
	MemSwapFree    uint64  `json:"memSwapFree"`
	MemSwapPercent float64 `json:"memSwapPercent"`
	//磁盘集
	Disk []DataGlobDisk `json:"disk"`
	//网络情况
	NetSendBytes     uint64 `json:"netSendBytes"`
	NetReadBytes     uint64 `json:"netReadBytes"`
	NetSendPackCount uint64 `json:"netSendPackCount"`
	NetReadPackCount uint64 `json:"netReadPackCount"`
	//进程信息
	NowPID DataGlobPID `json:"nowPID"`
	//其他进程信息
	OtherPIDs []DataGlobPID `json:"otherPIDs"`
}

type DataGlobDisk struct {
	//磁盘名称
	Name string `json:"name"`
	//总空间
	Total uint64 `json:"total"`
	//剩余空间
	Free uint64 `json:"free"`
	//已经使用
	Used uint64 `json:"used"`
	//使用占比
	UsedPercent float64 `json:"usedPercent"`
	//读次数
	ReadCount uint64 `json:"readCount"`
	//写次数
	WriteCount uint64 `json:"writeCount"`
	//读总量
	ReadBytes uint64 `json:"readBytes"`
	//写总量
	WriteBytes uint64 `json:"writeBytes"`
	//读时间
	ReadTime uint64 `json:"readTime"`
	//写时间
	WriteTime uint64 `json:"writeTime"`
}

type DataGlobPID struct {
	//名称
	Name string `json:"name"`
	//PID
	PID int32 `json:"pid"`
	//创建时间
	CreateTime int64 `json:"createTime"`
	//CPU占用率
	CPUPercent float64 `json:"cpuPercent"`
	//内存占用率
	MemoryPercent float32 `json:"memoryPercent"`
	//读次数
	IOReadCount uint64 `json:"ioReadCount"`
	//写次数
	IOWriteCount uint64 `json:"ioWriteCount"`
	//读总量
	IOReadBytes uint64 `json:"ioReadBytes"`
	//写总量
	IOWriteBytes uint64 `json:"ioWriteBytes"`
}
