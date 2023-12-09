package BaseLog

import (
	"fmt"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/golang-module/carbon"
	"strings"
)

// 将日志存储到数据库
func runSave() {
	//检索log目录
	dirSrc := CoreFile.BaseSrc + CoreFile.Sep + "log"
	//获取子目录
	// 存在2种模式，低速模式、高速模式。前者一层目录，后者二层目录
	// 1 -> 202104/gin.20210406.log
	// 2 -> 202104/06/gin.2021040621.log
	// 此处不管几层，直接递归查询
	runSaveDir(dirSrc)
}

func runSaveDir(dirSrc string) {
	files, err := CoreFile.GetFileList(dirSrc, []string{}, true)
	if err != nil {
		//没有文件，反馈
		return
	}
	for _, v := range files {
		if CoreFile.IsFolder(v) {
			runSaveDir(v)
			continue
		}
		runSaveFile(Router2SystemConfig.ServerName, Router2SystemConfig.ServerIP, v)
	}
}

func runSaveFile(mark, ip, fileSrc string) {
	//捕捉异常，大部分异常为map指向失败的异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("base log run, save file, ", r)
		}
	}()
	//跳过超大文件
	fileSize, err := CoreFile.GetFileSize(fileSrc)
	if err != nil {
		CoreLog.Error("base log run, get file size: ", fileSrc, ", err: ", err)
		return
	}
	if fileSize > 4*1024*1024 {
		return
	}
	//跳过非log文件格式
	fileInfos, err := CoreFile.GetFileNames(fileSrc)
	if err != nil {
		CoreLog.Error("base log run, get file info src: ", fileSrc)
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
	var timeType string
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
		timeType = "YYYY-MM-DD_H-M"
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
		timeType = "YYYY-MM-DD_HH"
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
		timeType = "YYYY-MM-DD"
	} else {
		//无法识别长度，退出
		return
	}
	//可能需跳过当下的数据
	skipNow, err := Router2SystemConfig.Cfg.Section("core").Key("log_save_now").Bool()
	if err != nil {
		skipNow = true
	}
	if skipNow {
		if insertTimeAt.Time.Unix() >= CoreFilter.GetNowTimeCarbon().SubHour().Time.Unix() {
			return
		}
	}
	//开始装载数据，写入数据库
	fileByte, err := CoreFile.LoadFile(fileSrc)
	if err != nil {
		CoreLog.Error("base log run, load file: ", fileSrc, ", err: ", err)
		return
	}
	if fileSize < 1 {
		if !isNowTime {
			if err := CoreFile.DeleteF(fileSrc); err != nil {
				CoreLog.Error("base log run, create new data, delete file: ", fileSrc, ", err: ", err)
				//继续运行后续
			}
		}
		return
	}
	// 检查该时间是否存在数据？
	haveData := false
	var data FieldsLog
	if err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM core_log WHERE create_at = $1 AND mark = $2 AND ip = $3 AND log_type = $4", insertTimeAt.Time, mark, ip, names[0]); err == nil {
		haveData = data.ID > 0
	}
	if haveData {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_log SET size = :size, contents = :contents WHERE id = :id", map[string]interface{}{
			"id":        data.ID,
			"create_at": insertTimeAt.Time,
			"size":      fileSize,
			"contents":  string(fileByte),
		})
	} else {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO core_log (create_at, mark, ip, log_type, time_type, size, contents) VALUES (:create_at, :mark, :ip, :log_type, :time_type, :size, :contents)", map[string]interface{}{
			"create_at": insertTimeAt.Time,
			"mark":      mark,
			"ip":        ip,
			"log_type":  names[0],
			"time_type": timeType,
			"size":      fileSize,
			"contents":  string(fileByte),
		})
	}
	if err != nil {
		CoreLog.Error("base log run, create or update data, ", err)
		return
	}
	if !isNowTime {
		if err := CoreFile.DeleteF(fileSrc); err != nil {
			CoreLog.Error("base log run, create new data, delete file: ", fileSrc, ", err: ", err)
			return
		}
	}
	if err != nil {
		CoreLog.Error("base log run, create new data, ", err)
		return
	}
}
