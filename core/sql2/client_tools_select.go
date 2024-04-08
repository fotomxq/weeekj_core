package CoreSQL2

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"time"
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
	//预占err
	globErr []error
	//是否经过SelectList处理
	haveSelectList bool
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

// SetTimeBetweenByArgQuery 设置时间范围
func (t *ClientListCtx) SetTimeBetweenByArgQuery(field string, betweenAt ArgsTimeBetween) *ClientListCtx {
	betweenTimeAt, err := betweenAt.GetFields()
	if err != nil {
		t.globErr = append(t.globErr, err)
		return t
	}
	return t.SetTimeBetweenQuery(field, betweenTimeAt.MinTime, betweenTimeAt.MaxTime)
}

// SetTimeBetweenQuery 设置时间范围
func (t *ClientListCtx) SetTimeBetweenQuery(field string, startAt time.Time, endAt time.Time) *ClientListCtx {
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("(", field, " >= $", t.preemptionNum, ")"), startAt)
	t.addPreemptionNum()
	t.addPreemption(fmt.Sprint("(", field, " <= $", t.preemptionNum, ")"), endAt)
	return t
}

// SetSearchQuery 设置搜索查询
// 如果启动此设定，请注意基于查询条件的$顺序，叠加后使用，否则讲造成条件和参数不匹配
func (t *ClientListCtx) SetSearchQuery(fields []string, search string) *ClientListCtx {
	if search == "" {
		return t
	}
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

// SetIDsQuery 常规IDs判断查询
func (t *ClientListCtx) SetIDsQuery(field string, param pq.Int64Array) *ClientListCtx {
	if len(param) > 0 {
		t.addPreemptionNum()
		t.addPreemption(fmt.Sprint("(", field, " = ANY($", t.preemptionNum, "))"), param)
	}
	return t
}

// SetIntQuery 常规Int判断查询
func (t *ClientListCtx) SetIntQuery(field string, param int) *ClientListCtx {
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

// SelectList
// 如果启用了自动组合方法，请尽可能不要使用本方法where和args，否则请在where条件中明确一共注入了几个参数，并从对应参数为起点计算，避免$x顺序不匹配
func (t *ClientListCtx) SelectList(where string, args ...interface{}) *ClientListCtx {
	t.haveSelectList = true
	step := 0
	if t.pages.Page > 0 {
		step = int((t.pages.Page - 1) * t.pages.Max)
	}
	haveNewWhere := where != ""
	var newArgs []any
	for k := 0; k < len(t.preemptionData); k++ {
		if !haveNewWhere && k == 0 {
			where = fmt.Sprint(t.preemptionData[k].Query)
		} else {
			where = fmt.Sprint(t.preemptionData[k].Query, " AND ", where)
		}
		newArgs = append(newArgs, t.preemptionData[k].Param)
	}
	t.clientCtx.query = t.getSQLSelect(where, step, int(t.pages.Max), t.pages.Sort, t.pages.Desc)
	t.queryCount = t.getSQLSelectCount(where)
	if len(args) > 0 {
		newArgs = append(newArgs, args...)
		t.clientCtx.appendArgs = newArgs
	} else {
		t.clientCtx.appendArgs = newArgs
	}
	return t
}

func (t *ClientListCtx) Result(data interface{}) error {
	if !t.haveSelectList {
		t.SelectList("")
	}
	err := t.clientCtx.Select(data, t.clientCtx.query, t.clientCtx.appendArgs...)
	appendLog("select", t.clientCtx.query, false, t.clientCtx.client.startAt, data, err)
	//fmt.Println("query: ", t.clientCtx.query)
	return err
}

func (t *ClientListCtx) ResultAndCount(data interface{}) (count int64, err error) {
	if !t.haveSelectList {
		t.SelectList("")
	}
	if len(t.globErr) > 0 {
		var errStr string
		for _, v := range t.globErr {
			errStr = fmt.Sprint(errStr, ";", v.Error())
		}
		err = errors.New(errStr)
		return
	}
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

func (t *ClientListCtx) ResultCount() (count int64, err error) {
	if !t.haveSelectList {
		t.SelectList("")
	}
	err = t.clientCtx.Get(&count, t.queryCount, t.clientCtx.appendArgs...)
	if err != nil {
		appendLog("analysis", t.clientCtx.query, false, t.clientCtx.client.startAt, 1, err)
		return
	}
	appendLog("analysis", t.clientCtx.query, false, t.clientCtx.client.startAt, 1, err)
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
