package MallCoreMod

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	"github.com/lib/pq"
	"time"
)

// FieldsOtherOption 附加选项
type FieldsOtherOption struct {
	//分类1
	Sort1 FieldsOtherOptionSort `db:"sort1" json:"sort1"`
	//分类2
	Sort2 FieldsOtherOptionSort `db:"sort2" json:"sort2"`
	//数据集合
	DataList FieldsOtherOptionChildList `db:"data_list" json:"dataList"`
}

// Value sql底层处理器
func (t FieldsOtherOption) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsOtherOption) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// FieldsOtherOptionSort 分类数据集合
type FieldsOtherOptionSort struct {
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//选项组
	Options pq.StringArray `db:"options" json:"options"`
}

// Value sql底层处理器
func (t FieldsOtherOptionSort) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsOtherOptionSort) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsOtherOptionChildList []FieldsOtherOptionChild

// Value sql底层处理器
func (t FieldsOtherOptionChildList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsOtherOptionChildList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsOtherOptionChild struct {
	//分类1的选项key
	Sort1 int `db:"sort1" json:"sort1" check:"intThan0" empty:"true"`
	//分类2的选项key
	Sort2 int `db:"sort2" json:"sort2" check:"intThan0" empty:"true"`
	//商品ID
	// 可以给0，则必须声明其他项目内容
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//选项key
	Key string `db:"key" json:"key" check:"mark"`
	//实际费用
	PriceReal int64 `db:"price_real" json:"priceReal" check:"price" empty:"true"`
	//折扣截止
	PriceExpireAt time.Time `db:"price_expire_at" json:"priceExpireAt" check:"isoTime" empty:"true"`
	//折扣前费用
	Price int64 `db:"price" json:"price" check:"price" empty:"true"`
	//封面ID
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//总库存
	Count int `db:"count" json:"count" check:"intThan0" empty:"true"`
	//商品条形码编码
	Code string `db:"code" json:"code"`
}

// Value sql底层处理器
func (t FieldsOtherOptionChild) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsOtherOptionChild) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// DataOtherOptions 附加选项
type DataOtherOptions struct {
	//分类1
	Sort1 DataOtherOptionSort `db:"sort1" json:"sort1"`
	//分类2
	Sort2 DataOtherOptionSort `db:"sort2" json:"sort2"`
	//数据集合
	DataList DataOtherOptionChildList `db:"data_list" json:"dataList"`
}

// Value sql底层处理器
func (t DataOtherOptions) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *DataOtherOptions) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// DataOtherOptionSort 分类数据集合
type DataOtherOptionSort struct {
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//选项组
	Options pq.StringArray `db:"options" json:"options"`
}

// Value sql底层处理器
func (t DataOtherOptionSort) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *DataOtherOptionSort) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// GetFields 转化数据为数据表结构
func (t DataOtherOptions) GetFields() (FieldsOtherOption, error) {
	var data FieldsOtherOption
	var err error
	data.Sort1 = FieldsOtherOptionSort{
		Name:    t.Sort1.Name,
		Options: t.Sort1.Options,
	}
	data.Sort2 = FieldsOtherOptionSort{
		Name:    t.Sort2.Name,
		Options: t.Sort2.Options,
	}
	for _, v := range t.DataList {
		var vExpire time.Time
		if v.PriceExpireAt != "" {
			vExpire, err = CoreFilter.GetTimeByISO(v.PriceExpireAt)
			if err != nil {
				return FieldsOtherOption{}, err
			}
		}
		data.DataList = append(data.DataList, FieldsOtherOptionChild{
			Sort1:         v.Sort1,
			Sort2:         v.Sort2,
			ProductID:     v.ProductID,
			Key:           v.Key,
			PriceReal:     v.PriceReal,
			PriceExpireAt: vExpire,
			Price:         v.Price,
			CoverFileID:   v.CoverFileID,
			Count:         v.Count,
			Code:          v.Code,
		})
	}
	return data, nil
}

type DataOtherOptionChildList []DataOtherOptionChild

// Value sql底层处理器
func (t DataOtherOptionChildList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *DataOtherOptionChildList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type DataOtherOptionChild struct {
	//分类1的选项key
	Sort1 int `db:"sort1" json:"sort1" check:"intThan0" empty:"true"`
	//分类2的选项key
	Sort2 int `db:"sort2" json:"sort2" check:"intThan0" empty:"true"`
	//选项名称
	Name string `db:"name" json:"name" check:"name"`
	//商品ID
	// 可以给0，则必须声明其他项目内容
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//选项key
	Key string `db:"key" json:"key" check:"mark"`
	//实际费用
	PriceReal int64 `db:"price_real" json:"priceReal" check:"price" empty:"true"`
	//折扣截止
	PriceExpireAt string `db:"price_expire_at" json:"priceExpireAt" check:"isoTime" empty:"true"`
	//折扣前费用
	Price int64 `db:"price" json:"price" check:"price" empty:"true"`
	//封面ID
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//总库存
	Count int `db:"count" json:"count" check:"intThan0" empty:"true"`
	//商品条形码编码
	Code string `db:"code" json:"code"`
}

// Value sql底层处理器
func (t DataOtherOptionChild) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *DataOtherOptionChild) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
