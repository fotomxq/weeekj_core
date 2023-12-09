package CoreSQLFrom

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

//FieldsFrom 来源设计
type FieldsFrom struct {
	System string `db:"system" json:"system" check:"system" empty:"true"`
	ID     int64  `db:"id" json:"id" check:"id" empty:"true"`
	Mark   string `db:"mark" json:"mark" check:"mark" empty:"true"`
	Name   string `db:"name" json:"name" check:"name" empty:"true"`
}

//Value sql底层处理器
func (t FieldsFrom) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsFrom) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

//GetFromOne 直接sql db数据
// 可用于获取唯一的来源结构体时
func (t *FieldsFrom) GetFromOne(db *sqlx.DB, tableName string, fields string, fieldName string, data interface{}) (err error) {
	if t.System == "" {
		err = errors.New("system is empty")
		return
	}
	var createData string
	createData, err = t.GetRaw()
	if err != nil {
		return
	}
	err = db.Get(data, fmt.Sprint("SELECT ", fields, " FROM ", tableName, " WHERE ", fieldName, " @> $1 LIMIT 1;"), createData)
	return
}

//CheckEg 检查两个数据是否相同？
func (t *FieldsFrom) CheckEg(b FieldsFrom) bool {
	return t.System == b.System && t.ID == b.ID && t.Mark == b.Mark
}

//GetSQL 获取sql语句带maps数据包
func (t *FieldsFrom) GetSQL(fieldName string, fieldMark string, maps map[string]interface{}) (query string, newMaps map[string]interface{}, err error) {
	query = fmt.Sprint(fieldName, " @> :", fieldMark)
	newMaps, err = t.GetMaps(fieldMark, maps)
	return
}

//GetMaps 获取maps组合数据包
func (t *FieldsFrom) GetMaps(fieldMark string, maps map[string]interface{}) (newMaps map[string]interface{}, err error) {
	if maps == nil {
		maps = map[string]interface{}{}
	}
	maps[fieldMark], err = t.GetRaw()
	return maps, err
}

//GetRaw 直接获取json处理后的筛选数据
func (t FieldsFrom) GetRaw() (data string, err error) {
	var jsonByte []byte
	jsonByte, err = json.Marshal(t)
	if err != nil {
		return
	}
	data = string(jsonByte)
	return
}

//GetString 直接获取json处理后的筛选数据
func (t FieldsFrom) GetString() string {
	return fmt.Sprint(t.System, ".", t.Mark, ".", t.ID, ".", t.Name)
}

//GetRawNoName 获取json不带name的结构
// mark=*，将自动忽略该参数
func (t FieldsFrom) GetRawNoName() (data string, err error) {
	var jsonByte []byte
	if t.Mark == "*" {
		type arg struct {
			System string `db:"system" json:"system"`
			ID     int64  `db:"id" json:"id"`
		}
		jsonByte, err = json.Marshal(arg{
			System: t.System,
			ID:     t.ID,
		})
		if err != nil {
			return
		}
	} else {
		type arg struct {
			System string `db:"system" json:"system"`
			ID     int64  `db:"id" json:"id"`
			Mark   string `db:"mark" json:"mark"`
		}
		jsonByte, err = json.Marshal(arg{
			System: t.System,
			ID:     t.ID,
			Mark:   t.Mark,
		})
		if err != nil {
			return
		}
	}
	data = string(jsonByte)
	return
}

//GetList 获取列表组合专用
// 自动判断是否存在参数，如果不存在则反馈原数据，否则写入需要的参数集合
func (t *FieldsFrom) GetList(fieldName string, fieldMark string, maps map[string]interface{}) (query string, newMap map[string]interface{}, err error) {
	if maps == nil {
		maps = map[string]interface{}{}
	}
	var jsonByte []byte
	if t.System != "" {
		if t.ID > 0 {
			if t.Mark != "" {
				type arg struct {
					System string `db:"system" json:"system"`
					ID     int64  `db:"id" json:"id"`
					Mark   string `db:"mark" json:"mark"`
				}
				jsonByte, err = json.Marshal(arg{
					System: t.System,
					ID:     t.ID,
					Mark:   t.Mark,
				})
				if err != nil {
					return
				}
			} else {
				type arg struct {
					System string `db:"system" json:"system"`
					ID     int64  `db:"id" json:"id"`
				}
				jsonByte, err = json.Marshal(arg{
					System: t.System,
					ID:     t.ID,
				})
				if err != nil {
					return
				}
			}
		} else {
			if t.Mark != "" {
				type arg struct {
					System string `db:"system" json:"system"`
					Mark   string `db:"mark" json:"mark"`
				}
				jsonByte, err = json.Marshal(arg{
					System: t.System,
					Mark:   t.Mark,
				})
				if err != nil {
					return
				}
			} else {
				type arg struct {
					System string `db:"system" json:"system"`
				}
				jsonByte, err = json.Marshal(arg{
					System: t.System,
				})
				if err != nil {
					return
				}
			}
		}
		maps[fieldMark] = string(jsonByte)
		query = fmt.Sprint(fieldName, " @> :", fieldMark)
	} else {
		return "", maps, nil
	}
	return query, maps, err
}

//GetListAnd 获取列表，自动束缚为AND连接
// 该方法主要为节约代码，上述getList获取数据后，还需要做多种判断，本方法将这些判断整合在一起
/** 正常语句应该做如下事情，本方法将尽可能完成以下内容，自动完成AND的处理
var newWhere string
newWhere, maps, err = CoreSQLFrom.JSONB.GetList("payment_create", "payment_create", args.PaymentCreate, maps)
if err != nil {
	return
} else {
	if newWhere != "" {
		where = where + " AND " + newWhere
	}
}
*/
func (t *FieldsFrom) GetListAnd(fieldName string, fieldMark string, query string, maps map[string]interface{}) (newQuery string, newMap map[string]interface{}, err error) {
	var newWhere string
	newWhere, maps, err = t.GetList(fieldName, fieldMark, maps)
	if err != nil {
		return
	} else {
		if newWhere != "" {
			if query == "" {
				newQuery = newWhere
			} else {
				newQuery = query + " AND " + newWhere
			}
			return newQuery, maps, err
		}
	}
	return query, maps, err
}

//FieldsFromOnlyID 简约设计，方便检索
type FieldsFromOnlyID struct {
	System string `db:"system" json:"system" check:"system" empty:"true"`
	ID     int64  `db:"id" json:"id" check:"id" empty:"true"`
}

//Value sql底层处理器
func (t FieldsFromOnlyID) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsFromOnlyID) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
