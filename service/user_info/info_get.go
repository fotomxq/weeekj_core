package ServiceUserInfo

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLIDs "gitee.com/weeekj/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"strings"
)

// ArgsGetInfoList 获取文件列表参数
type ArgsGetInfoList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 允许为0，则该信息不属于任何用户，或不和任何用户关联
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//从属关系
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//负责人
	Director int64 `db:"director" json:"director" check:"id" empty:"true"`
	//是否死亡
	IsDie bool `db:"is_die" json:"isDie" check:"bool"`
	//是否出院
	IsOut bool `db:"is_out" json:"isOut" check:"bool"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetInfoList 获取文件列表
func GetInfoList(args *ArgsGetInfoList) (dataList []FieldsInfo, dataCount int64, err error) {
	//获取数据
	maps := map[string]interface{}{}
	where := ""
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	where = CoreSQL.GetDeleteSQLField(args.IsDie, where, "die_at")
	where = CoreSQL.GetDeleteSQLField(args.IsOut, where, "out_at")
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.BindID > -1 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.Country > -1 {
		where = where + " AND country = :country"
		maps["country"] = args.Country
	}
	if args.SortID > 0 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if len(args.Tags) > 0 {
		where = where + " AND tags @> :tags"
		maps["tags"] = args.Tags
	}
	if args.Director > 0 {
		where = where + " AND (director_1 = :director OR director_2 = :director)"
		maps["director"] = args.Director
	}
	if args.Search != "" {
		args.Search = strings.ReplaceAll(args.Search, " ", "")
		//where = where + " AND (name ILIKE '%' || :search || '%' OR address ILIKE '%' || :search || '%' OR phone ILIKE '%' || :search || '%' OR profession ILIKE '%' || :search || '%' OR emergency_contact ILIKE '%' || :search || '%' OR emergency_contact_phone ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		where = where + " AND (REPLACE(name, ' ', '') ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "service_user_info"
	var rawList []FieldsInfo
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id "+"FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "date_of_birth", "level", "die_at", "out_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getInfoID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// GetInfoAllByFilter 获取符合条件的组织下所有老人
func GetInfoAllByFilter(orgID int64, sortID int64, tags pq.Int64Array, director int64, needIsDie bool, isDie bool, needIsOut bool, isOut bool, timeType int, betweenAt CoreSQLTime.DataCoreTime) (dataList []FieldsInfo, err error) {
	//获取数据
	maps := map[string]interface{}{}
	where := ""
	where = CoreSQL.GetDeleteSQL(false, where)
	if orgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = orgID
	}
	if sortID > 0 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = sortID
	}
	if len(tags) > 0 {
		where = where + " AND tags @> :tags"
		maps["tags"] = tags
	}
	if director > 0 {
		where = where + " AND (director_1 = :director OR director_2 = :director)"
		maps["director"] = director
	}
	var betweenAt2 CoreSQLTime.FieldsCoreTime
	betweenAt2, err = CoreSQLTime.GetBetweenByISO(betweenAt)
	if err != nil {
		return
	}
	switch timeType {
	case 0:
		//全部档案，时间范围无效
		// 标记时间范围无限大
		betweenAt2 = CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubYears(100).Time,
			MaxTime: CoreFilter.GetNowTimeCarbon().AddYears(100).Time,
		}
	case 1:
		//指定时间范围
	case 2:
		//结束日期之前
		betweenAt2.MinTime = CoreFilter.GetNowTimeCarbon().SubYears(100).Time
	case 3:
		//开始日期之后
		betweenAt2.MaxTime = CoreFilter.GetNowTimeCarbon().AddYears(100).Time
	}
	if needIsDie || needIsOut {
		if isDie && needIsDie {
			where, maps = CoreSQLTime.GetBetweenByTimeAnd("die_at", betweenAt2, where, maps)
		} else {
			if isOut && needIsOut {
				where, maps = CoreSQLTime.GetBetweenByTimeAnd("out_at", betweenAt2, where, maps)
			} else {
				where, maps = CoreSQLTime.GetBetweenByTimeAnd("create_at", betweenAt2, where, maps)
			}
		}
	} else {
		where, maps = CoreSQLTime.GetBetweenByTimeAnd("create_at", betweenAt2, where, maps)
	}
	where = " WHERE " + where
	var rawList []FieldsInfo
	err = CoreSQL.GetList(Router2SystemConfig.MainDB.DB, &rawList, "SELECT id FROM service_user_info"+where, maps)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getInfoID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// GetInfoByTagOrSort 获取符合标签和分类的所有人员列表
func GetInfoByTagOrSort(orgID int64, tagIDs pq.Int64Array, sortIDs pq.Int64Array, isExist bool) (dataList []FieldsInfo, err error) {
	//获取数据
	var rawList []FieldsInfo
	if isExist {
		err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM service_user_info WHERE org_id = $1 AND (tags @> $2 OR sort_id = ANY($3)) AND delete_at < to_timestamp(1000000) AND out_at < to_timestamp(1000000) AND die_at < to_timestamp(1000000)", orgID, tagIDs, sortIDs)
	} else {
		err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM service_user_info WHERE org_id = $1 AND (tags @> $2 OR sort_id = ANY($3))", orgID, tagIDs, sortIDs)
	}
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	//补充数据
	for _, v := range rawList {
		vData := getInfoID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

func GetInfoByTag(orgID int64, tagIDs pq.Int64Array, isExist bool) (dataList []FieldsInfo, err error) {
	//获取数据
	var rawList []FieldsInfo
	if isExist {
		err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM service_user_info WHERE org_id = $1 AND tags @> $2 AND delete_at < to_timestamp(1000000) AND out_at < to_timestamp(1000000) AND die_at < to_timestamp(1000000)", orgID, tagIDs)
	} else {
		err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM service_user_info WHERE org_id = $1 AND tags @> $2", orgID, tagIDs)
	}
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	//补充数据
	for _, v := range rawList {
		vData := getInfoID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

func GetInfoBySort(orgID int64, sortIDs pq.Int64Array, isExist bool) (dataList []FieldsInfo, err error) {
	//获取数据
	var rawList []FieldsInfo
	if isExist {
		err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM service_user_info WHERE org_id = $1 AND sort_id = ANY($2) AND delete_at < to_timestamp(1000000) AND out_at < to_timestamp(1000000) AND die_at < to_timestamp(1000000)", orgID, sortIDs)
	} else {
		err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM service_user_info WHERE org_id = $1 AND sort_id = ANY($2)", orgID, sortIDs)
	}
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	//补充数据
	for _, v := range rawList {
		vData := getInfoID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// ArgsGetInfoSearch 搜索房间参数
type ArgsGetInfoSearch struct {
	//最大个数
	Max int64 `db:"max" json:"max" check:"max"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

type DataGetInfoSearch struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//名称
	Name string `db:"name" json:"name"`
}

// GetInfoSearch 搜索房间
func GetInfoSearch(args *ArgsGetInfoSearch) (dataList []DataGetInfoSearch, err error) {
	//修正参数
	if args.Max < 1 {
		args.Max = 1
	}
	if args.Max > 30 {
		args.Max = 30
	}
	//获取数据
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, name FROM service_user_info WHERE org_id = $1 AND die_at < to_timestamp(1000000) AND delete_at < to_timestamp(1000000) AND (name ILIKE '%' || $2 || '%' OR des ILIKE '%' || $2 || '%') ORDER BY name LIMIT $3", args.OrgID, args.Search, args.Max)
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsGetInfoID 获取ID参数
type ArgsGetInfoID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// GetInfoID 获取ID
func GetInfoID(args *ArgsGetInfoID) (data FieldsInfo, err error) {
	data = getInfoID(args.ID)
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

// GetInfoByUserID 获取用户自己的档案
func GetInfoByUserID(userID int64) (data FieldsInfo) {
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM service_user_info WHERE user_id = $1 AND delete_at < to_timestamp(1000000)", userID)
	if data.ID < 1 {
		return
	}
	data = getInfoID(data.ID)
	if data.ID < 1 {
		return
	}
	return
}

// 获取档案信息
func getInfoID(id int64) (data FieldsInfo) {
	//获取缓冲
	cacheMark := getInfoCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	//获取数据
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, die_at, out_at, org_id, user_id, bind_id, bind_type, name, country, gender, id_card, id_card_front_file_id, id_card_back_file_id, phone, cover_file_id, des_files, address, date_of_birth, marital_status, education_status, profession, level, emergency_contact, emergency_contact_phone, sort_id, tags, doc_id, des, director_1, director_2, params FROM service_user_info WHERE id = $1", id)
	if err != nil {
		return
	}
	//写入缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheTime)
	//反馈
	return
}

// ArgsGetInfoMore 获取多个信息ID参数
type ArgsGetInfoMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetInfoMore 获取多个信息ID
func GetInfoMore(args *ArgsGetInfoMore) (dataList []FieldsInfo, err error) {
	//获取数据
	var rawList []FieldsInfo
	err = CoreSQLIDs.GetIDsAndDelete(&rawList, "service_user_info", "id", args.IDs, args.HaveRemove)
	if err != nil {
		return
	}
	//补充数据
	for _, v := range rawList {
		vData := getInfoID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

func GetInfoMoreNames(args *ArgsGetInfoMore) (data map[int64]string, err error) {
	//获取数据
	data, err = CoreSQLIDs.GetIDsNameAndDelete("service_user_info", args.IDs, args.HaveRemove)
	if err != nil {
		return
	}
	//反馈
	return
}

func GetInfoName(id int64) string {
	data := getInfoID(id)
	if data.ID < 1 {
		return ""
	}
	return data.Name
}

// ArgsGetOrgInfoMore 获取多个信息ID带组织参数
type ArgsGetOrgInfoMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// GetOrgInfoMore 获取多个信息ID带组织
func GetOrgInfoMore(args *ArgsGetOrgInfoMore) (dataList []FieldsInfo, err error) {
	//获取数据
	var rawList []FieldsInfo
	err = CoreSQLIDs.GetIDsOrgAndDelete(&rawList, "service_user_info", "id", args.IDs, args.OrgID, args.HaveRemove)
	if err != nil {
		return
	}
	//补充数据
	for _, v := range rawList {
		vData := getInfoID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

func GetOrgInfoMoreNames(args *ArgsGetOrgInfoMore) (data map[int64]string, err error) {
	//获取数据
	data, err = CoreSQLIDs.GetIDsOrgNameAndDelete("service_user_info", args.IDs, args.OrgID, args.HaveRemove)
	if err != nil {
		return
	}
	//反馈
	return
}

// CheckInfoIsExist 检查信息档案是否存在
func CheckInfoIsExist(infoID int64) (err error) {
	//获取数据
	infoData := getInfoID(infoID)
	if infoData.ID < 1 {
		err = errors.New("no data")
		return
	}
	if CoreSQL.CheckTimeHaveData(infoData.DeleteAt) || CoreSQL.CheckTimeHaveData(infoData.DieAt) || CoreSQL.CheckTimeHaveData(infoData.OutAt) {
		err = errors.New("no data")
		return
	}
	//反馈
	return
}

// ArgsSearchNameOrIDCard 搜索姓名和手机号参数
type ArgsSearchNameOrIDCard struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

type DataSearchNameOrIDCard struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//姓名
	Name string `db:"name" json:"name"`
	//证件编号
	IDCard string `db:"id_card" json:"idCard"`
}

// SearchNameOrIDCard 搜索姓名和手机号
func SearchNameOrIDCard(args *ArgsSearchNameOrIDCard) (dataList []DataSearchNameOrIDCard, err error) {
	//获取数据
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, name, id_card FROM service_user_info WHERE org_id = $1 AND delete_at < to_timestamp(1000000) AND (name = $2 OR id_card = $2) LIMIT 10", args.OrgID, args.Search)
	if err != nil || len(dataList) < 1 {
		err = errors.New(fmt.Sprint("data not exist, ", err))
		return
	}
	//反馈
	return
}
