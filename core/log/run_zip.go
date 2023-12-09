package CoreLog

import (
	"context"
	"fmt"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	"github.com/golang-module/carbon"
	"github.com/mholt/archiver/v4"
	"os"
	"strings"
)

// 压缩旧的日志
func runZip() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			Error("log run zip, ", r)
		}
	}()
	//debug模式下跳出
	if debugOn {
		return
	}
	//准备压缩的文件列表
	zipFileList := map[string]string{}
	//当前时间，根据当前时间比对非当前时间的目录，对其进行自动压缩处理
	// 如果目录不存在文件，直接删除目录
	folderList, err := CoreFile.GetFileList(logDir, []string{}, true)
	if err != nil {
		return
	}
	for _, vFolder := range folderList {
		//跳过文件
		if CoreFile.IsFile(vFolder) {
			continue
		} else {
			//检查目录下子文件
			zipFileList = runZipLoadChildFile(vFolder, zipFileList)
		}
	}
	//开始压缩
	files, err := archiver.FilesFromDisk(nil, zipFileList)
	if err != nil {
		Error("log run zip, create archiver files from disk, ", err)
		return
	}
	out, err := os.Create(fmt.Sprint(logFileDir, CoreFile.Sep, CoreFilter.GetNowTime().Format("20060102_1504"), ".tar.gz"))
	if err != nil {
		Error("log run zip, create os, ", err)
		return
	}
	defer func() {
		_ = out.Close()
	}()
	format := archiver.CompressedArchive{
		Compression: archiver.Gz{},
		Archival:    archiver.Tar{},
	}
	err = format.Archive(context.Background(), out, files)
	if err != nil {
		Error("log run zip, create archive, ", err)
		return
	}
	//清理所有文件
	for _, v := range zipFileList {
		err = CoreFile.DeleteF(v)
		if err != nil {
			Error("log run zip, delete file src: ", v, ", err: ", err)
		}
	}
}

// 读取子文件数据
func runZipLoadChildFile(fileSrc string, saveFileList map[string]string) map[string]string {
	if CoreFile.IsFile(fileSrc) {
		if runZipIsNowFile(fileSrc) {
			return saveFileList
		}
		saveFileList[fileSrc] = ""
	} else {
		fileList, err := CoreFile.GetFileList(fileSrc, []string{}, true)
		if err == nil && len(fileList) > 0 {
			for _, v := range fileList {
				if runZipIsNowFile(fileSrc) {
					continue
				}
				saveFileList = runZipLoadChildFile(v, saveFileList)
			}
		}
	}
	return saveFileList
}

// runZipIsNowFile 是否为当前时间文件
func runZipIsNowFile(fileSec string) (b bool) {
	//跳过非log文件格式
	fileInfos, err := CoreFile.GetFileNames(fileSec)
	if err != nil {
		Error("base log run, get file info src: ", fileSec)
		return
	}
	if fileInfos["type"] != "log" {
		return
	}
	//分析日志名称的结构类型
	names := strings.Split(fileInfos["only-name"], ".")
	if len(names) < 1 {
		//文件名称无法识别，非本系统日志，跳过
		return
	}
	//识别存储类型
	var insertTimeAt carbon.Carbon
	insertTimeAt = carbon.CreateFromDate(CoreFilter.GetNowTimeCarbon().Year(), CoreFilter.GetNowTimeCarbon().Month(), CoreFilter.GetNowTimeCarbon().Day())
	//如果是当前时间，则不会删除改文件
	isNowTime := false
	// 时间格式
	if len(names[1]) == 12 {
		//小时级别
		year, err := CoreFilter.GetIntByString(fmt.Sprint(names[1][0:4]))
		if err != nil {
			return
		}
		month, err := CoreFilter.GetIntByString(fmt.Sprint(names[1][4:6]))
		if err != nil {
			return
		}
		day, err := CoreFilter.GetIntByString(fmt.Sprint(names[1][6:8]))
		if err != nil {
			return
		}
		hour, err := CoreFilter.GetIntByString(fmt.Sprint(names[1][8:10]))
		if err != nil {
			return
		}
		min, err := CoreFilter.GetIntByString(fmt.Sprint(names[1][10:12]))
		if err != nil {
			return
		}
		if insertTimeAt.Year() == year && insertTimeAt.Month() == month && insertTimeAt.Day() == day && insertTimeAt.Hour() == hour && insertTimeAt.Minute() == hour {
			isNowTime = true
		}
		insertTimeAt = insertTimeAt.SetYear(year)
		insertTimeAt = insertTimeAt.SetMonth(month)
		insertTimeAt = insertTimeAt.SetDay(day)
		insertTimeAt = insertTimeAt.SetHour(hour)
		insertTimeAt = insertTimeAt.SetMinute(min)
		insertTimeAt = insertTimeAt.SetSecond(0)
	} else if len(names[1]) == 10 {
		//小时级别
		year, err := CoreFilter.GetIntByString(fmt.Sprint(names[1][0:4]))
		if err != nil {
			return
		}
		month, err := CoreFilter.GetIntByString(fmt.Sprint(names[1][4:6]))
		if err != nil {
			return
		}
		day, err := CoreFilter.GetIntByString(fmt.Sprint(names[1][6:8]))
		if err != nil {
			return
		}
		hour, err := CoreFilter.GetIntByString(fmt.Sprint(names[1][8:10]))
		if err != nil {
			return
		}
		if insertTimeAt.Year() == year && insertTimeAt.Month() == month && insertTimeAt.Day() == day && insertTimeAt.Hour() == hour {
			isNowTime = true
		}
		insertTimeAt = insertTimeAt.SetYear(year)
		insertTimeAt = insertTimeAt.SetMonth(month)
		insertTimeAt = insertTimeAt.SetDay(day)
		insertTimeAt = insertTimeAt.SetHour(hour)
		insertTimeAt = insertTimeAt.SetMinute(0)
		insertTimeAt = insertTimeAt.SetSecond(0)
	} else if len(names[1]) == 8 {
		//日级别
		year, err := CoreFilter.GetIntByString(fmt.Sprint(names[1][0:4]))
		if err != nil {
			return
		}
		month, err := CoreFilter.GetIntByString(fmt.Sprint(names[1][4:6]))
		if err != nil {
			return
		}
		day, err := CoreFilter.GetIntByString(fmt.Sprint(names[1][6:8]))
		if err != nil {
			return
		}
		if insertTimeAt.Year() == year && insertTimeAt.Month() == month && insertTimeAt.Day() == day {
			isNowTime = true
		}
		insertTimeAt = insertTimeAt.SetYear(year)
		insertTimeAt = insertTimeAt.SetMonth(month)
		insertTimeAt = insertTimeAt.SetDay(day)
		insertTimeAt = insertTimeAt.SetHour(0)
		insertTimeAt = insertTimeAt.SetMinute(0)
		insertTimeAt = insertTimeAt.SetSecond(0)
	} else {
		//无法识别长度，退出
		return
	}
	b = isNowTime
	return
}
