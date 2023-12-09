package OrgCoreCore

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetBindLast 自动分配人员管理器参数
// 该模块可以由外部任意使用，确保可抽取最早一个分配任务或其他行为的人员，作为最终候选人
// 模块内置记录功能，可记录上次分配行为的参数数据包，用于其他综合考虑
type ArgsGetBindLast struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//包含的分组
	// 可选，不提供则在全组织搜索；提供则仅在包含分组内查询
	GroupID int64 `db:"group_id" json:"groupID" check:"id" empty:"true"`
	//绑定关系的标识码
	// 建议给予任务下的特定群体，例如维修人员上门维修前分配任务，可以给予maintain的标记
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// GetBindLast 自动分配人员管理器
func GetBindLast(args *ArgsGetBindLast) (bindData FieldsBind, err error) {
	//查询组织下符合条件的所有人员
	var bindList []FieldsBind
	err = Router2SystemConfig.MainDB.Select(&bindList, "SELECT id FROM org_core_bind WHERE org_id = $1 AND ($2 < 1 OR $2 = ANY(group_ids)) AND delete_at < to_timestamp(1000000);", args.OrgID, args.GroupID)
	if err != nil || len(bindList) < 1 {
		err = errors.New(fmt.Sprint("no bind list, ", err))
		return
	}
	//组合条件
	var findBind pq.Int64Array
	for _, v := range bindList {
		findBind = append(findBind, v.ID)
	}
	//查询所有记录
	var lastList []FieldsBindLast
	err = Router2SystemConfig.MainDB.Select(&lastList, "SELECT l.id as id, l.bind_id as bind_id FROM org_core_bind_last as l, org_core_bind as b WHERE l.bind_id = ANY($1) AND l.mark = $2 AND b.id = l.bind_id AND b.delete_at < to_timestamp(1000000) ORDER BY l.last_at;", findBind, args.Mark)
	if err != nil || len(lastList) < 1 {
		//没有任何人符合条件，抽取第一个处理
		err = Router2SystemConfig.MainDB.Get(&bindData, "SELECT id, create_at, update_at, delete_at, user_id, org_id, group_ids, manager, params FROM org_core_bind WHERE id = $1 AND delete_at < to_timestamp(1000000);", bindList[0].ID)
		if err != nil {
			err = errors.New(fmt.Sprint("get bind data by bind list 0, ", err))
			return
		}
		err = getBindLastUpdate(0, &bindList[0], args.Mark, args.Params, true)
		return
	}
	//遍历数据，检查是否存在不存在的数据，如果存在将直接反馈
	for _, v := range bindList {
		isFind := false
		for _, v2 := range lastList {
			if v.ID == v2.BindID {
				isFind = true
				break
			}
		}
		if !isFind {
			err = Router2SystemConfig.MainDB.Get(&bindData, "SELECT id FROM org_core_bind WHERE id = $1 AND delete_at < to_timestamp(1000000);", v.ID)
			if err != nil {
				err = errors.New(fmt.Sprint("get bind data by bind list v, v bind id: ", v.ID, ", err: ", err))
				return
			}
			bindData = getBindByID(bindData.ID)
			err = getBindLastUpdate(0, &v, args.Mark, args.Params, true)
			return
		}
	}
	//最后将第一个数据记录，作为候选人反馈
	err = Router2SystemConfig.MainDB.Get(&bindData, "SELECT id FROM org_core_bind WHERE id = $1 AND delete_at < to_timestamp(1000000);", lastList[0].BindID)
	if err != nil {
		err = errors.New(fmt.Sprint("get last list 0, last id: ", lastList[0].ID, ", err: ", err))
		return
	}
	bindData = getBindByID(bindData.ID)
	err = getBindLastUpdate(lastList[0].ID, &bindData, args.Mark, args.Params, false)
	return
}

// 处理抽取的结果
func getBindLastUpdate(lastID int64, bindData *FieldsBind, mark string, params CoreSQLConfig.FieldsConfigsType, isCreate bool) (err error) {
	if isCreate {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_core_bind_last (mark, bind_id, params) VALUES (:mark, :bind_id, :params)", map[string]interface{}{
			"mark":    mark,
			"bind_id": bindData.ID,
			"params":  params,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("insert bind last, ", err))
		}
	} else {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_core_bind_last SET last_at = NOW(), params = :params WHERE id = :id", map[string]interface{}{
			"id":     lastID,
			"params": params,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("update bind last, id: ", lastID, ", err: ", err))
		}
	}
	return
}

// 删除绑定人的所有记录
// 该方法用于内部表优化处理
func deleteAllLastByBind(bindID int64) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "org_core_bind_last", "bind_id = :bind_id", map[string]interface{}{
		"bind_id": bindID,
	})
	return
}
