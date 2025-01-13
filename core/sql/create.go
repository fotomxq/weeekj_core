package CoreSQL

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// CreateOne 快速创建
func CreateOne(db *sqlx.DB, query string, args interface{}) (result sql.Result, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			return
		}
	}()
	tx := db.MustBegin()
	if args == nil {
		result, err = tx.Exec(query)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				err = errors.New(fmt.Sprint(err, ",rollback error: ", err2))
				return
			}
			return
		}
	} else {
		result, err = tx.NamedExec(query, args)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				err = errors.New(fmt.Sprint(err, ",rollback error: ", err2))
				return
			}
			return
		}
	}
	err = LastRowsAffected(tx, result, err)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			err = errors.New(fmt.Sprint(err, ",rollback error: ", err2))
			return
		}
		return
	}
	err = tx.Commit()
	return
}

// CreateMore 快速创建多条
func CreateMore(db *sqlx.DB, query string, args []interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			return
		}
	}()
	tx := db.MustBegin()
	for _, v := range args {
		var result sql.Result
		result, err = tx.NamedExec(query, v)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				err = errors.New(fmt.Sprint(err, ",rollback error: ", err2))
				return
			}
			return
		}
		err = LastRowsAffected(tx, result, err)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				err = errors.New(fmt.Sprint(err, ",rollback error: ", err2))
				return
			}
			return
		}
	}
	err = tx.Commit()
	return
}

// CreateOneAndID 带有反馈id的创建
// 本方法将在query尾巴位置插入获取ID的方法，所以请确保尾部语法封闭
func CreateOneAndID(db *sqlx.DB, query string, args interface{}) (lastID int64, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			return
		}
	}()
	tx := db.MustBegin()
	var stmt *sqlx.NamedStmt
	stmt, err = tx.PrepareNamed(query + " RETURNING id;")
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			err = errors.New(fmt.Sprint(err, ",rollback error: ", err2))
			return
		}
		return
	}
	lastID, err = LastRowsAffectedCreate(tx, stmt, args, err)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			err = errors.New(fmt.Sprint(err, ",rollback error: ", err2))
			return
		}
		return
	}
	err = tx.Commit()
	return
}

// CreateOneAndData 获取并反射数据集
func CreateOneAndData(db *sqlx.DB, tableName string, query string, args interface{}, data interface{}) (err error) {
	var lastID int64
	lastID, err = CreateOneAndID(db, query, args)
	if err != nil {
		err = errors.New(fmt.Sprint("create one and id error: ", err))
		return
	}
	err = db.Get(data, "SELECT * FROM "+tableName+" WHERE id = $1", lastID)
	if err != nil {
		err = errors.New(fmt.Sprint("get data error: ", err, ",last id: ", lastID))
	}
	return
}
