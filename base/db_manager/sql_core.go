package BaseDBManager

import (
	"fmt"
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//DB管理工具
// 表结构管理器

type SQLClient struct {
	//DB数据库句柄
	DB CoreSQL2.Client
	//表名称
	// eg: base_db_manager
	TableName string
	//主表缓冲标识码前缀
	// eg: xxx:xxx:xxx
	CacheName string
	//默认缓冲时效性
	// eg: 3600
	CacheTime int
	//结构映射关系
	// eg: &FieldsXXX{}
	FieldsAny any
	//是否存在更新时间
	HaveUpdateTime bool
	//是否存在软删除
	HaveSoftDelete bool
}

// Init 初始化数据库
func (t *SQLClient) Init() (err error) {
	_, err = t.DB.Init2(&Router2SystemConfig.MainSQL, t.TableName, t.FieldsAny)
	if err != nil {
		return
	}
	if t.CacheTime < 1 {
		t.CacheTime = CoreCache.CacheTime1Day
	}
	return
}

// GetByID 获取指定的ID
func (t *SQLClient) GetByID(id int64, data any) (err error) {
	//缓冲处理
	cacheMark := t.getCacheMark(id)
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, data); err == nil {
		return
	}
	//获取数据
	err = t.DB.Get().SetDefaultFields().GetByID(id).Result(data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, t.CacheTime)
	return
}

// Create 创建数据
func (t *SQLClient) Create(params map[string]any) (id int64, err error) {
	//执行创建
	id, err = t.DB.Insert().SetDefaultInsertFields().Add(params).ExecAndResultID()
	if err != nil {
		return
	}
	return
}

// UpdateByID 通过ID更新数据
func (t *SQLClient) UpdateByID(id int64, setFields []string, params map[string]any) (err error) {
	//执行更新
	if t.HaveUpdateTime {
		err = t.DB.Update().NeedUpdateTime().NeedSoft(t.HaveSoftDelete).AddWhereID(id).SetFields(setFields).NamedExec(params)
	} else {
		err = t.DB.Update().NeedSoft(t.HaveSoftDelete).AddWhereID(id).SetFields(setFields).NamedExec(params)
	}
	if err != nil {
		return
	}
	t.deleteCache(id)
	return
}

// UpdateDefaultFieldsByID 通过ID更新全量数据
func (t *SQLClient) UpdateDefaultFieldsByID(id int64, params map[string]any) (err error) {
	//执行更新
	if t.HaveUpdateTime {
		err = t.DB.Update().NeedUpdateTime().NeedSoft(t.HaveSoftDelete).AddWhereID(id).SetDefaultFields().NamedExec(params)
	} else {
		err = t.DB.Update().NeedSoft(t.HaveSoftDelete).AddWhereID(id).SetDefaultFields().NamedExec(params)
	}
	if err != nil {
		return
	}
	t.deleteCache(id)
	return
}

// DeleteByID 通过ID删除数据
func (t *SQLClient) DeleteByID(id int64) (err error) {
	//执行删除
	err = t.DB.Delete().NeedSoft(t.HaveSoftDelete).AddWhereID(id).ExecNamed(nil)
	if err != nil {
		return
	}
	t.deleteCache(id)
	return
}

// 缓冲处理机制
func (t *SQLClient) getCacheMark(id int64) string {
	return fmt.Sprintf("%s.%d", t.CacheName, id)
}

func (t *SQLClient) deleteCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(t.getCacheMark(id))
}
