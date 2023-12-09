package MarketCore

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetBindList 获取绑定关系参数
type ArgsGetBindList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//绑定的用户
	BindUserID int64 `db:"bind_user_id" json:"bindUserID" check:"id" empty:"true"`
	//绑定的档案
	BindInfoID int64 `db:"bind_info_id" json:"bindInfoID" check:"id" empty:"true"`
	//建立关系的渠道
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetBindList 获取绑定关系
func GetBindList(args *ArgsGetBindList) (dataList []FieldsBind, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
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
	if args.OrgID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.BindID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.BindUserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "bind_user_id = :bind_user_id"
		maps["bind_user_id"] = args.BindUserID
	}
	if args.BindInfoID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "bind_info_id = :bind_info_id"
		maps["bind_info_id"] = args.BindInfoID
	}
	if args.FromInfo.System != "" {
		where, maps, err = args.FromInfo.GetListAnd("from_info", "from_info", where, maps)
		if err != nil {
			return
		}
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"market_core_bind",
		"id",
		"SELECT id, create_at, update_at, delete_at, sort_id, tags, org_id, bind_id, bind_user_id, bind_info_id, from_info, des, params FROM market_core_bind WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetBindGroupList 根据营销人员聚合数据参数
type ArgsGetBindGroupList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//建立关系的渠道
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// DataGetBindGroupList 根据营销人员聚合数据数据
type DataGetBindGroupList struct {
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//关系人数
	Count int64 `db:"count" json:"count"`
}

// GetBindGroupList 根据营销人员聚合数据
func GetBindGroupList(args *ArgsGetBindGroupList) (dataList []DataGetBindGroupList, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
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
	if args.OrgID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.BindID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.FromInfo.System != "" {
		where, maps, err = args.FromInfo.GetListAnd("from_info", "from_info", where, maps)
		if err != nil {
			return
		}
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"market_core_bind",
		"id",
		"SELECT bind_id, COUNT(id) as count FROM market_core_bind WHERE "+where+" GROUP BY bind_id",
		where,
		maps,
		&args.Pages,
		[]string{"bind_id", "count"},
	)
	return
}

// ArgsGetBindByUserID 获取指定用户的营销关系参数
type ArgsGetBindByUserID struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//绑定的用户
	BindUserID int64 `db:"bind_user_id" json:"bindUserID" check:"id" empty:"true"`
	//绑定的档案
	BindInfoID int64 `db:"bind_info_id" json:"bindInfoID" check:"id" empty:"true"`
}

// GetBindByUserID 获取指定用户的营销关系
func GetBindByUserID(args *ArgsGetBindByUserID) (data FieldsBind, err error) {
	//必须存在档案或用户
	if args.BindUserID < 1 && args.BindInfoID < 1 {
		err = errors.New("user id or info id less 1")
		return
	}
	//获取数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, sort_id, tags, org_id, bind_id, bind_user_id, bind_info_id, from_info, des, params FROM market_core_bind WHERE org_id = $1 AND delete_at < to_timestamp(1000000) AND (bind_user_id = $2 OR bind_info_id = $3)", args.OrgID, args.BindUserID, args.BindInfoID)
	return
}

// ArgsCreateBind 建立推荐人关系参数
type ArgsCreateBind struct {
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//绑定的用户
	BindUserID int64 `db:"bind_user_id" json:"bindUserID" check:"id" empty:"true"`
	//绑定的档案
	BindInfoID int64 `db:"bind_info_id" json:"bindInfoID" check:"id" empty:"true"`
	//建立关系的渠道
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//客户备注
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateBind 建立推荐人关系
func CreateBind(args *ArgsCreateBind) (data FieldsBind, err error) {
	//必须存在档案或用户
	if args.BindUserID < 1 && args.BindInfoID < 1 {
		err = errors.New("user id or info id less 1")
		return
	}
	//只能存在一个关系
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, sort_id, tags, org_id, bind_id, bind_user_id, des, params FROM market_core_bind WHERE org_id = $1 AND (($2 > 0 AND bind_user_id = $2) OR ($3 > 0 AND bind_info_id = $3)) AND delete_at < to_timestamp(1000000)", args.OrgID, args.BindUserID, args.BindInfoID)
	if err == nil && data.ID > 0 {
		//如果不符合，则退出
		err = errors.New("bind user have bind")
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "market_core_bind", "INSERT INTO market_core_bind (sort_id, tags, org_id, bind_id, bind_user_id, bind_info_id, from_info, des, params) VALUES (:sort_id,:tags,:org_id,:bind_id,:bind_user_id,:bind_info_id,:from_info,:des,:params)", args, &data)
	return
}

// ArgsUpdateBind 更新推荐人参数
type ArgsUpdateBind struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//客户备注
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateBind 更新推荐人
func UpdateBind(args *ArgsUpdateBind) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE market_core_bind SET update_at = NOW(), sort_id = :sort_id, tags = :tags, des = :des, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:bind_id < 1 OR bind_id = :bind_id)", args)
	return
}

// ArgsUpdateBindToNewBind 批量修改营销人员关系到新营销人员参数
type ArgsUpdateBindToNewBind struct {
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//旧配送人员
	OldBindID int64 `db:"old_bind_id" json:"oldBindID" check:"id"`
	//新配送员
	NewBindID int64 `db:"new_bind_id" json:"newBindID" check:"id"`
}

// UpdateBindToNewBind 批量修改营销人员关系到新营销人员
func UpdateBindToNewBind(args *ArgsUpdateBindToNewBind) (err error) {
	//禁止自己转自己
	if args.OldBindID == args.NewBindID {
		err = errors.New("old and new is same")
		return
	}
	//修改新营销人员
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE market_core_bind SET update_at = NOW(), bind_id = :new_bind_id WHERE bind_id = :old_bind_id AND delete_at < to_timestamp(1000000) AND org_id = :org_id", args)
	if err != nil {
		return
	}
	return
}

// ArgsDeleteBind 删除推荐人参数
type ArgsDeleteBind struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
}

// DeleteBind 删除推荐人参数
func DeleteBind(args *ArgsDeleteBind) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "market_core_bind", "id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:bind_id < 1 OR bind_id = :bind_id)", args)
	return
}
