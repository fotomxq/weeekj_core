package UserRole

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
	FinancePay "github.com/fotomxq/weeekj_core/v5/finance/pay"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetPayLogList 获取日志列表参数
type ArgsGetPayLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//角色ID
	RoleID int64 `db:"role_id" json:"roleID" check:"id" empty:"true"`
	//系统来源
	SystemFrom string `db:"system_from" json:"systemFrom" check:"mark" empty:"true"`
	//系统ID
	FromID int64 `db:"from_id" json:"fromID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetPayLogList 获取日志列表
func GetPayLogList(args *ArgsGetPayLogList) (dataList []FieldsPayLog, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.RoleID > -1 {
		where = where + "role_id = :role_id"
		maps["role_id"] = args.RoleID
	}
	if args.SystemFrom != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "system_from = :system_from"
		maps["system_from"] = args.SystemFrom
		if args.FromID > 0 {
			if where != "" {
				where = where + " AND "
			}
			where = where + "from_id = :from_id"
			maps["from_id"] = args.FromID
		}
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "user_role_pay_log"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, role_id, deposit_mark, currency, system_from, from_id, des, price FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsPayToRole 平台或商户给与角色资金参数
type ArgsPayToRole struct {
	//角色ID
	RoleID int64 `db:"role_id" json:"roleID" check:"id"`
	//关联组织
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//储蓄配置
	DepositMark string `db:"deposit_mark" json:"depositMark" check:"mark"`
	//货币
	// eg: 86
	Currency int `db:"currency" json:"currency" check:"currency"`
	//支付金额
	Price int64 `db:"price" json:"price" check:"price"`
	//系统来源
	SystemFrom string `db:"system_from" json:"systemFrom" check:"mark"`
	//系统ID
	FromID int64 `db:"from_id" json:"fromID" check:"mark" empty:"true"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
}

// PayToRole 平台或商户给与角色资金
func PayToRole(args *ArgsPayToRole) (err error) {
	//检查金额
	if args.Price < 1 {
		return
	}
	//检查是否已经抽成？
	var logID int64
	err = Router2SystemConfig.MainDB.Get(&logID, "SELECT id FROM user_role_pay_log WHERE system_from = $1 AND from_id = $2", args.SystemFrom, args.FromID)
	if err == nil && logID > 0 {
		err = errors.New("replace data")
		return
	}
	//获取角色信息
	var roleData FieldsRole
	roleData, err = GetRoleID(&ArgsGetRoleID{
		ID: args.RoleID,
	})
	if err != nil {
		return
	}
	//创建转账
	var payID int64 = 0
	if args.OrgID > 0 {
		var payData FinancePay.FieldsPayType
		payData, _, err = FinancePay.CreateQuickPay(&FinancePay.ArgsCreate{
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: args.SystemFrom,
				ID:     args.FromID,
				Mark:   "",
				Name:   "",
			},
			PaymentCreate: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     args.OrgID,
				Mark:   "",
				Name:   "",
			},
			PaymentChannel: CoreSQLFrom.FieldsFrom{
				System: "deposit",
				ID:     0,
				Mark:   "merchant",
				Name:   "",
			},
			PaymentFrom: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     args.OrgID,
				Mark:   "",
				Name:   "",
			},
			TakeCreate: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     roleData.UserID,
				Mark:   "",
				Name:   "",
			},
			TakeChannel: CoreSQLFrom.FieldsFrom{
				System: "deposit",
				ID:     0,
				Mark:   args.DepositMark,
				Name:   "",
			},
			TakeFrom: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     args.OrgID,
				Mark:   "",
				Name:   "",
			},
			Des:      args.Des,
			ExpireAt: time.Time{},
			Currency: args.Currency,
			Price:    args.Price,
			Params:   CoreSQLConfig.FieldsConfigsType{},
		})
		if err != nil {
			return
		}
		payID = payData.ID
	} else {
		_, _, err = FinanceDeposit.SetByFrom(&FinanceDeposit.ArgsSetByFrom{
			UpdateHash: "",
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     roleData.UserID,
				Mark:   "",
				Name:   "",
			},
			FromInfo: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     0,
				Mark:   "",
				Name:   "",
			},
			ConfigMark:      args.DepositMark,
			AppendSavePrice: args.Price,
		})
		if err != nil {
			return
		}
	}
	//创建日志
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_role_pay_log (role_id, deposit_mark, currency, system_from, from_id, des, pay_id, price) VALUES (:role_id,:deposit_mark,:currency,:system_from,:from_id,:des,:pay_id,:price)", map[string]interface{}{
		"role_id":      roleData.ID,
		"deposit_mark": args.DepositMark,
		"currency":     args.Currency,
		"system_from":  args.SystemFrom,
		"from_id":      args.FromID,
		"des":          args.Des,
		"pay_id":       payID,
		"price":        args.Price,
	})
	if err != nil {
		return
	}
	//反馈
	return
}
