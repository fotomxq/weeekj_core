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
	//预占位参数序号
	preemptionNum int
	//预占位数据结构列
	preemptionData []clientGetCtxPreemption
}

type clientGetCtxPreemption struct {
	//sql条件语句，注意必须用()包裹
	Query string
	//参数值
	Param any
	//参数序号
	Num int
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
func (t *ClientGetCtx) SetDefaultFields() *ClientGetCtx {
	t.fieldOne = []string{}
	for k := 0; k < len(t.clientCtx.client.fieldNameList); k++ {
		t.fieldOne = append(t.fieldOne, t.clientCtx.client.fieldNameList[k].DBName)
	}
	return t
}

func (t *ClientGetCtx) SetDefaultListFields() *ClientGetCtx {
	t.fieldOne = []string{}
	for k := 0; k < len(t.clientCtx.client.fieldNameList); k++ {
		if !t.clientCtx.client.fieldNameList[k].IsList {
			continue
		}
		t.fieldOne = append(t.fieldOne, t.clientCtx.client.fieldNameList[k].DBName)
	}
	return t
}

func (t *ClientGetCtx) getFieldsOne() string {
	return t.clientCtx.GetFields(t.fieldOne)
}

// getSQLGet 组合sql命令
func (t *ClientGetCtx) getSQLGet(where string) string {
	if where == "" {
		haveNewWhere := where != ""
		for k := 0; k < len(t.preemptionData); k++ {
			if !haveNewWhere && k == 0 {
				where = fmt.Sprint(t.preemptionData[k].Query)
			} else {
				where = fmt.Sprint(t.preemptionData[k].Query, " AND ", where)
			}
			t.clientCtx.appendArgs = append(t.clientCtx.appendArgs, t.preemptionData[k].Param)
		}
	}
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
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("id = $", t.preemptionNum), id)
	return t
}

func (t *ClientGetCtx) GetByIDAndOrgID(id int64, orgID int64) *ClientGetCtx {
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("id = $", t.preemptionNum), id)
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("($", t.preemptionNum, " < 0 OR org_id = $", t.preemptionNum, ")"), orgID)
	return t
}

func (t *ClientGetCtx) GetByMarkAndOrgID(mark string, orgID int64) *ClientGetCtx {
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("mark = $", t.preemptionNum), mark)
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("($", t.preemptionNum, " < 0 OR org_id = $", t.preemptionNum, ")"), orgID)
	return t
}

func (t *ClientGetCtx) GetByCodeAndOrgID(code string, orgID int64) *ClientGetCtx {
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("code = $", t.preemptionNum), code)
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("($", t.preemptionNum, " < 0 OR org_id = $", t.preemptionNum, ")"), orgID)
	return t
}

func (t *ClientGetCtx) GetByIDAndUserID(id int64, userID int64) *ClientGetCtx {
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("id = $", t.preemptionNum), id)
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("($", t.preemptionNum, " < 0 OR user_id = $", t.preemptionNum, ")"), userID)
	return t
}

func (t *ClientGetCtx) SetIDQuery(field string, param int64) *ClientGetCtx {
	return t.SetInt64Query(field, param)
}

func (t *ClientGetCtx) SetStringQuery(field string, param string) *ClientGetCtx {
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("(", field, " = $", t.preemptionNum, " OR $", t.preemptionNum, " = '')"), param)
	return t
}

func (t *ClientGetCtx) SetInt64Query(field string, param int64) *ClientGetCtx {
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("(", field, " = $", t.preemptionNum, " OR $", t.preemptionNum, " < 0)"), param)
	return t
}

func (t *ClientGetCtx) SetIntQuery(field string, param int) *ClientGetCtx {
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("(", field, " = $", t.preemptionNum, " OR $", t.preemptionNum, " < 0)"), param)
	return t
}

// SetDeleteQuery 设置删除查询
// 如果启动此设定，请注意基于查询条件的$顺序，叠加后使用，否则讲造成条件和参数不匹配
func (t *ClientGetCtx) SetDeleteQuery(field string, param bool) *ClientGetCtx {
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("((", field, " < to_timestamp(1000000) AND $", t.preemptionNum, " = false) OR (", field, " >= to_timestamp(1000000) AND $", t.preemptionNum, " = true))"), param)
	return t
}

// AppendWhere 直接覆盖where
func (t *ClientGetCtx) AppendWhere(where string, args ...interface{}) *ClientGetCtx {
	t.clientCtx.query = t.getSQLGet(where)
	t.clientCtx.appendArgs = append(t.clientCtx.appendArgs, args...)
	return t
}

func (t *ClientGetCtx) addPreemptionNum() {
	if t.preemptionNum < 1 {
		t.preemptionNum = 0
	}
	t.preemptionNum += 1
}

func (t *ClientGetCtx) addPreemption(query string, param any) {
	t.preemptionData = append(t.preemptionData, clientGetCtxPreemption{
		Query: query,
		Param: param,
		Num:   t.preemptionNum,
	})
}

func (t *ClientGetCtx) Result(data interface{}) error {
	if t.clientCtx.query == "" {
		t.clientCtx.query = t.getSQLGet("")
	}
	err := t.clientCtx.Get(data, t.clientCtx.query, t.clientCtx.appendArgs...)
	appendLog("get", t.clientCtx.query, false, t.clientCtx.client.startAt, data, err)
	return err
}
