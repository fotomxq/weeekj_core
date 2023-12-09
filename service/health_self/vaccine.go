package ServiceHealthSelf

import (
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetVaccineList 获取疫苗接种记录参数
type ArgsGetVaccineList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// 必填
	OrgID int64 `json:"orgID" check:"id"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetVaccineList 获取疫苗接种记录
func GetVaccineList(args *ArgsGetVaccineList) (dataList []FieldsVaccine, dataCount int64, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.OrgBindID > -1 {
		where = where + " AND :org_bind_id = org_bind_id"
		maps["org_bind_id"] = args.OrgBindID
	}
	if args.UserID > -1 {
		where = where + " AND :user_id = user_id"
		maps["user_id"] = args.UserID
	}
	if args.Search != "" {
		where = where + " AND (address ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	var rawList []FieldsVaccine
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"service_health_self_vaccine",
		"id",
		"SELECT id FROM service_health_self_vaccine WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getVaccineByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsAppendVaccine 添加新的记录参数
type ArgsAppendVaccine struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//疫苗名称
	Name string `db:"name" json:"name" check:"name"`
	//接种地点
	Address string `db:"address" json:"address" check:"address" empty:"true"`
}

// AppendVaccine 添加新的记录
func AppendVaccine(args *ArgsAppendVaccine) (err error) {
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_health_self_vaccine (org_id, org_bind_id, user_id, name, address) VALUES (:org_id,:org_bind_id,:user_id,:name,:address)", args)
	if err != nil {
		return
	}
	return
}

// ArgsDeleteVaccine 删除记录参数
type ArgsDeleteVaccine struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteVaccine 删除记录
func DeleteVaccine(args *ArgsDeleteVaccine) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "service_health_self_vaccine", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteVaccineCache(args.ID)
	return
}

func getVaccineByID(id int64) (data FieldsVaccine) {
	cacheMark := getVaccineCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, org_bind_id, user_id, name, address FROM service_health_self_vaccine WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 10800)
	return
}

// 缓冲
func getVaccineCacheMark(id int64) string {
	return fmt.Sprint("service:health:self:vaccine:id:", id)
}

func deleteVaccineCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getVaccineCacheMark(id))
}
