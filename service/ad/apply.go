package ServiceAD

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsGetApplyList 获取请求列表参数
type ArgsGetApplyList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可以由用户发起，用户发起的广告主要向商户投放，由商户进行审核处理
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//是否审核
	NeedIsAudit bool `json:"needIsAudit" check:"bool"`
	IsAudit     bool `json:"isAudit" check:"bool"`
	//是否开始
	NeedIsStart bool `json:"needIsStart" check:"bool"`
	IsStart     bool `json:"isStart" check:"bool"`
	//是否结束
	NeedIsEnd bool `json:"needIsEnd" check:"bool"`
	IsEnd     bool `json:"isEnd" check:"bool"`
	//广告标识码
	Mark string `json:"mark" check:"mark" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetApplyList 获取请求列表
func GetApplyList(args *ArgsGetApplyList) (dataList []FieldsApply, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.NeedIsAudit {
		if args.IsAudit {
			where = where + " AND audit_at >= to_timestamp(1000000)"
		} else {
			where = where + " AND audit_at < to_timestamp(1000000)"
		}
	}
	if args.NeedIsStart {
		if args.IsStart {
			where = where + " AND start_at <= NOW()"
		} else {
			where = where + " AND start_at > NOW()"
		}
	}
	if args.NeedIsEnd {
		if args.IsEnd {
			where = where + " AND end_at <= NOW()"
		} else {
			where = where + " AND end_at > NOW()"
		}
	}
	if args.Mark != "" {
		where = where + " AND mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "service_ad_apply"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, start_at, end_at, audit_at, audit_des, audit_ban_des, org_id, user_id, area_ids, ad_id, mark, name, cover_file_id, count, click_count FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "start_at", "end_at", "audit_at"},
	)
	return
}

// ArgsGetApplyID 获取请求记录参数
type ArgsGetApplyID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 验证用
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可以由用户发起，用户发起的广告主要向商户投放，由商户进行审核处理
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// GetApplyID 获取请求记录
func GetApplyID(args *ArgsGetApplyID) (data FieldsApply, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, start_at, end_at, audit_at, audit_des, audit_ban_des, org_id, user_id, area_ids, ad_id, mark, name, des, cover_file_id, des_files, count, click_count, params FROM service_ad_apply WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND ($3 < 1 OR user_id = $3)", args.ID, args.OrgID, args.UserID)
	return
}

// ArgsCreateApply 创建请求参数
type ArgsCreateApply struct {
	//启动时间
	StartAt string `db:"start_at" json:"startAt" check:"isoTime"`
	//结束时间
	EndAt string `db:"end_at" json:"endAt" check:"isoTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//用户ID
	// 可以由用户发起，用户发起的广告主要向商户投放，由商户进行审核处理
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//投放分区ID列
	AreaIDs pq.Int64Array `db:"area_ids" json:"areaIDs" check:"ids" empty:"true"`
	//分区标识码
	// 作为前端抽取数据类型使用，可以重复指定多个
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300" empty:"true"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id"`
	//描述组图
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params" check:"params"`
}

// CreateApply 创建请求
// 申请的新记录将进入列队，在N秒内审核后绑定到对应的空间上
// 当过期后，将自动删除绑定关系，并删除对应的广告ID
func CreateApply(args *ArgsCreateApply) (data FieldsApply, err error) {
	//检查是否存在待审核请求
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM service_ad_apply WHERE org_id = $1 AND user_id = $2 AND delete_at < to_timestamp(1000000) AND audit_at < to_timestamp(1000000)", args.OrgID, args.UserID)
	if err == nil && id > 0 {
		err = errors.New("have same data")
		return
	}
	//构建时间
	var startAt, endAt time.Time
	startAt, err = CoreFilter.GetTimeByISO(args.StartAt)
	if err != nil {
		return
	}
	endAt, err = CoreFilter.GetTimeByISO(args.EndAt)
	if err != nil {
		return
	}
	//创建请求
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "service_ad_apply", "INSERT INTO service_ad_apply (start_at, end_at, audit_at, audit_des, audit_ban_des, org_id, user_id, area_ids, ad_id, mark, name, des, cover_file_id, des_files, count, click_count, params) VALUES (:start_at,:end_at,to_timestamp(0),'','',:org_id,:user_id,:area_ids,0,:mark,:name,:des,:cover_file_id,:des_files,0,0,:params)", map[string]interface{}{
		"start_at":      startAt,
		"end_at":        endAt,
		"org_id":        args.OrgID,
		"user_id":       args.UserID,
		"area_ids":      args.AreaIDs,
		"mark":          args.Mark,
		"name":          args.Name,
		"des":           args.Des,
		"cover_file_id": args.CoverFileID,
		"des_files":     args.DesFiles,
		"params":        args.Params,
	}, &data)
	return
}

// ArgsAuditApply 审核请求参数
type ArgsAuditApply struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 验证用
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否审核
	IsAudit bool `json:"isAudit" check:"bool"`
	//申请描述
	AuditDes string `db:"audit_des" json:"auditDes" check:"des" min:"1" max:"600" empty:"true"`
	//拒绝原因
	AuditBanDes string `db:"audit_ban_des" json:"auditBanDes" check:"des" min:"1" max:"600" empty:"true"`
	//投放分区ID列
	AreaIDs pq.Int64Array `db:"area_ids" json:"areaIDs" check:"ids" empty:"true"`
	//分区标识码
	// 作为前端抽取数据类型使用，可以重复指定多个
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
}

// AuditApply 审核请求
// 商户可以审核用户提交的请求
func AuditApply(args *ArgsAuditApply) (err error) {
	if args.IsAudit {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_ad_apply SET audit_at = NOW(), audit_des = :audit_des, area_ids = :area_ids, mark = :mark WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND audit_at < to_timestamp(1000000)", map[string]interface{}{
			"id":        args.ID,
			"org_id":    args.OrgID,
			"audit_des": args.AuditDes,
			"area_ids":  args.AreaIDs,
			"mark":      args.Mark,
		})
		if err != nil {
			return
		}
		//建立广告和绑定关系
		var data FieldsApply
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, start_at, end_at, org_id, user_id, area_ids, ad_id, mark, name, des, cover_file_id, des_files, params FROM service_ad_apply WHERE id = $1", args.ID)
		if err != nil {
			return
		}
		//检查是否存在广告及绑定关系
		if data.AdID > 0 {
			err = UpdateAD(&ArgsUpdateAD{
				ID:          data.ID,
				OrgID:       data.OrgID,
				Mark:        data.Mark,
				Name:        data.Name,
				Des:         data.Des,
				CoverFileID: data.CoverFileID,
				DesFiles:    data.DesFiles,
				Params:      data.Params,
			})
			if err != nil {
				return
			}
		} else {
			var adData FieldsAD
			adData, err = CreateAD(&ArgsCreateAD{
				OrgID:       data.OrgID,
				Mark:        data.Mark,
				Name:        data.Name,
				Des:         data.Des,
				CoverFileID: data.CoverFileID,
				DesFiles:    data.DesFiles,
				Params:      data.Params,
			})
			if err != nil {
				return
			}
			for _, v := range args.AreaIDs {
				_, err = SetBind(&ArgsSetBind{
					StartAt: data.StartAt,
					EndAt:   data.EndAt,
					OrgID:   data.OrgID,
					AreaID:  v,
					AdID:    adData.ID,
					Factor:  1,
					Params:  CoreSQLConfig.FieldsConfigsType{},
				})
				if err != nil {
					return
				}
			}
			_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_ad_apply SET ad_id = :ad_id WHERE id = :id", map[string]interface{}{
				"id":    data.ID,
				"ad_id": adData.ID,
			})
			if err != nil {
				return
			}
		}
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_ad_apply SET audit_at = to_timestamp(0), audit_ban_des = :audit_ban_des WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
			"id":            args.ID,
			"org_id":        args.OrgID,
			"audit_ban_des": args.AuditBanDes,
		})
		if err != nil {
			return
		}
		//删除广告和绑定关系
		var data FieldsApply
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, org_id, ad_id FROM service_ad_apply WHERE id = $1", args.ID)
		if err != nil {
			return
		}
		if data.AdID > 0 {
			//删除广告数据
			_ = DeleteAD(&ArgsDeleteAD{
				ID:    data.AdID,
				OrgID: data.OrgID,
			})
		}
	}
	return
}

// ArgsDeleteApply 删除请求参数
type ArgsDeleteApply struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 验证用
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可以由用户发起，用户发起的广告主要向商户投放，由商户进行审核处理
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// DeleteApply 删除请求
func DeleteApply(args *ArgsDeleteApply) (err error) {
	//获取数据
	var data FieldsApply
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, ad_id FROM service_ad_apply WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND ($3 < 1 OR user_id = $3) AND delete_at < to_timestamp(1000000)", args.ID, args.OrgID, args.UserID)
	if err != nil || data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//获取广告ID
	var adData FieldsAD
	adData, err = GetADByID(&ArgsGetADByID{
		ID:    data.AdID,
		OrgID: -1,
	})
	//如果存在数据，则继续处理
	if err == nil && adData.ID > 0 {
		//删除广告数据
		_ = DeleteAD(&ArgsDeleteAD{
			ID:    adData.ID,
			OrgID: adData.OrgID,
		})
	}
	//删除数据
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_ad_apply", "id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:user_id < 1 OR user_id = :user_id)", args)
	return
}

// 更新统计数据
func updateApplyAnalysis(adID int64, count int64, clickCount int64) {
	_, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_ad_apply SET count = count + :count, click_count = click_count + :click_count WHERE ad_id = :ad_id", map[string]interface{}{
		"ad_id":       adID,
		"count":       count,
		"click_count": clickCount,
	})
	if err != nil {
		//不记录错误
	}
}
