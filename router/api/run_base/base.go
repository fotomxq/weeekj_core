package RouterAPIRunBase

import (
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BasePedometer "github.com/fotomxq/weeekj_core/v5/base/pedometer"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

var (
	//配置维护监测
	runTime = time.Minute * 10
	//本服务内的配置更新时间
	updateTime int64 = 0
)

// Init 初始化设置
func Init() {
}

// Run 配置主动式维护模块
// 本配置维护主要针对于core层级的模块，该类模块构架不支持内部对modules的base及以上层级访问，所以需外部写入方式更新配置数据
func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("service config run error, ", r)
		}
	}()
	//本函数巡逻程序
	for {
		//如果需要更新再执行，否则延迟
		if updateTime >= BaseConfig.GetLastUpdateTime() {
			//延迟后再次执行
			time.Sleep(runTime)
			continue
		}
		//计数器
		if err := runConfigPedometer(); err != nil {
			CoreLog.Error("service config run error, pedometer config, ", err)
		}
		//更新时间
		updateTime = CoreFilter.GetNowTime().Unix()
		//检查postgresql连接有效性
		if err := Router2SystemConfig.MainDB.Ping(); err != nil {
			if err2 := Router2SystemConfig.LoadPostgres(); err2 != nil {
				CoreLog.Error("service config run error, postgresql connect failed, ", err, ", try connect, ", err2)
			}
		}
		//下一步
		time.Sleep(runTime)
	}
}

// 设置pedometer
func runConfigPedometer() error {
	// 初始化计步器的sms部分
	verificationCodeSMSSafeLimit, err := BaseConfig.GetDataInt("VerificationCodeSMSSafeLimit")
	if err != nil {
		verificationCodeSMSSafeLimit = 30
	}
	verificationCodeSMSSafeLimitExpire, err := BaseConfig.GetDataString("VerificationCodeSMSSafeLimitExpire")
	if err != nil {
		verificationCodeSMSSafeLimitExpire = "24h"
	}
	if err := BasePedometer.SetConfig(&BasePedometer.ArgsSetConfig{
		Mark: "sms", ExpireAdd: verificationCodeSMSSafeLimitExpire, MaxCount: verificationCodeSMSSafeLimit, IsAdd: true,
	}); err != nil {
		return err
	}
	// 初始化计步器安全事件部分
	// 对应逻辑在token头部的TokenAddSafety
	SafetyTokenLimit, err := BaseConfig.GetDataInt("SafetyTokenLimit")
	if err != nil {
		SafetyTokenLimit = 10
	}
	SafetyTokenLimitExpire, err := BaseConfig.GetDataString("SafetyTokenLimitExpire")
	if err != nil {
		SafetyTokenLimitExpire = "3m"
	}
	SafetyIPLimit, err := BaseConfig.GetDataInt("SafetyIPLimit")
	if err != nil {
		SafetyIPLimit = 10
	}
	SafetyIPLimitExpire, err := BaseConfig.GetDataString("SafetyIPLimitExpire")
	if err != nil {
		SafetyIPLimitExpire = "10m"
	}
	SafetyUserLimit, err := BaseConfig.GetDataInt("SafetyUserLimit")
	if err != nil {
		SafetyUserLimit = 30
	}
	SafetyUserLimitExpire, err := BaseConfig.GetDataString("SafetyUserLimitExpire")
	if err != nil {
		SafetyUserLimitExpire = "24h"
	}
	if err := BasePedometer.SetConfig(&BasePedometer.ArgsSetConfig{
		Mark: "safe-token", ExpireAdd: SafetyTokenLimitExpire, MaxCount: SafetyTokenLimit, IsAdd: true,
	}); err != nil {
		return err
	}
	if err := BasePedometer.SetConfig(&BasePedometer.ArgsSetConfig{
		Mark: "safe-ip", ExpireAdd: SafetyIPLimitExpire, MaxCount: SafetyIPLimit, IsAdd: true,
	}); err != nil {
		return err
	}
	if err := BasePedometer.SetConfig(&BasePedometer.ArgsSetConfig{
		Mark: "safe-user", ExpireAdd: SafetyUserLimitExpire, MaxCount: SafetyUserLimit, IsAdd: true,
	}); err != nil {
		return err
	}
	return nil
}
