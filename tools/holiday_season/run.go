package ToolsHolidaySeason

import (
	"encoding/json"
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreHttp "github.com/fotomxq/weeekj_core/v5/core/http"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSystemClose "github.com/fotomxq/weeekj_core/v5/core/system_close"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/robfig/cron"
	"time"
)

// Run 维护包
func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("base holiday season run, ", r)
		}
		//关闭时间
		runTimer.Stop()
	}()
	//初始化
	runTimer = cron.New()
	//是否启动
	if !Router2SystemConfig.GlobConfig.OtherAPI.OpenSyncHolidaySeason {
		return
	}
	//该方法不会做任何延迟，以避免影响其他前置模块获取异常
	//自动化处理
	if err := runTimer.AddFunc("@every 5m", func() {
		if runAPILock {
			return
		}
		runAPILock = true
		//调用API维护数据
		runAPI()
		runAPILock = false
	}); err != nil {
		CoreLog.Error("base holiday season run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}

// 检查最近30的数据，如果不存在将从远程拉取数据
// 除此之外，将检查最近修改且撤销锁定的数据，和远程进行同步
// api eg: http://timor.tech/api/holiday/info/2021-03-09
func runAPI() {
	//检查最近1分钟取消锁定的数据
	var dataList []FieldsHolidaySeason
	updateAt := CoreFilter.GetNowTimeCarbon().SubMinute()
	if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, update_at, date_at, status, is_holiday, name, wage, is_force FROM tools_holiday_season WHERE update_at > $1 AND is_force = false", updateAt.Time); err == nil {
		for _, v := range dataList {
			//调用API获取数据
			apiData, err := getAPIData(v.DateAt)
			if err != nil {
				CoreLog.Error("base holiday season run, get api data, ", err)
				//发生该异常后，直接退出执行程序
				//同时将自动发起延迟处理，避免异常问题
				time.Sleep(time.Minute * 10)
				return
			}
			//根据反馈结果，修改数据
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE tools_holiday_season SET status = :status, is_holiday = :is_holiday, name = :name, wage = :wage WHERE id = :id", map[string]interface{}{
				"id":         v.ID,
				"status":     apiData.Type.Type,
				"is_holiday": apiData.Holiday.Holiday,
				"name":       apiData.Type.Name,
				"wage":       apiData.Holiday.Wage,
			})
			if err != nil {
				CoreLog.Error("base holiday season run, update data, id: ", v.ID, ", err: ", err)
				time.Sleep(time.Minute * 5)
			}
		}
	}
	//检查今天到未来30天
	// 此处只要创建了一条数据，则会退出处理，直到下一轮数据维护
	step := 0
	nowDay := CoreFilter.GetNowTimeCarbon()
	nowDay = nowDay.SetHour(0)
	nowDay = nowDay.SetMinute(0)
	nowDay = nowDay.SetSecond(0)
	for {
		//超出41跳出
		if step > 41 {
			break
		}
		//获取数据，如果存在则跳过
		var data FieldsHolidaySeason
		if err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM tools_holiday_season WHERE date_at = $1", nowDay.Time); err == nil {
			if data.ID > 0 {
				nowDay = nowDay.AddDay()
				step += 1
				continue
			}
		}
		//获取API数据
		apiData, err := getAPIData(nowDay.Time)
		if err != nil {
			CoreLog.Error("base holiday season run, get api data, ", err)
			//发生该异常后，直接退出执行程序
			//同时将自动发起延迟处理，避免异常问题
			time.Sleep(time.Minute * 10)
			return
		}
		//写入新的数据
		err = Set(&ArgsSet{
			DateAt:    nowDay.Time,
			Status:    apiData.Type.Type,
			IsHoliday: apiData.Holiday.Holiday,
			Name:      apiData.Type.Name,
			Wage:      apiData.Holiday.Wage,
			IsForce:   false,
		})
		if err == nil {
			return
		} else {
			CoreLog.Error("base holiday season run, create new data, ", err)
		}
		//叠加时间
		nowDay = nowDay.AddDay()
		step += 1
	}
}

type dataAPI struct {
	//状态代码
	// 0 正常 -1 服务异常
	Code int `json:"code"`
	//数据集合
	Type dataAPIChild `json:"type"`
	//节假日数据
	Holiday dataAPIChildHoliday `json:"holiday"`
}

type dataAPIChild struct {
	//状态
	// enum(0, 1, 2, 3) 节假日类型，分别表示 工作日、周末、节日、调休。
	Type int `json:"type"`
	//名称
	Name string `json:"name"`
	//类型
	// 一周中的第几天。值为 1 - 7，分别表示 周一 至 周日
	Week int `json:"week"`
}

type dataAPIChildHoliday struct {
	//是否为节假日
	Holiday bool `json:"holiday"`
	//名称
	Name string `json:"name"`
	//工资倍数
	Wage int `json:"wage"`
	//只在调休下有该字段。true表示放完假后调休，false表示先调休再放假
	After bool `json:"after"`
	//只在调休下有该字段。表示调休的节假日
	Target string `json:"target"`
}

func getAPIData(date time.Time) (data dataAPI, err error) {
	//请求数据
	var dataByte []byte
	dataByte, err = CoreHttp.GetData(fmt.Sprint("https://timor.tech/api/holiday/info/", date.Format("2006-01-02")), nil, "", false)
	if err != nil {
		return
	}
	//解析数据
	if err = json.Unmarshal(dataByte, &data); err != nil {
		return
	}
	//反馈数据
	if data.Code != 0 {
		err = errors.New("api server failed")
		return
	}
	return
}
