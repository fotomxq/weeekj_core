package ERPCore

import (
	"errors"
	"fmt"
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

type ComponentVal struct {
	//缓冲名称
	CacheName string
	//表名称
	TableName string
}

// GetAllVal 获取绑定的所有内容
func (t *ComponentVal) GetAllVal(bindID int64) (dataList []FieldsComponentVal) {
	cacheMark := t.getCacheBindMark(bindID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &dataList); err == nil && len(dataList) > 0 {
		return
	}
	var rawList []FieldsComponentVal
	err := Router2SystemConfig.MainDB.Select(&rawList, fmt.Sprint("SELECT id, bind_id ", "FROM "+t.TableName+" WHERE bind_id = $1 ORDER BY sort"), bindID)
	if err != nil || len(rawList) < 1 {
		return
	}
	for k := 0; k < len(rawList); k++ {
		dataList = append(dataList, t.getValByID(rawList[k].BindID, rawList[k].ID))
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, dataList, CoreCache.CacheTime1Day)
	return
}

// ArgsSetMore 批量设置内容参数
type ArgsSetMore struct {
	//所属
	BindID int64 `json:"bindID"`
	//内容
	DataList FieldsComponentDefineList `json:"dataList"`
}

// SetMore 批量设置内容
func (t *ComponentVal) SetMore(args *ArgsSetMore) (err error) {
	for _, v := range args.DataList {
		err = t.setVal(args.BindID, &v)
		if err != nil {
			err = errors.New(fmt.Sprint("set component val, key: ", v.Key, ", err: ", err))
			return
		}
	}
	return
}

func (t *ComponentVal) setVal(bindID int64, args *FieldsComponentDefine) (err error) {
	//识别类型分析数字值
	var valInt64 int64
	var valFloat64 float64
	valInt64, _ = CoreFilter.GetInt64ByString(args.Val)
	valFloat64, _ = CoreFilter.GetFloat64ByString(args.Val)
	//获取数据
	var data FieldsComponentVal
	err = Router2SystemConfig.MainDB.Get(&data, fmt.Sprint("SELECT id, bind_id ", "FROM ", t.TableName, " WHERE bind_id = $1 AND key = $2"), bindID, args.Key)
	if err == nil {
		//如果存在数据则编辑
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, fmt.Sprint("UPDATE ", t.TableName, " SET val = :val, val_int64 = :val_int64, val_float64 = :val_float64, params = :params WHERE id = :id"), map[string]interface{}{
			"id":          data.ID,
			"val":         args.Val,
			"val_int64":   valInt64,
			"val_float64": valFloat64,
			"params":      args.Params,
		})
		if err != nil {
			return
		}
		//删除缓冲
		t.deleteCache(data.BindID, data.ID)
	} else {
		//如果不存在则创建
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, fmt.Sprint("INSERT ", "INTO ", t.TableName, " (bind_id, key, sort, component_type, name, help_des, val, val_int64, val_float64, check_val, is_require, params) VALUES (:bind_id,:key,:sort,:component_type,:name,:help_des,:val,:val_int64,:val_float64,:check_val,:is_require,:params)"), map[string]interface{}{
			"bind_id":        bindID,
			"key":            args.Key,
			"sort":           args.Sort,
			"component_type": args.ComponentType,
			"name":           args.Name,
			"help_des":       args.HelpDes,
			"val":            args.Val,
			"val_int64":      valInt64,
			"val_float64":    valFloat64,
			"check_val":      args.CheckVal,
			"is_require":     args.IsRequire,
			"params":         args.Params,
		})
		if err != nil {
			return
		}
	}
	//反馈
	return
}

// ArgsComponentValMoreSetOnlyUpdate 批量仅编辑操作参数
type ArgsComponentValMoreSetOnlyUpdate struct {
	//所属
	BindID int64 `json:"bindID"`
	//内容
	DataList []ArgsComponentValSetOnlyUpdate `json:"dataList"`
}

// SetValMoreOnlyUpdate 批量仅编辑操作
func (t *ComponentVal) SetValMoreOnlyUpdate(args *ArgsComponentValMoreSetOnlyUpdate) (err error) {
	for _, v := range args.DataList {
		err = t.setValOnlyUpdate(args.BindID, &v)
		if err != nil {
			return
		}
	}
	return
}

// ArgsComponentValSetOnlyUpdate 仅编辑操作参数
type ArgsComponentValSetOnlyUpdate struct {
	//组件key
	// 单个节点内必须唯一
	Key string `db:"key" json:"key"`
	//组件默认值
	Val string `db:"val" json:"val"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// setValOnlyUpdate 仅编辑操作
func (t *ComponentVal) setValOnlyUpdate(bindID int64, args *ArgsComponentValSetOnlyUpdate) (err error) {
	//获取数据
	var data FieldsComponentVal
	err = Router2SystemConfig.MainDB.Get(&data, fmt.Sprint("SELECT id, bind_id ", "FROM ", t.TableName, " WHERE bind_id = $1 AND key = $2"), bindID, args.Key)
	//不存在数据则退出
	if err != nil || data.ID < 1 {
		err = errors.New(fmt.Sprint("no data, bind id: ", bindID, ", key: ", args.Key, ", err: ", err))
		return
	}
	//识别类型分析数字值
	var valInt64 int64
	var valFloat64 float64
	valInt64, _ = CoreFilter.GetInt64ByString(args.Val)
	valFloat64, _ = CoreFilter.GetFloat64ByString(args.Val)
	//写入扩展参数
	for _, v := range args.Params {
		data.Params = CoreSQLConfig.Set(data.Params, v.Mark, v.Val)
	}
	//更新数据
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, fmt.Sprint("UPDATE ", t.TableName, " SET val = :val, val_int64 = :val_int64, val_float64 = :val_float64, params = :params WHERE id = :id"), map[string]interface{}{
		"id":          data.ID,
		"val":         args.Val,
		"val_int64":   valInt64,
		"val_float64": valFloat64,
		"params":      data.Params,
	})
	if err != nil {
		return
	}
	//删除缓冲
	t.deleteCache(data.BindID, data.ID)
	//反馈
	return
}

// DeleteByBindID 删除绑定的所有内容
func (t *ComponentVal) DeleteByBindID(bindID int64) (err error) {
	//删除数据
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, t.TableName, "bind_id = :bind_id", map[string]interface{}{
		"bind_id": bindID,
	})
	if err != nil {
		return
	}
	//删除缓冲
	t.deleteCache(bindID, -1)
	//反馈
	return
}

// getValByID 获取内容
func (t *ComponentVal) getValByID(bindID, id int64) (data FieldsComponentVal) {
	cacheMark := t.getCacheMark(bindID, id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, fmt.Sprint("SELECT id, bind_id, key, sort, component_type, name, help_des, val, val_int64, val_float64, check_val, is_require, params ", "FROM ", t.TableName, " WHERE id = $1"), id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Day)
	return
}

func (t *ComponentVal) getCacheMark(bindID, id int64) string {
	return fmt.Sprint(t.CacheName, bindID, ".", id)
}

func (t *ComponentVal) getCacheBindMark(bindID int64) string {
	return fmt.Sprint(t.CacheName, bindID)
}

func (t *ComponentVal) deleteCache(bindID int64, id int64) {
	Router2SystemConfig.MainCache.DeleteSearchMark(t.getCacheBindMark(bindID))
	if id > 0 {
		//Router2SystemConfig.MainCache.DeleteMark(t.getCacheMark(bindID, id))
	}
}
