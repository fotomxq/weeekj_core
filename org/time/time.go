package OrgTime

import (
	"database/sql"
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsGetList 获取列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// 必填
	OrgID int64 `json:"orgID" check:"id"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 获取列表
func GetList(args *ArgsGetList) (dataList []FieldsWorkTime, dataCount int64, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	var rawList []FieldsWorkTime
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"org_work_time",
		"id",
		"SELECT id FROM org_work_time WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "expire_at", "name"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getConfigByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetAllByBind 获取组织分组或成员参加的所有考勤安排参数
type ArgsGetAllByBind struct {
	//组织分组
	GroupIDs pq.Int64Array `json:"groupIDs" check:"ids"`
	//组织成员
	BindID int64 `json:"bindID" check:"id"`
}

// GetAllByBind 获取组织分组或成员参加的所有考勤安排
func GetAllByBind(args *ArgsGetAllByBind) (dataList []FieldsWorkTime, err error) {
	var rawList []FieldsWorkTime
	if err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM org_work_time WHERE (expire_at > NOW() OR expire_at < to_timestamp(1000000)) AND ($1 && groups OR $2 = ANY(binds));", args.GroupIDs, args.BindID); err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getConfigByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetOne 获取某一个数据参数
type ArgsGetOne struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//组织ID
	// 用于验证
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
}

// GetOne 获取某一个数据
func GetOne(args *ArgsGetOne) (data FieldsWorkTime, err error) {
	data = getConfigByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreate 创建新的数据参数
type ArgsCreate struct {
	//过期时间
	// 过期后自动失效，用于排临时班
	ExpireAt time.Time `db:"expire_at" json:"expireAt" check:"isoTime" empty:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//分组ID列
	Groups pq.Int64Array `db:"groups" json:"groups" check:"ids" empty:"true"`
	//绑定人列
	Binds pq.Int64Array `db:"binds" json:"binds" check:"ids" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//时间配置组
	Configs FieldsConfigs `db:"configs" json:"configs"`
	//轮动任务
	RotConfig FieldsConfigRot `db:"rot_config" json:"rotConfig"`
}

// Create 创建新的数据
func Create(args *ArgsCreate) (data FieldsWorkTime, err error) {
	if args.OrgID < 1 {
		err = errors.New("org id is error")
		return
	}
	if err = checkMonth(args.Configs.Month); err != nil {
		return
	}
	if err = checkMonthWeek(args.Configs.MonthWeek); err != nil {
		return
	}
	if err = checkMonthDay(args.Configs.MonthDay); err != nil {
		return
	}
	if err = checkWeek(args.Configs.Week); err != nil {
		return
	}
	if err = checkTime(args.Configs.WorkTime); err != nil {
		return
	}
	if err = checkTime(args.RotConfig.WorkTime); err != nil {
		return
	}
	//写入数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_work_time", "INSERT INTO org_work_time (expire_at, org_id, groups, binds, name, is_work, configs, rot_config) VALUES (:expire_at, :org_id, :groups, :binds, :name, false, :configs, :rot_config)", args, &data)
	if err != nil {
		return
	}
	//记录是否编辑过数据
	needUpdateConfigMemData = true
	//反馈
	return
}

// ArgsUpdateByID 修改数据参数
type ArgsUpdateByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//过期时间
	// 过期后自动失效，用于排临时班
	ExpireAt time.Time `db:"expire_at" json:"expireAt" check:"isoTime" empty:"true"`
	//分组ID列
	Groups pq.Int64Array `db:"groups" json:"groups" check:"ids" empty:"true"`
	//绑定人列
	Binds pq.Int64Array `db:"binds" json:"binds" check:"ids" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//时间配置组
	Configs FieldsConfigs `db:"configs" json:"configs"`
	//轮动任务
	RotConfig FieldsConfigRot `db:"rot_config" json:"rotConfig"`
}

// UpdateByID 修改数据
func UpdateByID(args *ArgsUpdateByID) (err error) {
	if err = checkMonth(args.Configs.Month); err != nil {
		return
	}
	if err = checkMonthWeek(args.Configs.MonthWeek); err != nil {
		return
	}
	if err = checkMonthDay(args.Configs.MonthDay); err != nil {
		return
	}
	if err = checkWeek(args.Configs.Week); err != nil {
		return
	}
	if err = checkTime(args.Configs.WorkTime); err != nil {
		return
	}
	if err = checkTime(args.RotConfig.WorkTime); err != nil {
		return
	}
	var result sql.Result
	result, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_work_time SET update_at = NOW(), expire_at = :expire_at, groups = :groups, binds = :binds, name = :name, configs = :configs, rot_config = :rot_config WHERE id = :id AND org_id = :org_id", args)
	if err != nil {
		return
	}
	aff, _ := result.RowsAffected()
	if aff < 1 {
		err = errors.New("aff is not exist")
		return
	}
	//记录是否编辑过数据
	needUpdateConfigMemData = true
	//删除缓冲
	deleteConfigCache(args.ID)
	//反馈
	return
}

// ArgsDeleteByID 删除数据参数
type ArgsDeleteByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteByID 删除数据
func DeleteByID(args *ArgsDeleteByID) (err error) {
	//执行操作
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "org_work_time", "id = :id AND org_id = :org_id", args)
	if err != nil {
		return
	}
	//记录是否编辑过数据
	needUpdateConfigMemData = true
	//删除缓冲
	deleteConfigCache(args.ID)
	//反馈
	return
}

// checkMonth 检查月份是否正确
func checkMonth(month []int) error {
	for _, v := range month {
		if v < 1 || v > 12 {
			return errors.New(fmt.Sprint("month is error, month: ", v))
		}
	}
	return nil
}

// checkMonthWeek 检查月份周
func checkMonthWeek(week []int) error {
	for _, v := range week {
		if v < 1 || v > 6 {
			return errors.New(fmt.Sprint("month week is error, month week: ", v))
		}
	}
	return nil
}

// checkMonthDay 检查月份天
func checkMonthDay(day []int) error {
	for _, v := range day {
		if v < 1 || v > 31 {
			return errors.New(fmt.Sprint("month day is error, month day: ", v))
		}
	}
	return nil
}

// checkWeek 检查周合法性
func checkWeek(week []int) error {
	for _, v := range week {
		if v < 0 || v > 7 {
			return errors.New(fmt.Sprint("week is error, week: ", v))
		}
	}
	return nil
}

// checkTime 检查工作时间的安排
func checkTime(data []FieldsWorkTimeTime) error {
	for _, v := range data {
		if v.StartHour < 0 || v.StartMinute < 0 || v.EndHour < 0 || v.EndMinute < 0 {
			return errors.New("task time is error")
		}
		if v.StartHour > 24 || v.StartMinute > 60 || v.EndHour > 24 || v.EndMinute > 60 {
			return errors.New("task time is error")
		}
		//if v.StartHour+v.StartMinute >= v.EndHour+v.EndMinute {
		//	return errors.New("task time is error")
		//}
	}
	return nil
}

// 获取配置
func getConfigByID(id int64) (data FieldsWorkTime) {
	cacheMark := getConfigCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, expire_at, org_id, groups, binds, name, is_work, configs, rot_config FROM org_work_time WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// 缓存
func getConfigCacheMark(id int64) string {
	return fmt.Sprint("org:time:config:id:", id)
}

func deleteConfigCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getConfigCacheMark(id))
}
