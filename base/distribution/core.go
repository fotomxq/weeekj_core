package BaseDistribution

import (
	"encoding/json"
	"errors"
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/robfig/cron"
	"sync"
	"time"
)

var (
	//定时器
	runTimer       *cron.Cron
	runDeleteLock  = false
	runSaveLock    = false
	runTestAvgLock = false
	//缓冲锁定
	cacheLock         sync.Mutex
	cacheChildLock    sync.Mutex
	cacheChildRunLock sync.Mutex
	//缓冲
	cacheData         []FieldsDistribution
	cacheChildData    []FieldsDistributionChild
	cacheChildRunData []FieldsDistributionChildRun
	//本地持久化数据的所在文件
	localDataSrc = "data.json"
	//全局顶级配置
	globMongodb DataGlobMongodb
	//全局postgres配置
	globPostgres DataGlobPostgres
	//获取配置表结构
	dataConfigs []DataService
	//总的锁定期
	dataConfigsLock sync.Mutex
	//内存存储配置最大时间
	maxSaveTime int64 = 300
	//全局debug配置
	globDebug = false
)

// Init 初始化
func Init() (err error) {
	//将所有数据加载到内存中
	var configData DataType
	//修正路径
	localDataSrc = Router2SystemConfig.RootDir + CoreFile.Sep + "conf" + CoreFile.Sep + "distribution" + CoreFile.Sep + localDataSrc
	//读取配置文件
	var dataByte []byte
	dataByte, err = CoreFile.LoadFile(localDataSrc)
	if err == nil {
		//解析数据
		err = json.Unmarshal(dataByte, &configData)
		if err != nil {
			CoreLog.Error("load file and get json data, file src: ", localDataSrc, ", ", err)
			err = nil
		}
	}
	//写入缓冲
	cacheLock.Lock()
	cacheData = configData.Service
	cacheLock.Unlock()
	cacheChildLock.Lock()
	cacheChildData = configData.Child
	cacheChildLock.Unlock()
	cacheChildRunLock.Lock()
	cacheChildRunData = configData.ChildRun
	cacheChildRunLock.Unlock()
	//反馈
	return
}

// 获取服务列表
type ArgsGetServiceList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList
	//搜索
	Search string
}

func GetServiceList(args *ArgsGetServiceList) (data []FieldsDistribution, dataCount int64, err error) {
	//遍历对应的数据，第一次为搜索
	var searchData []FieldsDistribution
	for _, v := range cacheData {
		if b := CoreFilter.MatchStr(args.Search, v.Mark) || CoreFilter.MatchStr(args.Search, v.Name); !b {
			continue
		}
		searchData = append(searchData, v)
		dataCount += 1
	}
	//开始计数
	skip := int((args.Pages.Page - 1) * args.Pages.Max)
	//遍历搜索后的数据，找出截取的值
	for k, v := range searchData {
		//跳过前缀部分
		if k < skip {
			continue
		}
		//达到目标值跳出
		if len(data) >= int(args.Pages.Max) {
			return
		}
		//写入反馈数据集合
		data = append(data, v)
	}
	return
}

// 设置服务
type ArgsSetService struct {
	//标识码
	Mark string
	//名称
	Name string
	//过期时间
	ExpireInterval int64
	//默认触发函数
	DefaultAction string
	DefaultFunc   string
}

func SetService(args *ArgsSetService) error {
	data, err := getService(args.Mark)
	//创建数据
	if err != nil || data.Mark == "" {
		newData := FieldsDistribution{
			Mark:           args.Mark,
			CreateAt:       CoreFilter.GetNowTime(),
			UpdateAt:       time.Time{},
			Name:           args.Name,
			ExpireInterval: args.ExpireInterval,
			DefaultAction:  args.DefaultAction,
			DefaultFunc:    args.DefaultFunc,
		}
		saveCache(newData)
		return nil
	}
	//修改数据
	data.UpdateAt = CoreFilter.GetNowTime()
	data.Name = args.Name
	data.ExpireInterval = args.ExpireInterval
	data.DefaultAction = args.DefaultAction
	data.DefaultFunc = args.DefaultFunc
	saveCache(data)
	return nil
}

// 获取服务配置
type ArgsGetServiceConfig struct {
	//标识码
	Mark string
}

func GetServiceConfig(args *ArgsGetServiceConfig) (data DataServiceConfig, err error) {
	//遍历获取数据
	for _, v := range dataConfigs {
		if v.Config.ServerMark != args.Mark {
			continue
		}
		if v.ExpireTime < CoreFilter.GetNowTime().Unix() {
			break
		}
		data = v.Config
		if data.NeedGlobDebug {
			data.Debug = globDebug
		}
		return
	}
	//没有找到或过期，则从本地文件拉取
	var configData DataService
	var dataByte []byte
	fileSrc := Router2SystemConfig.RootDir + CoreFile.Sep + "conf" + CoreFile.Sep + "glob" + CoreFile.Sep + args.Mark + ".json"
	dataByte, err = CoreFile.LoadFile(fileSrc)
	if err != nil {
		data = DataServiceConfig{
			ServerName:    args.Mark,
			ServerMark:    args.Mark,
			Debug:         globDebug,
			NeedGlobDebug: true,
		}
		err = nil
		return
	}
	//解析数据
	err = json.Unmarshal(dataByte, &configData.Config)
	if err != nil {
		err = errors.New(fmt.Sprint("load ", args.Mark, " json data, ", err))
		return
	}
	configData.ExpireTime = CoreFilter.GetNowTime().Unix() + maxSaveTime
	dataConfigsLock.Lock()
	dataConfigs = append(dataConfigs, configData)
	dataConfigsLock.Unlock()
	data = configData.Config
	if data.NeedGlobDebug {
		data.Debug = globDebug
	}
	//反馈
	return
}

// 获取指定服务的mongodb数据库配置
type ArgsGetServiceDBMongodbConfig struct {
	//服务标识码
	Mark string
}

func GetServiceDBMongodbConfig(_ *ArgsGetServiceDBMongodbConfig) (data DataGlobMongodb, err error) {
	//根据mark反馈该服务自定义设置
	//不存在则反馈全局设置
	data = globMongodb
	return
}

// 获取指定服务的postgres数据库配置
type ArgsGetServiceDBPostgresConfig struct {
	//服务标识码
	Mark string
}

func GetServiceDBPostgresConfig(_ *ArgsGetServiceDBPostgresConfig) (data DataGlobPostgres, err error) {
	//根据mark反馈该服务自定义设置
	//不存在则反馈全局设置
	data = globPostgres
	return
}

// 删除服务
type ArgsDeleteService struct {
	//标识码
	Mark string
}

func DeleteService(args *ArgsDeleteService) error {
	deleteCache(args.Mark)
	return deleteChildRun(args.Mark, "", "")
}

// 获取服务数据
func getService(mark string) (FieldsDistribution, error) {
	data, b := getCache(mark)
	if !b {
		return data, errors.New("service not exist")
	}
	return data, nil
}

// 获取缓冲
func getCache(mark string) (FieldsDistribution, bool) {
	for _, v := range cacheData {
		if v.Mark == mark {
			return v, true
		}
	}
	return FieldsDistribution{}, false
}

// 保存缓冲
func saveCache(data FieldsDistribution) {
	cacheLock.Lock()
	isFind := false
	for k, v := range cacheData {
		if v.Mark == data.Mark {
			cacheData[k] = data
			isFind = true
			break
		}
	}
	if !isFind {
		cacheData = append(cacheData, data)
	}
	cacheLock.Unlock()
}

// 删除缓冲
func deleteCache(mark string) {
	cacheLock.Lock()
	var newData []FieldsDistribution
	for _, v := range cacheData {
		if v.Mark == mark {
			continue
		}
		newData = append(newData, v)
	}
	cacheData = newData
	cacheLock.Unlock()
}
