package UserAddress

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetList 获取列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//上级
	// 0则表示不是历史数据；否则为历史数据，且必须指定上级地址ID
	ParentID int64 `json:"parent_id" check:"id" empty:"true"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country" empty:"true"`
	//省份 编码
	// eg: 710000
	Province int `db:"province" json:"province" check:"province" empty:"true"`
	//所属城市
	City int `db:"city" json:"city" check:"city" empty:"true"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//直接查询电话
	// 注意和普通搜索不会同时生效
	SearchPhone string `json:"searchPhone" check:"search" empty:"true"`
	//搜索
	// 昵称、姓名、电话、地址
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 获取列表
func GetList(args *ArgsGetList) (dataList []FieldsAddress, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.ParentID > -1 {
		where = where + "parent_id = :parent_id"
		maps["parent_id"] = args.ParentID
	}
	if args.UserID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.Country > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "country = :country"
		maps["country"] = args.Country
	}
	if args.Province > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "province = :province"
		maps["province"] = args.Province
	}
	if args.City > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "city = :city"
		maps["city"] = args.City
	}
	if args.IsRemove {
		if where != "" {
			where = where + " AND "
		}
		where = where + "delete_at > to_timestamp(1000000)"
	} else {
		if where != "" {
			where = where + " AND "
		}
		where = where + "delete_at < to_timestamp(1000000)"
	}
	if args.SearchPhone != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(phone ILIKE '%' || :search || '%')"
		maps["search"] = args.SearchPhone
	} else {
		if args.Search != "" {
			if where != "" {
				where = where + " AND "
			}
			where = where + "(nice_name ILIKE '%' || :search || '%' OR name ILIKE '%' || :search || '%' OR address ILIKE '%' || :search || '%' OR name ILIKE '%' || :search || '%' OR phone ILIKE '%' || :search || '%' OR email ILIKE '%' || :search || '%')"
			maps["search"] = args.Search
		}
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"user_address",
		"id",
		"SELECT id, create_at, update_at, delete_at, parent_id, user_id, nice_name, country, province, city, address, map_type, longitude, latitude, name, nation_code, phone, email, infos FROM user_address WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// GetAddressByUserID 获取用户的所有地址
func GetAddressByUserID(userID int64, limit int) (dataList []FieldsAddress, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, parent_id, user_id, nice_name, country, province, city, address, map_type, longitude, latitude, name, nation_code, phone, email, infos FROM user_address WHERE user_id = $1 AND delete_at < to_timestamp(1000000) AND parent_id = 0 LIMIT $2", userID, limit)
	if err != nil || len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetCount 查询用户的地址总数参数
type ArgsGetCount struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// GetCount 查询用户的地址总数
// 会自动剔除已删除部分
func GetCount(args *ArgsGetCount) (count int64, err error) {
	count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "user_address", "id", "delete_at < to_timestamp(1000000) AND user_id = $1 AND parent_id = 0", args.UserID)
	if err != nil {
		count = 0
	}
	return
}

// ArgsGetIDTop 追溯到该ID的最上级参数
type ArgsGetIDTop struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//验证用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// GetIDTop 追溯到该ID的最上级
// 例如在默认地址中，可以通过此方法追溯该用户最高级
func GetIDTop(args *ArgsGetIDTop) (data FieldsAddress, err error) {
	//找到该ID，无视是否删除
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, parent_id, user_id, nice_name, country, province, city, address, map_type, longitude, latitude, name, nation_code, phone, email, infos FROM user_address WHERE id = $1 AND ($2 < 1 OR user_id = $2)", args.ID, args.UserID)
	//检查是否存在上级？
	if data.ParentID < 1 {
		//不存在，则直接返回数据
		if data.DeleteAt.Unix() > 0 {
			err = errors.New("no data")
			return
		}
		return
	}
	//如果存在上级，开始找到上级反馈
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, parent_id, user_id, nice_name, country, province, city, address, map_type, longitude, latitude, name, nation_code, phone, email, infos FROM user_address WHERE id = $1", data.ParentID)
	return
}

// ArgsGetID 查看指定ID参数
type ArgsGetID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//验证用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//是否反馈被删除数据？
	IsRemove bool `db:"is_remove" check:"bool"`
}

// GetID 查看指定ID
func GetID(args *ArgsGetID) (data FieldsAddress, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, parent_id, user_id, nice_name, country, province, city, address, map_type, longitude, latitude, name, nation_code, phone, email, infos FROM user_address WHERE id = $1 AND ($2 < 1 OR user_id = $2) AND ($3 = false OR delete_at > to_timestamp(1000000))", args.ID, args.UserID, args.IsRemove)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetIDs 查看一组数据参数
type ArgsGetIDs struct {
	//ID组
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否反馈被删除数据？
	IsRemove bool `db:"is_remove" check:"bool"`
}

// GetIDs 查看一组数据
func GetIDs(args *ArgsGetIDs) (dataList []FieldsAddress, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, parent_id, user_id, nice_name, country, province, city, address, map_type, longitude, latitude, name, nation_code, phone, email, infos FROM user_address WHERE id = ANY($1) AND ($2 = false OR delete_at > to_timestamp(1000000))", args.IDs, args.IsRemove)
	if len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	return
}
