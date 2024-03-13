package CoreSQL2

import (
	"fmt"
)

type ClientGetCtx struct {
	//对象
	clientCtx *ClientCtx
	//单条读取字段列
	fieldOne []string
	//是否需要限制
	needLimit bool
	//是否需要排序
	needSort bool
	//排序字段
	sortField string
	//排序规则
	sortIsDesc bool
}

// NeedLimit 限制避免条件过于宽松，造成获取多条数据的异常问题
func (t *ClientGetCtx) NeedLimit() *ClientGetCtx {
	t.needLimit = true
	return t
}

// NeedSort 启动排序，用于反馈仅有数据时，获得最新或最早一条数据的设计
func (t *ClientGetCtx) NeedSort(needSort bool, sortField string, sortIsDesc bool) *ClientGetCtx {
	t.needSort = needSort
	t.sortField = sortField
	t.sortIsDesc = sortIsDesc
	return t
}

func (t *ClientGetCtx) SetFieldsOne(fields []string) *ClientGetCtx {
	t.fieldOne = fields
	return t
}

func (t *ClientGetCtx) getFieldsOne() string {
	return t.clientCtx.GetFields(t.fieldOne)
}

// getSQLGet 组合sql命令
func (t *ClientGetCtx) getSQLGet(where string) string {
	query := t.clientCtx.getSQLWhere(fmt.Sprint("SELECT ", t.getFieldsOne(), " FROM ", t.clientCtx.client.TableName), where)
	if t.needSort {
		query = fmt.Sprint(query, " ORDER BY ", t.sortField)
		if t.sortIsDesc {
			query = fmt.Sprint(query, " DESC")
		}
	}
	if t.needLimit {
		query = fmt.Sprint(query, " LIMIT 1")
	}
	return query
}

func (t *ClientGetCtx) GetByID(id int64) *ClientGetCtx {
	t.clientCtx.query = t.getSQLGet("id = $1")
	t.clientCtx.appendArgs = append(t.clientCtx.appendArgs, id)
	return t
}

func (t *ClientGetCtx) GetByIDAndOrgID(id int64, orgID int64) *ClientGetCtx {
	t.clientCtx.query = t.getSQLGet("id = $1 AND ($2 < 0 OR org_id = $2)")
	t.clientCtx.appendArgs = append(t.clientCtx.appendArgs, id)
	t.clientCtx.appendArgs = append(t.clientCtx.appendArgs, orgID)
	return t
}

func (t *ClientGetCtx) GetByMarkAndOrgID(mark string, orgID int64) *ClientGetCtx {
	t.clientCtx.query = t.getSQLGet("mark = $1 AND ($2 < 0 OR org_id = $2)")
	t.clientCtx.appendArgs = append(t.clientCtx.appendArgs, mark)
	t.clientCtx.appendArgs = append(t.clientCtx.appendArgs, orgID)
	return t
}

func (t *ClientGetCtx) GetByCodeAndOrgID(code string, orgID int64) *ClientGetCtx {
	t.clientCtx.query = t.getSQLGet("code = $1 AND ($2 < 0 OR org_id = $2)")
	t.clientCtx.appendArgs = append(t.clientCtx.appendArgs, code)
	t.clientCtx.appendArgs = append(t.clientCtx.appendArgs, orgID)
	return t
}

func (t *ClientGetCtx) GetByIDAndUserID(id int64, userID int64) *ClientGetCtx {
	t.clientCtx.query = t.getSQLGet("id = $1 AND ($2 < 0 OR user_id = $2)")
	t.clientCtx.appendArgs = append(t.clientCtx.appendArgs, id)
	t.clientCtx.appendArgs = append(t.clientCtx.appendArgs, userID)
	return t
}

func (t *ClientGetCtx) AppendWhere(where string, args ...interface{}) *ClientGetCtx {
	t.clientCtx.query = t.getSQLGet(where)
	t.clientCtx.appendArgs = append(t.clientCtx.appendArgs, args...)
	return t
}

func (t *ClientGetCtx) Result(data interface{}) error {
	err := t.clientCtx.Get(data, t.clientCtx.query, t.clientCtx.appendArgs...)
	appendLog("get", t.clientCtx.query, false, t.clientCtx.client.startAt, data, err)
	return err
}
