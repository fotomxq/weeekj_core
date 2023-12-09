package ToolsCommunication

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "gitee.com/weeekj/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsGetRoomList 获取房间列表参数
type ArgsGetRoomList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//分类
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//链接方式
	// 0 系统自带TCP握手方式; 1 系统自带RTC方式; 2 第三方agora服务
	ConnectType int `db:"connect_type" json:"connectType" check:"intThan0" empty:"true"`
	//通讯类型
	DataType int `db:"data_type" json:"dataType" check:"intThan0" empty:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否公开房间？
	// 私有化房间只允许特定链接链接，否则可以通过公共列表查询到
	IsPublic bool `db:"is_public" json:"isPublic" check:"bool"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetRoomList 获取房间列表
func GetRoomList(args *ArgsGetRoomList) (dataList []FieldsRoom, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.SortID > -1 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if len(args.Tags) > 0 {
		where = where + " AND tags @> :tags"
		maps["tags"] = args.Tags
	}
	if args.ConnectType > -1 {
		where = where + " AND connect_type = :connect_type"
		maps["connect_type"] = args.ConnectType
	}
	if args.DataType > -1 {
		where = where + " AND data_type = :data_type"
		maps["data_type"] = args.DataType
	}
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.IsPublic {
		where = where + " AND is_public = :is_public"
		maps["is_public"] = args.IsPublic
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "tools_communication_room"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, expire_at, connect_type, data_type, sort_id, tags, org_id, name, des, cover_file_id, max_count, params, is_public, password FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "expire_at", "max_count"},
	)
	return
}

// ArgsGetRoom 获取指定房间信息参数
type ArgsGetRoom struct {
	//房间ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetRoom 获取指定房间信息
func GetRoom(args *ArgsGetRoom) (data FieldsRoom, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, expire_at, connect_type, data_type, sort_id, tags, org_id, name, des, cover_file_id, max_count, params, is_public, password FROM tools_communication_room WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000)", args.ID, args.OrgID)
	if err == nil && data.ID < 1 {
		err = errors.New("room not exist")
	}
	return
}

// ArgsGetRoomMore 获取多个房间参数
type ArgsGetRoomMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetRoomMore 获取多个房间
func GetRoomMore(args *ArgsGetRoomMore) (dataList []FieldsRoom, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "tools_communication_room", "id, create_at, update_at, delete_at, expire_at, connect_type, data_type, sort_id, tags, org_id, name, des, cover_file_id, max_count, params, is_public", args.IDs, args.HaveRemove)
	return
}

// ArgsCheckOrgAndRoom 检查房间和组织参数
type ArgsCheckOrgAndRoom struct {
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

func CheckOrgAndRoom(args *ArgsCheckOrgAndRoom) (err error) {
	var data FieldsRoom
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM tools_communication_room WHERE org_id = $1 AND id = $2", args.OrgID, args.RoomID)
	if err == nil && data.ID < 1 {
		err = errors.New("not exist")
	}
	return
}

// ArgsCreateRoom 创建房间参数
type ArgsCreateRoom struct {
	//到期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt" check:"isoTime"`
	//链接方式
	// 0 系统自带TCP握手方式; 1 系统自带RTC方式; 2 第三方agora服务
	ConnectType int `db:"connect_type" json:"connectType" check:"intThan0" empty:"true"`
	//通讯类型
	DataType int `db:"data_type" json:"dataType" check:"intThan0"`
	//分类
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//房间名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//封面ID
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//是否公开房间？
	// 私有化房间只允许特定链接链接，否则可以通过公共列表查询到
	IsPublic bool `db:"is_public" json:"isPublic" check:"bool"`
	//房间链接密码
	Password string `db:"password" json:"password" check:"password" empty:"true"`
	//最大人数
	MaxCount int `db:"max_count" json:"maxCount" check:"intThan0"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateRoom 创建房间
func CreateRoom(args *ArgsCreateRoom) (data FieldsRoom, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "tools_communication_room", "INSERT INTO tools_communication_room (sort_id, tags, expire_at, connect_type, data_type, org_id, name, des, cover_file_id, max_count, params, is_public, password) VALUES (:sort_id,:tags,:expire_at,:connect_type,:data_type,:org_id,:name,:des,:cover_file_id,:max_count,:params,:is_public,:password)", args, &data)
	return
}

// 更新房间过期时间
// 继续延续5分钟
type ArgsUpdateRoomExpire struct {
	//房间ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//来源
	FromID int64 `db:"from_id" json:"fromID" check:"id" empty:"true"`
}

func UpdateRoomExpire(args *ArgsUpdateRoomExpire) (err error) {
	if args.FromID > 0 {
		var data FieldsFrom
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM tools_communication_from WHERE room_id = $1 AND id = $2 AND (role = 1 OR role = 2)", args.ID, args.FromID)
		if err != nil || data.ID < 1 {
			err = errors.New("from not exist")
			return
		}
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tools_communication_room SET update_at = NOW(), expire_at = :expire_at WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
		"id":        args.ID,
		"org_id":    args.OrgID,
		"expire_at": CoreFilter.GetNowTimeCarbon().AddMinutes(5).Time,
	})
	if err == nil {
		var roomData FieldsRoom
		roomData, err = GetRoom(&ArgsGetRoom{
			ID:    args.ID,
			OrgID: -1,
		})
		if err == nil {
			_ = pushRoomInfo(roomData)
		}
	} else {
		err = errors.New("update room info, " + err.Error())
	}
	return
}

// ArgsUpdateRoom 修改房间参数
type ArgsUpdateRoom struct {
	//房间ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//来源
	FromID int64 `db:"from_id" json:"fromID" check:"id" empty:"true"`
	//分类
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//到期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt" check:"isoTime"`
	//房间名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//封面ID
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//是否公开房间？
	// 私有化房间只允许特定链接链接，否则可以通过公共列表查询到
	IsPublic bool `db:"is_public" json:"isPublic" check:"bool"`
	//房间链接密码
	Password string `db:"password" json:"password" check:"password" empty:"true"`
	//最大人数
	MaxCount int `db:"max_count" json:"maxCount" check:"intThan0"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateRoom 修改房间
func UpdateRoom(args *ArgsUpdateRoom) (err error) {
	if args.FromID > 0 {
		var data FieldsFrom
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM tools_communication_from WHERE room_id = $1 AND id = $2 AND (role = 1 OR role = 2)", args.ID, args.FromID)
		if err != nil || data.ID < 1 {
			err = errors.New("from not exist")
			return
		}
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tools_communication_room SET update_at = NOW(), sort_id = :sort_id, tags = :tags, expire_at = :expire_at, name = :name, des = :des, cover_file_id = :cover_file_id, max_count = :max_count, params = :params, is_public = :is_public, password = :password WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
		"id":            args.ID,
		"org_id":        args.OrgID,
		"sort_id":       args.SortID,
		"tags":          args.Tags,
		"expire_at":     args.ExpireAt,
		"name":          args.Name,
		"des":           args.Des,
		"cover_file_id": args.CoverFileID,
		"is_public":     args.IsPublic,
		"password":      args.Password,
		"max_count":     args.MaxCount,
		"params":        args.Params,
	})
	if err == nil {
		var roomData FieldsRoom
		roomData, err = GetRoom(&ArgsGetRoom{
			ID:    args.ID,
			OrgID: -1,
		})
		if err == nil {
			_ = pushRoomInfo(roomData)
		}
	} else {
		err = errors.New("update room info, " + err.Error())
	}
	return
}

// ArgsDeleteRoom 删除房间参数
type ArgsDeleteRoom struct {
	//房间ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//来源
	FromID int64 `db:"from_id" json:"fromID" check:"id" empty:"true"`
}

// DeleteRoom 删除房间
func DeleteRoom(args *ArgsDeleteRoom) (err error) {
	if args.FromID > 0 {
		var data FieldsFrom
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM tools_communication_from WHERE room_id = $1 AND id = $2 AND role = 1", args.ID, args.FromID)
		if err != nil || data.ID < 1 {
			err = errors.New("from not exist")
			return
		}
	}
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "tools_communication_room", "id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
		"id":     args.ID,
		"org_id": args.OrgID,
	})
	if err == nil {
		var data FieldsRoom
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, org_id FROM tools_communication_room WHERE id = $1", args.ID)
		if err == nil {
			_ = pushRoomOrFromDelete(args.ID, 0)
		}
	}
	return
}
