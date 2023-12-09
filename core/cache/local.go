package CoreCache

import (
	"sync"
)

//TODO: 需增加分区划定，将内容拆分为不同分区存储，这样提高运行效率。读取和写入过程每个分区内所有具体的锁

type cacheLocal struct {
	cacheList []cacheLocalData
	writeLock sync.RWMutex
}

type cacheLocalData struct {
	//注意约定区间，必须是:分割，例如：erp:core:config:id:1.2
	// 其中，erp:core:config:id这几个部分会被划定为不同的存储分区，后缀的参数将划定用于识别分区内具体内容
	mark     string
	val      any
	expireAt int64
}

// GetString 读取数据
func (t *cacheLocal) GetString(mark string) (data string, err error) {
	return
}

func (t *cacheLocal) GetInt(mark string) (data int, err error) {
	return
}

func (t *cacheLocal) GetInt64(mark string) (data int64, err error) {
	return
}

func (t *cacheLocal) GetBool(mark string) (data bool, err error) {
	return
}

func (t *cacheLocal) GetFloat64(mark string) (data float64, err error) {
	return
}

func (t *cacheLocal) GetByte(mark string) (data []byte, err error) {
	return
}

func (t *cacheLocal) GetScanStruct(mark string, data []any) (err error) {
	return
}

func (t *cacheLocal) FindKeys(mark string) (data []string, err error) {
	return
}

func (t *cacheLocal) GetStruct(mark string, data interface{}) (err error) {
	return
}

// SetString 写入数据
func (t *cacheLocal) SetString(mark string, val string, expire int) {
}

func (t *cacheLocal) SetInt(mark string, val int, expire int) {
}

func (t *cacheLocal) SetInt64(mark string, val int64, expire int) {
}

func (t *cacheLocal) SetAny(mark string, val any, expire int) {
}

func (t *cacheLocal) SetBool(mark string, val bool, expire int) {
}

func (t *cacheLocal) SetStruct(mark string, val any, expire int) {
}

// DeleteMark 删除数据
func (t *cacheLocal) DeleteMark(mark string) {
}

func (t *cacheLocal) DeleteSearchMark(mark string) {
}

// DeleteAll 清理所有数据
func (t *cacheLocal) DeleteAll() {
	t.writeLock.Lock()
	defer t.writeLock.Unlock()
	t.cacheList = []cacheLocalData{}
}

// GetListAll 获取列表所有数据
func (t *cacheLocal) GetListAll(key string, data any) (err error) {
	return
}

// AppendList 向列表末尾插入数据
func (t *cacheLocal) AppendList(key string, data ...any) {
	return
}

// DeleteListFirst 删除列表第一条数据
func (t *cacheLocal) DeleteListFirst(key string) {
	return
}

// GetListLen 获取列表长度
func (t *cacheLocal) GetListLen(key string) (count int, err error) {
	return
}
