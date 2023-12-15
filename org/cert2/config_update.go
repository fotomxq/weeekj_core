package OrgCert2

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsUpdateConfig 修改配置参数
type ArgsUpdateConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定来源
	// user 用户 / org 商户 / finance_assets 财务资产 /
	BindFrom string `db:"bind_from" json:"bindFrom" check:"mark"`
	//默认过期时间长度
	DefaultExpire string `db:"default_expire" json:"defaultExpire"`
	//标识码
	// 用于程序化识别处理机制
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//证件名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//审核模式
	// none 无需审核; wait 人工审核; auto 自动审核(依赖其他模块，根据扩展参数具体识别方案);
	AuditType string `db:"audit_type" json:"auditType" check:"mark" empty:"true"`
	//审核费用
	// 如果为0则无效
	Currency int   `db:"currency" json:"currency" check:"currency" empty:"true"`
	Price    int64 `db:"price" json:"price" check:"price" empty:"true"`
	//序列号长度
	// 0则不会限制，但数据表最多存储300位，超出需使用扩展参数存储
	SNLen int `db:"sn_len" json:"snLen" check:"intThan0" empty:"true"`
	//通知类型
	// none 无通知; audit 审核通过后通知; expire 过期前通知; all 全部通知;
	TipType string `db:"tip_type" json:"tipType" check:"mark" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateConfig 修改配置
func UpdateConfig(args *ArgsUpdateConfig) (err error) {
	//核对auditType
	if args.AuditType == "" {
		args.AuditType = "none"
	}
	if err = checkAuditType(args.AuditType); err != nil {
		return
	}
	//更新的mark如果和之前不一致
	data := getConfigByMark(args.OrgID, args.Mark)
	if data.ID > 0 && !CoreSQL.CheckTimeHaveData(data.DeleteAt) && data.ID != args.ID {
		err = errors.New("mark is exist")
		return
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert_config2 SET update_at = NOW(), bind_from = :bind_from, default_expire = :default_expire, mark = :mark, name = :name, des = :des, cover_file_id = :cover_file_id, des_files = :des_files, audit_type = :audit_type, currency = :currency, price = :price, sn_len = :sn_len, tip_type = :tip_type, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteConfigCache(args.ID)
	//反馈
	return
}
