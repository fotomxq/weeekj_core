package CoreCache

import "github.com/gomodule/redigo/redis"

//服务缓冲体系
/**
1. 支持本地独立服务的缓冲支持
2. 支持切换为redis缓冲支持
3. 未来可支持其他元素方案

CoreRedis.
Router2SystemConfig.MainCache.
*/

var (
	//CacheTime6Sec 6秒过期时间
	CacheTime6Sec = 6
	//CacheTime1Minutes 1分钟过期时间
	CacheTime1Minutes = 60
	//CacheTime1Hour 1小时过期时间
	CacheTime1Hour = 3600
	//CacheTime1Day 1天过期时间
	CacheTime1Day = CacheTime1Hour * 24
	//CacheTime2Day 2天过期时间
	CacheTime2Day = CacheTime1Day * 2
	//CacheTime3Day 3天过期时间
	CacheTime3Day = CacheTime1Day * 3
	//CacheTime1Week 1周过期时间
	CacheTime1Week = CacheTime1Day * 7
	//CacheTime1Month 1个月过期时间
	CacheTime1Month = CacheTime1Day * 30
	//CacheTime1Year 一年过期时间
	CacheTime1Year = CacheTime1Month * 12
)

type CacheData struct {
	//缓冲模式
	// local 本地缓冲器; redis 缓冲器
	globMode string
	//本地缓冲器
	localCache cacheLocal
	//redis
	redisCache cacheRedis
}

func (t *CacheData) Init(setMode string) {
	t.globMode = setMode
}

// InitRedis 连接服务
func (t *CacheData) InitRedis(url string, password string, num int) (err error) {
	t.redisCache.SetDatabaseNum(num)
	err = t.redisCache.Init(url, password)
	if err != nil {
		return
	}
	return
}

// GetRedisDB 获取缓冲redis数据库原始数据
func (t *CacheData) GetRedisDB() (db *redis.Pool) {
	return t.redisCache.db
}

// GetString 读取数据
func (t *CacheData) GetString(mark string) (data string, err error) {
	switch t.globMode {
	case "local":
		return t.localCache.GetString(mark)
	case "redis":
		return t.redisCache.GetString(mark)
	}
	return
}

func (t *CacheData) GetInt(mark string) (data int, err error) {
	switch t.globMode {
	case "local":
		return t.localCache.GetInt(mark)
	case "redis":
		return t.redisCache.GetInt(mark)
	}
	return
}

func (t *CacheData) GetInt64(mark string) (data int64, err error) {
	switch t.globMode {
	case "local":
		return t.localCache.GetInt64(mark)
	case "redis":
		return t.redisCache.GetInt64(mark)
	}
	return
}

func (t *CacheData) GetBool(mark string) (data bool, err error) {
	switch t.globMode {
	case "local":
		return t.localCache.GetBool(mark)
	case "redis":
		return t.redisCache.GetBool(mark)
	}
	return
}

func (t *CacheData) GetFloat64(mark string) (data float64, err error) {
	switch t.globMode {
	case "local":
		return t.localCache.GetFloat64(mark)
	case "redis":
		return t.redisCache.GetFloat64(mark)
	}
	return
}

func (t *CacheData) GetByte(mark string) (data []byte, err error) {
	switch t.globMode {
	case "local":
		return t.localCache.GetByte(mark)
	case "redis":
		return t.redisCache.GetByte(mark)
	}
	return
}

func (t *CacheData) GetScanStruct(mark string, data []any) (err error) {
	switch t.globMode {
	case "local":
		return t.localCache.GetScanStruct(mark, data)
	case "redis":
		return t.redisCache.GetScanStruct(mark, data)
	}
	return
}

func (t *CacheData) FindKeys(mark string) (data []string, err error) {
	switch t.globMode {
	case "local":
		return t.localCache.FindKeys(mark)
	case "redis":
		return t.redisCache.FindKeys(mark)
	}
	return
}

func (t *CacheData) GetStruct(mark string, data interface{}) (err error) {
	switch t.globMode {
	case "local":
		return t.localCache.GetStruct(mark, data)
	case "redis":
		return t.redisCache.GetStruct(mark, data)
	}
	return
}

// SetString 写入数据
func (t *CacheData) SetString(mark string, val string, expire int) {
	switch t.globMode {
	case "local":
		t.localCache.SetString(mark, val, expire)
	case "redis":
		t.redisCache.SetString(mark, val, expire)
	}
}

func (t *CacheData) SetInt(mark string, val int, expire int) {
	switch t.globMode {
	case "local":
		t.localCache.SetInt(mark, val, expire)
	case "redis":
		t.redisCache.SetInt(mark, val, expire)
	}
}

func (t *CacheData) SetInt64(mark string, val int64, expire int) {
	switch t.globMode {
	case "local":
		t.localCache.SetInt64(mark, val, expire)
	case "redis":
		t.redisCache.SetInt64(mark, val, expire)
	}
}

func (t *CacheData) SetAny(mark string, val any, expire int) {
	switch t.globMode {
	case "local":
		t.localCache.SetAny(mark, val, expire)
	case "redis":
		t.redisCache.SetAny(mark, val, expire)
	}
}

func (t *CacheData) SetBool(mark string, val bool, expire int) {
	switch t.globMode {
	case "local":
		t.localCache.SetBool(mark, val, expire)
	case "redis":
		t.redisCache.SetBool(mark, val, expire)
	}
}

func (t *CacheData) SetStruct(mark string, val any, expire int) {
	switch t.globMode {
	case "local":
		t.localCache.SetStruct(mark, val, expire)
	case "redis":
		t.redisCache.SetStruct(mark, val, expire)
	}
}

func (t *CacheData) DeleteMark(mark string) {
	switch t.globMode {
	case "local":
		t.localCache.DeleteMark(mark)
	case "redis":
		t.redisCache.DeleteMark(mark)
	}
}

func (t *CacheData) DeleteSearchMark(mark string) {
	switch t.globMode {
	case "local":
		t.localCache.DeleteSearchMark(mark)
	case "redis":
		t.redisCache.DeleteSearchMark(mark)
	}
}

// DeleteAll 清理所有数据
func (t *CacheData) DeleteAll() {
	switch t.globMode {
	case "local":
		t.localCache.DeleteAll()
	case "redis":
		t.redisCache.DeleteAll()
	}
}

// GetListAll 获取列表所有数据
func (t *CacheData) GetListAll(key string, data any) (err error) {
	switch t.globMode {
	case "local":
		return t.localCache.GetListAll(key, data)
	case "redis":
		return t.redisCache.GetListAll(key, data)
	}
	return
}

// AppendList 向列表末尾插入数据
func (t *CacheData) AppendList(key string, data ...any) {
	switch t.globMode {
	case "local":
		t.localCache.AppendList(key, data...)
	case "redis":
		t.redisCache.AppendList(key, data...)
	}
	return
}

// DeleteListFirst 删除列表第一条数据
func (t *CacheData) DeleteListFirst(key string) {
	switch t.globMode {
	case "local":
		t.localCache.DeleteListFirst(key)
	case "redis":
		t.redisCache.DeleteListFirst(key)
	}
	return
}

// GetListLen 获取列表长度
func (t *CacheData) GetListLen(key string) (count int, err error) {
	switch t.globMode {
	case "local":
		return t.localCache.GetListLen(key)
	case "redis":
		return t.redisCache.GetListLen(key)
	}
	return
}
