package OrgCert2

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsCreateConfig 创建新的配置参数
type ArgsCreateConfig struct {
	//默认过期时间长度
	DefaultExpire string `db:"default_expire" json:"defaultExpire"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定来源
	// user 用户 / org 商户 / finance_assets 财务资产 /
	BindFrom string `db:"bind_from" json:"bindFrom" check:"mark"`
	//标识码
	// 用于程序化识别处理机制
	Mark string `db:"mark" json:"mark" check:"mark"`
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
	AuditType string `db:"audit_type" json:"auditType" check:"mark"`
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

// CreateConfig 创建新的配置
func CreateConfig(args *ArgsCreateConfig) (err error) {
	//核对auditType
	if args.AuditType == "" {
		args.AuditType = "none"
	}
	if err = checkAuditType(args.AuditType); err != nil {
		return
	}
	//核对TipType
	if args.TipType == "" {
		args.TipType = "none"
	}
	if err = checkTipType(args.TipType); err != nil {
		return
	}
	//检查mark
	data := getConfigByMark(args.OrgID, args.Mark)
	if data.ID > 0 && !CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		err = errors.New("config mark is exist")
		return
	}
	//创建数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_cert_config2 (default_expire, org_id, bind_from, mark, name, des, cover_file_id, des_files, audit_type, currency, price, sn_len, tip_type, params) VALUES (:default_expire,:org_id,:bind_from,:mark,:name,:des,:cover_file_id,:des_files,:audit_type,:currency,:price,:sn_len,:tip_type,:params)", args)
	if err != nil {
		return
	}
	return
}
