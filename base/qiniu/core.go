package BaseQiniu

import (
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/robfig/cron"
)

var (
	//定时器
	runTimer      = cron.New()
	runExpireLock = false
	//管理对象
	bucketManager *storage.BucketManager
)
