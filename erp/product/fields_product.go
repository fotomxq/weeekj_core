package ERPProduct

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsProduct ERP商品信息
type FieldsProduct struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//保质期小时
	ExpireHour int `db:"expire_hour" json:"expireHour"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//选择供应商
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//供货商名称
	CompanyName string `db:"company_name" json:"companyName" check:"title" min:"1" max:"300" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签ID列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//SN
	SN string `db:"sn" json:"sn"`
	//商品条形码编码
	Code string `db:"code" json:"code"`
	//拼音助记码
	PinYin string `db:"pin_yin" json:"pinYin" check:"des" min:"1" max:"300" empty:"true"`
	//英文名称
	EnName string `db:"en_name" json:"enName" check:"des" min:"1" max:"300" empty:"true"`
	//生产厂商名称
	ManufacturerName string `db:"manufacturer_name" json:"manufacturerName" check:"des" min:"1" max:"300" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"title" min:"1" max:"300"`
	//副标题
	TitleDes string `db:"title_des" json:"titleDes" check:"des" min:"1" max:"300" empty:"true"`
	//商品描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//封面ID
	CoverFileIDs pq.Int64Array `db:"cover_file_ids" json:"coverFileIDs"`
	//货物重量
	// 单位g
	Weight int `db:"weight" json:"weight" check:"intThan0" empty:"true"`
	//存储尺寸
	SizeW int `db:"size_w" json:"sizeW" check:"intThan0" empty:"true"`
	SizeH int `db:"size_h" json:"sizeH" check:"intThan0" empty:"true"`
	SizeZ int `db:"size_z" json:"sizeZ" check:"intThan0" empty:"true"`
	//规格类型
	// 0 盒装; 1 袋装; 3 散装; 4 瓶装
	PackType int `db:"pack_type" json:"packType"`
	//包装单位名称
	PackUnitName string `db:"pack_unit_name" json:"packUnitName" check:"des" min:"1" max:"100" empty:"true"`
	//包装内部含有产品数量
	PackUnit int `db:"pack_unit" json:"packUnit"`
	//建议零售价（不含税）
	TipPrice int64 `db:"tip_price" json:"tipPrice" check:"price" empty:"true"`
	//建议零售价（含税）
	// 该建议价格用于可直接用于最终零售价填入
	TipTaxPrice int64 `db:"tip_tax_price" json:"tipTaxPrice" check:"price" empty:"true"`
	//是否允许打折
	IsDiscount bool `db:"is_discount" json:"isDiscount" check:"bool" empty:"true"`
	//货币
	Currency int `db:"currency" json:"currency" check:"currency"`
	//单价成本（不含税）
	CostPrice int64 `db:"cost_price" json:"costPrice" check:"price" empty:"true"`
	//税率
	// 实际税率=tax/100000
	Tax int64 `db:"tax" json:"tax"`
	//单价成本（含税）
	TaxCostPrice int64 `db:"tax_cost_price" json:"taxCostPrice" check:"price" empty:"true"`
	//返利设计
	RebatePrice FieldsProductRebateList `db:"rebate_price" json:"rebatePrice" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// FieldsProductRebateList 返利设计列
type FieldsProductRebateList []FieldsProductRebate

// Value sql底层处理器
func (t FieldsProductRebateList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsProductRebateList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// FieldsProductRebate 返利设计
type FieldsProductRebate struct {
	//返利条件销量达到
	SellCount int64 `db:"sell_count" json:"sellCount" check:"intThan0"`
	//返利金额
	ReturnPrice int64 `db:"return_price" json:"returnPrice" check:"price"`
}

// Value sql底层处理器
func (t FieldsProductRebate) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsProductRebate) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
