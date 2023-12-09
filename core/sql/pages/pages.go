package CoreSQLPages

import (
	"errors"
	"fmt"
	"strconv"
)

// ArgsDataList 列表数据结构
type ArgsDataList struct {
	Page int64  `json:"page" check:"page"`
	Max  int64  `json:"max" check:"max"`
	Sort string `json:"sort" check:"sort"`
	Desc bool   `json:"desc" check:"desc"`
}

// GetCacheMark 获取缓冲名称
func (t *ArgsDataList) GetCacheMark() string {
	return fmt.Sprint(t.Page, ".", t.Max, ".", t.Sort, ".", t.Desc)
}

// GetArgsDataList 列表快速生成默认参数
func GetArgsDataList() ArgsDataList {
	return ArgsDataList{
		Page: 1,
		Max:  10,
		Sort: "_id",
		Desc: false,
	}
}

// GetSQL 处理分页
// page > 0才会处理分页，否则只有排序
func GetSQL(args *ArgsDataList) (query string) {
	query = "ORDER BY " + args.Sort
	if args.Desc {
		query = query + " DESC"
	} else {
		query = query + " ASC"
	}
	if args.Page > 0 && args.Max > 0 {
		query = query + " LIMIT " + strconv.FormatInt(args.Max, 10) + " OFFSET " + strconv.FormatInt(args.Max*(args.Page-1), 10)
	}
	return
}

func FilterSQLSort(sortName string, filters []string) (query string) {
	for _, v := range filters {
		if v == sortName {
			return v
		}
	}
	return filters[0]
}

func GetSQLDesc(desc bool) (query string) {
	if desc {
		return "DESC"
	} else {
		return "ASC"
	}
}

// GetMaps 获取maps层级的分页设计
// 注意部分mark不能继续使用，会发生冲突：page_sort / page_limit / page_offset
func GetMaps(args *ArgsDataList, maps map[string]interface{}) (query string, newMaps map[string]interface{}) {
	query = "ORDER BY " + args.Sort
	if maps == nil {
		maps = map[string]interface{}{}
	}
	if args.Desc {
		query = query + " DESC"
	} else {
		query = query + " ASC"
	}
	if args.Page > 0 && args.Max > 0 {
		query = query + " LIMIT :page_limit OFFSET :page_offset"
		maps["page_limit"] = args.Max
		maps["page_offset"] = (args.Page - 1) * args.Max
	}
	return query, maps
}

func GetMapsArgs(args *ArgsDataList, pageLimit int, offset int) (query string) {
	query = "ORDER BY " + args.Sort
	if args.Desc {
		query = query + " DESC"
	} else {
		query = query + " ASC"
	}
	if args.Page > 0 && args.Max > 0 {
		query = query + fmt.Sprint(" LIMIT ", pageLimit, " OFFSET ", offset)
	}
	return query
}

func GetMapsAnd(args *ArgsDataList, query string, maps map[string]interface{}) (newQuery string, newMaps map[string]interface{}) {
	var newWhere string
	newWhere, newMaps = GetMaps(args, maps)
	if newWhere != "" {
		newQuery = query + " " + newWhere
		return
	} else {
		return query, maps
	}
}

func GetMapsAndArgs(args *ArgsDataList, query string) (newQuery string) {
	var newWhere string
	newWhere = GetMapsArgs(args, int(args.Max), int((args.Page-1)*args.Max))
	if newWhere != "" {
		newQuery = query + " " + newWhere
		return
	} else {
		return query
	}
}

// GetMapsAndFilter 带有约束的排序
func GetMapsAndFilter(args *ArgsDataList, query string, maps map[string]interface{}, filterSort []string) (newQuery string, newMaps map[string]interface{}, err error) {
	isFind := false
	for _, v := range filterSort {
		if v == args.Sort {
			isFind = true
			break
		}
	}
	if !isFind {
		err = errors.New("sort not support")
		return
	}
	newQuery, newMaps = GetMapsAnd(args, query, maps)
	return
}

func GetMapsAndFilterArgs(args *ArgsDataList, query string, filterSort []string) (newQuery string, err error) {
	isFind := false
	for _, v := range filterSort {
		if v == args.Sort {
			isFind = true
			break
		}
	}
	if !isFind {
		err = errors.New("sort not support")
		return
	}
	newQuery = GetMapsAndArgs(args, query)
	return
}
