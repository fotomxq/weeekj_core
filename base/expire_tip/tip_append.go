package BaseExpireTip

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsAppendTip 添加新的通知参数
type ArgsAppendTip struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//系统标识码
	SystemMark string `db:"system_mark" json:"systemMark"`
	//关联ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//hash
	// 用于额外的数据对比，避免异常
	// 如果不给予，则本模块自动生成，方便对照
	Hash string `db:"hash" json:"hash"`
	//过期时间
	// 请将该时间和内部时间做对比，避免没有及时通知更新造成异常行为
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
}

// AppendTip 添加新的通知
func AppendTip(args *ArgsAppendTip) (err error) {
	//生成hash
	if args.Hash == "" {
		args.Hash = CoreFilter.GetRandStr4(30)
	}
	//根据来源查询是否存在相同数据？
	var newID int64
	var data FieldsTip
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM core_expire_tip WHERE org_id = $1 AND user_id = $2 AND system_mark = $3 AND bind_id = $4", args.OrgID, args.UserID, args.SystemMark, args.BindID)
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_expire_tip SET hash = :hash, expire_at = :expire_at WHERE id = :id", map[string]interface{}{
			"id":        data.ID,
			"hash":      args.Hash,
			"expire_at": args.ExpireAt,
		})
		if err != nil {
			return
		}
		newID = data.ID
	} else {
		//添加数据
		newID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO core_expire_tip (create_at, org_id, user_id, system_mark, bind_id, hash, expire_at) VALUES (:create_at,:org_id,:user_id,:system_mark,:bind_id,:hash,:expire_at)", map[string]interface{}{
			"create_at":   CoreFilter.GetNowTime(),
			"org_id":      args.OrgID,
			"user_id":     args.UserID,
			"system_mark": args.SystemMark,
			"bind_id":     args.BindID,
			"hash":        args.Hash,
			"expire_at":   args.ExpireAt,
		})
		if err != nil {
			return
		}
	}
	//通知
	CoreNats.PushDataNoErr("/base/expire_tip/new", "", newID, "", map[string]interface{}{
		"expireAt": args.ExpireAt,
	})
	//反馈
	return
}

// AppendTipNoErr 无错误推送
// 通知需订阅nats：/base/expire_tip/expire
func AppendTipNoErr(args *ArgsAppendTip) {
	if err := AppendTip(args); err != nil {
		CoreLog.Error("base expire tip append tip failed, system: ", args.SystemMark, ", bind id: ", args.BindID, ", err: ", err)
	}
}

// 通知写入新的数据
func appendHaveNewData(newID int64, hash string, expireAt time.Time) {
	if expireAt.Unix() > CoreFilter.GetNowTimeCarbon().AddHour().Time.Unix() {
		return
	}
	waitExpire1HourLock.Lock()
	defer waitExpire1HourLock.Unlock()
	for k, v := range waitExpire1HourList {
		if v.ID == newID {
			waitExpire1HourList[k].Hash = hash
			waitExpire1HourList[k].ExpireAt = expireAt
			return
		}
	}
	data, err := getID(newID)
	if err != nil {
		return
	}
	waitExpire1HourList = append(waitExpire1HourList, data)
}
