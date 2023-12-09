package CoreSQL

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

//UpdateOne 快速更新
func UpdateOne(db *sqlx.DB, query string, args interface{}) (result sql.Result, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint("recover, ", e))
			return
		}
	}()
	tx := db.MustBegin()
	if args == nil {
		result, err = tx.Exec(query)
	} else {
		result, err = tx.NamedExec(query, args)
	}
	err = LastRowsAffected(tx, result, err)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

//UpdateOneSoft 支持软删除的更新
// 请确保query where部分为封闭设计，避免影响到追加的软删除指令标记
func UpdateOneSoft(db *sqlx.DB, query string, args interface{}) (result sql.Result, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint("recover, ", e))
			return
		}
	}()
	tx := db.MustBegin()
	if args == nil {
		result = tx.MustExec(fmt.Sprint(query, " AND delete_at < to_timestamp(1000000)"))
	} else {
		result, err = tx.NamedExec(fmt.Sprint(query, " AND delete_at < to_timestamp(1000000)"), args)
	}
	err = LastRowsAffected(tx, result, err)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}
