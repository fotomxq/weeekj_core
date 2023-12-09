package UserFocus2

import (
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	BlogCoreMod "gitee.com/weeekj/weeekj_core/v5/blog/core/mod"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceInfoExchangeMod "gitee.com/weeekj/weeekj_core/v5/service/info_exchange/mod"
	UserCoreMod "gitee.com/weeekj/weeekj_core/v5/user/core/mod"
	UserMessageMod "gitee.com/weeekj/weeekj_core/v5/user/message/mod"
	"time"
)

// ArgsSetFocus 设置是否关注参数
type ArgsSetFocus struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//关注类型
	// focus 关注; like 喜欢
	Mark string `db:"mark" json:"mark" check:"mark"`
	//关注来源
	System string `db:"system" json:"system" check:"mark"`
	//来源ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//是否删除关注
	IsFocus bool `json:"isFocus" check:"bool"`
}

// SetFocus 设置是否关注
func SetFocus(args ArgsSetFocus) (err error) {
	//检查行为
	if err = checkMark(args.Mark); err != nil {
		return
	}
	//检查系统
	if err = checkSystem(args.System); err != nil {
		return
	}
	//获取是否关注了该数据
	cacheMark := getFocusUserCacheMark(args.UserID, args.Mark, args.System, args.BindID)
	var id int64
	id, err = Router2SystemConfig.MainCache.GetInt64(cacheMark)
	if err != nil || id < 1 {
		_ = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM user_focus2 WHERE user_id = $1 AND mark = $2 AND system = $3 AND bind_id = $4", args.UserID, args.Mark, args.System, args.BindID)
		err = nil
	}
	if id > 0 {
		if !args.IsFocus {
			_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "user_focus2", "id", map[string]interface{}{
				"id": id,
			})
			if err != nil {
				return
			}
			deleteCache(args.UserID, args.Mark, args.System, args.BindID)
		}
	} else {
		if args.IsFocus {
			_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_focus2(user_id, mark, system, bind_id) VALUES(:user_id, :mark, :system, :bind_id)", map[string]interface{}{
				"user_id": args.UserID,
				"mark":    args.Mark,
				"system":  args.System,
				"bind_id": args.BindID,
			})
			if err != nil {
				return
			}
		}
	}
	//核对和给用户发送消息
	if args.IsFocus {
		userFocusNewFocusSendMessage := BaseConfig.GetDataBoolNoErr("UserFocusNewFocusSendMessage")
		userFocusNewLikeSendMessage := BaseConfig.GetDataBoolNoErr("UserFocusNewLikeSendMessage")
		userFocusNewCollectionSendMessage := BaseConfig.GetDataBoolNoErr("UserFocusNewCollectionSendMessage")
		needSendMsg := userFocusNewFocusSendMessage || userFocusNewLikeSendMessage || userFocusNewCollectionSendMessage
		if needSendMsg {
			var targetUserID int64 = 0
			var targetUserName = ""
			var msgDes = ""
			switch args.System {
			case "blog_content":
				vData := BlogCoreMod.GetContentByIDNoErr(args.BindID, -1)
				targetUserID = vData.UserID
				targetUserName = UserCoreMod.GetUserNiceNameByID(targetUserID)
				msgDes = fmt.Sprint("的博客发布的文章《", vData.Title, "》")
			case "mall_product":
				//vData := MallCoreMod.GetProductNoErr(args.BindID, -1)
				//targetUserID = 0
			case "user_core":
				vData := UserCoreMod.GetUserByID(args.BindID, -1)
				targetUserID = vData.ID
				targetUserName = vData.Name
			case "info_exchange":
				vData := ServiceInfoExchangeMod.GetInfoID(args.BindID, -1, -1)
				targetUserID = vData.UserID
				targetUserName = UserCoreMod.GetUserNiceNameByID(targetUserID)
				msgDes = fmt.Sprint("的帖子《", vData.Title, "》")
			}
			if targetUserID > 0 {
				if targetUserName == "" {
					targetUserName = UserCoreMod.GetUserNiceNameByID(targetUserID)
				}
				switch args.Mark {
				case "focus":
					if userFocusNewFocusSendMessage {
						UserMessageMod.CreateSystemToUser(time.Time{}, targetUserID, "有人关注了你", fmt.Sprint("用户(", targetUserName, ")关注了你", msgDes), nil, nil)
					}
				case "like":
					if userFocusNewLikeSendMessage {
						UserMessageMod.CreateSystemToUser(time.Time{}, targetUserID, "有人为你点赞", fmt.Sprint("用户(", targetUserName, ")为你", msgDes, "点赞"), nil, nil)
					}
				case "collection":
					if userFocusNewLikeSendMessage {
						UserMessageMod.CreateSystemToUser(time.Time{}, targetUserID, "有人收藏你发布的内容", fmt.Sprint("用户(", targetUserName, ")收藏了你", msgDes), nil, nil)
					}
				}
			}
		}
	}
	//反馈
	return
}
