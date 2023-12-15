package FinanceDeposit

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	OrgUserMod "github.com/fotomxq/weeekj_core/v5/org/user/mod"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetList 获取储蓄数据列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom
	//存储来源
	FromInfo CoreSQLFrom.FieldsFrom
	//存储标识码
	ConfigMark string
	//最小金额
	// 0忽略
	MinPrice int64
	//最大金额
	// 0忽略
	MaxPrice int64
}

// GetList 获取储蓄数据列表
func GetList(args *ArgsGetList) (dataList []FieldsDepositType, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where, maps, err = args.CreateInfo.GetListAnd("create_info", "create_info", where, maps)
	if err != nil {
		return
	}
	where, maps, err = args.FromInfo.GetListAnd("from_info", "from_info", where, maps)
	if err != nil {
		return
	}
	if args.ConfigMark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "config_mark = :config_mark"
		maps["config_mark"] = args.ConfigMark
	}
	if args.MinPrice > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "save_price >= :min_price"
		maps["min_price"] = args.MinPrice
	}
	if args.MaxPrice > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "save_price <= :max_price"
		maps["max_price"] = args.MaxPrice
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"finance_deposit",
		"id",
		"SELECT id, create_at, update_at, update_hash, create_info, from_info, config_mark, save_price FROM finance_deposit WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "save_price"},
	)
	return
}

// ArgsGetByID 获取某个ID参数
type ArgsGetByID struct {
	//ID
	ID int64 `db:"id" json:"id"`
}

// GetByID 获取某个ID
func GetByID(args *ArgsGetByID) (data FieldsDepositType, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, update_hash, create_info, from_info, config_mark, save_price FROM finance_deposit WHERE id = $1", args.ID)
	return
}

// ArgsGetByFrom 获取来源参数
type ArgsGetByFrom struct {
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom
	//存储来源
	FromInfo CoreSQLFrom.FieldsFrom
	//存储标识码
	ConfigMark string
}

// GetByFrom 获取来源
func GetByFrom(args *ArgsGetByFrom) (data FieldsDepositType, err error) {
	var createInfo string
	createInfo, err = args.CreateInfo.GetRawNoName()
	if err != nil {
		return
	}
	var fromInfo string
	fromInfo, err = args.FromInfo.GetRawNoName()
	if err != nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, update_hash, create_info, from_info, config_mark, save_price FROM finance_deposit WHERE config_mark = $1 AND create_info @> $2 AND from_info @> $3", args.ConfigMark, createInfo, fromInfo)
	if err != nil {
		return
	}
	return
}

// GetPriceByFrom 检查来源金额是否足够？
func GetPriceByFrom(args *ArgsGetByFrom) (price int64) {
	var createInfo string
	var err error
	createInfo, err = args.CreateInfo.GetRawNoName()
	if err != nil {
		return
	}
	var fromInfo string
	fromInfo, err = args.FromInfo.GetRawNoName()
	if err != nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&price, "SELECT save_price FROM finance_deposit WHERE config_mark = $1 AND create_info @> $2 AND from_info @> $3", args.ConfigMark, createInfo, fromInfo)
	if err != nil {
		price = 0
	}
	return
}

// ArgsSetByFrom 修改或创建数据参数
type ArgsSetByFrom struct {
	//更新hash
	UpdateHash string
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom
	//存储来源
	FromInfo CoreSQLFrom.FieldsFrom
	//配置标识码
	ConfigMark string
	//增减金额
	AppendSavePrice int64
}

// SetByFrom 修改或创建数据
func SetByFrom(args *ArgsSetByFrom) (data FieldsDepositType, errCode string, err error) {
	//允许构建平台级别储蓄
	//if args.CreateInfo.System == "" || args.CreateInfo.ID < 1 {
	//	errCode = "create_empty"
	//	err = errors.New("create from not exist")
	//	return
	//}
	data, err = GetByFrom(&ArgsGetByFrom{
		CreateInfo: args.CreateInfo,
		FromInfo:   args.FromInfo,
		ConfigMark: args.ConfigMark,
	})
	if err != nil || data.ID < 1 {
		//创建新数据
		data, errCode, err = setByData(true, &FieldsDepositType{
			ID:         0,
			CreateAt:   time.Time{},
			UpdateAt:   time.Time{},
			UpdateHash: "",
			CreateInfo: args.CreateInfo,
			FromInfo:   args.FromInfo,
			ConfigMark: "",
			SavePrice:  0,
		}, args.ConfigMark, args.AppendSavePrice, "")
		if err != nil {
			return
		}
	} else {
		//修改旧数据
		if args.UpdateHash == "" {
			args.UpdateHash = data.UpdateHash
		}
		data, errCode, err = setByData(false, &data, args.ConfigMark, args.AppendSavePrice, args.UpdateHash)
		if err != nil {
			return
		}
	}
	//修改新储蓄
	//syncChangeNewSavings(&data, args.AppendSavePrice, "", false)
	//反馈
	return
}

// ArgsSetByID 修改必定存在的储蓄参数
type ArgsSetByID struct {
	//更新hash
	UpdateHash string
	//ID
	ID int64
	//配置标识码
	ConfigMark string
	//增减金额
	AppendSavePrice int64
}

// SetByID 修改必定存在的储蓄
func SetByID(args *ArgsSetByID) (data FieldsDepositType, errCode string, err error) {
	data, err = GetByID(&ArgsGetByID{
		ID: args.ID,
	})
	if err != nil {
		errCode = "not_exist"
		err = errors.New("cannot find data by id, " + err.Error())
		return
	} else {
		if data.UpdateHash != args.UpdateHash {
			errCode = "hash"
			err = errors.New("hash is error")
			return
		}
		return setByData(false, &data, args.ConfigMark, args.AppendSavePrice, args.UpdateHash)
	}
}

// setByData 设置储蓄
func setByData(isCreate bool, data *FieldsDepositType, configMark string, appendSavePrice int64, updateHash string) (newData FieldsDepositType, errCode string, err error) {
	//如果为0，则退出
	if !isCreate && appendSavePrice == 0 {
		errCode = "err_finance_add_0"
		err = errors.New("set deposit price 0")
		return
	}
	//获取配置
	var configData FieldsConfigType
	configData, err = GetConfigByMark(&ArgsGetConfigByMark{
		Mark: configMark,
	})
	if err != nil {
		errCode = "err_config"
		err = errors.New("get deposit config, mark: " + configMark + ", " + err.Error())
		return
	}
	//检查变动金额是否不足？
	if appendSavePrice > 0 {
		if configData.OnceSaveMinLimit > 0 && appendSavePrice <= configData.OnceSaveMinLimit {
			errCode = "err_finance_deposit_once_save_min_limit"
			err = errors.New("deposit config limit by once save min limit")
			return
		}
		if configData.OnceSaveMaxLimit > 0 && appendSavePrice >= configData.OnceSaveMaxLimit {
			errCode = "err_finance_deposit_once_save_max_limit"
			err = errors.New("deposit config limit by once save max limit")
			return
		}
	}
	if appendSavePrice < 0 {
		if configData.OnceTakeMinLimit > 0 && 0-appendSavePrice <= configData.OnceTakeMinLimit {
			errCode = "err_finance_deposit_once_take_min_limit"
			err = errors.New("deposit config limit by once take min limit")
			return
		}
		if configData.OnceTakeMaxLimit > 0 && 0-appendSavePrice >= configData.OnceTakeMaxLimit {
			errCode = "err_finance_deposit_once_take_max_limit"
			err = errors.New("deposit config limit by once take max limit")
			return
		}
	}
	//获取hash随机值
	var newUpdateHash string
	newUpdateHash, err = CoreFilter.GetRandStr3(10)
	if err != nil {
		errCode = "err_hash"
		err = errors.New("deposit rand hash, " + err.Error())
		return
	}
	if !isCreate {
		newPrice := data.SavePrice + appendSavePrice
		if newPrice < 0 {
			errCode = "err_finance_deposit_save_price_enough"
			err = errors.New("deposit save price not enough")
			return
		}
		if data.SavePrice < configData.TakeLimit {
			errCode = "err_finance_deposit_take_limit"
			err = errors.New("deposit save price not enough by take limit")
			return
		}
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE finance_deposit SET update_at = NOW(), update_hash = :new_update_hash, save_price = :save_price WHERE id = :id AND update_hash = :update_hash", map[string]interface{}{
			"new_update_hash": newUpdateHash,
			"save_price":      newPrice,
			"id":              data.ID,
			"update_hash":     updateHash,
		})
		if err != nil {
			errCode = "err_insert"
			err = errors.New("deposit update data by id, " + err.Error())
			return
		}
		newData, err = GetByID(&ArgsGetByID{
			ID: data.ID,
		})
		if err != nil {
			errCode = "err_no_data"
			return
		}
	} else {
		if appendSavePrice < 0 {
			errCode = "err_finance_deposit_save_price_enough"
			err = errors.New("deposit price less 1")
			return
		}
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "finance_deposit", "INSERT INTO finance_deposit (update_hash, create_info, from_info, config_mark, save_price) VALUES (:update_hash,:create_info,:from_info,:config_mark,:save_price)", map[string]interface{}{
			"update_hash": newUpdateHash,
			"create_info": data.CreateInfo,
			"from_info":   data.FromInfo,
			"config_mark": configMark,
			"save_price":  appendSavePrice,
		}, &newData)
		if err != nil {
			errCode = "err_insert"
			err = errors.New("deposit insert data, " + err.Error())
			return
		}
	}
	//第二代处理模块同步数据
	syncChangeNewSavings(data, appendSavePrice, "", false)
	//请求更新组织用户聚合数据
	if data.FromInfo.System == "org" && data.FromInfo.ID > 0 && data.CreateInfo.System == "user" && data.CreateInfo.ID > 0 {
		OrgUserMod.PushUpdateUserData(data.FromInfo.ID, data.CreateInfo.ID)
	}
	//反馈
	return
}
