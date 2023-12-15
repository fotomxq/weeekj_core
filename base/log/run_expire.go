package BaseLog

import (
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 将数据库数据，归档到文件内
func runExpire() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("base log expire run, ", r)
		}
	}()
	//构建存储目录
	fileLogDir = fmt.Sprint(CoreFile.BaseDir(), CoreFile.Sep, "log_file")
	if !CoreFile.IsFolder(fileLogDir) {
		if err := CoreFile.CreateFolder(fileLogDir); err != nil {
			CoreLog.Error("base log expire run, create folder src: ", fileLogDir, ", err: ", err)
			return
		}
	}
	//构建临时目录
	fileLogDirTemp := fmt.Sprint(fileLogDir, CoreFile.Sep, "temp")
	if !CoreFile.IsFolder(fileLogDirTemp) {
		if err := CoreFile.CreateFolder(fileLogDirTemp); err != nil {
			CoreLog.Error("base log expire run, create folder src: ", fileLogDirTemp, ", err: ", err)
			return
		}
	} else {
		//删除临时目录
		if err := CoreFile.DeleteF(fileLogDirTemp); err != nil {
			CoreLog.Error("base log expire run, delete folder src: ", fileLogDirTemp, ", err: ", err)
			return
		}
		if err := CoreFile.CreateFolder(fileLogDirTemp); err != nil {
			CoreLog.Error("base log expire run, create folder src: ", fileLogDirTemp, ", err: ", err)
			return
		}
	}
	//开始时间
	startAt := ""
	endAt := ""
	//找出需要归档的数据
	var step int64 = 0
	for {
		var dataList []FieldsLog
		if err := Router2SystemConfig.MainDB.Select(&dataList, fmt.Sprint("SELECT id, create_at, mark, ip, log_type, time_type, size, contents FROM core_log WHERE create_at < $1 LIMIT 5 OFFSET ", step), CoreFilter.GetNowTimeCarbon().SubDays(3).Time); err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		//写入文件
		for _, v := range dataList {
			if startAt == "" {
				startAt = dataList[0].CreateAt.Format("2006010215")
			}
			endAt = dataList[len(dataList)-1].CreateAt.Format("2006010215")
			vSrc := fmt.Sprint(fileLogDirTemp, CoreFile.Sep, dataList[0].CreateAt.Format("2006010215"), "_", dataList[len(dataList)-1].CreateAt.Format("2006010215"), "_", step, ".log")
			if CoreFile.IsExist(vSrc) {
				if err := CoreFile.DeleteF(vSrc); err != nil {
					CoreLog.Error("base log expire run, delete file: ", vSrc, ", err: ", err)
					return
				}
			}
			if err := CoreFile.WriteFile(vSrc, []byte(v.Contents)); err != nil {
				CoreLog.Error("base log expire run, write file: ", vSrc, ", err: ", err)
				return
			}
			//删除数据
			_, err := CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_log", "id", map[string]interface{}{
				"id": v.ID,
			})
			if err != nil {
				CoreLog.Error("base log expire run, delete log, id: ", v.ID, ", err: ", err)
				continue
			}
		}
		step += 5
	}
	if startAt == "" || endAt == "" {
		return
	}
	//构建zip文件路径
	zipFileSrc := fmt.Sprint(fileLogDir, CoreFile.Sep, startAt, "_", endAt, ".zip")
	//如果存在重复的文件，则删除
	if CoreFile.IsExist(zipFileSrc) {
		if err := CoreFile.DeleteF(zipFileSrc); err != nil {
			CoreLog.Error("base log expire run, delete file: ", zipFileSrc, ", err: ", err)
			return
		}
	}
	//压缩目录
	if err := CoreFile.ZipDir(fileLogDirTemp, zipFileSrc); err != nil {
		CoreLog.Error("base log expire run, zip temp: ", zipFileSrc, ", err: ", err)
		return
	}
	//删除临时目录
	if err := CoreFile.DeleteF(fileLogDirTemp); err != nil {
		CoreLog.Error("base log expire run, delete file: ", fileLogDirTemp, ", err: ", err)
		return
	}
}
