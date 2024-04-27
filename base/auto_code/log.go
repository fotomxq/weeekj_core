package BaseAutoCode

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// ArgsCreateNewCode 生成新的编码参数
type ArgsCreateNewCode struct {
	//系统模块标识码
	// eq: user_core 用户模块
	ModuleCode string `db:"module_code" json:"moduleCode" check:"code"`
	//模块内分支标识码
	// eq: user_core_core 用户模块的用户表
	BranchCode string `db:"branch_code" json:"branchCode" check:"code"`
}

// CreateNewCode 生成新的编码
// rawParams 本次计划提交的数据包，如果配置包含自定义规则，将从该数据包中取数据
func CreateNewCode(args *ArgsCreateNewCode, rawParams any) (newCode string, isReplaceData bool) {
	//进程锁定
	lockLog(args.ModuleCode, args.BranchCode)
	defer func() {
		unLockLog(args.ModuleCode, args.BranchCode)
	}()
	//获取配置
	configData := getConfigByCode(args.ModuleCode, args.BranchCode)
	if configData.ID < 1 {
		return
	}
	//检查是否启用
	if !configData.IsEnable {
		return
	}
	//生成编码前缀
	newCode = fmt.Sprint(configData.Prefix)
	//根据ID生成序号
	if configData.AutoNumber {
		newNum := getLogCount(args.ModuleCode, args.BranchCode) + 1
		if configData.AutoNumberLen < 1 {
			newCode = fmt.Sprint(newCode, newNum)
		} else {
			//填充序号前置部分
			newCode = fmt.Sprint(newCode, fmt.Sprintf("%0"+fmt.Sprint(configData.AutoNumberLen)+"d", newNum))
		}
	} else {
		if configData.CustomRule == "" {
			newNum := getLogCount(args.ModuleCode, args.BranchCode) + 1
			newCode = fmt.Sprint(newCode, newNum)
		} else {
			//解析自定义规则
			paramsVals := strings.Split(configData.CustomRule, ",")
			//获取数据包
			rawParamsFields := reflect.TypeOf(rawParams).Elem()
			rawParamsVals := reflect.ValueOf(rawParams).Elem()
			step := 0
			for step < rawParamsFields.NumField()-1 {
				//捕捉结构
				vField := rawParamsFields.Field(step)
				vValueType := rawParamsVals.Field(step)
				//下一步
				step += 1
				//找到匹配项
				isFind := false
				for _, vConfigVal := range paramsVals {
					if vConfigVal == vField.Name {
						isFind = true
						break
					}
				}
				if !isFind {
					continue
				}
				vVal := vValueType.Interface()
				//拼接
				newCode = fmt.Sprint(newCode, vVal)
			}
		}
	}
	//检查是否存在或创建日志
	if configData.IsLog {
		var haveData bool
		var err error
		if configData.IsGlobalUnique {
			haveData, _, err = createLog(&argsCreateLog{
				ModuleCode: "",
				BranchCode: "",
				ConfigID:   -1,
				Code:       newCode,
			})
			if haveData {
				isReplaceData = true
				return
			}
			if err != nil {
				return
			}
		}
		if configData.IsBranchUnique {
			haveData, _, err = createLog(&argsCreateLog{
				ModuleCode: configData.ModuleCode,
				BranchCode: configData.BranchCode,
				ConfigID:   configData.ID,
				Code:       newCode,
			})
			if haveData {
				isReplaceData = true
				return
			}
			if err != nil {
				return
			}
		}
	}
	//反馈
	return
}

// getLogCount 获取日志总数
func getLogCount(moduleCode string, branchCode string) (count int64) {
	count = logDB.Analysis().Count("module_code = $1 AND branch_code = $2", moduleCode, branchCode)
	return
}

// argsCreateLog 创建日志
type argsCreateLog struct {
	//系统模块标识码
	// eq: user_core 用户模块
	ModuleCode string `db:"module_code" json:"moduleCode" check:"code"`
	//模块内分支标识码
	// eq: user_core_core 用户模块的用户表
	BranchCode string `db:"branch_code" json:"branchCode" check:"code"`
	//采用配置ID
	ConfigID int64 `db:"config_id" json:"configId" check:"id"`
	//编码
	Code string `db:"code" json:"code" check:"code"`
}

// createLog 创建日志参数
func createLog(args *argsCreateLog) (haveData bool, newID int64, err error) {
	//检查是否创建日志
	var data FieldsLog
	_ = logDB.Get().SetFieldsOne([]string{"id"}).SetStringQuery("module_code", args.ModuleCode).SetStringQuery("branch_code", args.BranchCode).SetIDQuery("config_id", args.ConfigID).SetStringQuery("code", args.Code).Result(&data)
	if data.ID > 0 {
		haveData = true
		return
	}
	//创建日志
	newID, err = logDB.Insert().SetFields([]string{"module_code", "branch_code", "config_id", "code"}).Add(map[string]any{
		"module_code": args.ModuleCode,
		"branch_code": args.BranchCode,
		"config_id":   args.ConfigID,
		"code":        args.Code,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// lockLog 日志进程锁定
func lockLog(moduleCode string, branchCode string) {
	for k, v := range logLock {
		if v.ModuleCode == moduleCode && v.BranchCode == branchCode {
			logLock[k].Lock.Lock()
			return
		}
	}
	logLock = append(logLock, logLockData{
		ModuleCode: moduleCode,
		BranchCode: branchCode,
		Lock:       new(sync.Mutex),
	})
	return
}

// unLockLog 日志进程解锁
func unLockLog(moduleCode string, branchCode string) {
	for k, v := range logLock {
		if v.ModuleCode == moduleCode && v.BranchCode == branchCode {
			logLock[k].Lock.Unlock()
			return
		}
	}
	return
}
