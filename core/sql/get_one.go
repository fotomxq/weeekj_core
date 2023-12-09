package CoreSQL

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

//GetOne 通用解析程序
func GetOne(db *sqlx.DB, data interface{}, query string, maps map[string]interface{}) (err error) {
	if maps == nil{
		err = db.Get(data, query)
		return
	}
	var rows *sqlx.Rows
	rows, err = db.NamedQuery(query, maps)
	if err != nil {
		return
	}
	defer rows.Close()
	noData := true
	for rows.Next() {
		err = rows.StructScan(data)
		if err != nil {
			err = nil
			continue
		}
		noData = false
		break
	}
	if noData{
		err = errors.New("no data")
	}
	return
}