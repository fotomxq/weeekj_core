package BlogStuRead

import (
	"errors"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserSubscription "github.com/fotomxq/weeekj_core/v5/user/subscription"
	"time"
)

// ArgsAppendLog 添加新的数据参数
type ArgsAppendLog struct {
	//创建时间
	CreateAt string `json:"createAt" check:"defaultTime"`
	//结束时间
	EndAt string `json:"endAt" check:"defaultTime"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `json:"userID" check:"id"`
	//学习内容
	ContentID int64 `json:"contentID" check:"id"`
}

// AppendLog 添加新的数据
func AppendLog(args *ArgsAppendLog) (logData FieldsLog, errCode string, err error) {
	//检查是否启动了BlogUserReadMustUserSub，该设置为必须是会员，才会记录阅读记录
	configBlogUserReadMustUserSub := BaseConfig.GetDataBoolNoErr("BlogUserReadMustUserSub")
	if configBlogUserReadMustUserSub {
		//判断用户是否拥有会员
		if !UserSubscription.CheckHaveAnySub(args.UserID) {
			err = errors.New("user no user sub")
			errCode = "err_no_sub"
			return
		}
	}
	//获取参数
	var startAt, endAt time.Time
	startAt, err = CoreFilter.GetTimeByDefault(args.CreateAt)
	if err != nil {
		errCode = "err_time"
		return
	}
	endAt, err = CoreFilter.GetTimeByDefault(args.EndAt)
	if err != nil {
		errCode = "err_time"
		return
	}
	runTime := endAt.Unix() - startAt.Unix()
	//写入数据
	var logID int64
	logID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO blog_stu_read_log (create_at, end_at, run_time, org_id, user_id, content_id) VALUES (:create_at,:end_at,:run_time,:org_id,:user_id,:content_id)", map[string]interface{}{
		"create_at":  startAt,
		"end_at":     endAt,
		"run_time":   runTime,
		"org_id":     args.OrgID,
		"user_id":    args.UserID,
		"content_id": args.ContentID,
	})
	if err != nil {
		errCode = "err_insert"
		return
	}
	logData = getLog(logID)
	//推送nats
	CoreNats.PushDataNoErr("/blog/stu_read/log", "new", logData.UserID, "", logData)
	//反馈
	return
}
