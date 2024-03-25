package CoreSQL2

import (
	"fmt"
	"strings"
)

type ClientListCtx struct {
	//对象
	clientCtx *ClientCtx
	//列表读取字段列
	fieldsList []string
	//限定排序
	// 如果没有指定，则按照列表字段执行
	fieldsSort []string
	//分页
	pages ArgsPages
	//分页长度极限
	limitMax int
	//获取总数query
	queryCount string
	//预占位参数序号
	preemptionNum int
	//预占位数据结构列
	preemptionData []clientListCtxPreemption
}

type clientListCtxPreemption struct {
	//sql条件语句，注意必须用()包裹
	Query string
	//参数值
	Param any
	//参数序号
	Num int
}

func (t *ClientListCtx) SetFieldsList(fields []string) *ClientListCtx {
	t.fieldsList = fields
	if len(t.fieldsSort) < 1 {
		t.SetFieldsSort(fields)
	}
	return t
}

func (t *ClientListCtx) SetFieldsSort(fields []string) *ClientListCtx {
	t.fieldsSort = fields
	return t
}

func (t *ClientListCtx) getFieldsList() string {
	return t.clientCtx.GetFields(t.fieldsList)
}

func (t *ClientListCtx) SetPages(pages ArgsPages) *ClientListCtx {
	t.pages = pages
	return t
}

// SetDeleteQuery 设置删除查询
// 如果启动此设定，请注意基于查询条件的$顺序，叠加后使用，否则讲造成条件和参数不匹配
func (t *ClientListCtx) SetDeleteQuery(field string, param bool) *ClientListCtx {
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("((", field, " < to_timestamp(1000000) AND $", t.preemptionNum, " = false) OR (", field, " >= to_timestamp(1000000) AND $", t.preemptionNum, " = true))"), param)
	return t
}

// SetSearchQuery 设置搜索查询
// 如果启动此设定，请注意基于查询条件的$顺序，叠加后使用，否则讲造成条件和参数不匹配
func (t *ClientListCtx) SetSearchQuery(fields []string, search string) *ClientListCtx {
	t.addPreemptionNum()
	var newQuerys []string
	for _, v := range fields {
		newQuerys = append(newQuerys, fmt.Sprint(v, " ILIKE $", t.preemptionNum))
	}
	newQuery := fmt.Sprint("((", strings.Join(newQuerys, " OR "), ") OR $", t.preemptionNum, " = '')")
	t.addPreemption(newQuery, "%"+search+"%")
	return t
}

// SetIDQuery 常规ID判断查询
func (t *ClientListCtx) SetIDQuery(field string, param int64) *ClientListCtx {
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("(", field, " = $", t.preemptionNum, " OR $", t.preemptionNum, " < 0)"), param)
	return t
}

// SetStringQuery 常规字符串判断查询
func (t *ClientListCtx) SetStringQuery(field string, param string) *ClientListCtx {
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("(", field, " = $", t.preemptionNum, " OR $", t.preemptionNum, " = '')"), param)
	return t
}

func (t *ClientListCtx) SelectList(where string, args ...interface{}) *ClientListCtx {
	step := 0
	if t.pages.Page > 0 {
		step = int((t.pages.Page - 1) * t.pages.Max)
	}
	var newArgs []any
	for k := 0; k < len(t.preemptionData); k++ {
		where = fmt.Sprint(t.preemptionData[k].Query, " AND ", where)
		newArgs = append(newArgs, t.preemptionData[k].Param)
	}
	t.clientCtx.query = t.getSQLSelect(where, step, int(t.pages.Max), t.pages.Sort, t.pages.Desc)
	t.queryCount = t.getSQLSelectCount(where)
	if len(args) > 0 {
		newArgs = append(newArgs, args...)
		t.clientCtx.appendArgs = newArgs
	} else {
		t.clientCtx.appendArgs = args
	}
	return t
}

func (t *ClientListCtx) Result(data interface{}) error {
	err := t.clientCtx.Select(data, t.clientCtx.query, t.clientCtx.appendArgs...)
	appendLog("select", t.clientCtx.query, false, t.clientCtx.client.startAt, data, err)
	//fmt.Println("query: ", t.clientCtx.query)
	return err
}

func (t *ClientListCtx) ResultAndCount(data interface{}) (count int64, err error) {
	err = t.clientCtx.Select(data, t.clientCtx.query, t.clientCtx.appendArgs...)
	if err != nil {
		appendLog("select", t.clientCtx.query, false, t.clientCtx.client.startAt, data, err)
		return
	}
	err = t.clientCtx.Get(&count, t.queryCount, t.clientCtx.appendArgs...)
	if err != nil {
		appendLog("select", t.clientCtx.query, false, t.clientCtx.client.startAt, data, err)
		return
	}
	appendLog("select", t.clientCtx.query, false, t.clientCtx.client.startAt, data, err)
	return
}

func (t *ClientListCtx) addPreemptionNum() {
	if t.preemptionNum < 1 {
		t.preemptionNum = 0
	}
	t.preemptionNum += 1
}

func (t *ClientListCtx) addPreemption(query string, param any) {
	t.preemptionData = append(t.preemptionData, clientListCtxPreemption{
		Query: query,
		Param: param,
		Num:   t.preemptionNum,
	})
}

func (t *ClientListCtx) getSQLSelect(where string, step int, limit int, sort string, desc bool) string {
	appendQuery := ""
	if sort != "" {
		isFind := false
		for _, v := range t.fieldsSort {
			if v == sort {
				isFind = true
			}
		}
		if !isFind {
			if len(t.fieldsSort) > 0 {
				sort = t.fieldsSort[0]
			} else {
				if len(t.fieldsList) > 0 {
					if t.fieldsList[0] == "*" {
						sort = ""
					} else {
						sort = t.fieldsList[0]
					}
				}
			}
		}
		if sort != "" {
			appendQuery = fmt.Sprint(appendQuery, "ORDER BY ", sort)
			if desc {
				appendQuery = fmt.Sprint(appendQuery, " DESC")
			}
		}
	}
	if limit > 0 {
		if appendQuery != "" {
			appendQuery = appendQuery + " "
		}
		if limit > t.limitMax {
			limit = t.limitMax
		}
		appendQuery = fmt.Sprint(appendQuery, "LIMIT ", limit)
	} else {
		//if appendQuery != "" {
		//	appendQuery = appendQuery + " "
		//}
		//appendQuery = fmt.Sprint(appendQuery, "LIMIT ", limit)
	}
	if step > 0 {
		if appendQuery != "" {
			appendQuery = appendQuery + " "
		}
		appendQuery = fmt.Sprint(appendQuery, "OFFSET ", step)
	}
	if appendQuery != "" {
		return fmt.Sprint(t.clientCtx.getSQLWhere(fmt.Sprint("SELECT ", t.getFieldsList(), " FROM ", t.clientCtx.client.TableName), where), " ", appendQuery)
	} else {
		return t.clientCtx.getSQLWhere(fmt.Sprint("SELECT ", t.getFieldsList(), " FROM ", t.clientCtx.client.TableName), where)
	}
}

func (t *ClientListCtx) getSQLSelectCount(where string) string {
	return t.clientCtx.getSQLWhere(fmt.Sprint("SELECT COUNT("+t.clientCtx.client.GetKey()+") FROM ", t.clientCtx.client.TableName), where)
}
