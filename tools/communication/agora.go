package ToolsCommunication

import (
	"fmt"
	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/RtcTokenBuilder"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"strconv"
	"time"
)

//声网服务组件

// ArgsMakeAgoraToken 生成声网token参数
type ArgsMakeAgoraToken struct {
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id"`
	//来源系统
	FromSystem int `db:"from_system" json:"fromSystem" check:"intThan0"`
	//来源ID
	FromID int64 `db:"from_id" json:"fromID" check:"id"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//链接方式
	// 0 系统自带TCP握手方式; 1 系统自带RTC方式; 2 第三方agora服务字符串; 3 第三方agora服务uint32; 4 第三方agora服务字符串trc; 5 第三方agora服务uint32 rtc
	ConnectType int `db:"connect_type" json:"connectType"`
}

// MakeAgoraToken 生成声网token
func MakeAgoraToken(args *ArgsMakeAgoraToken) (data string, err error) {
	var appID, appCert string
	appID, err = BaseConfig.GetDataString("ToolsCommunicationAgoraAppID")
	if err != nil {
		return
	}
	appCert, err = BaseConfig.GetDataString("ToolsCommunicationAgoraAppKey")
	if err != nil {
		return
	}
	switch args.ConnectType {
	case 2:
		data, err = rtctokenbuilder.BuildTokenWithUserAccount(appID, appCert, fmt.Sprint(args.RoomID), fmt.Sprintf(strconv.Itoa(args.FromSystem), args.FromID), rtctokenbuilder.RoleAttendee, uint32(args.ExpireAt.Unix()))
		if err != nil {
			return
		}
	case 3:
		var uid int
		uid, err = CoreFilter.GetIntByString(fmt.Sprintf(strconv.Itoa(args.FromSystem), args.FromID))
		data, err = rtctokenbuilder.BuildTokenWithUID(appID, appCert, fmt.Sprint(args.RoomID), uint32(uid), rtctokenbuilder.RoleAttendee, uint32(args.ExpireAt.Unix()))
		if err != nil {
			return
		}
	case 4:
		data, err = rtctokenbuilder.BuildTokenWithUserAccountAndPrivilege(appID, appCert, fmt.Sprint(args.RoomID), fmt.Sprintf(strconv.Itoa(args.FromSystem), args.FromID), uint32(args.ExpireAt.Unix()), uint32(args.ExpireAt.Unix()), uint32(args.ExpireAt.Unix()), uint32(args.ExpireAt.Unix()))
		if err != nil {
			return
		}
	case 5:
		var uid int
		uid, err = CoreFilter.GetIntByString(fmt.Sprintf(strconv.Itoa(args.FromSystem), args.FromID))
		data, err = rtctokenbuilder.BuildTokenWithUIDAndPrivilege(appID, appCert, fmt.Sprint(args.RoomID), uint32(uid), uint32(args.ExpireAt.Unix()), uint32(args.ExpireAt.Unix()), uint32(args.ExpireAt.Unix()), uint32(args.ExpireAt.Unix()))
		if err != nil {
			return
		}
	}
	return
}
