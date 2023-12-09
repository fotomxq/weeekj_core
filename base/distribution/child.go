package BaseDistribution

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	"time"
)

// 获取子服务列表
type ArgsGetChildAll struct {
	//标识码
	Mark string
}

func GetChildAll(args *ArgsGetChildAll) ([]FieldsDistributionChild, error) {
	var data []FieldsDistributionChild
	//从缓冲找到数据
	for _, v := range cacheChildData {
		if v.Mark == args.Mark {
			data = append(data, v)
		}
	}
	if len(data) < 1 {
		return nil, errors.New("data is empty")
	}
	return data, nil
}

// 为服务设置子服务
type ArgsSetChild struct {
	//标识码
	Mark string
	//名称
	Name string
	//IP
	IP string
	//端口
	Port string
}

func SetChild(args *ArgsSetChild) error {
	//获取主服务
	serviceData, err := getService(args.Mark)
	if err != nil {
		return err
	}
	//尝试获取数据
	data, err := getChild(args.Mark, args.IP, args.Port)
	//创建数据
	if err != nil || data.Mark == "" {
		newData := FieldsDistributionChild{
			Mark:         args.Mark,
			CreateAt:     CoreFilter.GetNowTime(),
			UpdateAt:     time.Time{},
			ServerName:   args.Name,
			ServerIP:     args.IP,
			ServerPort:   args.Port,
			ExpireTime:   CoreFilter.GetNowTime().Unix() + serviceData.ExpireInterval,
			RunCount:     0,
			RunHour:      0,
			RunHourCount: 0,
			LastCountAvg: 0,
			LastCount:    nil,
		}
		saveChildCache(newData)
		return deleteChildRun(args.Mark, args.IP, args.Port)
	}
	//修改数据
	data.UpdateAt = CoreFilter.GetNowTime()
	data.ExpireTime = CoreFilter.GetNowTime().Unix() + serviceData.ExpireInterval
	data.RunCount += 1
	if data.RunHour == CoreFilter.GetNowTime().Hour() {
		data.RunHourCount += 1
	} else {
		data.RunHour = CoreFilter.GetNowTime().Hour()
		data.RunHourCount = 1
	}
	saveChildCache(data)
	return nil
}

// 通过负载获取子服务
// 不会反馈IP为空的数据，如果全部IP为空，则会反馈失败
type ArgsGetBalancing struct {
	//标识码
	Mark string
}

func GetBalancing(args *ArgsGetBalancing) (FieldsDistributionChild, error) {
	//获取所有数据
	data, err := GetChildAll(&ArgsGetChildAll{
		Mark: args.Mark,
	})
	if err != nil {
		return FieldsDistributionChild{}, err
	}
	//游标
	var result FieldsDistributionChild
	//如果找到数据，则遍历
	for _, v := range data {
		//检查是否过期
		if v.ExpireTime < CoreFilter.GetNowTime().Unix() {
			continue
		}
		//检查是否具备IP
		if v.ServerIP == "" {
			continue
		}
		//如果是第一次遍历，则直接赋值
		if result.Mark == "" {
			result = v
			continue
		}
		//第二次遍历或以上
		// 检查游标数据使用次数，小时平均次数
		if result.RunHourCount > v.RunHourCount {
			result = v
			continue
		}
		//可能多个次数相同
		// 检查响应时间，取最快的
		if v.LastCountAvg < result.LastCountAvg {
			result = v
			continue
		}
	}
	if result.Mark == "" {
		return FieldsDistributionChild{}, errors.New("service not have any server for use")
	}
	return result, nil
}

// 删除一个子服务
type ArgsDeleteChild struct {
	//标识码
	Mark string
	//IP
	IP string
	//端口
	Port string
}

func DeleteChild(args *ArgsDeleteChild) error {
	deleteChildCache(args.Mark, args.IP, args.Port)
	return deleteChildRun(args.Mark, args.IP, args.Port)
}

// 获取某个数据节点
func getChild(mark, ip, port string) (FieldsDistributionChild, error) {
	data, b := getChildCache(mark, ip, port)
	if !b {
		return data, errors.New("child not exist")
	}
	return data, nil
}

// 获取缓冲
func getChildCache(mark, ip, port string) (FieldsDistributionChild, bool) {
	for _, v := range cacheChildData {
		if v.Mark == mark && v.ServerIP == ip && v.ServerPort == port {
			return v, true
		}
	}
	return FieldsDistributionChild{}, false
}

// 保存缓冲
func saveChildCache(data FieldsDistributionChild) {
	cacheChildLock.Lock()
	isFind := false
	for k, v := range cacheChildData {
		if v.Mark == data.Mark && v.ServerIP == data.ServerIP && v.ServerPort == data.ServerPort {
			cacheChildData[k] = data
			isFind = true
			break
		}
	}
	if !isFind {
		cacheChildData = append(cacheChildData, data)
	}
	cacheChildLock.Unlock()
}

// 删除缓冲
func deleteChildCache(mark, ip, port string) {
	cacheChildLock.Lock()
	var newData []FieldsDistributionChild
	for _, v := range cacheChildData {
		if v.Mark == mark && ((ip != "" && v.ServerIP == ip) || ip == "") && ((port != "" && v.ServerPort == port) || port == "") {
			continue
		}
		newData = append(newData, v)
	}
	cacheChildData = newData
	cacheChildLock.Unlock()
}
