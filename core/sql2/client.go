package CoreSQL2

import (
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"github.com/golang-module/carbon"
	"reflect"
	"sync"
)

// Client 操作表核心
// 具体的表定义对象即可使用
type Client struct {
	//sql操作核心
	DB *SQLClient
	//表名称
	TableName string
	//空结构体
	// *args 注意采用引用关系，否则无法获取到结构体的类型
	StructData any
	//关键索引
	Key string
	//是否启动缓冲器
	openCache bool
	//缓冲时效
	cacheExpireSec int
	//缓冲对象
	cacheObj *CoreCache.CacheData
	//查看是否启动锁
	openViewLock bool
	viewLock     sync.Mutex
	//编辑是否启动锁
	openEditLock bool
	editLock     sync.Mutex
	//写入是否启动锁
	openCreateLock bool
	createLock     sync.Mutex
	//更新是否启动锁
	openUpdateLock bool
	updateLock     sync.Mutex
	//删除是否启动锁
	openDeleteLock bool
	deleteLock     sync.Mutex
	//开始时间
	startAt carbon.Carbon
	////////////////////////////////////////////////////
	//install特殊变量
	////////////////////////////////////////////////////
	//是否已经运行过install
	installHaveRun bool
	//等待插入的sql数据
	installAppendSQLData []string
	//主键发生数量，用于报错
	installNunIndexKeyNum int
}

func (t *Client) Init(mainDB *SQLClient, tableName string) *Client {
	t.DB = mainDB
	t.TableName = tableName
	t.Key = "id"
	return t
}

func (t *Client) Init2(mainDB *SQLClient, tableName string, structData any) (client *Client, err error) {
	t.DB = mainDB
	t.TableName = tableName
	t.Key = "id"
	t.StructData = structData
	err = t.InstallSQL()
	if err != nil {
		return
	}
	return t, nil
}

func (t *Client) SetKey(key string) *Client {
	t.Key = key
	return t
}

func (t *Client) GetKey() string {
	if t.Key != "" {
		return t.Key
	}
	return "id"
}

func (t *Client) SetCache(obj *CoreCache.CacheData) *Client {
	t.cacheObj = obj
	t.openCache = true
	if t.cacheExpireSec < 1 {
		t.cacheExpireSec = 60
	}
	return t
}

func (t *Client) SetViewLock(b bool) *Client {
	t.openViewLock = b
	return t
}

func (t *Client) SetEditLock(b bool) *Client {
	t.openEditLock = b
	return t
}

func (t *Client) SetCreateLock(b bool) *Client {
	t.openCreateLock = b
	return t
}

func (t *Client) SetUpdateLock(b bool) *Client {
	t.openUpdateLock = b
	return t
}

func (t *Client) SetDeleteLock(b bool) *Client {
	t.openDeleteLock = b
	return t
}

func (t *Client) SetExpireSec(sec int) *Client {
	t.cacheExpireSec = sec
	return t
}

// GetSortNameByJsonStruct 通过json/db结构体获取排序字段
func (t *Client) GetSortNameByJsonStruct(paramSort string, structData any, defaultSort string) (result string) {
	paramsType := reflect.TypeOf(structData).Elem()
	step := 0
	for step <= paramsType.NumField() {
		vField := paramsType.Field(step)
		jsonVal := vField.Tag.Get("json")
		//下一步
		step += 1
		if paramSort != jsonVal {
			continue
		}
		dbVal := vField.Tag.Get("db")
		result = dbVal
		break
	}
	if paramSort == "" {
		result = defaultSort
	}
	if result == paramSort {
		result = defaultSort
	}
	return result
}

func (t *Client) Get() *ClientGetCtx {
	t.startAt = CoreFilter.GetNowTimeCarbon()
	return &ClientGetCtx{
		clientCtx: &ClientCtx{
			client:          t,
			sqlNeedNoDelete: false,
			query:           "",
			appendArgs:      nil,
		},
		fieldOne:  []string{"*"},
		needLimit: false,
	}
}

func (t *Client) Select() *ClientListCtx {
	t.startAt = CoreFilter.GetNowTimeCarbon()
	return &ClientListCtx{
		clientCtx: &ClientCtx{
			client:          t,
			sqlNeedNoDelete: false,
			query:           "",
			appendArgs:      nil,
		},
		fieldsList: []string{"*"},
		fieldsSort: []string{},
		pages:      ArgsPages{},
		limitMax:   9999,
	}
}

func (t *Client) Insert() *ClientInsertCtx {
	t.startAt = CoreFilter.GetNowTimeCarbon()
	return &ClientInsertCtx{
		clientCtx: &ClientCtx{
			client:          t,
			sqlNeedNoDelete: false,
			query:           "",
			appendArgs:      []interface{}{},
		},
		fields: []string{},
	}
}

func (t *Client) Update() *ClientUpdateCtx {
	t.startAt = CoreFilter.GetNowTimeCarbon()
	return &ClientUpdateCtx{
		clientCtx: &ClientCtx{
			client:          t,
			sqlNeedNoDelete: false,
			query:           "",
			appendArgs:      nil,
		},
		updateFields:   []string{},
		updateFieldStr: "",
		whereFields:    []string{},
		whereArgs:      map[string]interface{}{},
		needUpdateAt:   false,
		haveWhere:      false,
	}
}

func (t *Client) Delete() *ClientDeleteCtx {
	t.startAt = CoreFilter.GetNowTimeCarbon()
	return &ClientDeleteCtx{
		clientCtx: &ClientCtx{
			client:          t,
			sqlNeedNoDelete: false,
			query:           "",
			appendArgs:      nil,
		},
		whereFields:    []string{},
		whereArgs:      map[string]interface{}{},
		needSoftDelete: false,
		haveWhere:      false,
	}
}

func (t *Client) Analysis() *ClientAnalysisCtx {
	t.startAt = CoreFilter.GetNowTimeCarbon()
	return &ClientAnalysisCtx{
		clientCtx: &ClientCtx{
			client:          t,
			sqlNeedNoDelete: false,
			query:           "",
			appendArgs:      nil,
		},
	}
}
