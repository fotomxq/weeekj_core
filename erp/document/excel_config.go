package ERPDocument

import (
	"fmt"
	CoreCache "gitee.com/weeekj/weeekj_core/v5/core/cache"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// GetExcelConfigByConfigID 获取指定配置
func GetExcelConfigByConfigID(configID int64) (data FieldsExcelConfig) {
	//获取数据
	cacheMark := getExcelConfigCache(configID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, config_id, sheets FROM erp_document_excel_config WHERE config_id = $1", configID)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime3Day)
	//反馈
	return
}

// ArgsSetExcelConfig 设置excel配置参数
type ArgsSetExcelConfig struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//数据表
	Sheets FieldsExcelConfigSheetList `db:"sheets" json:"sheets"`
}

// SetExcelConfig 设置excel配置
func SetExcelConfig(args *ArgsSetExcelConfig) (err error) {
	data := GetExcelConfigByConfigID(args.ConfigID)
	if data.ID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE erp_document_excel_config SET sheets = :sheets WHERE id = :id", map[string]interface{}{
			"id":     data.ID,
			"sheets": args.Sheets,
		})
		if err != nil {
			return
		}
		deleteExcelConfigCache(data.ConfigID)
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO erp_document_excel_config (config_id, sheets) VALUES (:config_id, :sheets)", map[string]interface{}{
			"config_id": args.ConfigID,
			"sheets":    args.Sheets,
		})
		if err != nil {
			return
		}
	}
	return
}

func deleteExcelConfig(configID int64) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "erp_document_excel_config", "config_id = :config_id", map[string]interface{}{
		"config_id": configID,
	})
	if err != nil {
		return
	}
	deleteExcelConfigCache(configID)
	return
}

func getExcelConfigCache(configID int64) string {
	return fmt.Sprint("erp:document:excel:config:config:id:", configID)
}

func deleteExcelConfigCache(configID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getExcelConfigCache(configID))
}
