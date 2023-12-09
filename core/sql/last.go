package CoreSQL

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

//LastRows 事务关系尾部封装
func LastRows(tx *sqlx.Tx, err error) error {
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
		}
		return err
	}
	return nil
}

//LastRowsAffected 事务关系尾封装，带影响行判断
func LastRowsAffected(tx *sqlx.Tx, result sql.Result, err error) error {
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
		}
		return err
	}else{
		var rowsAffected int64
		rowsAffected, err = result.RowsAffected()
		if err == nil{
			if rowsAffected < 1{
				err2 := tx.Rollback()
				if err2 != nil {
					return errors.New("rollback failed: " + err2.Error())
				}
				return errors.New("no rows affected")
			}
		}
	}
	return nil
}

//LastRowsAffectedCreate 带有ID的创建尾部
func LastRowsAffectedCreate(tx *sqlx.Tx, stmt *sqlx.NamedStmt, args interface{}, err error) (int64, error) {
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return 0, errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
		}
		return 0, errors.New("prepare named " + err.Error())
	}
	defer func() {
		_ = stmt.Close()
	}()
	var lastID int64
	err = stmt.Get(&lastID, args)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return 0, errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
		}
		return 0, errors.New("get last id " + err.Error())
	}
	return lastID, nil
}