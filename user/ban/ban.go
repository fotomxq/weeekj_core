package UserBan

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetBanList 获取黑名单列表参数
type ArgsGetBanList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//数据来源
	System string `db:"system" json:"system" check:"mark"`
}

// GetBanList 获取黑名单列表
func GetBanList(args *ArgsGetBanList) (dataList []FieldsBan, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.UserID > -1 {
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.System != "" {
		if err = checkSystem(args.System); err != nil {
			return
		}
		if where != "" {
			where = where + " AND "
		}
		where = where + "system = :system"
		maps["system"] = args.System
	}
	if where == "" {
		where = "true"
	}
	tableName := "user_ban"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT create_at, org_id, user_id, system, bind_id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if err != nil {
		return
	}
	return
}

// DataGetBanByFrom 查询一组来源是否在黑名单数据
type DataGetBanByFrom struct {
	//创建时间
	CreateAt string `db:"create_at" json:"createAt"`
	//来源ID
	BindID int64 `db:"bind_id" json:"bindID"`
}

// GetBanByFrom 查询一组来源是否在黑名单
func GetBanByFrom(userID int64, system string, ids pq.Int64Array) (dataList []DataGetBanByFrom) {
	for _, v := range ids {
		vData := getBanByFrom(userID, system, v)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, DataGetBanByFrom{
			CreateAt: CoreFilter.GetTimeToDefaultTime(vData.CreateAt),
			BindID:   vData.BindID,
		})
	}
	return
}

// CheckBan 检查一组ID是否为黑名单
func CheckBan(userID int64, system string, ids []int64) []int64 {
	var result []int64
	for _, v := range ids {
		vData := getBanByFrom(userID, system, v)
		if vData.ID < 1 {
			continue
		}
		result = append(result, vData.BindID)
	}
	return result
}

// ArgsSetBan 设置黑名单参数
type ArgsSetBan struct {
	//绑定组织
	// 根据数据来源决定，只是用于统计和记录，组织没有具体记录的访问权限
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//数据来源
	System string `db:"system" json:"system" check:"mark"`
	//来源ID
	BindID int64 `db:"bind_id" json:"bindID"`
}

// SetBan 设置黑名单
func SetBan(args *ArgsSetBan) (err error) {
	if err = checkSystem(args.System); err != nil {
		return
	}
	data := getBanByFrom(args.UserID, args.System, args.BindID)
	if data.ID > 0 {
		return
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_ban (org_id, user_id, system, bind_id) VALUES (:org_id,:user_id,:system,:bind_id)", args)
		if err != nil {
			return
		}
	}
	return
}

// ArgsDeleteBan 删除黑名单参数
type ArgsDeleteBan struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//数据来源
	System string `db:"system" json:"system" check:"mark"`
	//来源ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
}

// DeleteBan 删除黑名单
func DeleteBan(args *ArgsDeleteBan) (err error) {
	if err = checkSystem(args.System); err != nil {
		return
	}
	data := getBanByFrom(args.UserID, args.System, args.BindID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "user_ban", "id", map[string]interface{}{
		"id": data.ID,
	})
	if err != nil {
		return
	}
	deleteBanCache(args.UserID, args.System, args.BindID)
	return
}

func getBanByFrom(userID int64, system string, bindID int64) (data FieldsBan) {
	cacheMark := getBanCacheMark(userID, system, bindID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, user_id, system, bind_id FROM user_ban WHERE user_id = $1 AND system = $2 AND bind_id = $3", userID, system, bindID)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 7200)
	return
}

// 缓冲
func getBanCacheMark(userID int64, system string, bindID int64) string {
	return fmt.Sprint("user:ban:user:", userID, ".", system, ".", bindID)
}

func deleteBanCache(userID int64, system string, bindID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getBanCacheMark(userID, system, bindID))
}

// 检查system
func checkSystem(system string) (err error) {
	switch system {
	case "blog_core":
	case "user_core":
	default:
		err = errors.New("no support system")
	}
	return
}
