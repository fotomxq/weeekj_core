package Router2SystemInit

import (
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	"net/http"
	"os"
	"runtime/pprof"
	"time"
)

//Run
// 访问地址: http://localhost:9999/debug/pprof/
/**
profile（CPU Profiling）: $HOST/debug/pprof/profile，默认进行 30s 的 CPU Profiling，得到一个分析用的 profile 文件
allocs 查看过去所有内存分配样本，访问路径为 /debug/pprof/allocs
cmdline 当前程序命令行的完整调用路径
block（Block Profiling）：$HOST/debug/pprof/block，查看导致阻塞同步的堆栈跟踪
goroutine：$HOST/debug/pprof/goroutine，查看当前所有运行的 goroutines 堆栈跟踪
heap（Memory Profiling）: $HOST/debug/pprof/heap，查看活动对象的内存分配情况
mutex（Mutex Profiling）：$HOST/debug/pprof/mutex，查看导致互斥锁的竞争持有者的堆栈跟踪
threadcreate：$HOST/debug/pprof/threadcreate，查看创建新OS线程的堆栈跟踪
*/
func runDebug() {
	//go runCPU()
	err := http.ListenAndServe("0.0.0.0:9999", nil)
	if err != nil {
		fmt.Println("run debug http failed.")
	}
}

func runDebugCPU() {
	for {
		f, err := os.OpenFile(fmt.Sprint(CoreFile.BaseDir()+CoreFile.Sep, "/data/debug/files/cpu.profile"), os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			time.Sleep(time.Hour * 1)
			continue
		}
		//defer f.Close()
		err = pprof.StartCPUProfile(f)
		if err != nil {
			time.Sleep(time.Hour * 1)
			continue
		}
		defer pprof.StopCPUProfile()
		n := 10
		for i := 1; i <= 5; i++ {
			fmt.Printf("fib(%d)=%d\n", n, runDebugCPUFib(n))
			n += 3 * i
		}
		time.Sleep(time.Hour * 1)
	}
}

func runDebugCPUFib(n int) int {
	if n <= 1 {
		return 1
	}
	return runDebugCPUFib(n-1) + runDebugCPUFib(n-2)
}
