package BaseLog

import (
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

func runDownload() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("base log download run, ", r)
		}
	}()
	for k, v := range waitDownloads {
		//读取数据包
		var dataList []FieldsLog
		if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, mark, ip, log_type, time_type, size, contents FROM core_log WHERE create_at >= $1 AND create_at <= $2", v.TimeBetween.MinTime, v.TimeBetween.MaxTime); err != nil || len(dataList) < 1 {
			runDownloadDeleteKey(k)
			continue
		}
		//构建文件路径
		fileDir := fmt.Sprint(Router2SystemConfig.RootDir, CoreFile.Sep+"temp"+CoreFile.Sep, "log", CoreFile.Sep, v.FileSHA)
		if !CoreFile.IsFolder(fileDir) {
			if err := CoreFile.CreateFolder(fileDir); err != nil {
				CoreLog.Error("base log download run, create folder, ", err)
				runDownloadDeleteKey(k)
				continue
			}
		}
		//打包文件组
		if err := runDownloadSave(fileDir, dataList); err != nil {
			CoreLog.Error("base log download run, save to zip, ", err)
			runDownloadDeleteKey(k)
			continue
		}
		//全部完成后，删除key
		runDownloadDeleteKey(k)
		//避免阻塞，暂停线程1秒
		time.Sleep(time.Second * 1)
	}
}

func runDownloadSave(fileDir string, dataList []FieldsLog) error {
	waitDownloadsLock.Lock()
	defer waitDownloadsLock.Unlock()
	//构建zip文件路径
	zipFileSrc := fmt.Sprint(fileDir, ".zip")
	//如果存在重复的文件，则删除
	if CoreFile.IsExist(zipFileSrc) {
		if err := CoreFile.DeleteF(zipFileSrc); err != nil {
			return err
		}
	}
	//将数据包写入到临时文件夹目录下
	for _, v := range dataList {
		vSrc := fmt.Sprint(fileDir, CoreFile.Sep, v.CreateAt.Format("2006010215"), ".log")
		if CoreFile.IsExist(vSrc) {
			if err := CoreFile.DeleteF(vSrc); err != nil {
				return err
			}
		}
		if err := CoreFile.WriteFile(vSrc, []byte(v.Contents)); err != nil {
			return err
		}
	}
	//压缩目录
	if err := CoreFile.ZipDir(fileDir, zipFileSrc); err != nil {
		return err
	}
	//删除临时目录
	if err := CoreFile.DeleteF(fileDir); err != nil {
		return err
	}
	//反馈成功
	return nil
}

// 删除列队
func runDownloadDeleteKey(key int) {
	waitDownloadsLock.Lock()
	defer waitDownloadsLock.Unlock()
	var newData []DataWaitDownload
	for _, v := range waitDownloads {
		newData = append(newData, v)
	}
	waitDownloads = newData
}
