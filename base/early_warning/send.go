package BaseEarlyWarning

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// 获取目标人名下所有消息
type ArgsGetMessage struct {
	//分页
	Pages CoreSQLPages.ArgsDataList
	//用户ID
	UserID int64
	//送达ID
	ToID int64
	//模版ID
	TemplateID int64
	//是否需要已读参数
	NeedIsRead bool
	//是否已读
	IsRead bool
	//搜索
	Search string
}

func GetMessage(args *ArgsGetMessage) (dataList []FieldsWaitType, dataCount int64, err error) {
	where := "(content ILIKE '%')"
	maps := map[string]interface{}{
		"search": args.Search,
	}
	if args.UserID > 0 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.ToID > 0 {
		where = where + " AND to_id = :to_id"
		maps["user_id"] = args.ToID
	}
	if args.TemplateID > 0 {
		where = where + " AND template_id = :template_id"
		maps["template_id"] = args.TemplateID
	}
	if args.NeedIsRead {
		where = where + " AND is_read = :is_read"
		maps["is_read"] = args.IsRead
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_ew_wait",
		"id",
		"SELECT id, create_at, update_at, bind_id, level, to_id, is_send, is_read, expire_finish, expire_at, template_id, content, bind_data, need_phone, need_sms, need_email, need_app FROM core_ew_wait WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "level", "expire_at"},
	)
	return
}

// 给目标告警类目发送一个信息
// 关系是：任意模块可以使用该方法，该方法将抽取该类目下的目标人，按照预定的方案通知
type ArgsSendMod struct {
	//模版标识码
	Mark string
	//消息内容组
	Contents map[string]string
}

func SendMod(args *ArgsSendMod) (err error) {
	//通过mark获取模版信息
	var templateData FieldsTemplateType
	templateData, err = GetTemplateByMark(&ArgsGetTemplateByMark{
		Mark: args.Mark,
	})
	if err != nil {
		err = errors.New("cannot find template mark: " + args.Mark + ", " + err.Error())
		return
	}
	//获取该模版的最高级别
	var bindList []FieldsBindType
	var b bool
	bindList, b = getBindByLevel(&argsGetBindByLevel{
		TemplateID: templateData.ID,
		Level:      -1,
	})
	if !b {
		err = errors.New("cannot find bind data by template id or level 0")
		return
	}
	//遍历关系，发送数据
	for _, v := range bindList {
		err = send(v.ToID, v.TemplateID, args.Contents)
		if err != nil {
			err = errors.New(fmt.Sprint("cannot send message by template id: ", v.TemplateID, ", to id: ", v.ToID, ", err: ", err))
			return
		}
	}
	//反馈
	return
}

// 标记某个ID为已读
type ArgsUpdateSendIsRead struct {
	//ID
	ID int64 `db:"id"`
	//送达人ID
	// 可留空，用于检查
	ToID int64 `db:"to_id"`
}

func UpdateSendIsRead(args *ArgsUpdateSendIsRead) (err error) {
	where := "id = :id"
	if args.ToID > 0 {
		where = where + " AND to_id = :to_id"
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_ew_wait SET  update_at = NOW(), is_read = true WHERE "+where, args)
	return
}

// 标记某个mark全部为已读
type ArgsUpdateSendIsReadByMark struct {
	//标识码
	Mark string `db:"mark"`
}

func UpdateSendIsReadByMark(args *ArgsUpdateSendIsReadByMark) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_ew_wait w SET update_at = NOW(), is_read = true, expire_finish = false FROM core_ew_template t WHERE w.template_id = t.id AND t.mark = :mark", args)
	return
}

// 删除通知
type ArgsDeleteSendByID struct {
	//ID
	ID int64 `db:"id"`
}

func DeleteSendByID(args *ArgsDeleteSendByID) (err error) {
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_ew_wait", "id", args)
	if err != nil {
		err = errors.New("cannot delete wait by id, " + err.Error())
	}
	return
}

type ArgsDeleteSendByMark struct {
	//标识码
	Mark string `db:"mark"`
}

func DeleteSendByMark(args *ArgsDeleteSendByMark) (err error) {
	var templateData FieldsTemplateType
	templateData, err = GetTemplateByMark(&ArgsGetTemplateByMark{
		Mark: args.Mark,
	})
	if err != nil {
		err = errors.New("cannot find template data, " + err.Error())
		return
	}
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_ew_wait", "template_id = :template_id", map[string]interface{}{
		"template_id": templateData.ID,
	})
	if err != nil {
		err = errors.New("cannot delete wait by template id, " + err.Error())
	}
	return
}

// 给目标人发送一个消息
// 同一个模版、同一个人的未读、未过期消息将被强制过期处理
// params toID string 通知目标ID
// params templateID 消息模版，不同模版的级别、样式、通知方式都不同
// params contents 消息内容关系
func send(toID int64, templateID int64, contents map[string]string) (err error) {
	//标记此人旧的模版ID数据过期
	_, _ = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_ew_wait SET  expire_finish = true WHERE to_id = :to_id AND template_id = :template_id AND is_read = false AND expire_finish = false", map[string]interface{}{
		"to_id":       toID,
		"template_id": templateID,
	})
	//获取模版数据
	var templateData FieldsTemplateType
	templateData, err = GetTemplateByID(&ArgsGetTemplateByID{
		ID: templateID,
	})
	if err != nil {
		err = errors.New("cannot find template data, " + err.Error())
		return
	}
	//送达人数据
	var bindData FieldsBindType
	bindData, err = getBindByToAndTemplate(&argsGetBindByToAndTemplate{
		ToID: toID, TemplateID: templateID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("cannot find bind data by template id: ", templateID, ", to id: ", toID, ", err, ", err))
		return
	}
	//获取最短的过期时间
	var expireAt, expireBind, expireTemplate time.Time
	expireBind, err = CoreFilter.GetTimeByAdd(bindData.NextWaitTime)
	if err != nil {
		expireBind = CoreFilter.GetNowTime()
	}
	expireTemplate, err = CoreFilter.GetTimeByAdd(templateData.DefaultExpireTime)
	if err != nil {
		expireTemplate = CoreFilter.GetNowTime()
	}
	if expireBind.Unix() > expireTemplate.Unix() {
		expireAt = expireTemplate
	} else {
		expireAt = expireBind
	}
	//构建新的数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO core_ew_wait (bind_id, level, to_id, is_send, is_read, expire_finish, expire_at, template_id, content, bind_data, need_phone, need_sms, need_email, need_app) VALUES (:bind_id,:level,:to_id,false,false,false,:expire_at,:template_id,:content,:bind_data,:need_phone,:need_sms,:need_email,:need_app)", map[string]interface{}{
		"bind_id":     bindData.ID,
		"level":       bindData.Level,
		"to_id":       bindData.ToID,
		"expire_at":   expireAt,
		"template_id": templateData.ID,
		"content":     templateData.Content,
		"bind_data":   FieldsWaitBindData(contents),
		"need_phone":  bindData.NeedPhone,
		"need_sms":    bindData.NeedSMS,
		"need_email":  bindData.NeedEmail,
		"need_app":    bindData.NeedAPP,
	})
	if err != nil {
		err = errors.New("cannot create new data, " + err.Error())
		return
	}
	//反馈
	return
}

// 处理过期数据，并自动确保通知到下一个人
func clearSendExpire() (err error) {
	var data FieldsWaitType
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, bind_id, level, to_id, is_send, is_read, expire_finish, expire_at, template_id, content, bind_data, need_phone, need_sms, need_email, need_app FROM core_ew_wait WHERE  is_read = false AND expire_finish = false AND expire_at < NOW() ORDER BY id LIMIT 1")
	if err != nil {
		err = nil
		return
	}
	err = updateSendExpiredAndNext(&data)
	if err != nil {
		return
	}
	return
}

// 标记指定ID过期，并发送给下一个通知人
func updateSendExpiredAndNext(waitData *FieldsWaitType) (err error) {
	//作废当前消息
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_ew_wait SET  update_at = NOW(), expire_at = NOW(), expire_finish = true WHERE id = :id", map[string]interface{}{
		"id": waitData.ID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("update old msg, id: ", waitData.ID, ", err: ", err))
		return
	}
	//获取模版信息
	var templateData FieldsTemplateType
	templateData, err = GetTemplateByID(&ArgsGetTemplateByID{
		ID: waitData.TemplateID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("cannot find template cache by id: ", waitData.ID, ", err: "+err.Error()))
		return
	}
	//获取关系
	var bindData FieldsBindType
	bindData, err = GetBindID(&ArgsGetBindID{
		ID: waitData.BindID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("cannot find bind cache by id: ", waitData.BindID))
		return
	}
	//检查该绑定关系是否为最后一个级别？
	if bindData.Level < 1 {
		return
	}
	//获取当前级别所有消息
	switch bindData.LevelMode {
	case "or":
		//或关系，只有有一个已读，则返回
		var waitDataList []FieldsWaitType
		err = Router2SystemConfig.MainDB.Select(&waitDataList, "SELECT id, create_at, update_at, bind_id, level, to_id, is_send, is_read, expire_finish, expire_at, template_id, content, bind_data, need_phone, need_sms, need_email, need_app FROM core_ew_wait WHERE template_id = $1 AND level = :level", map[string]interface{}{
			"template_id": templateData.ID,
			"level":       bindData.Level,
		})
		if err != nil {
			//没有则跳过
			break
		}
		if len(waitDataList) < 1 {
			break
		}
		for _, v := range waitDataList {
			if v.IsRead {
				return
			}
		}
	case "and":
		//联合同意，只要有一个还未查看，则返回
		var waitDataList []FieldsWaitType
		err = Router2SystemConfig.MainDB.Select(&waitDataList, "SELECT id, create_at, update_at, bind_id, level, to_id, is_send, is_read, expire_finish, expire_at, template_id, content, bind_data, need_phone, need_sms, need_email, need_app FROM core_ew_wait WHERE template_id = $1 AND level = :level", map[string]interface{}{
			"template_id": templateData.ID,
			"level":       bindData.Level,
		})
		if err != nil {
			//没有则跳过
			break
		}
		if len(waitDataList) < 1 {
			break
		}
		for _, v := range waitDataList {
			if !v.IsRead && !v.ExpireFinish {
				return
			}
		}
	case "none":
		//直接跳过处理
	default:
		//直接跳过处理，按照none处理
	}
	//到达下一级进行通知
	var nextBindList []FieldsBindType
	var b bool
	nextBindList, b = getBindByLevel(&argsGetBindByLevel{
		TemplateID: templateData.ID,
		Level:      bindData.Level,
	})
	if b {
		//开始通知
		for _, v := range nextBindList {
			err = send(v.ToID, v.TemplateID, waitData.BindData)
			if err != nil {
				CoreLog.Error(fmt.Sprint("send next bind failed, to id: ", v.ToID, ", template id: ", v.TemplateID, ", err: ", err))
			}
		}
	}
	//反馈
	return
}

// 检查预警信息是否已读、已经过期
func checkWaitIsRead(data *FieldsWaitType) bool {
	if data.IsRead || data.ExpireFinish {
		return true
	}
	return false
}

// 检查某个模块是否存在未读预警
// 如果存在，就同时反馈该预警信息
func getLastSendMod(mark string) (data FieldsWaitType, err error) {
	var templateData FieldsTemplateType
	templateData, err = GetTemplateByMark(&ArgsGetTemplateByMark{
		Mark: mark,
	})
	if err != nil {
		return FieldsWaitType{}, errors.New("cannot find template data, " + err.Error())
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, bind_id, level, to_id, is_send, is_read, expire_finish, expire_at, template_id, content, bind_data, need_phone, need_sms, need_email, need_app FROM core_ew_wait WHERE template_id = $1 AND is_read = false AND expire_finish = false ORDER BY id LIMIT 1", templateData.ID)
	if err != nil {
		return
	}
	return
}

// 检查某人是否存在未读、未过期信息？
func getLastSendIsWaitByTemplate(toID int64, mark string) (data FieldsWaitType, err error) {
	var templateData FieldsTemplateType
	templateData, err = GetTemplateByMark(&ArgsGetTemplateByMark{
		Mark: mark,
	})
	if err != nil {
		err = errors.New("cannot find template data, " + err.Error())
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, bind_id, level, to_id, is_send, is_read, expire_finish, expire_at, template_id, content, bind_data, need_phone, need_sms, need_email, need_app FROM core_ew_wait WHERE template_id = $1 AND to_id = $2 AND is_read = false AND expire_finish = false ORDER BY id LIMIT 1", templateData.ID, toID)
	if err != nil {
		return
	}
	return
}

// 检查预警是否未读，且超出某个时间
// 过期时间为秒
// 可用于检查是否还需要发送预警信息
func checkWaitIsReadAndExpire(mark string, expireTime int64) bool {
	data, err := getLastSendMod(mark)
	if err != nil {
		return true
	}
	if !data.IsRead && data.CreateAt.Unix()+expireTime < CoreFilter.GetNowTime().Unix() {
		return true
	}
	return false
}

// 检查并发送
func sendBeforeCheck(mark string, contents map[string]string, expireTime int64) error {
	if b := checkWaitIsReadAndExpire(mark, expireTime); !b {
		return errors.New("check not pass")
	}
	return SendMod(&ArgsSendMod{
		Mark: mark, Contents: contents,
	})
}
