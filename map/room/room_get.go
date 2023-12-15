package MapRoom

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetRoomList 获取房屋列表参数
type ArgsGetRoomList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//状态
	// 0 空闲; 1 有人; 2 退房; 3 不可用
	Status int `db:"status" json:"status" check:"intThan0" empty:"true"`
	//入驻人员
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetRoomList 获取房屋列表参数
func GetRoomList(args *ArgsGetRoomList) (dataList []FieldsRoom, dataCount int64, err error) {
	//获取数据包
	maps := map[string]interface{}{}
	where := ""
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.SortID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if len(args.Tags) > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "tags @> :tags"
		maps["tags"] = args.Tags
	}
	if args.Status > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "status = :status"
		maps["status"] = args.Status
	}
	if args.InfoID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + ":info_id = ANY(infos)"
		maps["info_id"] = args.InfoID
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%' OR code ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "map_room"
	var rawList []FieldsRoom
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "code"},
	)
	if err != nil {
		return
	}
	//覆盖数据
	for _, v := range rawList {
		vData := getRoomID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// ArgsGetRoomSearch 搜索房间参数
type ArgsGetRoomSearch struct {
	//最大个数
	Max int64 `db:"max" json:"max" check:"max"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

type DataGetRoomSearch struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//名称
	Name string `db:"name" json:"name"`
}

// GetRoomSearch 搜索房间
func GetRoomSearch(args *ArgsGetRoomSearch) (dataList []DataGetRoomSearch, err error) {
	//修正参数
	if args.Max < 1 {
		args.Max = 1
	}
	if args.Max > 30 {
		args.Max = 30
	}
	//获取数据
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, concat(name, '(', code, ')') as name FROM map_room WHERE org_id = $1 AND delete_at < to_timestamp(1000000) AND (name ILIKE '%' || $2 || '%' OR des ILIKE '%' || $2 || '%') ORDER BY code LIMIT $3", args.OrgID, args.Search, args.Max)
	if err != nil {
		return
	}
	//反馈
	return
}

// DataGetRoomAndInfosSearch 搜索房间数据据
type DataGetRoomAndInfosSearch struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//名称
	Name string `db:"name" json:"name"`
	//入驻人员列
	Infos pq.Int64Array `db:"infos" json:"infos"`
}

// GetRoomAndInfosSearch 搜索房间
func GetRoomAndInfosSearch(args *ArgsGetRoomSearch) (dataList []DataGetRoomAndInfosSearch, err error) {
	//修正参数
	if args.Max < 1 {
		args.Max = 1
	}
	if args.Max > 30 {
		args.Max = 30
	}
	//获取数据
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, concat(name, '(', code, ')') as name, infos FROM map_room WHERE org_id = $1 AND delete_at < to_timestamp(1000000) AND (name ILIKE '%' || $2 || '%' OR des ILIKE '%' || $2 || '%') ORDER BY code LIMIT $3", args.OrgID, args.Search, args.Max)
	if err != nil || len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	//反馈
	return
}

// GetAllRoomListByOrgID 获取组织下所有房间
func GetAllRoomListByOrgID(orgID int64) (dataList []FieldsRoom, err error) {
	//获取数据
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, status, name, code, cover_file_id, status, service_status, infos FROM map_room WHERE org_id = $1 AND delete_at < to_timestamp(1000000) ORDER BY code", orgID)
	if err != nil || len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	//反馈
	return
}

// GetRoomCountByOrgID 获取房间总数
func GetRoomCountByOrgID(orgID int64) (count int64) {
	//获取数据
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM map_room WHERE org_id = $1 AND delete_at < to_timestamp(1000000)", orgID)
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsGetRoomByInfo 获取人员入驻的房间参数
type ArgsGetRoomByInfo struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//入驻人员
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//状态
	// -1 不过滤状态; 0 空闲; 1 有人; 2 退房; 3 不可用
	Status int `db:"status" json:"status" check:"intThan0" empty:"true"`
}

// GetRoomByInfo 获取人员入驻的房间参数
func GetRoomByInfo(args *ArgsGetRoomByInfo) (data FieldsRoom, err error) {
	//获取数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM map_room WHERE delete_at < to_timestamp(1000000) AND ($1 < 1 OR org_id = $1) AND $2 = ANY(infos) AND ($3 < 0 OR status = $3) ORDER BY update_at DESC LIMIT 1", args.OrgID, args.InfoID, args.Status)
	if err != nil {
		return
	}
	//覆盖数据
	if data.ID > 0 {
		data = getRoomID(data.ID)
		if data.ID < 1 {
			err = errors.New("no data")
			return
		}
	}
	//反馈
	return
}

func GetRoomIDByInfo(infoID int64) (roomID int64) {
	cacheMark := getRoomInfoCacheMark(infoID)
	roomID, _ = Router2SystemConfig.MainCache.GetInt64(cacheMark)
	if roomID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&roomID, "SELECT id FROM map_room WHERE delete_at < to_timestamp(1000000) AND $1 = ANY(infos) ORDER BY update_at DESC LIMIT 1", infoID)
	Router2SystemConfig.MainCache.SetInt64(cacheMark, roomID, cacheTime)
	return
}

// GetRoomListByInfo 获取信息档案入住的房间列表
func GetRoomListByInfo(infoID int64) (dataList []FieldsRoom, err error) {
	//获取数据
	var rawList []FieldsRoom
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM map_room WHERE delete_at < to_timestamp(1000000) AND $1 = ANY(infos) ORDER BY update_at DESC", infoID)
	if err != nil {
		return
	}
	//覆盖数据
	for _, v := range rawList {
		vData := getRoomID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// ArgsGetRoomID 获取指定房间ID参数
type ArgsGetRoomID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetRoomID 获取指定房间ID
func GetRoomID(args *ArgsGetRoomID) (data FieldsRoom, err error) {
	data = getRoomID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	if !CoreFilter.EqID2(args.OrgID, data.OrgID) || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		err = errors.New("no data")
		return
	}
	return
}

// 获取房间信息
func getRoomID(id int64) (data FieldsRoom) {
	//获取缓冲
	cacheMark := getRoomCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	//获取数据
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, sort_id, tags, status, service_status, service_bind_id, service_mission_id, infos, code, name, des, cover_file_id, des_files, bg_color, params FROM map_room WHERE id = $1", id)
	if err != nil || data.ID < 1 {
		return
	}
	//写入缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheTime)
	//反馈
	return
}

// ArgsGetRoomMore 获取一组房间参数
type ArgsGetRoomMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetRoomMore 获取一组房间
func GetRoomMore(args *ArgsGetRoomMore) (dataList []FieldsRoom, err error) {
	//获取数据
	var rawList []FieldsRoom
	err = CoreSQLIDs.GetIDsAndDelete(&rawList, "map_room", "id", args.IDs, args.HaveRemove)
	if err != nil {
		return
	}
	//覆盖数据
	for _, v := range rawList {
		vData := getRoomID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

func GetRoomNames(args *ArgsGetRoomMore) (dataList map[int64]string, err error) {
	//获取数据
	dataList, err = CoreSQLIDs.GetIDsNameAndDelete("map_room", args.IDs, args.HaveRemove)
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsGetRoomCode 获取指定房间ID参数
type ArgsGetRoomCode struct {
	//房间编号
	Code string `db:"code" json:"code" check:"mark"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetRoomCode 获取指定房间ID
func GetRoomCode(args *ArgsGetRoomCode) (data FieldsRoom, err error) {
	//获取数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM map_room WHERE code = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000)", args.Code, args.OrgID)
	if err != nil {
		return
	}
	//覆盖数据
	if data.ID > 0 {
		data = getRoomID(data.ID)
		if data.ID < 1 {
			err = errors.New("no data")
			return
		}
	}
	//反馈
	return
}

// ArgsGetRooms 获取一组房间参数
type ArgsGetRooms struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetRooms 获取一组房间
func GetRooms(args *ArgsGetRooms) (dataList []FieldsRoom, err error) {
	//获取数据
	var rawList []FieldsRoom
	err = CoreSQLIDs.GetIDsAndDelete(&rawList, "map_room", "id", args.IDs, args.HaveRemove)
	if err != nil {
		return
	}
	//覆盖数据
	for _, v := range rawList {
		vData := getRoomID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

func GetRoomsName(args *ArgsGetRooms) (data map[int64]string, err error) {
	//获取数据
	data, err = CoreSQLIDs.GetIDsNameAndDelete("map_room", args.IDs, args.HaveRemove)
	if err != nil {
		return
	}
	//反馈
	return
}

func GetRoomName(id int64) (name string) {
	data := getRoomID(id)
	name = data.Name
	return
}

// ArgsGetRoomByOrg 获取组织下所有可用房间参数
type ArgsGetRoomByOrg struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//是否存在呼叫
	NeedCall bool `json:"needCall" check:"bool"`
	IsCall   bool `json:"isCall" check:"bool"`
}

// GetRoomByOrg 获取组织下所有可用房间
func GetRoomByOrg(args *ArgsGetRoomByOrg) (dataList []FieldsRoom, err error) {
	//获取房间数据
	var rawList []FieldsRoom
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM map_room WHERE org_id = $1 AND delete_at < to_timestamp(1000000) AND ($2 < 1 OR sort_id = $2) AND ($3 < 1 OR (($3 = true AND service_status != 0) OR ($3 = false AND service_status = 0)))", args.OrgID, args.SortID)
	if err != nil || len(dataList) < 1 {
		err = errors.New(fmt.Sprint("no data, ", err))
		return
	}
	//覆盖数据
	for _, v := range rawList {
		vData := getRoomID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// 获取房间数据包
func getRunSensorRoomList(limit int, step int) (roomList []FieldsRoom) {
	//获取房间数据
	if err := Router2SystemConfig.MainDB.Select(&roomList, "SELECT id, org_id FROM map_room WHERE delete_at < to_timestamp(1000000) LIMIT $1 OFFSET $2", limit, step); err != nil {
		return
	}
	//反馈
	return
}

func getAppendWarningRoom(roomID int64) (infoIDs pq.Int64Array) {
	//获取数据
	_ = Router2SystemConfig.MainDB.Get(&infoIDs, "SELECT infos FROM map_room WHERE id = $1", roomID)
	//反馈
	return
}
