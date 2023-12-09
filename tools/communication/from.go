package ToolsCommunication

import (
	"errors"
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetFromList 获取房间的参与列表参数
type ArgsGetFromList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id" empty:"true"`
	//来源系统
	FromSystem int `db:"from_system" json:"fromSystem" check:"intThan0" empty:"true"`
	//来源ID
	FromID int64 `db:"from_id" json:"fromID" check:"id" empty:"true"`
	//是否允许发言
	NeedAllowSend bool `db:"need_allow_send" json:"needAllowSend" check:"bool"`
	AllowSend     bool `db:"allow_send" json:"allowSend" check:"bool"`
	//角色类型
	// 0 普通; 1 房主; 2 副房主
	Role int `db:"role" json:"role" check:"intThan0" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetFromList 获取房间的参与列表
func GetFromList(args *ArgsGetFromList) (dataList []FieldsFrom, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.RoomID > -1 {
		where = where + "room_id = :room_id"
		maps["room_id"] = args.RoomID
	}
	if args.FromSystem > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "from_system = :from_system"
		maps["from_system"] = args.FromSystem
		if args.FromID > -1 {
			where = where + " AND from_id = :from_id"
			maps["from_id"] = args.FromID
		}
	}
	if args.NeedAllowSend {
		if where != "" {
			where = where + " AND "
		}
		where = where + "allow_send = :allow_send"
		maps["allow_send"] = args.AllowSend
	}
	if args.Role > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "role = :role"
		maps["role"] = args.Role
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "tools_communication_from"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, expire_at, room_id, from_system, from_id, name, allow_send, role, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "expire_at"},
	)
	return
}

// ArgsGetFrom 检查指定来源数据包参数
type ArgsGetFrom struct {
	//来源系统
	FromSystem int `db:"from_system" json:"fromSystem" check:"intThan0"`
	//来源ID
	FromID int64 `db:"from_id" json:"fromID" check:"id"`
}

// GetFrom 检查指定来源数据包
func GetFrom(args *ArgsGetFrom) (dataList []FieldsFrom, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, expire_at, room_id, from_system, from_id, name, token, allow_send, role, params FROM tools_communication_from WHERE from_system = $1 AND from_id = $2", args.FromSystem, args.FromID)
	return
}

// ArgsGetFromRoomID 检查指定来源数据包参数
type ArgsGetFromRoomID struct {
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id"`
	//来源系统
	FromSystem int `db:"from_system" json:"fromSystem" check:"intThan0"`
	//来源ID
	FromID int64 `db:"from_id" json:"fromID" check:"id"`
}

// GetFromRoomID 检查指定来源数据包
func GetFromRoomID(args *ArgsGetFromRoomID) (data FieldsFrom, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, room_id, from_system, from_id, name, token, allow_send, role, params FROM tools_communication_from WHERE from_system = $1 AND from_id = $2 AND room_id = $3 AND (expire_at < to_timestamp(1000000) OR expire_at >= NOW())", args.FromSystem, args.FromID, args.RoomID)
	return
}

// ArgsCheckFromAndRoom 检查房间和来源关联性参数
type ArgsCheckFromAndRoom struct {
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id"`
	//来源系统
	FromSystem int `db:"from_system" json:"fromSystem" check:"intThan0"`
	//来源ID
	FromID int64 `db:"from_id" json:"fromID" check:"id"`
}

// CheckFromAndRoom 检查房间和来源关联性
func CheckFromAndRoom(args *ArgsCheckFromAndRoom) (data FieldsFrom, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, room_id, from_system, from_id, name, token, allow_send, role, params FROM tools_communication_from WHERE from_system = $1 AND from_id = $2 AND room_id = $3", args.FromSystem, args.FromID, args.RoomID)
	return
}

// ArgsAppendRoom 加入房间参数
type ArgsAppendRoom struct {
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id"`
	//链接方式
	// 0 系统自带TCP握手方式; 1 系统自带RTC方式; 2 第三方agora服务
	ConnectType int `db:"connect_type" json:"connectType" check:"intThan0" empty:"true"`
	//来源系统
	FromSystem int `db:"from_system" json:"fromSystem" check:"intThan0"`
	//来源ID
	FromID int64 `db:"from_id" json:"fromID" check:"id"`
	//昵称
	Name string `db:"name" json:"name" check:"name"`
	//链接token
	// 用于第三方链接用
	Token string `db:"token" json:"token"`
	//是否允许发言
	AllowSend bool `db:"allow_send" json:"allowSend" check:"bool"`
	//角色类型
	// 0 普通; 1 房主; 2 副房主
	Role int `db:"role" json:"role" check:"intThan0" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// AppendRoom 加入房间
func AppendRoom(args *ArgsAppendRoom) (data FieldsFrom, err error) {
	//获取房间数据
	var roomData FieldsRoom
	roomData, err = GetRoom(&ArgsGetRoom{
		ID:    args.RoomID,
		OrgID: -1,
	})
	if err != nil {
		err = errors.New("room not exist, " + err.Error())
		return
	}
	//不能超出房间人数限制
	var roomCount int64
	roomCount, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "tools_communication_from", "id", "room_id = :room_id", map[string]interface{}{
		"room_id": roomData.ID,
	})
	if err != nil {
		roomCount = 0
	}
	if roomCount+1 > int64(roomData.MaxCount) {
		err = errors.New("room too many from")
		return
	}
	//写入房间ID
	args.RoomID = roomData.ID
	//检查房间和来源是否存在数据？
	var fromData FieldsFrom
	if err = Router2SystemConfig.MainDB.Get(&fromData, "SELECT id FROM tools_communication_from WHERE from_system = $1 AND from_id = $2 AND room_id = $3", args.FromSystem, args.FromID, args.RoomID); err == nil && fromData.ID > 0 {
		err = errors.New("already in the room")
		return
	}
	//获取过期时间
	expireTime := getDefaultFromExpire()
	expireAt := CoreFilter.GetNowTimeCarbon().AddSeconds(expireTime)
	//识别连接方式
	if args.ConnectType > -1 {
		roomData.ConnectType = args.ConnectType
	}
	switch roomData.ConnectType {
	case 0:
	case 1:
	case 2:
		args.Token, err = MakeAgoraToken(&ArgsMakeAgoraToken{
			RoomID:      roomData.ID,
			FromSystem:  args.FromSystem,
			FromID:      args.FromID,
			ExpireAt:    expireAt.Time,
			ConnectType: roomData.ConnectType,
		})
		if err != nil {
			err = errors.New("make agora token, " + err.Error())
			return
		}
	case 3:
		args.Token, err = MakeAgoraToken(&ArgsMakeAgoraToken{
			RoomID:      roomData.ID,
			FromSystem:  args.FromSystem,
			FromID:      args.FromID,
			ExpireAt:    expireAt.Time,
			ConnectType: roomData.ConnectType,
		})
		if err != nil {
			err = errors.New("make agora token, " + err.Error())
			return
		}
	case 4:
		args.Token, err = MakeAgoraToken(&ArgsMakeAgoraToken{
			RoomID:      roomData.ID,
			FromSystem:  args.FromSystem,
			FromID:      args.FromID,
			ExpireAt:    expireAt.Time,
			ConnectType: roomData.ConnectType,
		})
		if err != nil {
			err = errors.New("make agora token, " + err.Error())
			return
		}
	case 5:
		args.Token, err = MakeAgoraToken(&ArgsMakeAgoraToken{
			RoomID:      roomData.ID,
			FromSystem:  args.FromSystem,
			FromID:      args.FromID,
			ExpireAt:    expireAt.Time,
			ConnectType: roomData.ConnectType,
		})
		if err != nil {
			err = errors.New("make agora token, " + err.Error())
			return
		}
	default:
		err = errors.New("room connect type error")
		return
	}
	//创建数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "tools_communication_from", "INSERT INTO tools_communication_from (expire_at, room_id, from_system, from_id, name, token, allow_send, role, params) VALUES (:expire_at,:room_id,:from_system,:from_id,:name,:token,:allow_send,:role,:params)", map[string]interface{}{
		"expire_at":   expireAt,
		"room_id":     args.RoomID,
		"from_system": args.FromSystem,
		"from_id":     args.FromID,
		"name":        args.Name,
		"token":       args.Token,
		"allow_send":  args.AllowSend,
		"role":        args.Role,
		"params":      args.Params,
	}, &data)
	return
}

// ArgsAppendRoomTwo 检查并建立双向聊天房间参数
type ArgsAppendRoomTwo struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//链接方式
	// 0 系统自带TCP握手方式; 1 系统自带RTC方式; 2 第三方agora服务
	ConnectType int `db:"connect_type" json:"connectType" check:"intThan0" empty:"true"`
	//通讯类型
	DataType int `db:"data_type" json:"dataType" check:"intThan0"`
	//分类
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//房间名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//封面ID
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//房间链接密码
	Password string `db:"password" json:"password" check:"password" empty:"true"`
	//来源系统
	FromSystem int `db:"from_system" json:"fromSystem" check:"intThan0"`
	//来源ID
	FromID int64 `db:"from_id" json:"fromID" check:"id"`
	//昵称
	FromName string `db:"from_name" json:"fromName" check:"name"`
	//链接方式
	// 0 系统自带TCP握手方式; 1 系统自带RTC方式; 2 第三方agora服务
	FromConnectType int `db:"from_connect_type" json:"fromConnectType" check:"intThan0" empty:"true"`
	//到达系统
	ToSystem int   `db:"to_system" json:"toSystem" check:"intThan0"`
	ToID     int64 `db:"to_id" json:"toID" check:"id"`
	//到达昵称
	ToName string `db:"to_name" json:"toName" check:"name"`
	//链接方式
	// 0 系统自带TCP握手方式; 1 系统自带RTC方式; 2 第三方agora服务
	ToConnectType int `db:"to_connect_type" json:"toConnectType" check:"intThan0" empty:"true"`
}

// AppendRoomTwo 检查并建立双向聊天房间
// 两个来源建立，第一个为房主
func AppendRoomTwo(args *ArgsAppendRoomTwo) (roomData FieldsRoom, oneData FieldsFrom, twoData FieldsFrom, err error) {
	//检查是否存在已建立的数据？
	// 如果存在将直接反馈数据
	var findList []FieldsFrom
	err = Router2SystemConfig.MainDB.Select(&findList, "SELECT id, create_at, expire_at, room_id, from_system, from_id, name, token, allow_send, role, params FROM tools_communication_from WHERE from_system = $1 AND from_id = $2", args.FromSystem, args.FromID)
	if err == nil {
		for _, v := range findList {
			var vFind FieldsFrom
			err = Router2SystemConfig.MainDB.Get(&vFind, "SELECT id, create_at, expire_at, room_id, from_system, from_id, name, token, allow_send, role, params FROM tools_communication_from WHERE from_system = $1 AND from_id = $2 AND room_id = $3", args.ToSystem, args.ToID, v.RoomID)
			if err == nil && vFind.ID > 0 {
				roomData, err = GetRoom(&ArgsGetRoom{
					ID:    v.RoomID,
					OrgID: -1,
				})
				if err == nil && roomData.ID > 0 {
					oneData = v
					twoData = vFind
					return
				} else {
					_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "tools_communication_from", "room_id = :room_id", map[string]interface{}{
						"room_id": v.RoomID,
					})
				}
			}
		}
	}
	//构建新的房间
	roomData, err = CreateRoom(&ArgsCreateRoom{
		ExpireAt:    CoreFilter.GetNowTimeCarbon().AddMinutes(5).Time,
		ConnectType: args.ConnectType,
		DataType:    args.DataType,
		SortID:      args.SortID,
		Tags:        args.Tags,
		OrgID:       args.OrgID,
		Name:        args.Name,
		Des:         args.Des,
		CoverFileID: args.CoverFileID,
		IsPublic:    false,
		Password:    args.Password,
		MaxCount:    2,
		Params:      []CoreSQLConfig.FieldsConfigType{},
	})
	if err != nil {
		return
	}
	//加入房间开始聊天
	//获取过期时间
	expireTime := getDefaultFromExpire()
	expireAt := CoreFilter.GetNowTimeCarbon().AddSeconds(expireTime)
	//识别连接方式
	if args.FromConnectType < 0 {
		args.FromConnectType = roomData.ConnectType
	}
	var fromToken string
	switch args.FromConnectType {
	case 0:
	case 1:
	case 2:
		fromToken, err = MakeAgoraToken(&ArgsMakeAgoraToken{
			RoomID:      roomData.ID,
			FromSystem:  args.FromSystem,
			FromID:      args.FromID,
			ExpireAt:    expireAt.Time,
			ConnectType: args.FromConnectType,
		})
		if err != nil {
			err = errors.New("make agora token, " + err.Error())
			return
		}
	case 3:
		fromToken, err = MakeAgoraToken(&ArgsMakeAgoraToken{
			RoomID:      roomData.ID,
			FromSystem:  args.FromSystem,
			FromID:      args.FromID,
			ExpireAt:    expireAt.Time,
			ConnectType: args.FromConnectType,
		})
		if err != nil {
			err = errors.New("make agora token, " + err.Error())
			return
		}
	case 4:
		fromToken, err = MakeAgoraToken(&ArgsMakeAgoraToken{
			RoomID:      roomData.ID,
			FromSystem:  args.FromSystem,
			FromID:      args.FromID,
			ExpireAt:    expireAt.Time,
			ConnectType: args.FromConnectType,
		})
		if err != nil {
			err = errors.New("make agora token, " + err.Error())
			return
		}
	case 5:
		fromToken, err = MakeAgoraToken(&ArgsMakeAgoraToken{
			RoomID:      roomData.ID,
			FromSystem:  args.FromSystem,
			FromID:      args.FromID,
			ExpireAt:    expireAt.Time,
			ConnectType: args.FromConnectType,
		})
		if err != nil {
			err = errors.New("make agora token, " + err.Error())
			return
		}
	default:
		err = errors.New("room connect type error")
		return
	}
	//创建数据One
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "tools_communication_from", "INSERT INTO tools_communication_from (expire_at, room_id, from_system, from_id, name, token, allow_send, role, params) VALUES (:expire_at,:room_id,:from_system,:from_id,:name,:token,:allow_send,:role,:params)", map[string]interface{}{
		"expire_at":   expireAt,
		"room_id":     roomData.ID,
		"from_system": args.FromSystem,
		"from_id":     args.FromID,
		"name":        args.FromName,
		"token":       fromToken,
		"allow_send":  true,
		"role":        1,
		"params": CoreSQLConfig.FieldsConfigsType{
			{
				Mark: "connectType",
				Val:  fmt.Sprint(args.FromConnectType),
			},
		},
	}, &oneData)
	if err != nil {
		return
	}
	//识别连接方式Two
	if args.ToConnectType < 0 {
		args.ToConnectType = roomData.ConnectType
	}
	var toToken string
	switch args.ToConnectType {
	case 0:
	case 1:
	case 2:
		toToken, err = MakeAgoraToken(&ArgsMakeAgoraToken{
			RoomID:      roomData.ID,
			FromSystem:  args.ToSystem,
			FromID:      args.ToID,
			ExpireAt:    expireAt.Time,
			ConnectType: args.ToConnectType,
		})
		if err != nil {
			err = errors.New("make agora token, " + err.Error())
			return
		}
	case 3:
		toToken, err = MakeAgoraToken(&ArgsMakeAgoraToken{
			RoomID:      roomData.ID,
			FromSystem:  args.ToSystem,
			FromID:      args.ToID,
			ExpireAt:    expireAt.Time,
			ConnectType: args.ToConnectType,
		})
		if err != nil {
			err = errors.New("make agora token, " + err.Error())
			return
		}
	case 4:
		toToken, err = MakeAgoraToken(&ArgsMakeAgoraToken{
			RoomID:      roomData.ID,
			FromSystem:  args.ToSystem,
			FromID:      args.ToID,
			ExpireAt:    expireAt.Time,
			ConnectType: args.ToConnectType,
		})
		if err != nil {
			err = errors.New("make agora token, " + err.Error())
			return
		}
	case 5:
		toToken, err = MakeAgoraToken(&ArgsMakeAgoraToken{
			RoomID:      roomData.ID,
			FromSystem:  args.ToSystem,
			FromID:      args.ToID,
			ExpireAt:    expireAt.Time,
			ConnectType: args.ToConnectType,
		})
		if err != nil {
			err = errors.New("make agora token, " + err.Error())
			return
		}
	default:
		err = errors.New("room connect type error")
		return
	}
	//创建数据Two
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "tools_communication_from", "INSERT INTO tools_communication_from (expire_at, room_id, from_system, from_id, name, token, allow_send, role, params) VALUES (:expire_at,:room_id,:from_system,:from_id,:name,:token,:allow_send,:role,:params)", map[string]interface{}{
		"expire_at":   expireAt,
		"room_id":     roomData.ID,
		"from_system": args.ToSystem,
		"from_id":     args.ToID,
		"name":        args.ToName,
		"token":       toToken,
		"allow_send":  true,
		"role":        0,
		"params": CoreSQLConfig.FieldsConfigsType{
			{
				Mark: "connectType",
				Val:  fmt.Sprint(args.ToConnectType),
			},
		},
	}, &twoData)
	if err != nil {
		return
	}
	return
}

// ArgsUpdateFromExpire 更新来源的到期时间参数
type ArgsUpdateFromExpire struct {
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id"`
	//来源系统
	FromSystem int `db:"from_system" json:"fromSystem" check:"intThan0"`
	//来源ID
	FromID int64 `db:"from_id" json:"fromID" check:"id"`
}

// UpdateFromExpire 更新来源的到期时
func UpdateFromExpire(args *ArgsUpdateFromExpire) (err error) {
	expireAt := CoreFilter.GetNowTimeCarbon().AddSeconds(getDefaultFromExpire())
	var roomData FieldsRoom
	roomData, err = GetRoom(&ArgsGetRoom{
		ID:    args.RoomID,
		OrgID: -1,
	})
	if err != nil {
		return
	}
	var fromData FieldsFrom
	err = Router2SystemConfig.MainDB.Get(&fromData, "SELECT id, expire_at, from_system, from_id, params FROM tools_communication_from WHERE room_id = $1 AND from_system = $2 AND from_id = $3", args.RoomID, args.FromSystem, args.FromID)
	if err != nil || fromData.ID < 1 {
		err = errors.New(fmt.Sprint("not find from data, ", err))
		return
	}
	for _, v := range fromData.Params {
		if v.Mark == "connectType" {
			var c int
			c, err = CoreFilter.GetIntByString(v.Val)
			if err != nil {
				break
			}
			roomData.ConnectType = c
			break
		}
	}
	if fromData.ExpireAt.Unix()-60 < CoreFilter.GetNowTime().Unix() {
		var newToken string
		newToken, err = MakeAgoraToken(&ArgsMakeAgoraToken{
			RoomID:      roomData.ID,
			FromSystem:  fromData.FromSystem,
			FromID:      fromData.FromID,
			ExpireAt:    expireAt.Time,
			ConnectType: roomData.ConnectType,
		})
		if err != nil {
			err = errors.New("make agora token, " + err.Error())
			return
		}
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE tools_communication_from SET expire_at = :expire_at, token = :token WHERE from_system = :from_system AND from_id = :from_id", map[string]interface{}{
			"expire_at":   expireAt.Time,
			"token":       newToken,
			"from_system": args.FromSystem,
			"from_id":     args.FromID,
		})
		if err != nil {
			return
		}
	}
	if roomData.ExpireAt.Unix()-60 < CoreFilter.GetNowTime().Unix() {
		err = UpdateRoomExpire(&ArgsUpdateRoomExpire{
			ID:     roomData.ID,
			OrgID:  -1,
			FromID: -1,
		})
		if err != nil {
			return
		}
	}
	if fromData.FromSystem == 2 {
		roomData, err = GetRoom(&ArgsGetRoom{
			ID:    args.RoomID,
			OrgID: -1,
		})
		if err != nil {
			return
		}
		err = Router2SystemConfig.MainDB.Get(&fromData, "SELECT id, create_at, expire_at, room_id, from_system, from_id, name, token, allow_send, role, params FROM tools_communication_from WHERE room_id = $1 AND from_system = $2 AND from_id = $3", args.RoomID, args.FromSystem, args.FromID)
		if err != nil || fromData.ID < 1 {
			err = errors.New(fmt.Sprint("not find from data, ", err))
			return
		}
		pushNew(fromData.ID, roomData, fromData)
	}
	return
}

// ArgsOutRoom 退出房间参数
type ArgsOutRoom struct {
	//来源ID
	ID int64 `db:"id" json:"id" check:"id"`
	//来源系统
	FromSystem int `db:"from_system" json:"fromSystem" check:"intThan0"`
	//来源ID
	FromID int64 `db:"from_id" json:"fromID" check:"id"`
}

// OutRoom 退出房间
func OutRoom(args *ArgsOutRoom) (err error) {
	var fromData FieldsFrom
	err = Router2SystemConfig.MainDB.Get(&fromData, "SELECT id, room_id FROM tools_communication_from WHERE id = $1 AND from_system = $2 AND from_id = $3", args.ID, args.FromSystem, args.FromID)
	if err != nil {
		err = errors.New(fmt.Sprint("get room from by id: ", args.ID, ", from system: ", args.FromSystem, ", from id: ", args.FromID, ", err: ", err))
		return
	}
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "tools_communication_from", "id", map[string]interface{}{
		"id": args.ID,
	})
	if err != nil {
		return
	} else {
		_ = pushRoomOrFromDelete(0, fromData.ID)
	}
	var count int64
	count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "tools_communication_from", "id", "room_id = :room_id", map[string]interface{}{
		"room_id": fromData.RoomID,
	})
	if err != nil || count < 2 {
		err = DeleteRoom(&ArgsDeleteRoom{
			ID:     fromData.RoomID,
			OrgID:  -1,
			FromID: -1,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("delete room, ", err))
		}
	}
	return
}

// 获取默认过期时间
func getDefaultFromExpire() int {
	data, err := BaseConfig.GetDataInt64("ToolsCommunicationFromExpire")
	if err != nil {
		data = 180
	}
	return int(data)
}
