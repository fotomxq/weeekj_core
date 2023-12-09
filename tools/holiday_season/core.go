package ToolsHolidaySeason

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/golang-module/carbon"
	"github.com/robfig/cron"
	"time"
)

//节假日模块
// 用于记录和获取外部API的节假日数据，方便其他模块进行使用处理
// 后台可以手动调整节假日数据，手动调整后的数据将无法被覆盖，否则将通过API查询节假日数据集合
// API来源: http://timor.tech/api/holiday

var (
	//定时器
	runTimer   *cron.Cron
	runAPILock = false
)

// ArgsGetList 获取列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//查询范围
	DateMin time.Time `json:"dateMin" check:"isoTime"`
	DateMax time.Time `json:"dateMax" check:"isoTime"`
	//是否包含上班？
	HaveWork bool `json:"haveWork" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 获取列表
func GetList(args *ArgsGetList) (dataList []FieldsHolidaySeason, dataCount int64, err error) {
	if args.DateMin.Unix() >= args.DateMax.Unix() {
		err = errors.New(fmt.Sprint("date is error, date min: ", args.DateMin.String(), ", date max: ", args.DateMax.String()))
		return
	}
	where := "date_at >= :date_min AND date_at <= :date_max"
	maps := map[string]interface{}{
		"date_min": args.DateMin,
		"date_max": args.DateMax,
	}
	if !args.HaveWork {
		where = where + " AND status != 0"
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"tools_holiday_season",
		"id",
		"SELECT id, update_at, date_at, status, is_holiday, name, wage, is_force FROM tools_holiday_season WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "update_at", "date_at"},
	)
	return
}

// ArgsCheckIsWork 检查指定的时间是否上班？参数
type ArgsCheckIsWork struct {
	//对应的日期
	DateAt time.Time `db:"date_at" json:"dateAt"`
}

// CheckIsWork 检查指定的时间是否上班？
func CheckIsWork(args *ArgsCheckIsWork) (b bool) {
	var data FieldsHolidaySeason
	if err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, status FROM tools_holiday_season WHERE date_at = $1", args.DateAt); err == nil {
		if data.ID > 0 {
			return data.Status == 0
		}
	}
	return
}

// 设置指定天的记录
type ArgsSet struct {
	//对应的日期
	DateAt time.Time `db:"date_at" json:"dateAt" check:"isoTime"`
	//节假日类型
	// 0 工作日、1 周末、2 节日、3 调休
	Status int `db:"status" json:"status"`
	//是否房价
	IsHoliday bool `db:"is_holiday" json:"isHoliday"`
	//名称
	// eg: 周二
	Name string `db:"name" json:"name" check:"des" min:"1" max:"100"`
	//薪资倍数
	Wage int `db:"wage" json:"wage"`
	//是否强制修改
	// 如果启动，则只认准修改数据，将不同步API
	IsForce bool `db:"is_force" json:"isForce"`
}

func Set(args *ArgsSet) (err error) {
	//修正时间数据
	dateAt := carbon.CreateFromTimestamp(args.DateAt.Unix())
	dateAt = dateAt.SetHour(0)
	dateAt = dateAt.SetMinute(0)
	dateAt = dateAt.SetSecond(0)
	var data FieldsHolidaySeason
	if err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM tools_holiday_season WHERE date_at = $1", args.DateAt); err == nil {
		if data.ID > 0 {
			//修改数据
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE tools_holiday_season SET  update_at = NOW(), date_at = :date_at, status = :status, is_holiday = :is_holiday, name = :name, wage = :wage, is_force = :is_force WHERE id = :id", map[string]interface{}{
				"id":         data.ID,
				"date_at":    dateAt,
				"status":     args.Status,
				"is_holiday": args.IsHoliday,
				"name":       args.Name,
				"wage":       args.Wage,
				"is_force":   args.IsForce,
			})
			if err != nil {
				err = errors.New("update failed, " + err.Error())
				return
			}
			return
		}
	}
	//创建新的数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO tools_holiday_season (date_at, status, is_holiday, name, wage, is_force) VALUES (:date_at, :status, :is_holiday, :name, :wage, :is_force)", map[string]interface{}{
		"date_at":    dateAt,
		"status":     args.Status,
		"is_holiday": args.IsHoliday,
		"name":       args.Name,
		"wage":       args.Wage,
		"is_force":   args.IsForce,
	})
	return
}

// 删除记录
type ArgsDelete struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

func Delete(args *ArgsDelete) (err error) {
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "tools_holiday_season", "id", args)
	return
}
