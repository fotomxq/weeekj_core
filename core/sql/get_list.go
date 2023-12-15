package CoreSQL

import (
	"fmt"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	"github.com/jmoiron/sqlx"
)

// GetList 遍历和解析多结构体
func GetList(db *sqlx.DB, data interface{}, query string, maps map[string]interface{}) (err error) {
	if maps == nil {
		err = db.Select(data, query)
		return
	}
	var rows *sqlx.Rows
	rows, err = db.NamedQuery(query, maps)
	if err != nil {
		return
	}
	defer rows.Close()
	err = sqlx.StructScan(rows, data)
	return
}

// GetListAndCount 带count的处理机制
func GetListAndCount(db *sqlx.DB, data interface{}, tableName string, fieldName string, query string, where string, maps map[string]interface{}) (count int64, err error) {
	err = GetList(db, data, query, maps)
	if err != nil {
		return
	}
	count, err = GetAllCountMap(db, tableName, fieldName, where, maps)
	return
}

// GetListPageAndCount 带有分页方法的设计GetListPageAndCount
func GetListPageAndCount(db *sqlx.DB, data interface{}, tableName string, fieldName string, query string, where string, maps map[string]interface{}, pages *CoreSQLPages.ArgsDataList, filterSort []string) (count int64, err error) {
	var newMaps map[string]interface{}
	query, newMaps, err = CoreSQLPages.GetMapsAndFilter(pages, query, maps, filterSort)
	if err != nil {
		return
	}
	err = GetList(db, data, query, newMaps)
	if err != nil {
		return
	}
	if where == "" {
		where = "true"
	}
	count, err = GetAllCountMap(db, tableName, fieldName, where, maps)
	return
}

func GetListPageAndCountArgs(db *sqlx.DB, data interface{}, tableName string, fieldName string, query string, where string, pages *CoreSQLPages.ArgsDataList, filterSort []string, args ...interface{}) (count int64, err error) {
	count, err = GetAllCount(db, tableName, fieldName, where, args...)
	if err != nil {
		return
	}
	query, err = CoreSQLPages.GetMapsAndFilterArgs(pages, query, filterSort)
	if err != nil {
		return
	}
	err = db.Select(data, query, args...)
	if err != nil {
		return
	}
	return
}

// GetWhereInt64 检查某个数字参数
func GetWhereInt64(where string, maps map[string]interface{}, fieldName string, param int64) (string, map[string]interface{}) {
	if param > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + fmt.Sprint(fieldName, " = :", fieldName)
		maps[fieldName] = param
	}
	return where, maps
}

// GetNeedChange 组合某个开关类选项
func GetNeedChange(where string, fieldName string, need bool, open bool) string {
	if need {
		if where != "" {
			where = where + " AND "
		}
		if open {
			where = where + fieldName + " > to_timestamp(1000000)"
		} else {
			where = where + fieldName + " <= to_timestamp(1000000)"
		}
	}
	return where
}

func GetNeedChangeNow(where string, fieldName string, need bool, open bool) string {
	if need {
		if where != "" {
			where = where + " AND "
		}
		if open {
			where = where + fieldName + " > NOW()"
		} else {
			where = where + fieldName + " <= NOW()"
		}
	}
	return where
}
