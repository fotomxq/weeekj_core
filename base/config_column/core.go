package BaseConfigColumn

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

//前端列头存储专用模块
/**
该模块可以匹配到系统、组织、用户、成员四个级别，自动识别目标终端类型，并呈现不同的自定义列头信息
前端可根据该信息，呈现不同页面的列头、布局等内容。
*/

// 获取列表
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//系统类型
	System int `json:"system" check:"intThan0" empty:"true"`
	//来源ID
	BindID int64 `json:"bindID" check:"id" empty:"true"`
	//标识码
	Mark string `json:"mark" check:"mark" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

func GetList(args *ArgsGetList) (dataList []FieldsColumn, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.System > -1 {
		where = where + "system = :system"
		maps["system"] = args.System
	}
	if args.BindID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.Mark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(mark ILIKE '%' || :search || '%' OR data -> 'name' ? :search)"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_config_column",
		"id",
		"SELECT id, create_at, update_at, system, bind_id, mark, data FROM core_config_column WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "send_at"},
	)
	return
}

// 获取数据
type ArgsGetMark struct {
	//获取的Mark
	Mark string `json:"mark" check:"mark"`
	//获取的所属组织ID
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//获取的用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
}

func GetMark(args *ArgsGetMark) (data FieldsColumn, err error) {
	if args.UserID > 0 {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, system, bind_id, mark, data FROM core_config_column WHERE system = 2 AND mark = $1 AND bind_id = $32", args.Mark, args.UserID)
		if err == nil {
			return
		}
	}
	if args.OrgID > 0 {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, system, bind_id, mark, data FROM core_config_column WHERE system = 1 AND mark = $1 AND bind_id = $2", args.Mark, args.OrgID)
		if err == nil {
			return
		}
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, system, bind_id, mark, data FROM core_config_column WHERE system = 0 AND mark = $1 AND bind_id = 0", args.Mark)
	if err == nil {
		return
	}
	return
}

// 设数据
type ArgsSet struct {
	//来源系统
	// 0 系统层 / 1 组织层 / 2 用户层
	// 系统层影响所有系统配置设计，该设计全系统通用，但用户层可自定义覆盖设定
	// 组织层用于声明组织内部的所有列头，用于覆盖系统层的设计
	// 用户层可直接覆盖系统层或组织层的设定
	System int `db:"system" json:"system"`
	//获取的Mark
	Mark string `db:"mark" json:"mark" check:"mark"`
	//来源ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//保存数据集
	// 前后顺序将按照该顺序一致
	Data FieldsChildList `db:"data" json:"data"`
}

func Set(args *ArgsSet) (err error) {
	var data FieldsColumn
	if args.System > 0 {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM core_config_column WHERE system = 0 AND mark = $1", args.Mark)
		if err != nil {
			return
		}
		if data.ID < 1 {
			err = errors.New("mark not exist")
			return
		}
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM core_config_column WHERE system = $1 AND mark = $2 AND bind_id = $3", args.System, args.Mark, args.BindID)
	if err == nil {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_config_column SET update_at = NOW(), data = :data WHERE id = :id", map[string]interface{}{
			"id":   data.ID,
			"data": args.Data,
		})
		return
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO core_config_column (system, bind_id, mark, data) VALUES (:system,:bind_id,:mark,:data)", args)
	return
}

// 恢复用户数据到组织层
type ArgsReturnUser struct {
	//获取的Mark
	Mark string `json:"mark" check:"mark"`
	//用户ID
	UserID int64 `json:"userID" check:"id"`
}

func ReturnUser(args *ArgsReturnUser) (err error) {
	//检查用户是否存在数据
	var data FieldsColumn
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM core_config_column WHERE system = 2 AND mark = $1 AND bind_id = $2", args.Mark, args.UserID)
	if err != nil {
		err = nil
		return
	}
	//删除数据
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_config_column", "id", map[string]interface{}{
		"id": data.ID,
	})
	return
}

// 恢复组织层数据到系统层
type ArgsReturnOrg struct {
	//获取的Mark
	Mark string `json:"mark" check:"mark"`
	//组织ID
	OrgID int64 `json:"orgID" check:"id"`
}

func ReturnOrg(args *ArgsReturnOrg) (err error) {
	//检查组织是否存在数据
	var data FieldsColumn
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM core_config_column WHERE system = 1 AND mark = $1 AND bind_id = $2", args.Mark, args.OrgID)
	if err != nil {
		err = nil
		return
	}
	//删除数据
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_config_column", "id", map[string]interface{}{
		"id": data.ID,
	})
	return
}

// 删除指定数据
// 将删除所有关联的数据
type ArgsDeleteMark struct {
	//获取的Mark
	Mark string `db:"mark" json:"mark" check:"mark"`
}

func DeleteMark(args *ArgsDeleteMark) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_config_column", "mark = :mark", args)
	return
}
