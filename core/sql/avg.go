package CoreSQL

import "github.com/jmoiron/sqlx"

//GetAllAvg 获取平均行数
func GetAllAvg(db *sqlx.DB, tableName string, fieldsName string, where string, whereArgs ...interface{}) (count float64, err error) {
	if where != "" {
		where = " WHERE " + where
	}
	err = db.QueryRow("SELECT AVG("+fieldsName+") FROM "+tableName+where, whereArgs...).Scan(&count)
	return
}

//GetAllAvgMap 获取平均行数，map方式
func GetAllAvgMap(db *sqlx.DB, tableName string, fieldsName string, where string, whereArgs interface{}) (count float64, err error) {
	if where != "" {
		where = " WHERE " + where
	}
	var rows *sqlx.Rows
	rows, err = db.NamedQuery("SELECT AVG("+fieldsName+") FROM "+tableName+where, whereArgs)
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

//GetAllAvgMapTables 混合多表查询方法集合
func GetAllAvgMapTables(db *sqlx.DB, fieldsName string, query string, whereArgs interface{}) (count float64, err error) {
	var rows *sqlx.Rows
	rows, err = db.NamedQuery("SELECT AVG("+fieldsName+") FROM "+query, whereArgs)
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
