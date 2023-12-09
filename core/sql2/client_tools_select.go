package CoreSQL2

import "fmt"

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

func (t *ClientListCtx) SelectList(where string, args ...interface{}) *ClientListCtx {
	step := 0
	if t.pages.Page > 0 {
		step = int((t.pages.Page - 1) * t.pages.Max)
	}
	t.clientCtx.query = t.getSQLSelect(where, step, int(t.pages.Max), t.pages.Sort, t.pages.Desc)
	t.queryCount = t.getSQLSelectCount(where)
	t.clientCtx.appendArgs = args
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
