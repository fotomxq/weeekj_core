package CoreLog

import (
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type logConfig struct {
	//日志句柄
	LogHandle *logrus.Logger
	//上一次生成的目录路径
	// 切换目录时，会压缩目录
	LastDirSrc string
	//上一次日志的文件路径
	// 切换时，如果上一次文件尺寸为0会自动删除
	LastFileSrc string
	//日志前缀
	// eg: gin.
	Prefix string
	//日志文件操作句柄
	fileHandle *os.File
	// 第二个句柄，避免混淆写入
	fileHandle2 *os.File
	// 当前所在句柄
	fileHandleNow int
}

// Init 初始化
func (t *logConfig) Init() {
	t.LogHandle = logrus.New()
	t.LogHandle.SetFormatter(&logrus.JSONFormatter{})
}

// MakeNextHandle 生成下一轮日志句柄
func (t *logConfig) MakeNextHandle() {
	if debugOn {
		t.LogHandle.SetOutput(os.Stdout)
		t.LogHandle.SetNoLock()
	} else {
		//删除上一个为空的文件
		allowHigh := false
		fileMore3M := false
		if fileSize, err := CoreFile.GetFileSize(t.LastFileSrc); err == nil {
			if fileSize < 1 {
				_ = CoreFile.DeleteF(t.LastFileSrc)
			}
			if fileSize >= 1*1024*1024 {
				allowHigh = true
			}
			if fileSize >= 3*1024*1024 {
				fileMore3M = true
			}
		}
		//检查是否高频模式
		if allowHigh {
			//生成下一个目录路径
			t.LastDirSrc = logDir + CoreFile.Sep + CoreFilter.GetNowTime().Format("200601"+CoreFile.Sep+"02")
			if fileMore3M {
				//生成下一个文件句柄
				t.LastFileSrc = t.LastDirSrc + CoreFile.Sep + t.Prefix + CoreFilter.GetNowTime().Format("200601021504") + ".log"
			} else {
				//生成下一个文件句柄
				t.LastFileSrc = t.LastDirSrc + CoreFile.Sep + t.Prefix + CoreFilter.GetNowTime().Format("2006010215") + ".log"
			}
		} else {
			//生成下一个目录路径
			t.LastDirSrc = logDir + CoreFile.Sep + CoreFilter.GetNowTime().Format("200601")
			//生成下一个文件句柄
			t.LastFileSrc = t.LastDirSrc + CoreFile.Sep + t.Prefix + CoreFilter.GetNowTime().Format("20060102") + ".log"
		}
		if !CoreFile.IsFolder(t.LastDirSrc) {
			if err := CoreFile.CreateFolder(t.LastDirSrc); err != nil {
				time.Sleep(time.Second * 10)
				return
			}
		}
		//启动新的句柄
		var err error
		var fileHandle *os.File
		fileHandle, err = os.OpenFile(t.LastFileSrc, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
		if err != nil {
			time.Sleep(time.Second * 10)
			return
		}
		//自动切换到下一个句柄
		if t.fileHandleNow == 1 {
			//关闭旧的句柄
			_ = t.fileHandle2.Close()
			//覆盖句柄
			t.fileHandle2 = fileHandle
			t.LogHandle.SetOutput(t.fileHandle2)
			t.LogHandle.SetNoLock()
			t.fileHandleNow = 2
			//关闭旧的句柄
			_ = t.fileHandle.Close()
		} else {
			//关闭旧的句柄
			_ = t.fileHandle.Close()
			//覆盖句柄
			t.fileHandle = fileHandle
			t.LogHandle.SetOutput(t.fileHandle)
			t.LogHandle.SetNoLock()
			t.fileHandleNow = 1
			//关闭旧的句柄
			_ = t.fileHandle2.Close()
		}
	}
}
