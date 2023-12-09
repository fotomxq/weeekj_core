package Router2SystemBomb

import (
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/golang-module/carbon"
	"time"
)

// StartBomb 设置定时炸弹时间
/**
1. 系统启动后，将在隐藏系统配置中存储相关内容，并自动完成读写
2. 首次写入时，如果发现为0则不写入；否则写入配置项中
*/
func StartBomb(startAt carbon.Carbon) {
	Router2SystemConfig.BombSec, _ = BaseConfig.GetDataInt64("HideBombSec")
	bombStartMonth, _ := BaseConfig.GetDataInt64("HideBombStartMonth")
	if bombStartMonth < 1 {
		Router2SystemConfig.BombStartMonth = startAt
		_ = BaseConfig.UpdateByMark(&BaseConfig.ArgsUpdateByMark{
			UpdateHash: "",
			Mark:       "HideBombStartMonth",
			Value:      fmt.Sprint(Router2SystemConfig.BombStartMonth.Time.Unix()),
		})
	} else {
		Router2SystemConfig.BombStartMonth = carbon.CreateFromTimestamp(bombStartMonth)
	}
	if Router2SystemConfig.BombSec < 1 {
		Router2SystemConfig.BombSec = 1
		_ = BaseConfig.UpdateByMark(&BaseConfig.ArgsUpdateByMark{
			UpdateHash: "",
			Mark:       "HideBombSec",
			Value:      fmt.Sprint(Router2SystemConfig.BombSec),
		})
	}
}

// GetBomb 获取定时炸弹
// 自带拦截器，将系统阻断沉默，沉默后继续执行后续代码
// 直接调用本方法即可
func GetBomb() {
	if Router2SystemConfig.BombSec < 1 {
		return
	}
	//计算间隔月份
	nowAt := CoreFilter.GetNowTimeCarbon()
	if Router2SystemConfig.BombStartMonth.Time.Unix() > nowAt.Time.Unix() {
		return
	}
	diffMonth := Router2SystemConfig.BombStartMonth.DiffInDaysWithAbs(nowAt) / 30
	//计算睡眠时间
	sleepSec := diffMonth * Router2SystemConfig.BombSec
	Router2SystemConfig.BombSec = diffMonth
	_ = BaseConfig.UpdateByMark(&BaseConfig.ArgsUpdateByMark{
		UpdateHash: "",
		Mark:       "HideBombSec",
		Value:      fmt.Sprint(Router2SystemConfig.BombSec),
	})
	if sleepSec > 9 {
		sleepSec = 9
	}
	//开始执行
	time.Sleep(time.Duration(sleepSec) * time.Second)
}
