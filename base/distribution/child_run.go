package BaseDistribution

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"time"
)

// ArgsGetChildRun 获取子服务的所有run参数
type ArgsGetChildRun struct {
	//标识码
	Mark string
	//IP
	IP string
	//端口
	Port string
}

// GetChildRun 获取子服务的所有run
func GetChildRun(args *ArgsGetChildRun) ([]FieldsDistributionChildRun, error) {
	var data []FieldsDistributionChildRun
	//从缓冲找到数据
	for _, v := range cacheChildRunData {
		if v.Mark == args.Mark && v.ServerIP == args.IP && v.ServerPort == args.Port {
			data = append(data, v)
		}
	}
	if len(data) < 1 {
		return nil, errors.New("data is empty")
	}
	return data, nil
}

// 设置run
type ArgsSetChildRun struct {
	//标识码
	Mark string
	//IP
	IP string
	//端口
	Port string
	//运行标识码
	RunMark string
	//过期时间
	ExpireAddTime int64
}

func SetChildRun(args *ArgsSetChildRun) error {
	//获取主服务
	_, err := getChild(args.Mark, args.IP, args.Port)
	if err != nil {
		return errors.New("child not exist, " + err.Error())
	}
	//尝试获取数据
	data, err := getChildRun(args.Mark, args.IP, args.Port, args.RunMark)
	//创建数据
	if err != nil {
		newData := FieldsDistributionChildRun{
			Mark:         args.Mark,
			CreateAt:     CoreFilter.GetNowTime(),
			UpdateAt:     time.Time{},
			RunMark:      args.RunMark,
			ServerIP:     args.IP,
			ServerPort:   args.Port,
			ExpireTime:   CoreFilter.GetNowTime().Unix() + args.ExpireAddTime,
			RunCount:     0,
			RunHour:      0,
			RunHourCount: 0,
		}
		saveChildRunCache(newData)
		return nil
	}
	//修改数据
	data.UpdateAt = CoreFilter.GetNowTime()
	data.ExpireTime = CoreFilter.GetNowTime().Unix() + args.ExpireAddTime
	data.RunCount += 1
	if data.RunHour == CoreFilter.GetNowTime().Hour() {
		data.RunHourCount += 1
	} else {
		data.RunHour = CoreFilter.GetNowTime().Hour()
		data.RunHourCount = 1
	}
	saveChildRunCache(data)
	return nil
}

// 删除mark对应的所有run
func deleteChildRun(mark, ip, port string) error {
	deleteChildRunCache(mark, ip, port)
	return nil
}

// 获取子run
func getChildRun(mark, ip, port, runMark string) (FieldsDistributionChildRun, error) {
	data, b := getChildRunCache(mark, ip, port, runMark)
	if !b {
		return data, errors.New("child run not exist")
	}
	return data, nil
}

// 获取缓冲
func getChildRunCache(mark, ip, port, runMark string) (FieldsDistributionChildRun, bool) {
	for _, v := range cacheChildRunData {
		if v.Mark == mark && v.ServerIP == ip && v.ServerPort == port && v.RunMark == runMark {
			return v, true
		}
	}
	return FieldsDistributionChildRun{}, false
}

// 保存缓冲
func saveChildRunCache(data FieldsDistributionChildRun) {
	cacheChildRunLock.Lock()
	isFind := false
	for k, v := range cacheChildRunData {
		if v.Mark == data.Mark && v.ServerIP == data.ServerIP && v.ServerPort == data.ServerPort {
			cacheChildRunData[k] = data
			isFind = true
			break
		}
	}
	if !isFind {
		cacheChildRunData = append(cacheChildRunData, data)
	}
	cacheChildRunLock.Unlock()
}

// 删除缓冲
func deleteChildRunCache(mark, ip, port string) {
	cacheChildRunLock.Lock()
	var newData []FieldsDistributionChildRun
	for _, v := range cacheChildRunData {
		if v.Mark == mark && ((ip != "" && v.ServerIP == ip) || ip == "") && ((port != "" && v.ServerPort == port) || port == "") {
			continue
		}
		newData = append(newData, v)
	}
	cacheChildRunData = newData
	cacheChildRunLock.Unlock()
}
