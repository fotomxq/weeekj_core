package CoreSQL2

import (
	"errors"
	"fmt"
	"strings"
)

type ClientDeleteCtx struct {
	//对象
	clientCtx *ClientCtx
	//条件字段
	// a = :a
	whereFields []string
	// key: a
	// val: interface{}
	whereArgs map[string]interface{}
	//是否为软删除
	needSoftDelete bool
	//是否已经追加了where
	haveWhere bool
}

// NeedSoft 是否软删除
func (t *ClientDeleteCtx) NeedSoft(b bool) *ClientDeleteCtx {
	t.needSoftDelete = b
	return t
}

// SetWhereAnd 添加等于关系判断
func (t *ClientDeleteCtx) SetWhereAnd(name string, val interface{}) *ClientDeleteCtx {
	t.whereFields = append(t.whereFields, fmt.Sprint(name, " = :", name))
	t.whereArgs[name] = val
	return t
}

// SetWhereOrThan 设置条件或查询关系
// 可用于负数跳过、0以上等于的判断机制
func (t *ClientDeleteCtx) SetWhereOrThan(name string, val interface{}) *ClientDeleteCtx {
	t.whereFields = append(t.whereFields, fmt.Sprint("(", ":", name, " < 0 OR ", name, " = :", name, ")"))
	t.whereArgs[name] = val
	return t
}

// AddWhereID 添加ID条件
func (t *ClientDeleteCtx) AddWhereID(id int64) *ClientDeleteCtx {
	t.SetWhereAnd("id", id)
	return t
}

// AddWhereOrgID 添加组织ID条件
func (t *ClientDeleteCtx) AddWhereOrgID(orgID int64) *ClientDeleteCtx {
	t.SetWhereOrThan("org_id", orgID)
	return t
}

// AddWhereUserID 添加用户ID条件
func (t *ClientDeleteCtx) AddWhereUserID(userID int64) *ClientDeleteCtx {
	t.SetWhereOrThan("user_id", userID)
	return t
}

// SetWhereStr 追加覆盖条件部分
func (t *ClientDeleteCtx) SetWhereStr(where string, arg map[string]interface{}) *ClientDeleteCtx {
	if !t.haveWhere {
		t.clientCtx.query = fmt.Sprint(t.clientCtx.query, " WHERE ")
		t.haveWhere = true
	}
	t.clientCtx.query = fmt.Sprint(t.clientCtx.query, where)
	for k, v := range arg {
		t.whereArgs[k] = v
	}
	return t
}

func (t *ClientDeleteCtx) AddQuery(fieldName string, param any) *ClientDeleteCtx {
	t.SetWhereOrThan(fieldName, param)
	return t
}

func (t *ClientDeleteCtx) makeWhere() {
	if !t.haveWhere {
		t.clientCtx.query = fmt.Sprint(t.clientCtx.query, " WHERE ")
		t.haveWhere = true
	}
	if len(t.whereFields) > 0 {
		t.clientCtx.query = fmt.Sprint(t.clientCtx.query, strings.Join(t.whereFields, " AND "))
	}
}

func (t *ClientDeleteCtx) makeArgs(arg map[string]interface{}) map[string]interface{} {
	if arg == nil {
		return t.whereArgs
	}
	for k, v := range t.whereArgs {
		arg[k] = v
	}
	return arg
}

// Exec 执行删除
func (t *ClientDeleteCtx) Exec(where string, args ...any) error {
	if t.needSoftDelete {
		t.clientCtx.query = fmt.Sprint("UPDATE ", t.clientCtx.client.TableName, " SET delete_at = NOW() WHERE ", where)
	} else {
		t.clientCtx.query = fmt.Sprint("DELETE ", "FROM ", t.clientCtx.client.TableName, " WHERE ", where)
	}
	t.makeWhere()
	result, err := t.clientCtx.Exec(t.clientCtx.query, args...)
	if err == nil {
		var affected int64
		affected, err = result.RowsAffected()
		if err == nil {
			if affected < 1 {
				err = errors.New(fmt.Sprint("no affected rows"))
			}
		}
	}
	if err != nil {
		err = errors.New(fmt.Sprint(err, ", ", t.clientCtx.getErrorQueryByArgs(args)))
	}
	appendLog("delete", t.clientCtx.query, false, t.clientCtx.client.startAt, nil, err)
	return err
}

// ExecAny 执行删除
func (t *ClientDeleteCtx) ExecAny(arg interface{}) error {
	if t.needSoftDelete {
		t.clientCtx.query = fmt.Sprint("UPDATE ", t.clientCtx.client.TableName, " SET delete_at = NOW() ", t.clientCtx.query)
	} else {
		t.clientCtx.query = fmt.Sprint("DELETE ", "FROM ", t.clientCtx.client.TableName, t.clientCtx.query)
	}
	t.makeWhere()
	result, err := t.clientCtx.NamedExec(t.clientCtx.query, arg)
	if err == nil {
		var affected int64
		affected, err = result.RowsAffected()
		if err == nil {
			if affected < 1 {
				err = errors.New(fmt.Sprint("no affected rows"))
			}
		}
	}
	if err != nil {
		err = errors.New(fmt.Sprint(err, ", ", t.clientCtx.getErrorQueryByArgs(arg)))
	}
	appendLog("delete", t.clientCtx.query, false, t.clientCtx.client.startAt, nil, err)
	return err
}

// ExecNamed 执行删除
// 需要给与map[string]interface{}的参数
// 如果没有，则可以给与nil，程序会自动跳过
func (t *ClientDeleteCtx) ExecNamed(arg map[string]interface{}) error {
	if t.needSoftDelete {
		t.clientCtx.query = fmt.Sprint("UPDATE ", t.clientCtx.client.TableName, " SET delete_at = NOW() ", t.clientCtx.query)
	} else {
		t.clientCtx.query = fmt.Sprint("DELETE ", "FROM ", t.clientCtx.client.TableName, t.clientCtx.query)
	}
	t.makeWhere()
	result, err := t.clientCtx.NamedExec(t.clientCtx.query, t.makeArgs(arg))
	if err == nil {
		var affected int64
		affected, err = result.RowsAffected()
		if err == nil {
			if affected < 1 {
				err = errors.New(fmt.Sprint("no affected rows"))
			}
		}
	}
	if err != nil {
		err = errors.New(fmt.Sprint(err, ", ", t.clientCtx.getErrorQueryByArgs(arg)))
	}
	appendLog("delete", t.clientCtx.query, false, t.clientCtx.client.startAt, nil, err)
	return err
}

// ClearAll 清理表
func (t *ClientDeleteCtx) ClearAll() error {
	t.clientCtx.query = "TRUNCATE " + "TABLE " + t.clientCtx.client.TableName
	_, err := t.clientCtx.Exec(t.clientCtx.query)
	if err != nil {
		err = errors.New(fmt.Sprint(err, ", ", t.clientCtx.query))
	}
	return err
}
