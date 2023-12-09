package CoreSQLHistory

import (
	"errors"
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

//自动化归档和处理程序
// 使用方法：引用到模块的run内运行，注意外部需要自行实现分布式部分
// 本方法主要给与一系列参数，自动实现数据的归档重组
// 注意，新表的ID请勿设置自增，以及其他自动填写或自增设置，本方法将直接覆盖

type ArgsRun struct {
	//归档多久之前的数据
	BeforeTime time.Time
	//时间字段
	// 留空后的默认值: create_at
	TimeFieldName string
	//旧的表名
	OldTableName string
	//新的表名
	NewTableName string
}

func Run(args *ArgsRun) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			return
		}
	}()
	if args.TimeFieldName == "" {
		args.TimeFieldName = "create_at"
	}
	var count int64
	err = Router2SystemConfig.MainDB.Get(&count, fmt.Sprint("SELECT COUNT(id) FROM ", args.OldTableName, " WHERE ", args.TimeFieldName, " < $1"), args.BeforeTime)
	if count < 1 {
		return
	}
	tx := Router2SystemConfig.MainDB.MustBegin()
	_, err = tx.Exec("INSERT INTO "+args.NewTableName+" (SELECT * FROM "+args.OldTableName+" WHERE "+args.TimeFieldName+" < $1)", args.BeforeTime)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			err = errors.New(fmt.Sprint(err, ", rollback: ", err2))
		}
		return
	}
	_, err = tx.Exec("DELETE FROM "+args.OldTableName+" WHERE "+args.TimeFieldName+" < $1", args.BeforeTime)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			err = errors.New(fmt.Sprint(err, ", rollback: ", err2))
		}
		return
	}
	err = tx.Commit()
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			err = errors.New(fmt.Sprint(err, ", rollback: ", err2))
		}
		return
	}
	return
}
