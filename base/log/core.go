package BaseLog

import (
	"errors"
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/robfig/cron"
	"sync"
)

//本模块主要用于
// 1\将本地化文件存入数据库
// 2\提供检索、预警联动服务

var (
	//预备下载列队
	waitDownloads     []DataWaitDownload
	waitDownloadsLock sync.Mutex
	//定时器
	runTimer          = cron.New()
	runDeleteTempLock = false
	runDownloadLock   = false
	runExpireLock     = false
	runSaveLock       = false
	//日志重新归档目录
	fileLogDir = "log_file"
)

type DataWaitDownload struct {
	//时间戳
	FileSHA string
	//时间间隔
	TimeBetween CoreSQLTime.FieldsCoreTime
}

// 获取列表
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//上报主机名称
	Mark string `db:"mark" json:"mark"`
	//上报主机IP
	IP string `db:"ip" json:"ip"`
	//日志类型
	LogType string `db:"log_type" json:"logType"`
	//时间类型
	TimeType string `db:"time_type" json:"timeType"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

func GetList(args *ArgsGetList) (dataList []FieldsLog, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.Mark != "" {
		where = where + "mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.IP != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "ip = :ip"
		maps["ip"] = args.IP
	}
	if args.LogType != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "log_type = :log_type"
		maps["log_type"] = args.LogType
	}
	if args.TimeType != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "time_type = :time_type"
		maps["time_type"] = args.TimeType
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(contents ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_log",
		"id",
		"SELECT id, create_at, mark, ip, log_type, time_type, size FROM core_log WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "size"},
	)
	return
}

// 获取指定的数据ID
type ArgsGetByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

func GetByID(args *ArgsGetByID) (data FieldsLog, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, mark, ip, log_type, time_type, size, contents FROM core_log WHERE id = $1", args.ID)
	return
}

// 下载指定时间段的所有日志
// 本地会安排构建请求
type ArgsDownload struct {
	//时间段
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

func Download(args *ArgsDownload) (fileName string, dataByte []byte, err error) {
	//构建请求时间的字符串
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	if timeBetween.MinTime.Unix() > timeBetween.MaxTime.Unix() {
		err = errors.New("time between error")
		return
	}
	fileSHA := fmt.Sprint(timeBetween.MinTime.Format("2006010215"), "_", timeBetween.MaxTime.Format("2006010215"))
	fileName = fileSHA + ".zip"
	//检查该数据包是否存在？
	downloadFileSrc := fmt.Sprint(Router2SystemConfig.RootDir, CoreFile.Sep+"temp"+CoreFile.Sep, "log", CoreFile.Sep, fileSHA, ".zip")
	if CoreFile.IsFile(downloadFileSrc) {
		//直接下载该文件
		dataByte, err = CoreFile.LoadFile(downloadFileSrc)
		if err != nil {
			//读取文件失败，说明权限可能存在异常问题，退出
		}
		//尝试删除该文件
		_ = CoreFile.DeleteF(downloadFileSrc)
		return
	}
	//将请求交给列队处理
	waitDownloadsLock.Lock()
	defer waitDownloadsLock.Unlock()
	isFind := false
	for _, v := range waitDownloads {
		if v.FileSHA == fileSHA {
			isFind = true
			break
		}
	}
	if !isFind {
		waitDownloads = append(waitDownloads, DataWaitDownload{
			FileSHA:     fileSHA,
			TimeBetween: timeBetween,
		})
	}
	//反馈
	return
}
