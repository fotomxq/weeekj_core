package CoreSQL

import (
	"github.com/jmoiron/sqlx"
)

//GetAllSumMap 获取合计数，map方式
func GetAllSumMap(db *sqlx.DB, tableName string, fieldsName string, where string, whereArgs interface{}) (count int64, err error) {
	if where != "" {
		where = " WHERE " + where
	}
	var rows *sqlx.Rows
	rows, err = db.NamedQuery("SELECT SUM("+fieldsName+") FROM "+tableName+where, whereArgs)
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
