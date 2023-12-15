package MallCoreV3

import (
	BaseQiniu "github.com/fotomxq/weeekj_core/v5/base/qiniu"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"gopkg.in/errgo.v2/fmt/errors"
	"time"

	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
)

// ArgsGetShoppingList 获取购物车列表参数
type ArgsGetShoppingList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// DataGetShoppingList 获取商品列表数据
type DataGetShoppingList struct {
	//商品ID
	ProductID int64 `db:"product_id" json:"productID"`
	//商品key
	ProductKey string `db:"product_key" json:"productKey" check:"mark" empty:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//商品删除时间
	DeleteAt time.Time `json:"deleteAt"`
	//商品封面
	CoverFileURL string `json:"coverFileURL"`
	//商品名称
	ProductTitle string `json:"productTitle"`
	//商品原价
	Price int64 `json:"price"`
	//折扣截止
	PriceExpireAt time.Time `db:"price_expire_at" json:"priceExpireAt"`
	//商品优惠价
	RealPrice int64 `json:"realPrice"`
	//购买件数
	Count int64 `json:"count"`
	//可用票据
	UseTickets []int64 `json:"useTickets"`
}

// GetShoppingList 获取商品列表
func GetShoppingList(args *ArgsGetShoppingList) (dataList []DataGetShoppingList, dataCount int64, err error) {
	where := "user_id = :user_id"
	maps := map[string]interface{}{
		"user_id": args.UserID,
	}
	var rawList []FieldsShopping
	tableName := "mall_core_shopping"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id, create_at, product_id, product_key, count "+"FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if err != nil || len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		appendData := DataGetShoppingList{
			ProductID:     v.ProductID,
			ProductKey:    v.ProductKey,
			CreateAt:      v.CreateAt,
			DeleteAt:      time.Time{},
			CoverFileURL:  "",
			ProductTitle:  "",
			Price:         0,
			PriceExpireAt: time.Time{},
			RealPrice:     0,
			Count:         v.Count,
			UseTickets:    []int64{},
		}
		var productData FieldsProduct
		productData, _ = GetProduct(&ArgsGetProduct{
			ID:    v.ProductID,
			OrgID: -1,
		})
		if productData.ID > 0 {
			appendData.DeleteAt = productData.DeleteAt
			appendData.ProductTitle = productData.Title
			if v.ProductKey == "" {
				appendData.Price = productData.Price
				appendData.PriceExpireAt = productData.PriceExpireAt
				appendData.RealPrice = productData.PriceReal
			} else {
				isFind := false
				for _, v2 := range productData.OtherOptions.DataList {
					if v2.Key == v.ProductKey {
						isFind = true
						appendData.Price = v2.Price
						appendData.PriceExpireAt = v2.PriceExpireAt
						appendData.RealPrice = v2.PriceReal
						break
					}
				}
				if !isFind {
					appendData.Price = productData.Price
					appendData.PriceExpireAt = productData.PriceExpireAt
					appendData.RealPrice = productData.PriceReal
				}
			}
			appendData.UseTickets = productData.UserTicket
			if len(productData.CoverFileIDs) > 0 {
				appendData.CoverFileURL, _ = BaseQiniu.GetPublicURLStr(productData.CoverFileIDs[0])
			}
		}
		dataList = append(dataList, appendData)
	}
	return
}

// ArgsAddShopping 添加购物车参数
type ArgsAddShopping struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//商品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//商品key
	ProductKey string `db:"product_key" json:"productKey" check:"mark" empty:"true"`
	//数量
	Count int64 `db:"count" json:"count"`
}

// AddShopping 添加购物车
func AddShopping(args *ArgsAddShopping) (err error) {
	if args.Count < 1 {
		args.Count = 0
	}
	if args.Count > 9999 {
		args.Count = 9999
	}
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM mall_core_shopping WHERE user_id = $1 AND product_id = $2", args.UserID, args.ProductID)
	if err == nil && id > 0 {
		if args.Count < 1 {
			return DeleteShopping(&ArgsDeleteShopping{
				UserID:    args.UserID,
				ProductID: args.ProductID,
			})
		} else {
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE mall_core_shopping SET product_key = :product_key, count = :count WHERE id = :id", map[string]interface{}{
				"id":          id,
				"product_key": args.ProductKey,
				"count":       args.Count,
			})
			return
		}
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO mall_core_shopping (user_id, product_id, product_key, count) VALUES (:user_id,:product_id,:product_key,:count)", args)
	return
}

// ArgsDeleteShopping 删除购物车参数
type ArgsDeleteShopping struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//商品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
}

// DeleteShopping 删除购物车
func DeleteShopping(args *ArgsDeleteShopping) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "mall_core_shopping", "user_id = :user_id AND product_id = :product_id", args)
	return
}

// ArgsClearShopping 清空购物车参数
type ArgsClearShopping struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

func ClearShopping(args *ArgsClearShopping) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "mall_core_shopping", "user_id = :user_id", args)
	return
}
