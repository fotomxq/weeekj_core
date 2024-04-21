package CoreSQL2

import (
	"fmt"
	"strings"
)

type ClientCtx struct {
	//句柄
	client *Client
	//上下文链
	//是否剔除被删除数据 delete_at < to_timestamp(1000000)
	sqlNeedNoDelete bool
	//组装的语句
	query string
	//等待参数
	appendArgs []interface{}
}

func (t *ClientCtx) GetFields(l []string) string {
	return strings.Join(l, ",")
}

func (t *ClientCtx) DataNoDelete() *ClientCtx {
	t.sqlNeedNoDelete = true
	return t
}

func (t *ClientCtx) getSQLWhere(query string, where string) string {
	if t.sqlNeedNoDelete {
		if where != "" {
			where = fmt.Sprint("(", where, ") AND delete_at", t.getSQLTimeLessDefault())
		} else {
			where = fmt.Sprint("delete_at", t.getSQLTimeLessDefault())
		}
	}
	if where != "" {
		return fmt.Sprint(query, " WHERE ", where)
	} else {
		return fmt.Sprint(query)
	}
}

func (t *ClientCtx) getSQLTimeLessDefault() string {
	return " <= to_timestamp(1000000)"
}

func (t *ClientCtx) getSQLTimeThanDefault() string {
	return " > to_timestamp(1000000)"
}

func (t *ClientCtx) getSQLTimeLessNow() string {
	return " < NOW()"
}

func (t *ClientCtx) getSQLTimeThanNow() string {
	return " >= NOW()"
}

func (t *ClientCtx) getErrorQueryByArgs(args any) string {
	//检查参数长度，自动裁剪
	argsStr := fmt.Sprint(args)
	if len(argsStr) > 100 {
		argsStr = fmt.Sprint(argsStr[:100], "...")
	}
	return fmt.Sprint("query: ", t.query, ", args: ", argsStr)
}
