package CoreSystemClose

import "sync"

//系统关闭处理模块

var (
	//lockAllRun 阻塞所有定时器模块
	lockAllRun sync.WaitGroup
	//waitCount 启动的进程数量
	waitCount = 0
)

//Wait 等待进程
func Wait() {
	waitCount += 1
	lockAllRun.Add(1)
	lockAllRun.Wait()
}

//Close 关闭系统预备
func Close() {
	for k := 0; k < waitCount; k++ {
		lockAllRun.Done()
	}
}
