package CoreCache

import (
	"encoding/json"
	"fmt"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	"github.com/gomodule/redigo/redis"
	"time"
)

type cacheRedis struct {
	//连接句柄
	db *redis.Pool
	//数据库编号
	databaseNum int
}

// Init 连接服务
func (t *cacheRedis) Init(url string, password string) (err error) {
	t.db = &redis.Pool{ //实例化一个连接池
		MaxIdle:     100, //最初的连接数量
		IdleTimeout: 180 * time.Second,
		// MaxActive:1000000,    //最大连接数量
		MaxActive: 4000, //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		Wait:      true,
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			var c redis.Conn
			c, err = redis.Dial("tcp", url, redis.DialPassword(password), redis.DialDatabase(t.databaseNum))
			if err != nil {
				fmt.Println("connect redis failed, tcp, ", err)
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err = c.Do("PING")
			if err != nil {
				//fmt.Println("connect redis failed, ping, ", err)
				return err
			}
			return err
		},
	}
	var c redis.Conn
	defer func() {
		_ = c.Close()
	}()
	if t.databaseNum < 1 {
		t.databaseNum = 0
	}
	c, err = redis.Dial("tcp", url, redis.DialPassword(password), redis.DialDatabase(t.databaseNum))
	return
}

// SetDatabaseNum 设置数据库编号
func (t *cacheRedis) SetDatabaseNum(num int) {
	t.databaseNum = num
}

// GetString 读取数据
func (t *cacheRedis) GetString(mark string) (data string, err error) {
	return redis.String(t.getData(mark))
}

func (t *cacheRedis) GetInt(mark string) (data int, err error) {
	return redis.Int(t.getData(mark))
}

func (t *cacheRedis) GetInt64(mark string) (data int64, err error) {
	return redis.Int64(t.getData(mark))
}

func (t *cacheRedis) GetBool(mark string) (data bool, err error) {
	return redis.Bool(t.getData(mark))
}

func (t *cacheRedis) GetFloat64(mark string) (data float64, err error) {
	return redis.Float64(t.getData(mark))
}

func (t *cacheRedis) GetByte(mark string) (data []byte, err error) {
	return redis.Bytes(t.getData(mark))
}

func (t *cacheRedis) GetScanStruct(mark string, data []any) (err error) {
	var reply []any
	reply, err = redis.Values(t.do("keys", mark))
	if err != nil {
		return
	}
	for _, v := range reply {
		vMark := string(v.([]byte))
		var vVal any
		err = t.GetStruct(vMark, &vVal)
		if err != nil {
			continue
		}
		data = append(data, vVal)
	}
	err = nil
	return
}

func (t *cacheRedis) FindKeys(mark string) (result []string, err error) {
	var reply []any
	reply, err = redis.Values(t.do("keys", mark))
	if err != nil {
		return
	}
	for _, v := range reply {
		result = append(result, string(v.([]byte)))
	}
	return
}

func (t *cacheRedis) GetStruct(mark string, data interface{}) (err error) {
	var val []byte
	val, err = redis.Bytes(t.getData(mark))
	if err != nil {
		return
	}
	err = json.Unmarshal(val, data)
	return
}

func (t *cacheRedis) getData(mark string) (reply any, err error) {
	c := t.db.Get()
	defer func() {
		_ = c.Close()
	}()
	reply, err = c.Do("Get", mark)
	if err != nil {
		CoreLog.Error("redis get data, ", err)
		return
	}
	return
}

func (t *cacheRedis) getList(mark string) (reply []any, err error) {
	c := t.db.Get()
	defer func() {
		_ = c.Close()
	}()
	reply, err = redis.Values(c.Do("hgetall", mark))
	if err != nil {
		CoreLog.Error("redis get data, ", err)
		return
	}
	return
}

func (t *cacheRedis) do(commandName string, args ...any) (reply any, err error) {
	c := t.db.Get()
	defer func() {
		_ = c.Close()
	}()
	reply, err = c.Do(commandName, args...)
	if err != nil {
		CoreLog.Error("redis do, ", err)
		return
	}
	return
}

// SetString 写入数据
func (t *cacheRedis) SetString(mark string, val string, expire int) {
	t.setData(mark, val, expire)
}

func (t *cacheRedis) SetInt(mark string, val int, expire int) {
	t.setData(mark, val, expire)
}

func (t *cacheRedis) SetInt64(mark string, val int64, expire int) {
	t.setData(mark, val, expire)
}

func (t *cacheRedis) SetAny(mark string, val any, expire int) {
	t.setData(mark, val, expire)
}

func (t *cacheRedis) SetBool(mark string, val bool, expire int) {
	t.setData(mark, val, expire)
}

func (t *cacheRedis) SetStruct(mark string, val any, expire int) {
	jsonData, err := json.Marshal(val)
	if err != nil {
		CoreLog.Error("redis set string json, ", err)
		return
	}
	t.setData(mark, jsonData, expire)
}

func (t *cacheRedis) setData(mark string, val any, expire int) {
	c := t.db.Get()
	defer func() {
		_ = c.Close()
	}()
	_, err := c.Do("Set", mark, val)
	if err != nil {
		CoreLog.Error("redis set string, ", err)
		return
	}
	if expire < 1 {
		expire = 2592000
	}
	if expire > 0 {
		_, err = c.Do("expire", mark, expire)
		if err != nil {
			CoreLog.Error("redis set string expire, ", err)
			return
		}
	}
}

// DeleteMark 删除数据
func (t *cacheRedis) DeleteMark(mark string) {
	t.deleteData(mark)
}

func (t *cacheRedis) DeleteSearchMark(mark string) {
	c := t.db.Get()
	defer func() {
		_ = c.Close()
	}()
	vals, err := redis.Strings(c.Do("KEYS", fmt.Sprint(mark, "*")))
	if err != nil {
		return
	}
	if len(vals) < 1 {
		return
	}
	err = c.Send("MULTI")
	if err != nil {
		CoreLog.Error("delete redis search MULTI ", err)
		return
	}
	for i := 0; i < len(vals); i++ {
		err = c.Send("DEL", vals[i])
		if err != nil {
			CoreLog.Error("delete redis search delete ", vals[i])
		}
	}
	_, err = redis.Values(c.Do("EXEC"))
	if err != nil {
		CoreLog.Error("delete redis search exec ", err)
	}
}

// DeleteAll 清理所有数据
func (t *cacheRedis) DeleteAll() {
	c := t.db.Get()
	defer func() {
		_ = c.Close()
	}()
	_ = c.Send("flushall")
}

func (t *cacheRedis) deleteData(mark string) {
	c := t.db.Get()
	defer func() {
		_ = c.Close()
	}()
	_, err := c.Do("DEL", mark)
	if err != nil {
		CoreLog.Error("redis get data, ", err)
		return
	}
	return
}

// GetListAll 获取列表所有数据
func (t *cacheRedis) GetListAll(key string, data any) (err error) {
	var val []any
	val, err = redis.Values(t.do("LRANGE", key, 0, -1))
	if err != nil {
		return
	}
	err = redis.ScanStruct(val, data)
	if err != nil {
		return
	}
	return
}

// AppendList 向列表末尾插入数据
func (t *cacheRedis) AppendList(key string, data ...any) {
	data = append([]any{key}, data...)
	_, _ = t.do("RPUSH", data...)
}

// DeleteListFirst 删除列表第一条数据
func (t *cacheRedis) DeleteListFirst(key string) {
	_, _ = t.do("Lpop", key)
}

// GetListLen 获取列表长度
func (t *cacheRedis) GetListLen(key string) (count int, err error) {
	return redis.Int(t.do("LLEN", key))
}
