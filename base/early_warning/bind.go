package BaseEarlyWarning

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 通过关联人找到关系结构
type ArgsGetBindID struct {
	//绑定ID
	ID int64
}

func GetBindID(args *ArgsGetBindID) (data FieldsBindType, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, level, level_mode, next_wait_time, to_id, template_id, need_phone, need_sms, need_email, need_app FROM core_ew_bind WHERE id = $1", args.ID)
	return
}

// 通过关联人找到关系结构
type ArgsGetBindByToID struct {
	//送达人ID
	ToID int64
}

func GetBindByToID(args *ArgsGetBindByToID) (data []FieldsBindType, err error) {
	err = Router2SystemConfig.MainDB.Select(&data, "SELECT id, create_at, update_at, level, level_mode, next_wait_time, to_id, template_id, need_phone, need_sms, need_email, need_app FROM core_ew_bind WHERE to_id = $1", args.ToID)
	return
}

// 通过模版找到关系结构
type ArgsGetBindByTemplateID struct {
	//模版ID
	TemplateID int64
}

func GetBindByTemplateID(args *ArgsGetBindByTemplateID) (data []FieldsBindType, err error) {
	err = Router2SystemConfig.MainDB.Select(&data, "SELECT id, create_at, update_at, level, level_mode, next_wait_time, to_id, template_id, need_phone, need_sms, need_email, need_app FROM core_ew_bind WHERE template_id = $1", args.TemplateID)
	return
}

// 添加或设定关联
type ArgsSetBind struct {
	//送达ID
	ToID int64 `db:"to_id"`
	//模版ID
	TemplateID int64 `db:"template_id"`
	//级别
	Level int `db:"level"`
	//级别模式
	LevelMode string `db:"level_mode"`
	//下一个等待时间
	NextWaitTime string `db:"next_wait_time"`
	//通知方式
	NeedPhone bool `db:"need_phone"`
	NeedSMS   bool `db:"need_sms"`
	NeedEmail bool `db:"need_email"`
	NeedAPP   bool `db:"need_app"`
}

func SetBind(args *ArgsSetBind) (data FieldsBindType, err error) {
	//查询是否存在关联
	data, err = getBindByToAndTemplate(&argsGetBindByToAndTemplate{
		ToID:       args.ToID,
		TemplateID: args.TemplateID,
	})
	if err != nil {
		var lastID int64
		lastID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO core_ew_bind (level, level_mode, next_wait_time, to_id, template_id, need_phone, need_sms, need_email, need_app) VALUES (:level,:level_mode,:next_wait_time,:to_id,:template_id,:need_phone,:need_sms,:need_email,:need_app)", args)
		if err == nil {
			data, err = GetBindID(&ArgsGetBindID{
				ID: lastID,
			})
		}
	} else {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_ew_bind SET update_at = NOW(), level = :level, level_mode = :levelMode, next_wait_time = :next_wait_time, to_id = :to_id, template_id = :template_id, need_phone = :need_phone, need_sms = :need_sms, need_email = :need_email, need_app = :need_app WHERE id = :id", map[string]interface{}{
			"id":             data.ID,
			"to_id":          args.ToID,
			"template_id":    args.TemplateID,
			"level":          args.Level,
			"level_mode":     args.LevelMode,
			"next_wait_time": args.NextWaitTime,
			"need_phone":     args.NeedPhone,
			"need_sms":       args.NeedSMS,
			"need_email":     args.NeedEmail,
			"need_app":       args.NeedAPP,
		})
		if err == nil {
			data, err = GetBindID(&ArgsGetBindID{
				ID: args.TemplateID,
			})
		}
	}
	if err != nil {
		return
	}
	//检查模版下所有关系人，同一个级别的信息，强制修改为一致的条件模式
	_, _ = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_ew_bind SET update_at = NOW(), level_mode = :level_mode, next_wait_time = :next_wait_time, to_id = :to_id, need_phone = :need_phone, need_sms = :need_sms, need_email = :need_email, need_app = :need_app WHERE template_id = :template_id AND level = :level", map[string]interface{}{
		"id":             data.ID,
		"to_id":          args.ToID,
		"template_id":    args.TemplateID,
		"level":          args.Level,
		"level_mode":     args.LevelMode,
		"next_wait_time": args.NextWaitTime,
		"need_phone":     args.NeedPhone,
		"need_sms":       args.NeedSMS,
		"need_email":     args.NeedEmail,
		"need_app":       args.NeedAPP,
	})
	return
}

// 解除关系
type ArgsSetUnBind struct {
	//送达人ID
	// 可以留空
	ToID int64 `db:"to_id"`
	//模版ID
	// 可以留空
	TemplateID int64 `db:"template_id"`
}

func SetUnBind(args *ArgsSetUnBind) (err error) {
	if args.ToID > 0 {
		if args.TemplateID > 0 {
			_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_ew_bind", "to_id = :to_id AND template_id = :template_id", args)
		} else {
			_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_ew_bind", "to_id = :to_id", args)
		}
	} else {
		_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_ew_bind", "template_id = :template_id", args)
	}
	return
}

// 检查关系人和模版，是否存在关联
type argsCheckBindAndTemplate struct {
	//送达人ID
	ToID int64 `db:"to_id"`
	//模版ID
	TemplateID int64 `db:"template_id"`
}

func checkBindAndTemplate(args *argsCheckBindAndTemplate) bool {
	count, err := CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "core_ew_bind", "id", "to_id = :to_id AND template_id = :template_id", args)
	if err != nil {
		count = 0
	}
	return count > 0
}

// 获取某人的模版关系
type argsGetBindByToAndTemplate struct {
	//送达人ID
	ToID int64
	//模版ID
	TemplateID int64
}

func getBindByToAndTemplate(args *argsGetBindByToAndTemplate) (data FieldsBindType, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, level, level_mode, next_wait_time, to_id, template_id, need_phone, need_sms, need_email, need_app FROM core_ew_bind WHERE to_id = $1 AND template_id = $2", args.ToID, args.TemplateID)
	return
}

// 找到模版的，关系人和模版最优先关系
// 根据要求的最大级别以下查询，例如查询1，则不包含1的向下查询
// return []FieldsBindType 绑定关系序列
// return bool 是否存在数据
type argsGetBindByLevel struct {
	//模版ID
	TemplateID int64
	//级别
	// 如果为-1，则自动走模版绑定的最大级别
	Level int
}

func getBindByLevel(args *argsGetBindByLevel) (dataList []FieldsBindType, b bool) {
	//抽取该模版的所有关系人
	var findAllData []FieldsBindType
	var err error
	err = Router2SystemConfig.MainDB.Select(&findAllData, "SELECT * FROM core_ew_bind WHERE level > $1 AND template_id = $2", args.Level, args.TemplateID)
	if err != nil {
		return
	}
	//如果level < 1，则找出该绑定关系最大值
	if args.Level < 0 {
		for _, v := range findAllData {
			if v.Level > args.Level {
				args.Level = v.Level
			}
		}
	}
	//遍历数据，抽取出最大level的级别作为最后选择
	for _, v := range findAllData {
		//如果级别小于0，则无效退出
		if args.Level < 0 {
			return
		}
		//如果级别相同，则加入列队
		if v.Level == args.Level {
			dataList = append(dataList, v)
			continue
		}
		//递减
		args.Level = args.Level - 1
	}
	b = len(dataList) > 0
	//反馈
	return
}
