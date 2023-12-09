package CoreSQL

import (
	"github.com/jmoiron/sqlx"
)

//GetAllCount 获取总行数
func GetAllCount(db *sqlx.DB, tableName string, fieldsName string, where string, whereArgs ...interface{}) (count int64, err error) {
	if where != "" {
		where = " WHERE " + where
	}
	err = db.QueryRow("SELECT COUNT("+fieldsName+") FROM "+tableName+where, whereArgs...).Scan(&count)
	return
}

//GetAllCountMap 获取总行数，map方式
func GetAllCountMap(db *sqlx.DB, tableName string, fieldsName string, where string, whereArgs interface{}) (count int64, err error) {
	if where != "" {
		where = " WHERE " + where
	}
	var rows *sqlx.Rows
	rows, err = db.NamedQuery("SELECT COUNT("+fieldsName+") FROM "+tableName+where, whereArgs)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&count)
		break
	}
	return
}

//GetAllCountMapTables 混合多表查询方法集合
func GetAllCountMapTables(db *sqlx.DB, fieldsName string, query string, whereArgs interface{}) (count int64, err error) {
	var rows *sqlx.Rows
	rows, err = db.NamedQuery("SELECT COUNT("+fieldsName+") FROM "+query, whereArgs)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&count)
		break
	}
	return
}
