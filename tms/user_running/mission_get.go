package TMSUserRunning

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetMissionList 获取任务列表参数
type ArgsGetMissionList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//跑腿单类型
	// 0 帮我送 ; 1 帮我买; 2 帮我取
	RunType int `db:"run_type" json:"runType" check:"intThan0" empty:"true"`
	//关联组织
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 允许为0，则该信息不属于任何用户，或不和任何用户关联
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//是否完成取货
	NeedIsTake bool `json:"needIsTake" check:"bool"`
	IsTake     bool `json:"isTake" check:"bool"`
	//是否完成
	NeedIsFinish bool `json:"needIsFinish" check:"bool"`
	IsFinish     bool `json:"isFinish" check:"bool"`
	//关联订单ID
	// 可能没有关联订单
	OrderID int64 `db:"order_id" json:"orderID" check:"id" empty:"true"`
	//履约跑腿员
	// 用户角色ID
	RoleID int64 `db:"role_id" json:"roleID" check:"id" empty:"true"`
	//是否支付跑腿费
	// 当前费用部分
	NeedIsRunPay bool `json:"needIsRunPay" check:"bool"`
	IsRunPay     bool `json:"isRunPay" check:"bool"`
	//是否已经支付过跑腿费
	// 存在支付过的费用
	NeedHaveRunPay bool `json:"needHaveRunPay" check:"bool"`
	HaveRunPay     bool `json:"haveRunPay" check:"bool"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetMissionList 获取任务列表
func GetMissionList(args *ArgsGetMissionList) (dataList []FieldsMission, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.RunType > -1 {
		where = where + " AND run_type = :run_type"
		maps["run_type"] = args.RunType
	}
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.NeedIsTake {
		if args.IsTake {
			where = where + " AND take_at > to_timestamp(1000000)"
		} else {
			where = where + " AND take_at <= to_timestamp(1000000)"
		}
	}
	if args.NeedIsFinish {
		if args.IsFinish {
			where = where + " AND finish_at > to_timestamp(1000000)"
		} else {
			where = where + " AND finish_at <= to_timestamp(1000000)"
		}
	}
	if args.OrderID > -1 {
		where = where + " AND order_id = :order_id"
		maps["order_id"] = args.OrderID
	}
	if args.RoleID > -1 {
		where = where + " AND role_id = :role_id"
		maps["role_id"] = args.RoleID
	}
	if args.NeedIsRunPay {
		if args.IsRunPay {
			where = where + " AND run_pay_at > to_timestamp(1000000)"
		} else {
			where = where + " AND run_pay_at <= to_timestamp(1000000)"
		}
	}
	if args.NeedHaveRunPay {
		if args.HaveRunPay {
			where = where + " AND run_price > 0"
		} else {
			where = where + " AND run_price <= 0"
		}
	}
	if args.Search != "" {
		where = where + " AND (order_wait_price ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%' OR order_des ILIKE '%' || :search || '%' OR from_address ->> 'address' ILIKE '%' || :search || '%' OR from_address ->> 'name' ILIKE '%' || :search || '%' OR from_address ->> 'phone' ILIKE '%' || :search || '%' OR to_address ->> 'address' ILIKE '%' || :search || '%' OR to_address ->> 'name' ILIKE '%' || :search || '%' OR to_address ->> 'phone' ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "tms_user_running_mission"
	var rawList []FieldsMission
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id "+"FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "take_at", "finish_at"},
	)
	if err != nil || len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getMissionID(v.ID)
		if vData.ID < 1 {
			continue
		}
		vData.TakeCode = ""
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetMissionID 获取指定任务信息参数
type ArgsGetMissionID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//关联组织
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//关联用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//履约跑腿员
	// 用户角色ID
	RoleID int64 `db:"role_id" json:"roleID" check:"id" empty:"true"`
}

// GetMissionID 获取指定任务信息
func GetMissionID(args *ArgsGetMissionID) (data FieldsMission, err error) {
	data = getMissionID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	if !CoreFilter.EqID2(args.OrgID, data.OrgID) || !CoreFilter.EqID2(args.UserID, data.UserID) || !CoreFilter.EqID2(args.RoleID, data.RoleID) {
		err = errors.New("no data")
		return
	}
	data.TakeCode = ""
	return
}

func getMissionID(id int64) (data FieldsMission) {
	cacheMark := getMissionCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, wait_at, good_type, take_at, finish_at, take_code, run_type, org_id, user_id, order_id, role_id, run_pay_at, run_pay_id, run_price, run_wait_price, run_pay_list, run_pay_after, order_pay_after, order_wait_price, order_price, order_pay_at, order_pay_id, des, order_des_files, order_des, good_widget, from_address, to_address, logs, params FROM tms_user_running_mission WHERE id = $1", id)
	if err != nil || data.ID < 1 {
		err = errors.New("no data")
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheTime)
	return
}

// GetMissionAllInfoID 获取指定任务全部信息
func GetMissionAllInfoID(args *ArgsGetMissionID) (data FieldsMission, err error) {
	data = getMissionID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	if !CoreFilter.EqID2(args.OrgID, data.OrgID) || !CoreFilter.EqID2(args.UserID, data.UserID) || !CoreFilter.EqID2(args.RoleID, data.RoleID) {
		err = errors.New("no data")
		return
	}
	return
}

// GetMissionTakeCodeByID 获取任务领取代码
func GetMissionTakeCodeByID(args *ArgsGetMissionID) (code string, err error) {
	var data FieldsMission
	data = getMissionID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	if data.TakeCode == "" {
		err = errors.New("no data")
		return
	}
	code = data.TakeCode
	return
}

// 获取关联订单数据集
func getMissionListByOrderID(orderID int64) (dataList []FieldsMission, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM tms_user_running_mission WHERE order_id = $1 AND delete_at < to_timestamp(1000000)", orderID)
	if err != nil {
		return
	}
	for k, v := range dataList {
		dataList[k] = getMissionID(v.ID)
	}
	return
}

// 获取支付ID的相关任务
func getMissionListByPayID(payID int64) (dataList []FieldsMission, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM tms_user_running_mission WHERE (order_pay_id = $1 OR run_pay_id = $1) AND delete_at < to_timestamp(1000000)", payID)
	if err != nil {
		return
	}
	for k, v := range dataList {
		dataList[k] = getMissionID(v.ID)
	}
	return
}

// 获取缓冲名称
func getMissionCacheMark(id int64) string {
	return fmt.Sprint("tms:user:running:id:", id)
}

// 删除缓冲
func deleteMissionCache(id int64) {
	cacheMark := getMissionCacheMark(id)
	Router2SystemConfig.MainCache.DeleteMark(cacheMark)
}
