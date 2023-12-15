package BaseEmail

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetServerList 获取发送列表参数
type ArgsGetServerList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetServerList 获取发送列表
func GetServerList(args *ArgsGetServerList) (dataList []FieldsEmailServerType, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(title ILIKE '%' || :search || '%' OR title_des ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "core_email_server"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, org_id, name, host, port, is_ssl, email, password, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at"},
	)
	return
}

// ArgsGetServerByID 获取发送方参数
type ArgsGetServerByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetServerByID 获取发送方
func GetServerByID(args *ArgsGetServerByID) (data FieldsEmailServerType, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, org_id, name, host, port, is_ssl, email, password, params FROM core_email_server WHERE id = $1 AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	return
}

// ArgsCreateServer 创建新的发送方参数
type ArgsCreateServer struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//Host
	Host string `db:"host" json:"host"`
	//端口
	// 可以留空，将自动指定默认
	Port string `db:"port" json:"port"`
	//是否为SSL方式连接
	IsSSL bool `db:"is_ssl" json:"isSSL" check:"bool"`
	//邮件地址
	Email string `db:"email" json:"email" check:"email"`
	//密码
	Password string `db:"password" json:"password"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateServer 创建新的发送方
// 如果启动ssl且没有给予端口，则自动采用默认465
// 如果没有启动ssl，且没有端口号，则默认采用25
func CreateServer(args *ArgsCreateServer) (data FieldsEmailServerType, err error) {
	if args.Port == "" {
		if args.IsSSL {
			args.Port = "465"
		} else {
			args.Port = "25"
		}
	}
	var lastID int64
	lastID, err = CoreSQL.CreateOneAndID(
		Router2SystemConfig.MainDB.DB,
		"INSERT INTO core_email_server(org_id, name, host, port, is_ssl, email, password, params) VALUES (:org_id, :name, :host, :port, :is_ssl, :email, :password, :params)",
		args,
	)
	if err == nil {
		data, err = GetServerByID(&ArgsGetServerByID{
			ID:    lastID,
			OrgID: args.OrgID,
		})
	}
	return
}

// ArgsUpdateServer 更新发送方信息参数
type ArgsUpdateServer struct {
	//ID
	ID int64 `db:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//Host
	Host string `db:"host" json:"host"`
	//端口
	// 可以留空，将自动指定默认
	Port string `db:"port" json:"port"`
	//是否为SSL方式连接
	IsSSL bool `db:"is_ssl" json:"isSSL" check:"bool"`
	//邮件地址
	Email string `db:"email" json:"email" check:"email"`
	//密码
	Password string `db:"password" json:"password"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateServer 更新发送方信息
func UpdateServer(args *ArgsUpdateServer) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_email_server SET update_at=NOW(), host=:host, port=:port, email=:email, password=:password, is_ssl=:is_ssl, name=:name, params = :params WHERE id=:id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteServerByID 删除发送方参数
type ArgsDeleteServerByID struct {
	//ID
	ID int64 `db:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteServerByID 删除发送方
func DeleteServerByID(args *ArgsDeleteServerByID) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_email_server", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
