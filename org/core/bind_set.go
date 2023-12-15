package OrgCoreCore

import (
	"errors"
	"fmt"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsSetBind 设置绑定关系参数
type ArgsSetBind struct {
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//用户头像
	Avatar int64 `db:"avatar" json:"avatar" check:"id" empty:"true"`
	//名称
	Name string `json:"name" check:"name"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//组织分组ID列
	GroupIDs pq.Int64Array `db:"group_ids" json:"groupIDs" check:"ids" empty:"true"`
	//角色配置列
	RoleConfigIDs pq.Int64Array `db:"role_config_ids" json:"roleConfigIDs" check:"ids" empty:"true"`
	//权利主张
	Manager pq.StringArray `db:"manager" json:"manager" check:"marks" empty:"true"`
	//联系电话
	NationCode string `db:"nation_code" json:"nationCode" check:"nationCode" empty:"true"`
	Phone      string `db:"phone" json:"phone" check:"phone" empty:"true"`
	//邮件地址
	Email string `db:"email" json:"email" check:"email" empty:"true"`
	//同步专用设计
	// 可用于同步其他系统来源
	SyncSystem string `db:"sync_system" json:"syncSystem"`
	SyncID     int64  `db:"sync_id" json:"syncID"`
	SyncHash   string `db:"sync_hash" json:"syncHash"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetBind 设置绑定关系
// 注意，也可以修改删除掉的数据，但不会恢复状态
func SetBind(args *ArgsSetBind) (data FieldsBind, err error) {
	//修正参数
	if len(args.GroupIDs) < 1 {
		args.GroupIDs = []int64{}
	}
	if len(args.RoleConfigIDs) < 1 {
		args.RoleConfigIDs = []int64{}
	}
	if len(args.Manager) < 1 {
		args.Manager = []string{}
	}
	//查询存在的绑定关系
	if args.UserID > 0 {
		//绑定用户权限
		if err = checkAndUpdateUserGroup(args.UserID); err != nil {
			return
		}
		//尝试获取用户关联的数据
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_core_bind WHERE org_id = $1 AND user_id = $2", args.OrgID, args.UserID)
	}
	if args.SyncSystem != "" {
		findBind := GetBindBySync(args.OrgID, args.SyncSystem, args.SyncID, args.SyncHash)
		if findBind.ID > 0 && findBind.ID != data.ID && data.ID > 0 {
			err = errors.New("have sync data")
			return
		}
		data = findBind
	}
	if args.NationCode != "" && args.Phone != "" {
		findBind := GetBindByPhone(args.OrgID, args.NationCode, args.Phone)
		if findBind.ID > 0 && findBind.ID != data.ID && data.ID > 0 {
			err = errors.New("have phone")
			return
		}
		data = findBind
	}
	if args.Email != "" {
		findBind := GetBindByEmail(args.OrgID, args.Email)
		if findBind.ID > 0 && findBind.ID != data.ID && data.ID > 0 {
			err = errors.New("have email")
			return
		}
		data = findBind
	}
	//绑定关系
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_core_bind SET update_at = NOW(), delete_at = to_timestamp(0), avatar = :avatar, name = :name, group_ids = :group_ids, role_config_ids = :role_config_ids, manager = :manager, nation_code = :nation_code, phone = :phone, email = :email, sync_system = :sync_system, sync_id = :sync_id, sync_hash = :sync_hash, params = :params WHERE id = :id", map[string]interface{}{
			"id":              data.ID,
			"avatar":          args.Avatar,
			"name":            args.Name,
			"group_ids":       args.GroupIDs,
			"role_config_ids": args.RoleConfigIDs,
			"manager":         args.Manager,
			"nation_code":     args.NationCode,
			"phone":           args.Phone,
			"email":           args.Email,
			"sync_system":     args.SyncSystem,
			"sync_id":         args.SyncID,
			"sync_hash":       args.SyncHash,
			"params":          args.Params,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("update org bind, ", err))
			return
		}
		data = getBindByID(data.ID)
		if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
			err = errors.New("no data")
			return
		}
		//清理缓冲
		deleteBindCache(data.ID)
		//推送NATS
		CoreNats.PushDataNoErr("/org/core/bind", "update", data.ID, "", nil)
	} else {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_core_bind", "INSERT INTO org_core_bind (user_id, avatar, name, org_id, group_ids, role_config_ids, manager, nation_code, phone, email, sync_system, sync_id, sync_hash, params) VALUES (:user_id, :avatar, :name, :org_id, :group_ids, :role_config_ids, :manager, :nation_code, :phone, :email, :sync_system, :sync_id, :sync_hash, :params)", map[string]interface{}{
			"user_id":         args.UserID,
			"avatar":          args.Avatar,
			"name":            args.Name,
			"org_id":          args.OrgID,
			"group_ids":       args.GroupIDs,
			"role_config_ids": args.RoleConfigIDs,
			"manager":         args.Manager,
			"nation_code":     args.NationCode,
			"phone":           args.Phone,
			"email":           args.Email,
			"sync_system":     args.SyncSystem,
			"sync_id":         args.SyncID,
			"sync_hash":       args.SyncHash,
			"params":          args.Params,
		}, &data)
		if err != nil {
			err = errors.New(fmt.Sprint("insert org bind, ", err))
			return
		}
		//推送NATS
		CoreNats.PushDataNoErr("/org/core/bind", "create", data.ID, "", nil)
	}
	//更新扩展数据表内容
	if args.NationCode != "" && args.Phone != "" {
		bindInfoData := getBindInfoByBindID(data.ID)
		if bindInfoData.ID > 0 {
			_ = updateBindInfoBaseInBind(data.ID, data.Phone)
		}
	}
	//更新统计
	updateOrgBindAnalysis(data.OrgID)
	//反馈
	return
}

// ArgsSetBindParams 设置绑定关系的某组参数
type ArgsSetBindParams struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//用户ID
	UserID int64 `json:"userID" check:"id"`
	//附加绑定参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetBindParams 设置绑定关系的某组
// 增量处理，如果希望删除请使用Set方法
func SetBindParams(args *ArgsSetBindParams) (err error) {
	var data FieldsBind
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, params FROM org_core_bind WHERE org_id = $1 AND user_id = $2", args.OrgID, args.UserID)
	if err != nil {
		err = errors.New(fmt.Sprint("get no data, id: ", args.OrgID, ", err: ", err))
		return
	}
	for _, v := range data.Params {
		isFind := false
		for _, v2 := range args.Params {
			if v.Mark == v2.Mark {
				isFind = true
				break
			}
		}
		if !isFind {
			args.Params = append(args.Params, v)
		}
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_core_bind SET update_at = NOW(), params = :params WHERE id = :id", map[string]interface{}{
		"id":     data.ID,
		"params": args.Params,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("update data, id: ", data.ID, ", err: ", err))
		return
	}
	//推送NATS
	CoreNats.PushDataNoErr("/org/core/bind", "update", data.ID, "", nil)
	//清理缓冲
	deleteBindCache(data.ID)
	//反馈
	return
}

// 清理所有指定的角色配置ID
func deleteAllBindRoleConfigID(roleConfigID int64) {
	var bindList []FieldsBind
	err := Router2SystemConfig.MainDB.Select(&bindList, "SELECT id FROM org_core_bind WHERE $1 = ANY(role_config_ids)", roleConfigID)
	if err != nil || len(bindList) < 1 {
		return
	}
	for _, v := range bindList {
		vBind := getBindByID(v.ID)
		if vBind.ID < 1 {
			continue
		}
		var newRoleConfigIDs pq.Int64Array
		for _, v2 := range vBind.RoleConfigIDs {
			if v2 == roleConfigID {
				continue
			}
			newRoleConfigIDs = append(newRoleConfigIDs, v2)
		}
		vBind.RoleConfigIDs = newRoleConfigIDs
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_core_bind SET update_at = NOW(), role_config_ids = :role_config_ids WHERE id = :id", map[string]interface{}{
			"id":              vBind.ID,
			"role_config_ids": vBind.RoleConfigIDs,
		})
		if err != nil {
			CoreLog.Warn("org core delete all bind role config id, ", err)
			continue
		}
		//清理缓冲
		deleteBindCache(vBind.ID)
		//推送NATS
		CoreNats.PushDataNoErr("/org/core/bind", "update", vBind.ID, "", nil)
	}
}

// 局部修改成员的手机号
func updateBindPhone(orgBindID int64, nationCode string, phone string) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_core_bind SET nation_code = :nation_code, phone = :phone WHERE id = :id", map[string]interface{}{
		"id":          orgBindID,
		"nation_code": nationCode,
		"phone":       phone,
	})
	if err != nil {
		return
	}
	deleteBindCache(orgBindID)
	return
}
