package CoreSQLIDs

import (
	"errors"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

//查询一组ID的通用方法设计

// GetIDsAndDelete 带有delete的方法设计
// 必须带有id / delete_at 字段设计
func GetIDsAndDelete(dataList interface{}, tableName string, fields string, ids pq.Int64Array, haveRemove bool) (err error) {
	//去重复
	newIDs := pq.Int64Array{}
	if len(ids) > 0 {
		for key := 0; key < len(ids); key++ {
			isFind := false
			if len(newIDs) > 0 {
				for key2 := 0; key2 < len(newIDs); key2++ {
					if ids[key] == newIDs[key2] {
						isFind = true
						break
					}
				}
				if isFind {
					continue
				}
			}
			newIDs = append(newIDs, ids[key])
		}
	} else {
		return
	}
	//获取数据
	err = Router2SystemConfig.MainDB.Select(dataList, "SELECT "+fields+" FROM "+tableName+" WHERE id = ANY($1) AND ($2 = true OR delete_at < to_timestamp(1000000)) LIMIT 100", newIDs, haveRemove)
	return
}

// GetIDsOrgAndDelete 带有delete的方法设计
// 必须带有id / org_id / delete_at 字段设计
func GetIDsOrgAndDelete(dataList interface{}, tableName string, fields string, ids pq.Int64Array, orgID int64, haveRemove bool) (err error) {
	//去重复
	newIDs := pq.Int64Array{}
	if len(ids) > 0 {
		for key := 0; key < len(ids); key++ {
			isFind := false
			if len(newIDs) > 0 {
				for key2 := 0; key2 < len(newIDs); key2++ {
					if ids[key] == newIDs[key2] {
						isFind = true
						break
					}
				}
				if isFind {
					continue
				}
			}
			newIDs = append(newIDs, ids[key])
		}
	} else {
		return
	}
	//获取数据
	err = Router2SystemConfig.MainDB.Select(dataList, "SELECT "+fields+" FROM "+tableName+" WHERE id = ANY($1) AND ($2 = true OR delete_at < to_timestamp(1000000)) AND ($3 < 1 OR org_id = $3) LIMIT 100", newIDs, haveRemove, orgID)
	return
}

// GetIDsNameAndDelete 必须带有id / delete_at / name字段设计
func GetIDsNameAndDelete(tableName string, ids pq.Int64Array, haveRemove bool) (data map[int64]string, err error) {
	//去重复
	newIDs := pq.Int64Array{}
	if len(ids) > 0 {
		for key := 0; key < len(ids); key++ {
			isFind := false
			if len(newIDs) > 0 {
				for key2 := 0; key2 < len(newIDs); key2++ {
					if ids[key] == newIDs[key2] {
						isFind = true
						break
					}
				}
				if isFind {
					continue
				}
			}
			newIDs = append(newIDs, ids[key])
		}
	} else {
		return
	}
	//获取数据
	type dataType struct {
		//ID
		ID int64 `db:"id" json:"id"`
		//名称
		Name string `db:"name" json:"name"`
	}
	var dataList []dataType
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, name FROM "+tableName+" WHERE id = ANY($1) AND ($2 = true OR delete_at < to_timestamp(1000000)) LIMIT 100", newIDs, haveRemove)
	if err == nil {
		data = map[int64]string{}
		for _, v := range dataList {
			data[v.ID] = v.Name
		}
	}
	return
}

// GetIDsOrgNameAndDelete 必须带有id / org_id / delete_at / name字段设计
func GetIDsOrgNameAndDelete(tableName string, ids pq.Int64Array, orgID int64, haveRemove bool) (data map[int64]string, err error) {
	//去重复
	newIDs := pq.Int64Array{}
	if len(ids) > 0 {
		for key := 0; key < len(ids); key++ {
			isFind := false
			if len(newIDs) > 0 {
				for key2 := 0; key2 < len(newIDs); key2++ {
					if ids[key] == newIDs[key2] {
						isFind = true
						break
					}
				}
				if isFind {
					continue
				}
			}
			newIDs = append(newIDs, ids[key])
		}
	} else {
		return
	}
	//获取数据
	type dataType struct {
		//ID
		ID int64 `db:"id" json:"id"`
		//名称
		Name string `db:"name" json:"name"`
	}
	var dataList []dataType
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, name FROM "+tableName+" WHERE id = ANY($1) AND ($2 = true OR delete_at < to_timestamp(1000000)) AND ($3 < 1 OR org_id = $3) LIMIT 100", newIDs, haveRemove, orgID)
	if err == nil {
		data = map[int64]string{}
		for _, v := range dataList {
			data[v.ID] = v.Name
		}
	}
	return
}

type DataGetIDsOrgNameAndDeleteAndAvatar struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//昵称
	Name string `db:"name" json:"name"`
	//头像
	Avatar int64 `db:"avatar" json:"avatar"`
}

// GetIDsOrgNameAndDeleteAndAvatar 必须带有id / org_id / delete_at / name字段设计
func GetIDsOrgNameAndDeleteAndAvatar(tableName string, ids pq.Int64Array, orgID int64, haveRemove bool) (data []DataGetIDsOrgNameAndDeleteAndAvatar, err error) {
	//去重复
	newIDs := pq.Int64Array{}
	if len(ids) > 0 {
		for key := 0; key < len(ids); key++ {
			isFind := false
			if len(newIDs) > 0 {
				for key2 := 0; key2 < len(newIDs); key2++ {
					if ids[key] == newIDs[key2] {
						isFind = true
						break
					}
				}
				if isFind {
					continue
				}
			}
			newIDs = append(newIDs, ids[key])
		}
	} else {
		return
	}
	//获取数据
	err = Router2SystemConfig.MainDB.Select(&data, "SELECT id, name, avatar FROM "+tableName+" WHERE id = ANY($1) AND ($2 = true OR delete_at < to_timestamp(1000000)) AND ($3 < 1 OR org_id = $3) LIMIT 100", newIDs, haveRemove, orgID)
	if err != nil || len(data) < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// GetIDsOrgTitleAndDelete 必须带有id / org_id / delete_at / title字段设计
func GetIDsOrgTitleAndDelete(tableName string, ids pq.Int64Array, orgID int64, haveRemove bool) (data map[int64]string, err error) {
	//去重复
	newIDs := pq.Int64Array{}
	if len(ids) > 0 {
		for key := 0; key < len(ids); key++ {
			isFind := false
			if len(newIDs) > 0 {
				for key2 := 0; key2 < len(newIDs); key2++ {
					if ids[key] == newIDs[key2] {
						isFind = true
						break
					}
				}
				if isFind {
					continue
				}
			}
			newIDs = append(newIDs, ids[key])
		}
	} else {
		return
	}
	//获取数据
	type dataType struct {
		//ID
		ID int64 `db:"id" json:"id"`
		//名称
		Title string `db:"title" json:"title"`
	}
	var dataList []dataType
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, title FROM "+tableName+" WHERE id = ANY($1) AND ($2 = true OR delete_at < to_timestamp(1000000)) AND ($3 < 1 OR org_id = $3) LIMIT 100", newIDs, haveRemove, orgID)
	if err == nil {
		data = map[int64]string{}
		for _, v := range dataList {
			data[v.ID] = v.Title
		}
	}
	return
}

// GetIDsTitleAndDelete 必须带有id / delete_at / title字段设计
func GetIDsTitleAndDelete(tableName string, ids pq.Int64Array, haveRemove bool) (data map[int64]string, err error) {
	//去重复
	newIDs := pq.Int64Array{}
	if len(ids) > 0 {
		for key := 0; key < len(ids); key++ {
			isFind := false
			if len(newIDs) > 0 {
				for key2 := 0; key2 < len(newIDs); key2++ {
					if ids[key] == newIDs[key2] {
						isFind = true
						break
					}
				}
				if isFind {
					continue
				}
			}
			newIDs = append(newIDs, ids[key])
		}
	} else {
		return
	}
	//获取数据
	type dataType struct {
		//ID
		ID int64 `db:"id" json:"id"`
		//标题
		Title string `db:"title" json:"title"`
	}
	var dataList []dataType
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, title FROM "+tableName+" WHERE id = ANY($1) AND ($2 = true OR delete_at < to_timestamp(1000000)) LIMIT 100", newIDs, haveRemove)
	if err == nil {
		data = map[int64]string{}
		for _, v := range dataList {
			data[v.ID] = v.Title
		}
	}
	return
}
