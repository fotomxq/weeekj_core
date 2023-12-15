package ServiceAD

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetADList 获取广告列表参数
type ArgsGetADList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分区标识码
	// 作为前端抽取数据类型使用，可以重复指定多个
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetADList 获取广告列表
func GetADList(args *ArgsGetADList) (dataList []FieldsAD, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Mark != "" {
		where = where + " AND mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"service_ad",
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, mark, name, cover_file_id, des_files, params FROM service_ad WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "name"},
	)
	return
}

// ArgsGetADByID 获取指定的ID参数
type ArgsGetADByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetADByID 获取指定的ID
func GetADByID(args *ArgsGetADByID) (data FieldsAD, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, mark, name, des, cover_file_id, des_files, params FROM service_ad WHERE id = $1 AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	return
}

// ArgsGetADMore 获取一组数据参数
type ArgsGetADMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
}

// GetADMore 获取一组数据
func GetADMore(args *ArgsGetADMore) (dataList []FieldsAD, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "service_ad", "id, create_at, update_at, delete_at, org_id, mark, name, des, cover_file_id, des_files, params", args.IDs, args.HaveRemove)
	return
}

func GetGetADMoreMap(args *ArgsGetADMore) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsNameAndDelete("service_ad", args.IDs, args.HaveRemove)
	return
}

// ArgsCreateAD 创建新的广告参数
type ArgsCreateAD struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分区标识码
	// 作为前端抽取数据类型使用，可以重复指定多个
	Mark string `db:"mark" json:"mark" check:"mark"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述组图
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateAD 创建新的广告
func CreateAD(args *ArgsCreateAD) (data FieldsAD, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "service_ad", "INSERT INTO service_ad(org_id, mark, name, des, cover_file_id, des_files, params) VALUES(:org_id, :mark, :name, :des, :cover_file_id, :des_files, :params)", args, &data)
	return
}

// ArgsUpdateAD 修改广告参数
type ArgsUpdateAD struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 验证用
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分区标识码
	// 作为前端抽取数据类型使用，可以重复指定多个
	Mark string `db:"mark" json:"mark" check:"mark"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述组图
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateAD 修改广告
func UpdateAD(args *ArgsUpdateAD) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_ad SET mark = :mark, name = :name, des = :des, cover_file_id = :cover_file_id, des_files = :des_files, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteAD 删除广告参数
type ArgsDeleteAD struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 验证用
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteAD 删除广告
func DeleteAD(args *ArgsDeleteAD) (err error) {
	//删除广告数据
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_ad", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err == nil {
		//尝试删除绑定关系
		_, _ = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_ad_bind", "ad_id = :ad_id", map[string]interface{}{
			"ad_id": args.ID,
		})
	}
	return
}
