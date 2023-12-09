package CoreSQL2

import (
	"fmt"
	"strings"
)

type ClientUpdateCtx struct {
	//对象
	clientCtx *ClientCtx
	//更新字段
	updateFields   []string
	updateFieldStr string
	//条件字段
	// a = :a
	whereFields []string
	// key: a
	// val: interface{}
	whereArgs map[string]interface{}
	//是否需要更新时间
	needUpdateAt bool
	//是否已经追加了where
	haveWhere bool
}

func (t *ClientUpdateCtx) SetWhereAnd(name string, val interface{}) *ClientUpdateCtx {
	t.whereFields = append(t.whereFields, fmt.Sprint(name, " = :", name))
	t.whereArgs[name] = val
	return t
}

func (t *ClientUpdateCtx) SetWhereOrThan(name string, val interface{}) *ClientUpdateCtx {
	t.whereFields = append(t.whereFields, fmt.Sprint("(", ":", name, " < 0 OR ", name, " = :", name, ")"))
	t.whereArgs[name] = val
	return t
}

func (t *ClientUpdateCtx) NeedUpdateTime() *ClientUpdateCtx {
	t.needUpdateAt = true
	return t
}

func (t *ClientUpdateCtx) SetFields(fields []string) *ClientUpdateCtx {
	t.updateFields = fields
	t.clientCtx.query = fmt.Sprint("UPDATE ", t.clientCtx.client.TableName, " SET ")
	if t.needUpdateAt {
		t.clientCtx.query = fmt.Sprint(t.clientCtx.query, "update_at = NOW(),")
	}
	for k := 0; k < len(t.updateFields); k++ {
		t.updateFields[k] = fmt.Sprint(t.updateFields[k], " = :", t.updateFields[k])
	}
	t.clientCtx.query = fmt.Sprint(t.clientCtx.query, strings.Join(t.updateFields, ","))
	return t
}

func (t *ClientUpdateCtx) SetFieldStr(fields string) *ClientUpdateCtx {
	t.updateFieldStr = fields
	t.clientCtx.query = fmt.Sprint("UPDATE ", t.clientCtx.client.TableName, " SET ")
	if t.needUpdateAt {
		t.clientCtx.query = fmt.Sprint(t.clientCtx.query, "update_at = NOW(),")
	}
	t.clientCtx.query = fmt.Sprint(t.clientCtx.query, t.updateFieldStr)
	return t
}

func (t *ClientUpdateCtx) AddWhereID(id int64) *ClientUpdateCtx {
	t.SetWhereAnd("id", id)
	return t
}

func (t *ClientUpdateCtx) AddWhereOrgID(orgID int64) *ClientUpdateCtx {
	t.SetWhereOrThan("org_id", orgID)
	return t
}

func (t *ClientUpdateCtx) AddWhereUserID(userID int64) *ClientUpdateCtx {
	t.SetWhereOrThan("user_id", userID)
	return t
}

// SetWhereStr 追加覆盖条件部分
func (t *ClientUpdateCtx) SetWhereStr(where string, arg map[string]interface{}) *ClientUpdateCtx {
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

func (t *ClientUpdateCtx) makeWhere() {
	if !t.haveWhere {
		t.clientCtx.query = fmt.Sprint(t.clientCtx.query, " WHERE ")
		t.haveWhere = true
	}
	if len(t.whereFields) > 0 {
		t.clientCtx.query = fmt.Sprint(t.clientCtx.query, strings.Join(t.whereFields, " AND "))
	}
}

func (t *ClientUpdateCtx) makeArgs(arg map[string]interface{}) map[string]interface{} {
	if arg == nil {
		return t.whereArgs
	}
	for k, v := range t.whereArgs {
		arg[k] = v
	}
	return arg
}

func (t *ClientUpdateCtx) NamedExec(arg map[string]interface{}) error {
	t.makeWhere()
	_, err := t.clientCtx.NamedExec(t.clientCtx.query, t.makeArgs(arg))
	appendLog("update", t.clientCtx.query, false, t.clientCtx.client.startAt, nil, err)
	return err
}
