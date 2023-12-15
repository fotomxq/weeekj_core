package CoreSQLMarks

import (
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// GetMarks 带有delete的方法设计
// 必须带有mark 字段设计
func GetMarks(dataList interface{}, tableName string, fields string, marks pq.StringArray) (err error) {
	//去重复
	newMarks := pq.StringArray{}
	for _, v := range marks {
		isFind := false
		for _, v2 := range newMarks {
			if v == v2 {
				isFind = true
				break
			}
		}
		if isFind {
			continue
		}
		newMarks = append(newMarks, v)
	}
	//获取数据
	err = Router2SystemConfig.MainDB.Select(dataList, "SELECT "+fields+" FROM "+tableName+" WHERE mark = ANY($1) LIMIT 100", newMarks)
	return
}

// GetMarksAndDelete 带有delete的方法设计
// 必须带有mark / delete_at 字段设计
func GetMarksAndDelete(dataList interface{}, tableName string, fields string, marks pq.StringArray, haveRemove bool) (err error) {
	//去重复
	newMarks := pq.StringArray{}
	for _, v := range marks {
		isFind := false
		for _, v2 := range newMarks {
			if v == v2 {
				isFind = true
				break
			}
		}
		if isFind {
			continue
		}
		newMarks = append(newMarks, v)
	}
	//获取数据
	err = Router2SystemConfig.MainDB.Select(dataList, "SELECT "+fields+" FROM "+tableName+" WHERE mark = ANY($1) AND ($2 = true OR delete_at < to_timestamp(1000000)) LIMIT 100", newMarks, haveRemove)
	return
}

// GetMarksName 必须带有mark / name字段设计
func GetMarksName(tableName string, marks pq.StringArray) (data map[string]string, err error) {
	//去重复
	newMarks := pq.StringArray{}
	for _, v := range marks {
		isFind := false
		for _, v2 := range newMarks {
			if v == v2 {
				isFind = true
				break
			}
		}
		if isFind {
			continue
		}
		newMarks = append(newMarks, v)
	}
	//获取数据
	type dataType struct {
		//mark
		Mark string `db:"mark" json:"mark"`
		//名称
		Name string `db:"name" json:"name"`
	}
	var dataList []dataType
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT mark, name FROM "+tableName+" WHERE mark = ANY($1) LIMIT 100", newMarks)
	if err == nil {
		data = map[string]string{}
		for _, v := range dataList {
			data[v.Mark] = v.Name
		}
	}
	return
}

// GetMarksNameAndDelete 必须带有mark / delete_at / name字段设计
func GetMarksNameAndDelete(tableName string, marks pq.StringArray, haveRemove bool) (data map[string]string, err error) {
	//去重复
	newMarks := pq.StringArray{}
	for _, v := range marks {
		isFind := false
		for _, v2 := range newMarks {
			if v == v2 {
				isFind = true
				break
			}
		}
		if isFind {
			continue
		}
		newMarks = append(newMarks, v)
	}
	//获取数据
	type dataType struct {
		//mark
		Mark string `db:"mark" json:"mark"`
		//名称
		Name string `db:"name" json:"name"`
	}
	var dataList []dataType
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT mark, name FROM "+tableName+" WHERE mark = ANY($1) AND ($2 = true OR delete_at < to_timestamp(1000000)) LIMIT 100", newMarks, haveRemove)
	if err == nil {
		data = map[string]string{}
		for _, v := range dataList {
			data[v.Mark] = v.Name
		}
	}
	return
}

// GetMarksTitleAndDelete 必须带有marks / delete_at / title字段设计
func GetMarksTitleAndDelete(tableName string, marks pq.StringArray, haveRemove bool) (data map[string]string, err error) {
	//去重复
	newMarks := pq.StringArray{}
	for _, v := range marks {
		isFind := false
		for _, v2 := range newMarks {
			if v == v2 {
				isFind = true
				break
			}
		}
		if isFind {
			continue
		}
		newMarks = append(newMarks, v)
	}
	//获取数据
	type dataType struct {
		//mark
		Mark string `db:"mark" json:"mark"`
		//标题
		Title string `db:"title" json:"title"`
	}
	var dataList []dataType
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT mark, title FROM "+tableName+" WHERE mark = ANY($1) AND ($2 = true OR delete_at < to_timestamp(1000000)) LIMIT 100", newMarks, haveRemove)
	if err == nil {
		data = map[string]string{}
		for _, v := range dataList {
			data[v.Mark] = v.Title
		}
	}
	return
}
