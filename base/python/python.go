package BasePython

import (
	"errors"
	"fmt"
	BaseExpireTip "github.com/fotomxq/weeekj_core/v5/base/expire_tip"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// PushSync 推送一个同步的数据请求
// timeout 超时时间，单位秒
func PushSync(system string, bindID int64, mark string, param []byte, fileSrc string, timeout int) (resultFileSrc string, err error) {
	//请求数据
	var newID int64
	newID, err = Push(system, bindID, mark, param, fileSrc, timeout)
	if err != nil {
		return
	}
	//循环等待
	nowSec := 0
	for {
		//超时反馈错误
		if nowSec >= timeout {
			err = errors.New("time out")
			return
		}
		//检查是否完成
		if b := checkResult(newID); b {
			resultFileSrc, err = getResult(newID)
			if err != nil {
				return
			}
			break
		}
		//递增
		nowSec += 1
		time.Sleep(time.Second * 1)
	}
	//反馈成功
	return
}

// PushCall 回调形式处理数据请求
func PushCall(system string, bindID int64, mark string, param []byte, timeout int, fileSrc string, handle func(resultFileSrc string)) (err error) {
	//请求数据
	var newID int64
	newID, err = Push(system, bindID, mark, param, fileSrc, timeout)
	if err != nil {
		return
	}
	//循环等待
	nowSec := 0
	for {
		//超时反馈错误
		if nowSec >= timeout {
			err = errors.New("time out")
			return
		}
		//检查是否完成
		if b := checkResult(newID); b {
			var resultFileSrc string
			resultFileSrc, err = getResult(newID)
			if err != nil {
				return
			}
			handle(resultFileSrc)
			break
		}
		//递增
		nowSec += 1
		time.Sleep(time.Second * 1)
	}
	//反馈成功
	return
}

// Push 推送一个数据请求
func Push(system string, bindID int64, mark string, param []byte, fileSrc string, timeout int) (newID int64, err error) {
	//检查是否存在相同的请求
	var data fieldsWait
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM core_python WHERE system = $1 AND bind_id = $2 AND mark = $3 AND expire_at >= NOW()", system, bindID, mark)
	if err == nil && data.ID > 0 {
		newID = data.ID
		return
	}
	//创建新的数据
	newID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO core_python (expire_at, is_finish, system, bind_id, mark, param) VALUES (:expire_at,false,:system,:bind_id,:mark,:param)", map[string]interface{}{
		"expire_at": CoreFilter.GetNowTimeCarbon().AddSeconds(timeout).Time,
		"system":    system,
		"bind_id":   bindID,
		"mark":      mark,
		"param":     param,
	})
	if err != nil || newID < 1 {
		err = errors.New("insert data")
		return
	}
	//过期提醒
	BaseExpireTip.AppendTipNoErr(&BaseExpireTip.ArgsAppendTip{
		OrgID:      0,
		UserID:     0,
		SystemMark: "core_python",
		BindID:     newID,
		Hash:       "",
		ExpireAt:   CoreFilter.GetNowTimeCarbon().AddSeconds(timeout).Time,
	})
	//推送nats请求
	CoreNats.PushDataNoErr("base_python_new", "/base/python/new", system, newID, mark, map[string]interface{}{
		"bindID":     bindID,
		"fileSrc":    fileSrc,
		"newFileSrc": getResultSrc(newID),
		"param":      param,
	})
	//反馈
	return
}

// GetAndCheck 检查并获取结果
func GetAndCheck(id int64) (result string, err error) {
	//检查结果
	if b := checkResult(id); !b {
		err = errors.New("no finish")
		return
	}
	//获取结果
	result, err = getResult(id)
	if err != nil {
		return
	}
	//反馈
	return
}

func GetAndCheckByFrom(system string, bindID int64, mark string) (resultFileSrc string, err error) {
	var data fieldsWait
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM core_python WHERE system = $1 AND bind_id = $2 AND mark = $3 AND expire_at >= NOW()", system, bindID, mark)
	if err != nil || data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return GetAndCheck(data.ID)
}

// 更新结果
func updateResult(id int64) {
	//更新结果
	_, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_python SET is_finish = true WHERE id = :id", map[string]interface{}{
		"id": id,
	})
	if err != nil {
		CoreLog.Error("core python update result, id: ", id, ", err: ", err)
		return
	}
	//删除缓冲
	deleteCache(id)
}

// 删除数据
func deleteByID(id int64) (err error) {
	data := getByID(id)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	fileSrc := getResultSrc(id)
	if CoreFile.IsFile(fileSrc) {
		err = CoreFile.DeleteF(fileSrc)
		if err != nil {
			return
		}
	}
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_python", "id", map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return
	}
	deleteCache(id)
	return
}

// 检查结果
// 需要在过期之前转移走，否则文件会被移除
func checkResult(id int64) bool {
	cacheMark := getCacheMark(id)
	var data fieldsWait
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return data.IsFinish
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, is_finish FROM core_python WHERE id = $1", id)
	if data.ID < 1 {
		return false
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 86400)
	return data.IsFinish
}

// 获取指定ID数据
func getByID(id int64) (data fieldsWait) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, is_finish, system, bind_id, mark, param FROM core_python WHERE id = $1", id)
	if err != nil {
		return
	}
	return
}

// getResult 获取处理结果
func getResult(id int64) (result string, err error) {
	data := getByID(id)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	if !data.IsFinish {
		err = errors.New("no finish")
		return
	}
	fileSrc := getResultSrc(id)
	if !CoreFile.IsFile(fileSrc) {
		err = errors.New("no file")
		return
	}
	result = fileSrc
	return
}

// 根据数据获取文件路径
func getResultSrc(id int64) string {
	data := getByID(id)
	if data.ID < 1 {
		return ""
	}
	dataDirSrc := fmt.Sprint(dirSrc, data.CreateAt.Format("200601"), CoreFile.Sep, data.CreateAt.Format("02"))
	if !CoreFile.IsFolder(dataDirSrc) {
		_ = CoreFile.CreateFolder(dataDirSrc)
	}
	return fmt.Sprint(dataDirSrc, CoreFile.Sep, data.ID)
}

// 缓冲
func getCacheMark(id int64) string {
	return fmt.Sprint("core:python:id:", id)
}

func deleteCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getCacheMark(id))
}
