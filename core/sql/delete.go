package CoreSQL

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

//GetDeleteSQL 获取删除或尚未删除的sql语句部分
func GetDeleteSQL(isRemove bool, where string) string {
	return GetDeleteSQLField(isRemove, where, "delete_at")
}

func GetDeleteSQLField(isRemove bool, where string, field string) string {
	if where == "" {
		if isRemove {
			where = field + " > to_timestamp(1000000)"
		} else {
			where = field + " < to_timestamp(1000000)"
		}
		return where
	} else {
		if isRemove {
			where = where + " AND " + field + " > to_timestamp(1000000)"
		} else {
			where = where + " AND " + field + " < to_timestamp(1000000)"
		}
		return where
	}
}

//DeleteOne 通用删除处理
func DeleteOne(db *sqlx.DB, tableName string, fieldName string, value interface{}) (result sql.Result, err error) {
	//捕捉错误
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			return
		}
	}()
	//删除操作
	tx := db.MustBegin()
	if value == nil {
		result, err = tx.Exec(fmt.Sprint("DELETE FROM ", tableName, " WHERE ", fieldName, "=:", fieldName))
	} else {
		result, err = tx.NamedExec(fmt.Sprint("DELETE FROM ", tableName, " WHERE ", fieldName, "=:", fieldName), value)
	}
	err = LastRowsAffected(tx, result, err)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

//DeleteOneSoft 通用软删除
func DeleteOneSoft(db *sqlx.DB, tableName string, fieldName string, value interface{}) (result sql.Result, err error) {
	//捕捉错误
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			return
		}
	}()
	//删除操作
	tx := db.MustBegin()
	if value == nil {
		result, err = tx.Exec(fmt.Sprint("UPDATE ", tableName, " SET delete_at=NOW() WHERE ", fieldName, "=:", fieldName))
	} else {
		result, err = tx.NamedExec(fmt.Sprint("UPDATE ", tableName, " SET delete_at=NOW() WHERE ", fieldName, "=:", fieldName), value)
	}
	err = LastRowsAffected(tx, result, err)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

//DeleteAllSoft 通用软删除
func DeleteAllSoft(db *sqlx.DB, tableName string, where string, value interface{}) (result sql.Result, err error) {
	//捕捉错误
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			return
		}
	}()
	//删除操作
	tx := db.MustBegin()
	if value == nil {
		result, err = tx.Exec(fmt.Sprint("UPDATE ", tableName, " SET delete_at=NOW() WHERE ", where))
	} else {
		result, err = tx.NamedExec(fmt.Sprint("UPDATE ", tableName, " SET delete_at=NOW() WHERE ", where), value)
	}
	err = LastRowsAffected(tx, result, err)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

//DeleteAll 清理所有符合条件的
func DeleteAll(db *sqlx.DB, tableName string, where string, args interface{}) (result sql.Result, err error) {
	//捕捉错误
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			return
		}
	}()
	//删除操作
	tx := db.MustBegin()
	if args == nil {
		if where == "" {
			result, err = tx.Exec(fmt.Sprint("DELETE FROM ", tableName))
		} else {
			result, err = tx.Exec(fmt.Sprint("DELETE FROM ", tableName, " WHERE ", where))
		}
	} else {
		if where == "" {
			result, err = tx.NamedExec(fmt.Sprint("DELETE FROM ", tableName), args)
		} else {
			result, err = tx.NamedExec(fmt.Sprint("DELETE FROM ", tableName, " WHERE ", where), args)
		}
	}
	err = LastRowsAffected(tx, result, err)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}
