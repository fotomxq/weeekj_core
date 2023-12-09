package UserLogin

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetSave 提取数据
type ArgsGetSave struct {
	//密钥
	Key string `db:"key" json:"key"`
}

func GetSave(args *ArgsGetSave) (data FieldsSaveReportData, err error) {
	var dataResult FieldsSave
	err = Router2SystemConfig.MainDB.Get(&dataResult, "SELECT id, data FROM user_login_save WHERE key = $1 AND expire_at >= NOW()", args.Key)
	if err != nil {
		return
	}
	if dataResult.ID < 1 {
		err = errors.New("not exist")
		return
	}
	data = dataResult.Data
	_, _ = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "user_login_save", "id", map[string]interface{}{
		"id": dataResult.ID,
	})
	return
}

// 添加数据到集合内
func appendSave(data FieldsSaveReportData) (key string, err error) {
	key, err = CoreFilter.GetRandStr3(30)
	if err != nil {
		return
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_login_save (expire_at, key, data) VALUES (:expire_at,:key,:data)", map[string]interface{}{
		"expire_at": CoreFilter.GetNowTimeCarbon().AddSeconds(30).Time,
		"key":       key,
		"data":      data,
	})
	return
}
