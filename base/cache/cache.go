package BaseCache

import (
	"encoding/json"
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	"github.com/robfig/cron"
	"sync"
)

//本方法集为底层实现

// 锁定
var (
	//定时器
	runTimer     *cron.Cron
	runCacheLock = false
	//缓冲锁定和数据集合
	cacheLock sync.Mutex
	cacheData []DataCache
)

// ArgsGetByMark 获取某个数据集合
type ArgsGetByMark struct {
	Mark string
}

// Deprecated
func GetByMark(args *ArgsGetByMark) (DataCache, bool) {
	for _, v := range cacheData {
		if v.Mark == args.Mark {
			if v.ExpireTime < CoreFilter.GetNowTime().Unix() {
				return DataCache{}, false
			}
			return v, true
		}
	}
	return DataCache{}, false
}

// GetByMarkInterface 获取缓冲模块interface数据
// Deprecated
func GetByMarkInterface(mark string, data interface{}) bool {
	rawData, b := GetByMark(&ArgsGetByMark{
		Mark: mark,
	})
	if !b {
		return b
	}
	data = rawData.ValueInterface
	return true
}

// GetByMarkInterfaceReturn 获取缓冲模块interface数据
// Deprecated
func GetByMarkInterfaceReturn(mark string) (data interface{}, b bool) {
	var rawData DataCache
	rawData, b = GetByMark(&ArgsGetByMark{
		Mark: mark,
	})
	if !b {
		return
	}
	data = rawData.ValueInterface
	return
}

// SetData 写入数据集合
// Deprecated
func SetData(args *DataCache) {
	cacheLock.Lock()
	args.ExpireTime = CoreFilter.GetNowTime().Unix() + args.ExpireTime
	isFind := false
	for k, v := range cacheData {
		if v.Mark == args.Mark {
			cacheData[k] = DataCache{
				CreateTime:     args.CreateTime,
				ExpireTime:     args.ExpireTime,
				Mark:           args.Mark,
				Value:          args.Value,
				ValueInt64:     args.ValueInt64,
				ValueFloat64:   args.ValueFloat64,
				ValueBool:      args.ValueBool,
				ValueByte:      args.ValueByte,
				ValueInterface: args.ValueInterface,
			}
			isFind = true
		}
	}
	if !isFind {
		cacheData = append(cacheData, DataCache{
			CreateTime:     args.CreateTime,
			ExpireTime:     args.ExpireTime,
			Mark:           args.Mark,
			Value:          args.Value,
			ValueInt64:     args.ValueInt64,
			ValueFloat64:   args.ValueFloat64,
			ValueBool:      args.ValueBool,
			ValueByte:      args.ValueByte,
			ValueInterface: args.ValueInterface,
		})
	}
	cacheLock.Unlock()
}

// GetDataByByte 将byte转为特殊结构体
// Deprecated
func GetDataByByte(mark string, data interface{}) error {
	cacheData, b := GetByMark(&ArgsGetByMark{
		Mark: mark,
	})
	if !b {
		return errors.New("no cache")
	}
	return json.Unmarshal(cacheData.ValueByte, &data)
}

// GetDataByInt64 将int64转为特殊结构体
// Deprecated
func GetDataByInt64(mark string) (val int64, b bool) {
	var cacheData DataCache
	cacheData, b = GetByMark(&ArgsGetByMark{
		Mark: mark,
	})
	if !b {
		return
	}
	val = cacheData.ValueInt64
	return
}

// SetDataByByte 将特殊结构体转为byte
// Deprecated
func SetDataByByte(mark string, addExpireTime int64, data interface{}) error {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	SetData(&DataCache{
		CreateTime:     0,
		ExpireTime:     addExpireTime,
		Mark:           mark,
		Value:          "",
		ValueInt64:     0,
		ValueFloat64:   0,
		ValueBool:      false,
		ValueByte:      dataByte,
		ValueInterface: nil,
	})
	return nil
}

// SetDataByBool 直接存储对应值
// Deprecated
func SetDataByBool(mark string, addExpireTime int64, b bool) error {
	SetData(&DataCache{
		CreateTime:     0,
		ExpireTime:     addExpireTime,
		Mark:           mark,
		Value:          "",
		ValueInt64:     0,
		ValueFloat64:   0,
		ValueBool:      b,
		ValueByte:      nil,
		ValueInterface: nil,
	})
	return nil
}

// SetDataByInt64 直接存储对应值
// Deprecated
func SetDataByInt64(mark string, addExpireTime int64, val int64) error {
	SetData(&DataCache{
		CreateTime:     0,
		ExpireTime:     addExpireTime,
		Mark:           mark,
		Value:          "",
		ValueInt64:     val,
		ValueFloat64:   0,
		ValueBool:      false,
		ValueByte:      nil,
		ValueInterface: nil,
	})
	return nil
}

// SetDataByFloat64 直接存储对应值
// Deprecated
func SetDataByFloat64(mark string, addExpireTime int64, val float64) error {
	SetData(&DataCache{
		CreateTime:     0,
		ExpireTime:     addExpireTime,
		Mark:           mark,
		Value:          "",
		ValueInt64:     0,
		ValueFloat64:   val,
		ValueBool:      false,
		ValueByte:      nil,
		ValueInterface: nil,
	})
	return nil
}

// SetDataByInterface 设置动态变量
// Deprecated
func SetDataByInterface(mark string, addExpireTime int64, val interface{}) error {
	SetData(&DataCache{
		CreateTime:     0,
		ExpireTime:     addExpireTime,
		Mark:           mark,
		Value:          "",
		ValueInt64:     0,
		ValueFloat64:   0,
		ValueBool:      false,
		ValueByte:      nil,
		ValueInterface: val,
	})
	return nil
}

// DeleteMark 删除缓冲
// Deprecated
func DeleteMark(mark string) {
	cacheLock.Lock()
	var newCacheData []DataCache
	for _, v := range cacheData {
		if v.Mark == mark {
			continue
		}
		newCacheData = append(newCacheData, v)
	}
	cacheData = newCacheData
	cacheLock.Unlock()
}
