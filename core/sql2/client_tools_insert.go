package CoreSQL2

import (
	"errors"
	"fmt"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	"strings"
)

type ClientInsertCtx struct {
	//对象
	clientCtx *ClientCtx
	//插入字段
	fields []string
}

func (t *ClientInsertCtx) Add(args interface{}) *ClientInsertCtx {
	t.clientCtx.appendArgs = append(t.clientCtx.appendArgs, args)
	return t
}

func (t *ClientInsertCtx) SetFields(fields []string) *ClientInsertCtx {
	t.fields = fields
	return t
}

func (t *ClientInsertCtx) SetDefaultFields() *ClientInsertCtx {
	return t.SetFields(t.clientCtx.client.GetFields())
}

func (t *ClientInsertCtx) SetDefaultInsertFields() *ClientInsertCtx {
	var result []string
	for k := 0; k < len(t.clientCtx.client.fieldNameList); k++ {
		v := t.clientCtx.client.fieldNameList[k]
		if !v.IsCreateRequired {
			continue
		}
		result = append(result, v.DBName)
	}
	return t.SetFields(result)
}

func (t *ClientInsertCtx) getFieldVal() string {
	return fmt.Sprint(":", strings.Join(t.fields, ",:"))
}

func (t *ClientInsertCtx) getSQL() {
	t.clientCtx.query = fmt.Sprint("INSERT ", "INTO ", t.clientCtx.client.TableName, "(", t.clientCtx.GetFields(t.fields), ") VALUES(", t.getFieldVal(), ")")
}

// Exec 执行
// TODO: 核对代码存在问题，参数没有正确写入
func (t *ClientInsertCtx) Exec() error {
	t.getSQL()
	for _, vArgs := range t.clientCtx.appendArgs {
		_, err := t.clientCtx.Exec(t.clientCtx.query, vArgs)
		appendLog("insert", t.clientCtx.query, false, t.clientCtx.client.startAt, nil, err)
		if OpenDebug {
			if err != nil {
				CoreLog.Error("sql exec failed: ", err, ", sql: ", t.clientCtx.query, ", args: ", t.clientCtx.appendArgs)
			}
		}
	}
	return nil
}

func (t *ClientInsertCtx) ExecAndCheckID() error {
	addID, err := t.ExecAndResultID()
	if err != nil {
		return err
	}
	if addID < 1 {
		return errors.New("insert no id")
	}
	return nil
}

// ExecAndResultData 执行并返回数据
// TODO: 部分场景下正常写入，但反馈失败信息
func (t *ClientInsertCtx) ExecAndResultData(data interface{}) error {
	id, err := t.ExecAndResultID()
	if err != nil {
		return errors.New(fmt.Sprint("exec result id failed, ", err))
	}
	if id < 1 {
		return errors.New("no id")
	}
	err = t.clientCtx.Get(data, fmt.Sprint("SELECT * ", "FROM ", t.clientCtx.client.TableName, " WHERE ", t.clientCtx.client.GetKey(), " = $1"), id)
	if err != nil {
		return errors.New(fmt.Sprint("get data failed exec after, ", err))
	}
	return nil
}

func (t *ClientInsertCtx) ExecAndResultID() (int64, error) {
	c, err := t.execAndResultIDChild()
	appendLog("insert", t.clientCtx.query, false, t.clientCtx.client.startAt, nil, err)
	return c, err
}

func (t *ClientInsertCtx) execAndResultIDChild() (int64, error) {
	//构建sql
	if len(t.clientCtx.appendArgs) < 1 {
		return 0, errors.New("no insert data")
	}
	t.getSQL()
	//defer func() {
	//	if e := recover(); e != nil {
	//		appendLog("insert", t.clientCtx.query, false, t.clientCtx.client.startAt, nil, errors.New(fmt.Sprint(e)))
	//		return
	//	}
	//}()
	////事务关系开始运行
	//tx := t.clientCtx.MustBegin()
	////运行sql
	//result, err := tx.NamedExec(t.clientCtx.query, t.clientCtx.appendArgs[0])
	//if err != nil {
	//	err2 := tx.Rollback()
	//	if err2 != nil {
	//		err = errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
	//		return 0, err
	//	}
	//	return 0, err
	//}
	////获取影响的行ID
	//resultID, err := result.RowsAffected()
	//if err != nil {
	//	resultID, err = result.LastInsertId()
	//	if err != nil {
	//		err2 := tx.Rollback()
	//		if err2 != nil {
	//			err = errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
	//			return 0, err
	//		}
	//		return 0, err
	//	}
	//}
	//if resultID < 1 {
	//	err2 := tx.Rollback()
	//	if err2 != nil {
	//		err = errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
	//		return 0, err
	//	} else {
	//		err = errors.New("result id is empty")
	//	}
	//	return 0, err
	//}
	////启动事务
	//err = tx.Commit()
	//if err != nil {
	//	return 0, err
	//}
	////反馈
	//return resultID, nil
	tx := t.clientCtx.MustBegin()
	stmt, err := tx.PrepareNamed(t.clientCtx.query + " RETURNING id")
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return 0, errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
		}
		return 0, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	var id int64
	err = stmt.Get(&id, t.clientCtx.appendArgs[0])
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return 0, errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
		}
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return id, nil
}
