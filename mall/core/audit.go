package MallCore

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsGetAuditList 获取审核列表参数
type ArgsGetAuditList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//商品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//是否已经过期
	NeedIsExpire bool `db:"need_is_expire" json:"needIsExpire" check:"bool"`
	IsExpire     bool `db:"is_expire" json:"isExpire" check:"bool"`
	//是否过审
	NeedIsPassing bool `db:"need_is_passing" json:"needIsPassing" check:"bool"`
	IsPassing     bool `db:"is_passing" json:"isPassing" check:"bool"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetAuditList 获取审核列表
func GetAuditList(args *ArgsGetAuditList) (dataList []FieldsAudit, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ProductID > -1 {
		where = where + " AND product_id = :product_id"
		maps["product_id"] = args.ProductID
	}
	if args.NeedIsExpire {
		if args.IsExpire {
			where = where + " AND expire_at < NOW()"
		} else {
			where = where + " AND expire_at >= NOW()"
		}
	}
	if args.NeedIsPassing {
		if args.IsPassing {
			where = where + " AND audit_at > to_timestamp(1000000)"
		} else {
			where = where + " AND audit_at <= to_timestamp(1000000)"
		}
	}
	if args.Search != "" {
		where = where + " AND (ban_des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "mall_core_audit"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, expire_at, org_id, product_id, audit_at, ban_des, ban_des_files FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "expire_at", "audit_at"},
	)
	return
}

// ArgsGetAuditByProduct 通过商品获取审核数据参数
type ArgsGetAuditByProduct struct {
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//商品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
}

// GetAuditByProduct 通过商品获取审核数据
func GetAuditByProduct(args *ArgsGetAuditByProduct) (data FieldsAudit, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, expire_at, org_id, product_id, audit_at, ban_des, ban_des_files FROM mall_core_audit WHERE org_id = $1 AND product_id = $2", args.OrgID, args.ProductID)
	return
}

// ArgsGetAuditByProducts 通过一组商品获取审核数据参数
type ArgsGetAuditByProducts struct {
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//商品ID
	ProductIDs pq.Int64Array `db:"product_ids" json:"productIDs" check:"ids"`
}

// GetAuditByProducts 通过一组商品获取审核数据
func GetAuditByProducts(args *ArgsGetAuditByProducts) (dataList []FieldsAudit, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, expire_at, org_id, product_id, audit_at, ban_des, ban_des_files FROM mall_core_audit WHERE org_id = $1 AND product_id = ANY($2)", args.OrgID, args.ProductIDs)
	return
}

// ArgsCreateAudit 创建审核请求参数
type ArgsCreateAudit struct {
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//商品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
}

// CreateAudit 创建审核请求
func CreateAudit(args *ArgsCreateAudit) (data FieldsAudit, err error) {
	//存在审核或已经审核通过，则跳出
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, expire_at, org_id, product_id, audit_at, ban_des, ban_des_files FROM mall_core_audit WHERE org_id = $1 AND product_id = $2 AND audit_at > to_timestamp(1000000) ORDER BY id DESC LIMIT 1", args.OrgID, args.ProductID)
	if err == nil && data.ID > 0 {
		err = updateProductPublish(&argsUpdateProductPublish{
			ID: data.ProductID,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("update product publish, ", err))
			return
		}
		return
	}
	//检查配置项
	var configAutoAudit bool
	configAutoAudit, err = BaseConfig.GetDataBool("MallCommodityAutoAudit")
	if err != nil {
		configAutoAudit = false
	}
	auditAt := time.Time{}
	if configAutoAudit {
		auditAt = CoreFilter.GetNowTime()
	}
	var configExpire int
	configExpire, err = BaseConfig.GetDataInt("MallCommodityAuditExpire")
	if err != nil {
		configExpire = 24
	}
	configExpireAt := CoreFilter.GetNowTimeCarbon().AddHours(configExpire)
	//创建新审核请求
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "mall_core_audit", "INSERT INTO mall_core_audit (expire_at, org_id, product_id, audit_at, ban_des, ban_des_files) VALUES (:expire_at,:org_id,:product_id,:audit_at,'','{}')", map[string]interface{}{
		"expire_at":  configExpireAt.Time,
		"org_id":     args.OrgID,
		"product_id": args.ProductID,
		"audit_at":   auditAt,
	}, &data)
	if err != nil {
		err = errors.New(fmt.Sprint("create mall core audit, insert, ", err))
		return
	} else {
		if configAutoAudit {
			err = updateProductPublish(&argsUpdateProductPublish{
				ID: data.ProductID,
			})
			if err != nil {
				err = errors.New(fmt.Sprint("update product publish, ", err))
				return
			}
		}
	}
	//反馈
	return
}

// ArgsUpdateAuditPassing 通过审核参数
type ArgsUpdateAuditPassing struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// UpdateAuditPassing 通过审核
func UpdateAuditPassing(args *ArgsUpdateAuditPassing) (err error) {
	var data FieldsAudit
	data, err = getAudit(&argsGetAudit{
		ID: args.ID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("audit not exist, ", err))
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE mall_core_audit SET update_at = NOW(), audit_at = NOW() WHERE id = :id", args)
	if err != nil {
		err = errors.New(fmt.Sprint("update audit passing, ", err))
		return
	}
	err = updateProductPublish(&argsUpdateProductPublish{
		ID: data.ProductID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("update product publish, ", err, ", product id: ", data.ProductID))
		return
	}
	return
}

// ArgsUpdateAuditBan 拒绝审核参数
type ArgsUpdateAuditBan struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//拒绝原因
	BanDes string `db:"ban_des" json:"banDes" check:"des" min:"1" max:"1000"`
	//拒绝附加文件
	BanDesFiles pq.Int64Array `db:"ban_des_files" json:"banDesFiles" check:"ids" empty:"true"`
}

// UpdateAuditBan 拒绝审核
func UpdateAuditBan(args *ArgsUpdateAuditBan) (err error) {
	var data FieldsAudit
	data, err = getAudit(&argsGetAudit{
		ID: args.ID,
	})
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE mall_core_audit SET update_at = NOW(), audit_at = to_timestamp(0), ban_des = :ban_des, ban_des_files = :ban_des_files WHERE id = :id", args)
	if err != nil {
		return
	}
	err = DeleteProduct(&ArgsDeleteProduct{
		ID:    data.ProductID,
		OrgID: data.OrgID,
	})
	if err != nil {
		return
	}
	return
}

// ArgsDeleteAudit 删除审核请求参数
type ArgsDeleteAudit struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteAudit 删除审核请求
func DeleteAudit(args *ArgsDeleteAudit) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "mall_core_audit", "id = :id AND org_id = :org_id AND expire_at >= NOW() AND audit_at < to_timestamp(1000000)", args)
	return
}

// argsGetAudit 获取指定的审核数据包参数
type argsGetAudit struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// getAudit 获取指定的审核数据包
func getAudit(args *argsGetAudit) (data FieldsAudit, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, expire_at, org_id, product_id, audit_at, ban_des, ban_des_files FROM mall_core_audit WHERE id = $1", args.ID)
	if data.ID < 1 {
		err = errors.New("data not exist")
	}
	return
}
