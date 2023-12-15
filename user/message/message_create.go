package UserMessage

import (
	"errors"
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BaseExpireTip "github.com/fotomxq/weeekj_core/v5/base/expire_tip"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsCreate 创建新的消息参数
type ArgsCreate struct {
	//预计发送时间
	WaitSendAt time.Time `db:"wait_send_at" json:"waitSendAt" check:"isoTime" empty:"true"`
	//发送人
	// 如果为0则为系统消息，同时自动跳过时间差验证机制
	SendUserID int64 `db:"send_user_id" json:"sendUserID" check:"id"`
	//接收人
	ReceiveUserID int64 `db:"receive_user_id" json:"receiveUserID" check:"id"`
	//标题
	Title string `db:"title" json:"title" check:"des" min:"1" max:"300"`
	//内容
	Content string `db:"content" json:"content" check:"des" min:"1" max:"1500"`
	//附件文件列
	Files pq.Int64Array `db:"files" json:"files" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// Create 创建新的消息
func Create(args *ArgsCreate) (data FieldsMessage, err error) {
	//分析参数
	if len(args.Params) < 1 {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	if len(args.Files) > 10 {
		err = errors.New("file too many")
		return
	}
	//创建消息
	var messageID int64
	messageID, err = create(args)
	if err != nil {
		return
	}
	//获取消息数据包
	data = getByID(messageID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//反馈
	return
}

func create(args *ArgsCreate) (messageID int64, err error) {
	//如果存在发件人
	if args.SendUserID > 0 {
		//禁止连续推送消息，如需推送系统消息，走系统层设计
		var userMessageTimeLimit int64
		userMessageTimeLimit, err = BaseConfig.GetDataInt64("UserMessageTimeLimit")
		if err != nil {
			userMessageTimeLimit = 5
		}
		err = Router2SystemConfig.MainDB.Get(&messageID, "SELECT id FROM user_message WHERE send_user_id = $1 AND create_at > $2 LIMIT 1", args.SendUserID, CoreFilter.GetNowTimeCarbon().SubSeconds(int(userMessageTimeLimit)).Time)
		if err == nil || messageID > 0 {
			err = errors.New("time too short")
			return
		}
	}
	//修正附加文件列
	if args.Files == nil {
		args.Files = []int64{}
	}
	//创建数据
	messageID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO user_message (wait_send_at, status, send_user_id, receive_user_id, title, content, files, params) VALUES (:wait_send_at, 0, :send_user_id, :receive_user_id, :title, :content, :files, :params)", args)
	if err != nil {
		return
	}
	if messageID < 1 {
		err = errors.New("insert data")
		return
	}
	//获取消息
	data := getByID(messageID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//统计数据
	if args.SendUserID > 0 {
		var count int64
		_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM user_message WHERE send_user_id = $1", args.SendUserID)
		AnalysisAny2.AppendData("re", "user_message_send_count", time.Time{}, 0, args.SendUserID, args.SendUserID, 0, 0, count)
	}
	if args.ReceiveUserID > 0 {
		var countAll int64
		_ = Router2SystemConfig.MainDB.Get(&countAll, "SELECT COUNT(id) FROM user_message WHERE receive_user_id = $1", args.ReceiveUserID)
		if countAll < 1 {
			countAll = 0
		}
		AnalysisAny2.AppendData("re", "user_message_receive_count", time.Time{}, 0, args.ReceiveUserID, args.ReceiveUserID, 0, 0, countAll)
		var countUnRead int64
		_ = Router2SystemConfig.MainDB.Get(&countUnRead, "SELECT COUNT(id) FROM user_message WHERE receive_user_id = $1 AND receive_delete_at < to_timestamp(1000000) AND receive_read_at < to_timestamp(1000000)", args.ReceiveUserID)
		AnalysisAny2.AppendData("re", "user_message_receive_unread_count", time.Time{}, 0, args.ReceiveUserID, args.ReceiveUserID, 0, 0, countUnRead)
	}
	//根据时间设置到期发送提醒
	BaseExpireTip.AppendTipNoErr(&BaseExpireTip.ArgsAppendTip{
		OrgID:      0,
		UserID:     0,
		SystemMark: "user_message",
		BindID:     data.ID,
		Hash:       "",
		ExpireAt:   data.WaitSendAt,
	})
	//清理缓冲
	if args.ReceiveUserID > 0 {
		deleteMessageReceiveCountCache(args.ReceiveUserID)
	}
	//反馈
	return
}
