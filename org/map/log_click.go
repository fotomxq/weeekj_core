package OrgMap

import (
	"errors"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	FinanceDeposit "gitee.com/weeekj/weeekj_core/v5/finance/deposit"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsClickMapAd 触发点击参数
type ArgsClickMapAd struct {
	//查看时间
	CreateAt string `db:"create_at" json:"createAt" check:"isoTime"`
	//完成时间
	FinishAt string `db:"finish_at" json:"finishAt" check:"isoTime"`
	//点击用户
	ClickUserID int64 `db:"click_user_id" json:"clickUserID" check:"id"`
	//地图ID
	MapID int64 `db:"map_id" json:"mapID" check:"id"`
}

// ClickMapAd 触发点击
func ClickMapAd(args *ArgsClickMapAd) (successMsg string, errCode string, err error) {
	//获取时间
	var createAt, finishAt time.Time
	createAt, err = CoreFilter.GetTimeByISO(args.CreateAt)
	if err != nil {
		errCode = "err_time"
		return
	}
	finishAt, err = CoreFilter.GetTimeByISO(args.FinishAt)
	if err != nil {
		errCode = "err_time"
		return
	}
	//获取数据集
	var mapData FieldsMap
	mapData, err = GetMapByID(&ArgsGetMapByID{
		ID: args.MapID,
	})
	if err != nil {
		errCode = "err_no_data"
		return
	}
	//检查最短时间
	//if mapData.ViewTimeLimit < finishAt.Unix()-createAt.Unix() {
	//	errCode = "err_limit"
	//	err = errors.New("limit time")
	//	return
	//}
	//检查点击量是否超出
	if mapData.AdCountLimit < mapData.AdCount {
		errCode = "err_limit"
		err = errors.New("ad count limit less")
		return
	}
	//检查记录
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM org_map_ad_log WHERE click_user_id = $1 AND map_id = $2", args.ClickUserID, args.MapID)
	if err == nil && id > 0 {
		errCode = "replace"
		err = errors.New("replace click ad")
		return
	}
	// 开始处理点击奖励
	// 平台设置用户点击奖励金额
	var orgMapClickUnitRandPrice int64
	orgMapClickUnitRandPrice, err = BaseConfig.GetDataInt64("OrgMapClickUnitRandPrice")
	if err != nil {
		errCode = "err_config"
		return
	}
	_, errCode, err = FinanceDeposit.SetByFrom(&FinanceDeposit.ArgsSetByFrom{
		UpdateHash: "",
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     args.ClickUserID,
			Mark:   "",
			Name:   "",
		},
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
		ConfigMark:      "savings",
		AppendSavePrice: orgMapClickUnitRandPrice,
	})
	if err != nil {
		errCode = "deposit_finance_err"
		return
	}
	//添加记录
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_map_ad_log(create_at, finish_at, org_id, user_id, click_user_id, map_id, count, integral_count, bonus) VALUES(:create_at, :finish_at, :org_id, :user_id, :click_user_id, :map_id, :count, :integral_count, :bonus)", map[string]interface{}{
		"create_at":      createAt,
		"finish_at":      finishAt,
		"org_id":         mapData.OrgID,
		"user_id":        mapData.UserID,
		"click_user_id":  args.ClickUserID,
		"map_id":         args.MapID,
		"count":          0,
		"integral_count": 1,
		"bonus":          orgMapClickUnitRandPrice,
	})
	if err != nil {
		errCode = "err_insert"
		return
	}
	//变更统计次数
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_map SET ad_count = ad_count + 1 WHERE id = :id", map[string]interface{}{
		"id": args.MapID,
	})
	if err != nil {
		errCode = "err_update"
		return
	}
	successMsg = "ok"
	//反馈
	return
}
