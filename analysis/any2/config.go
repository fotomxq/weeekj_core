package AnalysisAny2

import (
	"fmt"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 获取配置
func getConfigByMark(mark string, noCreate bool) (data FieldsConfig, err error) {
	cacheMark := getConfigCacheMark(mark)
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.Mark != "" {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, mark, particle, file_day FROM analysis_any2_config WHERE mark = $1", mark)
	if err != nil {
		if noCreate {
			return
		}
		//自动生成配置
		if err = setConfig(mark); err != nil {
			return
		}
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, mark, particle, file_day FROM analysis_any2_config WHERE mark = $1", mark)
		if err != nil {
			return
		}
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 2592000)
	return
}

// 生成配置
func setConfig(mark string) (err error) {
	cacheMark := getConfigCacheMark(mark)
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO analysis_any2_config (mark, particle, file_day) VALUES (:mark,:particle,:file_day)", map[string]interface{}{
		"mark":     mark,
		"particle": 3,
		"file_day": 30,
	})
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.DeleteMark(cacheMark)
	return
}

// SetConfigBefore 预先构建数据集
// particle 颗粒度参考: 0 小时（默认）/ 1 1天 / 2 1周 / 3 1月 / 4 1年
func SetConfigBefore(mark string, particle int, fileDay int) (err error) {
	data, _ := getConfigByMark(mark, true)
	if data.ID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE analysis_any2_config SET particle = :particle, file_day = :file_day WHERE id = :id", map[string]interface{}{
			"id":       data.ID,
			"particle": particle,
			"file_day": fileDay,
		})
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO analysis_any2_config (mark, particle, file_day) VALUES (:mark,:particle,:file_day)", map[string]interface{}{
			"mark":     mark,
			"particle": particle,
			"file_day": fileDay,
		})
	}
	if err != nil {
		return
	}
	cacheMark := getConfigCacheMark(mark)
	Router2SystemConfig.MainCache.DeleteMark(cacheMark)
	return
}

// SetConfigBeforeNoErr 预先构建数据集
// particle 颗粒度参考: 0 小时（默认）/ 1 1天 / 2 1周 / 3 1月 / 4 1年
func SetConfigBeforeNoErr(mark string, particle int, fileDay int) {
	if err := SetConfigBefore(mark, particle, fileDay); err != nil {
		CoreLog.Error("analysis any2 init config, ", err)
	}
}

// SetConfigFile 设置配置归档天数
func SetConfigFile(mark string, fileDay int) (err error) {
	if fileDay < 1 {
		fileDay = 3
	}
	data, _ := getConfigByMark(mark, true)
	if data.ID > 0 {
		if data.Mark == mark && data.FileDay == fileDay {
			return
		}
	}
	cacheMark := getConfigCacheMark(mark)
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE analysis_any2_config SET file_day = :file_day WHERE id = :id", map[string]interface{}{
		"id":       data.ID,
		"file_day": fileDay,
	})
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.DeleteMark(cacheMark)
	return
}

// SetConfigParticle 设置配置颗粒度
func SetConfigParticle(mark string, particle int) (err error) {
	switch particle {
	case 0:
	case 1:
	case 2:
	case 3:
	case 4:
	default:
		particle = 0
	}
	data, _ := getConfigByMark(mark, true)
	if data.ID > 0 {
		if data.Mark == mark && data.Particle == particle {
			return
		}
	}
	cacheMark := getConfigCacheMark(mark)
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE analysis_any2_config SET particle = :particle WHERE id = :id", map[string]interface{}{
		"id":       data.ID,
		"particle": particle,
	})
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.DeleteMark(cacheMark)
	return
}

// 获取配置标识码
func getConfigCacheMark(mark string) string {
	return fmt.Sprint("analysis:any2:mark:", mark)
}
