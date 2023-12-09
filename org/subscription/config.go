package OrgSubscription

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "gitee.com/weeekj/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// 查看配置列表
type ArgsGetConfigList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

func GetConfigList(args *ArgsGetConfigList) (dataList []FieldsConfig, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.Search != "" {
		where = where + " AND (title ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"org_sub_config",
		"id",
		"SELECT id, create_at, update_at, delete_at, mark, func_list, time_type, time_n, currency, price, price_old, title, cover_file_id, des_files, style_id, params FROM org_sub_config WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetConfigByID 查看指定配置ID参数
type ArgsGetConfigByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetConfigByID 查看指定配置ID
func GetConfigByID(args *ArgsGetConfigByID) (data FieldsConfig, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, mark, func_list, time_type, time_n, currency, price, price_old, title, des, cover_file_id, des_files, style_id, params FROM org_sub_config WHERE id = $1 AND delete_at < to_timestamp(1000000)", args.ID)
	return
}

type ArgsGetConfigByMark struct {
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
}

func GetConfigByMark(args *ArgsGetConfigByMark) (data FieldsConfig, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, mark, func_list, time_type, time_n, currency, price, price_old, title, des, cover_file_id, des_files, style_id, params FROM org_sub_config WHERE mark = $1 AND delete_at < to_timestamp(1000000)", args.Mark)
	return
}

// 批量查询
type ArgsGetConfigMore struct {
	//一组ID
	IDs pq.Int64Array `db:"ids" json:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
}

func GetConfigMore(args *ArgsGetConfigMore) (dataList []FieldsConfig, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "org_sub_config", "id, create_at, update_at, delete_at, mark, func_list, time_type, time_n, currency, price, price_old, title, cover_file_id, des_files, style_id, params", args.IDs, args.HaveRemove)
	return
}

// ArgsSetConfig 设置配置参数
type ArgsSetConfig struct {
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//对应组织功能
	FuncList pq.StringArray `db:"func_list" json:"funcList"`
	//时间类型
	// 0 小时 1 天 2 周 3 月 4 年
	TimeType int `db:"time_type" json:"timeType"`
	//时间长度
	TimeN int `db:"time_n" json:"timeN"`
	//开通价格
	Currency int   `db:"currency" json:"currency" check:"currency"`
	Price    int64 `db:"price" json:"price" check:"price"`
	//折扣前费用，用于展示
	PriceOld int64 `db:"price_old" json:"priceOld" check:"price"`
	//标题
	Title string `db:"title" json:"title" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//样式ID
	// 关联到样式库后，本记录的图片和文本将交给样式库布局实现
	StyleID int64 `db:"style_id" json:"styleID" check:"id" empty:"true"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetConfig 设置配置
func SetConfig(args *ArgsSetConfig) (data FieldsConfig, err error) {
	data, err = GetConfigByMark(&ArgsGetConfigByMark{
		Mark: args.Mark,
	})
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_sub_config SET update_at = NOW(), func_list = :func_list, time_type = :time_type, time_n = :time_n, currency = :currency, price = :price, price_old = :price_old, title = :title, des = :des, cover_file_id = :cover_file_id, des_files = :des_files, style_id = :style_id, params = :params WHERE id = :id", map[string]interface{}{
			"id":            data.ID,
			"func_list":     args.FuncList,
			"time_type":     args.TimeType,
			"time_n":        args.TimeN,
			"currency":      args.Currency,
			"price":         args.Price,
			"price_old":     args.PriceOld,
			"title":         args.Title,
			"des":           args.Des,
			"cover_file_id": args.CoverFileID,
			"des_files":     args.DesFiles,
			"style_id":      args.StyleID,
			"params":        args.Params,
		})
		if err == nil {
			data, err = GetConfigByID(&ArgsGetConfigByID{
				ID: data.ID,
			})
			if err != nil {
				err = errors.New("update data after get by id, " + err.Error())
			}
		} else {
			err = errors.New("update failed, " + err.Error())
		}
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_sub_config", "INSERT INTO org_sub_config (mark, func_list, time_type, time_n, currency, price, price_old, title, des, cover_file_id, des_files, style_id, params) VALUES (:mark,:func_list,:time_type,:time_n,:currency,:price,:price_old,:title,:des,:cover_file_id,:des_files,:style_id,:params)", args, &data)
	if err != nil {
		err = errors.New("insert failed, " + err.Error())
	}
	return
}

// ArgsDeleteConfig 删除配置参数
type ArgsDeleteConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteConfig 删除配置
func DeleteConfig(args *ArgsDeleteConfig) (err error) {
	_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "org_sub_config", "id", args)
	return
}
