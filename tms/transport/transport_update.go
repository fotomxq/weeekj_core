package TMSTransport

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLFields "gitee.com/weeekj/weeekj_core/v5/core/sql/fields"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsUpdateTransportOldToNewBind 修改配送员所有未完成配送单为新的配送员参数
type ArgsUpdateTransportOldToNewBind struct {
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作组织人员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//旧配送人员
	OldBindID int64 `db:"old_bind_id" json:"oldBindID" check:"id"`
	//新配送员
	NewBindID int64 `db:"new_bind_id" json:"newBindID" check:"id"`
}

// UpdateTransportOldToNewBind 修改配送员所有未完成配送单为新的配送员
func UpdateTransportOldToNewBind(args *ArgsUpdateTransportOldToNewBind) (err error) {
	//禁止自己转自己
	if args.OldBindID == args.NewBindID {
		err = errors.New("old and new is same")
		return
	}
	//获取所有配送单
	var dataList []CoreSQLFields.FieldsID
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM tms_transport WHERE bind_id = $1 AND delete_at < to_timestamp(1000000) AND status != 3 AND finish_at < to_timestamp(1000000) AND org_id = $2", args.OldBindID, args.OrgID)
	if err != nil || len(dataList) < 1 {
		err = nil
		return
	}
	//修改为新配送员
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport SET update_at = NOW(), bind_id = :new_bind_id WHERE bind_id = :old_bind_id AND delete_at < to_timestamp(1000000) AND status != 3 AND finish_at < to_timestamp(1000000) AND org_id = :org_id", args)
	if err != nil {
		return
	}
	//增加日志
	for _, v := range dataList {
		_ = appendLog(&argsAppendLog{
			OrgID:           args.OrgID,
			BindID:          args.BindID,
			TransportID:     v.ID,
			TransportBindID: args.NewBindID,
			Mark:            "update_new_bind",
			Des:             fmt.Sprint("批量修改配送员[", args.OldBindID, "]到新的配送员[", args.NewBindID, "]"),
		})
		pushMQTTTransportUpdate(args.OrgID, args.BindID, 0, 3)
		pushMQTTTransportUpdate(args.OrgID, args.NewBindID, 0, 4)
	}
	if args.OldBindID > 0 {
		pushNatsAnalysisBind(args.OldBindID)
	}
	if args.NewBindID > 0 {
		pushNatsAnalysisBind(args.NewBindID)
	}
	return
}

// ArgsUpdateTransportGPS 更新配送单定位信息参数
type ArgsUpdateTransportGPS struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//地图制式
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
}

// UpdateTransportGPS 更新配送单定位信息
func UpdateTransportGPS(args *ArgsUpdateTransportGPS) (err error) {
	var data FieldsTransport
	data, err = GetTransport(&ArgsGetTransport{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		return
	}
	err = appendTransportGPS(&argsAppendTransportGPS{
		OrgID:       data.OrgID,
		TransportID: data.ID,
		MapType:     args.MapType,
		Longitude:   args.Longitude,
		Latitude:    args.Latitude,
	})
	if err != nil {
		return
	}
	if data.BindID > 0 {
		err = appendBindGPS(&argsAppendBindGPS{
			OrgID:     data.OrgID,
			BindID:    data.BindID,
			MapType:   args.MapType,
			Longitude: args.Longitude,
			Latitude:  args.Latitude,
		})
		if err != nil {
			return
		}
	}
	return
}

// ArgsUpdateTransportTaskAt 修改上门时间参数
type ArgsUpdateTransportTaskAt struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作组织人员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//期望送货时间
	TaskAt string `db:"task_at" json:"taskAt" check:"isoTime" empty:"true"`
}

// UpdateTransportTaskAt 修改上门时间
func UpdateTransportTaskAt(args *ArgsUpdateTransportTaskAt) (err error) {
	//预期上门时间
	var taskAt time.Time
	taskAt, err = CoreFilter.GetTimeByISO(args.TaskAt)
	if err != nil {
		taskAt = time.Now()
	}
	//修改时间
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport SET update_at = NOW(), task_at = :task_at WHERE id = :id AND (org_id = :org_id OR :org_id < 1) AND (:bind_id < 1 OR bind_id = :bind_id)", map[string]interface{}{
		"id":      args.ID,
		"org_id":  args.OrgID,
		"bind_id": args.BindID,
		"task_at": taskAt,
	})
	if err != nil {
		return
	}
	var data FieldsTransport
	data, err = GetTransport(&ArgsGetTransport{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		return
	}
	_ = appendLog(&argsAppendLog{
		OrgID:           args.OrgID,
		BindID:          args.BindID,
		TransportID:     args.ID,
		TransportBindID: data.BindID,
		Mark:            "task_at",
		Des:             fmt.Sprint("更新新的配送时间(", args.TaskAt, ")"),
	})
	pushMQTTTransportUpdate(data.OrgID, data.BindID, data.ID, 1)
	return
}

// ArgsUpdateTransportPick 更新配送单状态到取货中参数
type ArgsUpdateTransportPick struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作组织人员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
}

// UpdateTransportPick 更新配送单状态到取货中
func UpdateTransportPick(args *ArgsUpdateTransportPick) (err error) {
	var data FieldsTransport
	data, err = GetTransport(&ArgsGetTransport{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport SET update_at = NOW(), status = 1 WHERE id = :id AND (org_id = :org_id OR :org_id < 1) AND status = 0", map[string]interface{}{
		"id":     args.ID,
		"org_id": args.OrgID,
	})
	if err != nil {
		return
	}
	//记录日志
	_ = appendLog(&argsAppendLog{
		OrgID:           args.OrgID,
		BindID:          args.BindID,
		TransportID:     args.ID,
		TransportBindID: data.BindID,
		Mark:            "pick",
		Des:             fmt.Sprint("更新状态为配送人员取件中"),
	})
	//推送nats
	pushNatsStatusUpdate("pick", data.ID, "更新状态为配送人员取件中")
	//推送mqtt
	pushMQTTTransportUpdate(data.OrgID, data.BindID, data.ID, 1)
	//反馈
	return
}

// ArgsUpdateTransportSend 更新配送单状态到送货中参数
type ArgsUpdateTransportSend struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作组织人员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
}

// UpdateTransportSend 更新配送单状态到送货中
func UpdateTransportSend(args *ArgsUpdateTransportSend) (err error) {
	var data FieldsTransport
	data, err = GetTransport(&ArgsGetTransport{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport SET update_at = NOW(), status = 2 WHERE id = :id AND (org_id = :org_id OR :org_id < 1) AND status = 1", map[string]interface{}{
		"id":     args.ID,
		"org_id": args.OrgID,
	})
	if err != nil {
		return
	}
	//记录日志
	_ = appendLog(&argsAppendLog{
		OrgID:           args.OrgID,
		BindID:          args.BindID,
		TransportID:     args.ID,
		TransportBindID: data.BindID,
		Mark:            "send",
		Des:             fmt.Sprint("更新状态为配送人员送货中"),
	})
	//推送nats
	pushNatsStatusUpdate("send", data.ID, "更新状态为配送人员送货中")
	//推送mqtt
	pushMQTTTransportUpdate(data.OrgID, data.BindID, data.ID, 1)
	//反馈
	return
}

// ArgsUpdateTransportComment 评价配送单参数
type ArgsUpdateTransportComment struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//客户档案ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//客户用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//评级
	// 1-5 级别
	Level int `db:"level" json:"level"`
}

// UpdateTransportComment 评价配送单
func UpdateTransportComment(args *ArgsUpdateTransportComment) (err error) {
	var data FieldsTransport
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, org_id, bind_id FROM tms_transport WHERE id = $1 AND info_id = $2 AND user_id = $3", args.ID, args.InfoID, args.UserID)
	if err != nil || data.ID < 1 {
		err = errors.New("data not exist")
		return
	}
	err = appendAnalysis(&argsAppendAnalysis{
		OrgID:       data.OrgID,
		BindID:      0,
		InfoID:      args.InfoID,
		UserID:      args.UserID,
		TransportID: data.ID,
		KM:          0,
		OverTime:    0,
		Level:       args.Level,
	})
	if err != nil {
		return
	}
	if data.BindID > 0 {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport_bind SET level_1_day = level_1_day + :level_1_day WHERE bind_id = :bind_id", map[string]interface{}{
			"bind_id":     data.BindID,
			"level_1_day": args.Level,
		})
		if err != nil {
			return
		}
	}
	return
}

// ArgsUpdateTransportBind 修改配送人员参数
type ArgsUpdateTransportBind struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作组织人员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//新的配送员
	NewBindID int64 `db:"new_bind_id" json:"newBindID"`
}

// UpdateTransportBind 修改配送人员
func UpdateTransportBind(args *ArgsUpdateTransportBind) (err error) {
	var data FieldsTransport
	data, err = GetTransport(&ArgsGetTransport{
		ID:     args.ID,
		OrgID:  args.OrgID,
		InfoID: 0,
		UserID: 0,
	})
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport SET update_at = NOW(), bind_id = :new_bind_id WHERE id = :id AND (org_id = :org_id OR :org_id < 1) AND status != 3", map[string]interface{}{
		"id":          args.ID,
		"org_id":      args.OrgID,
		"new_bind_id": args.NewBindID,
	})
	if err != nil {
		return
	}
	_ = appendLog(&argsAppendLog{
		OrgID:           args.OrgID,
		BindID:          args.BindID,
		TransportID:     args.ID,
		TransportBindID: data.BindID,
		Mark:            "send",
		Des:             fmt.Sprint("修改配送人员为[", args.NewBindID, "]"),
	})
	if err2 := updateTransportBindAnalysis(data, args.NewBindID); err2 != nil {
		CoreLog.Error("update transport bind analysis, ", err2)
		return
	}
	if data.BindID > 0 {
		pushMQTTTransportUpdate(data.OrgID, data.BindID, data.ID, 3)
		pushNatsAnalysisBind(data.BindID)
	}
	if args.NewBindID > 0 {
		pushMQTTTransportUpdate(data.OrgID, args.NewBindID, data.ID, 4)
		pushNatsAnalysisBind(args.NewBindID)
	}
	return
}

// 通知变更配送状态
func pushNatsStatusUpdate(action string, id int64, des string) {
	CoreNats.PushDataNoErr("/tms/transport/update", action, id, "", map[string]interface{}{
		"des": des,
	})
}
