package ServiceCompany

import (
	"errors"
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetBindList 获取绑定列表参数
type ArgsGetBindList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
}

// GetBindList 获取绑定列表
func GetBindList(args *ArgsGetBindList) (dataList []FieldsBind, dataCount int64, err error) {
	//获取数据
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.CompanyID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "company_id = :company_id"
		maps["company_id"] = args.CompanyID
	}
	if where == "" {
		where = "true"
	}
	var rawList []FieldsBind
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"service_company_bind",
		"id",
		"SELECT id FROM service_company_bind WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	//遍历重组数据
	for _, v := range rawList {
		vData := getBind(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// GetBindByOrgID 获取公司绑定关系列表
func GetBindByOrgID(orgID int64) (dataList []FieldsBind, err error) {
	var rawList []FieldsBind
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM service_company_bind WHERE org_id = $1 AND user_id = 0 AND phone = '' LIMIT 999", orgID)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getBind(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// CheckBindAndOrg 检查组织是否具备公司绑定关系
func CheckBindAndOrg(orgID int64, companyID int64) bool {
	var id int64
	err := Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM service_company_bind WHERE org_id = $1 AND company_id = $2 AND user_id = 0 AND phone = '' LIMIT 1", orgID, companyID)
	return err == nil && id > 0
}

// GetBindByUserID 获取用户现在的所有绑定关系
func GetBindByUserID(userID int64) (dataList []FieldsBind, err error) {
	var rawList []FieldsBind
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT b.id as id FROM service_company_bind as b, service_company as c WHERE b.user_id = $1 AND b.company_id = c.id AND c.delete_at < to_timestamp(1000000) LIMIT 999", userID)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getBind(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// CheckBindAndUser 获取用户和公司的绑定关系
func CheckBindAndUser(userID int64, companyID int64) bool {
	var id int64
	err := Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM service_company_bind WHERE user_id = $1 AND company_id = $2 LIMIT 1", userID, companyID)
	return err == nil && id > 0
}

// ArgsSetBind 设置新绑定关系参数
type ArgsSetBind struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//可以预设手机号，手续用户绑定后自动绑定对应用户
	//绑定手机号的国家代码
	NationCode string `db:"nation_code" json:"nationCode" check:"nationCode" empty:"true"`
	//手机号码，绑定后的手机
	Phone string `db:"phone" json:"phone" check:"phone" empty:"true"`
	//公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id"`
	//赋予能力
	// 系统约定的几个特定能力，平台无法编辑该范围，只能授权
	Managers pq.StringArray `db:"managers" json:"managers"`
}

// SetBind 设置新绑定关系
func SetBind(args *ArgsSetBind) (data FieldsBind, err error) {
	//锁定机制
	bindLock.Lock()
	defer bindLock.Unlock()
	//检查参数
	if args.OrgID < 1 && args.UserID < 1 && args.Phone == "" {
		err = errors.New("user and phone is empty")
		return
	}
	//检查限制数量
	if !checkBindLimit(args.OrgID, args.UserID, args.NationCode, args.Phone) {
		err = errors.New("bind limit")
		return
	}
	//单独绑定组织；绑定对应用户；预绑定手机号
	if args.OrgID > 0 && args.UserID == 0 && args.Phone == "" {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM service_company_bind WHERE company_id = $1 AND org_id = $2 AND user_id = 0 AND phone = '' LIMIT 1", args.CompanyID, args.OrgID)
	} else {
		if args.UserID > 0 {
			err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM service_company_bind WHERE company_id = $1 AND org_id = $2 AND user_id = $3 LIMIT 1", args.CompanyID, args.OrgID, args.UserID)
		} else {
			err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM service_company_bind WHERE company_id = $1 AND org_id = $2 AND nation_code = $3 AND phone = $4 LIMIT 1", args.CompanyID, args.OrgID, args.NationCode, args.Phone)
		}
	}
	//修改或创建数据
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_company_bind SET managers = :managers WHERE id = :id", map[string]interface{}{
			"id":       data.ID,
			"managers": args.Managers,
		})
		if err != nil {
			return
		}
	} else {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "service_company_bind", "INSERT INTO service_company_bind(org_id, user_id, nation_code, phone, company_id, managers) VALUES(:org_id, :user_id, :nation_code, :phone, :company_id, :managers)", args, &data)
		if err != nil {
			return
		}
	}
	//删除缓冲
	deleteBindCache(data.ID)
	//获取数据
	data = getBind(data.ID)
	//反馈
	return
}

// ArgsDeleteBind 删除绑定关系参数
type ArgsDeleteBind struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteBind 删除绑定关系
func DeleteBind(args *ArgsDeleteBind) (err error) {
	//删除数据
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "service_company_bind", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//删除缓冲
	deleteBindCache(args.ID)
	//反馈
	return
}

// 删除公司所有供应商关系
func deleteBindByCompanyID(orgID int64, companyID int64) {
	var dataList []FieldsBind
	_ = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM service_company_bind WHERE org_id = $1 AND company_id = $2", orgID, companyID)
	for _, v := range dataList {
		_ = DeleteBind(&ArgsDeleteBind{
			ID:    v.ID,
			OrgID: orgID,
		})
	}
}

// 检查用户绑定公司数量
func checkBindLimit(orgID int64, userID int64, nationCode, phone string) bool {
	//获取限制配置
	serviceCompanyBindLimit, _ := BaseConfig.GetDataInt64("ServiceCompanyBindLimit")
	//获取绑定人数
	var count int64
	var err error
	if orgID > 0 && userID == 0 && phone == "" {
		err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM service_company_bind WHERE org_id = $1 AND user_id = 0 AND phone = ''", orgID)
	} else {
		if userID > 0 {
			err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM service_company_bind WHERE org_id = $1 AND user_id = $2", orgID, userID)
		} else {
			err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM service_company_bind WHERE org_id = $1 AND nation_code = $2 AND phone = $3", orgID, nationCode, phone)
		}
	}
	//不存在数据反馈成功
	if (err != nil || count < 1) && serviceCompanyBindLimit > 0 {
		return true
	}
	//反馈数据
	return count < serviceCompanyBindLimit
}

func getBind(id int64) (data FieldsBind) {
	cacheMark := getBindCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, user_id, nation_code, phone, company_id, managers FROM service_company_bind WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// 缓冲
func getBindCacheMark(id int64) string {
	return fmt.Sprint("service:company:bind:", id)
}

func deleteBindCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getBindCacheMark(id))
}
