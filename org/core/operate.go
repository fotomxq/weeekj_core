package OrgCoreCore

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetOperateList 获取控制关系列参数
type ArgsGetOperateList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//绑定来源
	BindInfo CoreSQLFrom.FieldsFrom `json:"bindInfo"`
	//资源来源
	From CoreSQLFrom.FieldsFrom `json:"from"`
	//符合权限的
	Manager string `json:"manager"`
	//是否删除
	IsRemove bool `json:"isRemove"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetOperateList 获取控制关系列
func GetOperateList(args *ArgsGetOperateList) (dataList []FieldsOperate, dataCount int64, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	where, maps, err = args.From.GetListAnd("from_info", "from_info", where, maps)
	if err != nil {
		return
	}
	where, maps, err = args.BindInfo.GetListAnd("bind_info", "bind_info", where, maps)
	if err != nil {
		return
	}
	if args.Manager != "" {
		where = where + " AND :manager = ANY(manager)"
		maps["manager"] = args.Manager
	}
	if args.Search != "" {
		where = where + " AND (bind_info -> 'name' ? :search OR from_info -> 'name' ? :search)"
		maps["search"] = args.Search
	}
	var rawList []FieldsOperate
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"org_core_operate",
		"id",
		"SELECT id FROM org_core_operate WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getOperateByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetOperate 获取指定的控制关系参数
type ArgsGetOperate struct {
	//组织ID
	OrgID int64
	//绑定来源
	BindInfo CoreSQLFrom.FieldsFrom
	//资源来源
	FromInfo CoreSQLFrom.FieldsFrom
}

// GetOperate 获取指定的控制关系
func GetOperate(args *ArgsGetOperate) (data FieldsOperate, err error) {
	var bindInfo string
	bindInfo, err = args.BindInfo.GetRaw()
	if err != nil {
		return
	}
	var fromInfo string
	fromInfo, err = args.FromInfo.GetRaw()
	if err != nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_core_operate WHERE org_id = $1 AND bind_info @> $2 AND from_info @> $3 AND delete_at < to_timestamp(1000000) LIMIT 1;", args.OrgID, bindInfo, fromInfo)
	if err != nil {
		return
	}
	data = getOperateByID(data.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCheckOperate 检查控制关系以及权限参数
type ArgsCheckOperate struct {
	//组织ID
	OrgID int64
	//绑定来源
	BindInfo CoreSQLFrom.FieldsFrom
	//资源来源
	FromInfo CoreSQLFrom.FieldsFrom
	//检查的权限
	// 一组权限
	Managers pq.StringArray
	//一组权限的检查方式
	// or 或的关系，只要一个满足即可
	// and 全部满足才能放行
	Filter string
}

// CheckOperate  检查控制关系以及权限
func CheckOperate(args *ArgsCheckOperate) (data FieldsOperate, b bool, err error) {
	data, err = GetOperate(&ArgsGetOperate{
		OrgID:    args.OrgID,
		BindInfo: args.BindInfo,
		FromInfo: args.FromInfo,
	})
	if err != nil {
		return
	}
	for _, v := range data.Manager {
		if v == "all" {
			b = true
			return
		} else {
			if args.Filter == "or" {
				for _, v2 := range args.Managers {
					if v == v2 {
						b = true
						return
					}
				}
			} else {
				if args.Filter == "and" {
					isFind := false
					for _, v2 := range args.Managers {
						if v == v2 {
							isFind = true
							break
						}
					}
					if !isFind {
						return
					}
				}
			}
		}
	}
	if args.Filter == "and" {
		b = true
		return
	}
	return
}

// CheckOperateOnlyBool  检查控制关系以及权限
// 只反馈真假
func CheckOperateOnlyBool(args *ArgsCheckOperate) (b bool) {
	_, b, _ = CheckOperate(args)
	return
}

// ArgsSetOperate 设置控制关系参数
type ArgsSetOperate struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//绑定来源
	// system: bind | id: 成员ID
	// system: group | id: 分组ID
	BindInfo CoreSQLFrom.FieldsFrom `db:"bind_info" json:"bindInfo"`
	//控制资源的来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//控制权限列
	Manager pq.StringArray `db:"manager" json:"manager"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetOperate 设置控制关系
func SetOperate(args *ArgsSetOperate) (err error) {
	var data FieldsOperate
	data, err = GetOperate(&ArgsGetOperate{
		OrgID:    args.OrgID,
		BindInfo: args.BindInfo,
		FromInfo: args.FromInfo,
	})
	if err == nil {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_core_operate SET update_at = NOW(), manager = :manager, params = :params WHERE id = :id", map[string]interface{}{
			"id":      data.ID,
			"manager": args.Manager,
			"params":  args.Params,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("update no data, org id: ", args.OrgID, ", err: ", err))
			return
		}
		deleteOperateCache(data.ID)
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_core_operate (org_id, bind_info, from_info, manager, params) VALUES (:org_id, :bind_info, :from_info, :manager, :params)", args)
		if err != nil {
			err = errors.New(fmt.Sprint("create failed, org id: ", args.OrgID, ", err: ", err))
			return
		}
	}
	return
}

// ArgsDeleteOperate 删除控制权限参数
type ArgsDeleteOperate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteOperate 删除控制权限
func DeleteOperate(args *ArgsDeleteOperate) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_core_operate", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteOperateCache(args.ID)
	return
}

// ArgsReturnOperate 恢复绑定关系参数
type ArgsReturnOperate struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
}

// ReturnOperate 恢复绑定关系
func ReturnOperate(args *ArgsReturnOperate) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_core_operate SET delete_at = to_timestamp(0) WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteOperateCache(args.ID)
	return
}

func getOperateByID(id int64) (data FieldsOperate) {
	cacheMark := getOperateCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, bind_info, from_info, manager, params FROM org_core_operate WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, operateCacheTime)
	return
}

// 缓冲
func getOperateCacheMark(id int64) string {
	return fmt.Sprint("org:core:operate:id:", id)
}

func deleteOperateCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getOperateCacheMark(id))
}
