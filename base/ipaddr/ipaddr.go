package BaseIPAddr

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"sync"
	"time"
)

//ip处理模块

var (
	//是否启用全局黑名单
	allowOpenBan = true
	//是否启用全局白名单
	allowOpenWhite = true
	//所有带有正则表达式的片段数据
	// 无论是否启动，都会自动加载
	dataByMatch []FieldsIPAddr
	//正则表达式片段集合，加载锁定
	dataByMatchLock sync.Mutex
)

// Init 初始化
func Init() (err error) {
	//加载所有正则数据
	_ = getAllMatchData()
	return
}

// SetOpenBan 设置是否启用黑名单
func SetOpenBan(args bool) {
	allowOpenBan = args
}

// SetOpenWhite 设置是否启用白名单
func SetOpenWhite(args bool) {
	allowOpenWhite = args
}

// CheckAuto 自动化通过处理
// 不能是ban，同时必须white
// 可根据情况跳过相关设定，但注意是全局跳过，否则必须遵守上述规则
func CheckAuto(args string) bool {
	isBan := CheckIsBan(args)
	isWhite := CheckIsWhite(args)
	return isWhite && !isBan
}

// CheckIsBan 检查IP是否在黑名单
func CheckIsBan(args string) bool {
	if !allowOpenBan {
		return false
	}
	data, err := getIP(args)
	if err != nil {
		return false
	}
	//如果过期，则删除配置并返回不存在
	if data.ExpireAt.Unix() < CoreFilter.GetNowTime().Unix() {
		//删除此配置
		_ = ClearIP(&ArgsClearIP{
			IP: data.IP,
		})
		return false
	}
	if !data.IsBan {
		matchData := checkMatchData(args)
		for _, v := range matchData {
			if v.IsBan {
				return true
			}
		}
	}
	return data.IsBan
}

// CheckIsWhite 检查IP是否在白名单
func CheckIsWhite(args string) bool {
	if !allowOpenWhite {
		return true
	}
	data, err := getIP(args)
	if err != nil {
		return false
	}
	//如果过期，则删除配置并返回不存在
	if data.ExpireAt.Unix() < CoreFilter.GetNowTime().Unix() {
		//删除此配置
		_ = ClearIP(&ArgsClearIP{
			IP: data.IP,
		})
		return false
	}
	if !data.IsWhite {
		matchData := checkMatchData(args)
		for _, v := range matchData {
			if v.IsWhite {
				return true
			}
		}
	}
	return data.IsWhite
}

// ArgsGetList 获取列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList
	//搜索
	Search string
}

// GetList 获取列表
func GetList(args *ArgsGetList) (dataList []FieldsIPAddr, dataCount int64, err error) {
	maps := map[string]interface{}{
		"search": args.Search,
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_ipaddr",
		"id",
		"SELECT id, create_at, update_at, expire_at, ip, is_match, is_ban, is_white FROM core_ipaddr WHERE ip ILIKE '%' || :search || '%'",
		"ip ILIKE '%' || :search || '%'",
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "expire_at", "ip"},
	)
	return
}

// ArgsSetBan 设置IP在黑名单参数
type ArgsSetBan struct {
	//IP地址或正则表达式
	IP string
	//是否为正则表达式
	IsMatch bool
	//是否列入名单
	B bool
	//过期时间
	ExpireTime string
}

// SetBan 设置IP在黑名单
func SetBan(args *ArgsSetBan) (err error) {
	var data FieldsIPAddr
	data, err = getIP(args.IP)
	if err != nil {
		err = SetIP(&ArgsSetIP{
			IP:         args.IP,
			IsMatch:    args.IsMatch,
			IsBan:      args.B,
			IsWhite:    false,
			ExpireTime: args.ExpireTime,
		})
	} else {
		err = SetIP(&ArgsSetIP{
			IP:         args.IP,
			IsMatch:    args.IsMatch,
			IsBan:      args.B,
			IsWhite:    data.IsWhite,
			ExpireTime: args.ExpireTime,
		})
	}
	if args.IsMatch {
		_ = getAllMatchData()
	}
	return
}

// ArgsSetWhite 设置IP在白名单情况参数
type ArgsSetWhite struct {
	//IP地址或正则表达式
	IP string
	//是否为正则表达式
	IsMatch bool
	//是否列入名单
	B bool
	//过期时间
	ExpireTime string
}

// SetWhite 设置IP在白名单情况
func SetWhite(args *ArgsSetWhite) (err error) {
	var data FieldsIPAddr
	data, err = getIP(args.IP)
	if err != nil {
		err = SetIP(&ArgsSetIP{
			IP:         args.IP,
			IsMatch:    args.IsMatch,
			IsBan:      false,
			IsWhite:    args.B,
			ExpireTime: args.ExpireTime,
		})
	} else {
		err = SetIP(&ArgsSetIP{
			IP:         args.IP,
			IsMatch:    args.IsMatch,
			IsBan:      data.IsBan,
			IsWhite:    args.B,
			ExpireTime: args.ExpireTime,
		})
	}
	if args.IsMatch {
		_ = getAllMatchData()
	}
	return
}

// ArgsSetIP 设置数据参数
type ArgsSetIP struct {
	//IP地址或正则表达式
	IP string
	//是否为正则表达式
	IsMatch bool
	//是否列入黑名单
	IsBan bool
	//是否列入白名单
	IsWhite bool
	//过期时间
	// ISO时间格式
	ExpireTime string
}

// SetIP 设置数据
func SetIP(args *ArgsSetIP) (err error) {
	var expireAt time.Time
	expireAt, err = CoreFilter.GetTimeByISO(args.ExpireTime)
	if err != nil {
		err = errors.New("expire time, " + err.Error())
		return err
	}
	var data FieldsIPAddr
	data, err = getIP(args.IP)
	if err != nil {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO core_ipaddr(expire_at, ip, is_match, is_ban, is_white) VALUES(:expire_at, :ip, :is_match, :is_ban, :is_white)", map[string]interface{}{
			"expire_at": expireAt,
			"ip":        args.IP,
			"is_match":  args.IsMatch,
			"is_ban":    args.IsBan,
			"is_white":  args.IsWhite,
		})
		if err != nil {
			err = errors.New("insert data, " + err.Error())
			return err
		}
	} else {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_ipaddr SET expire_at=:expire_at, is_match=:is_match, is_ban=:is_ban, is_white=:is_white, update_at=NOW() WHERE id = :id;", map[string]interface{}{
			"expire_at": expireAt,
			"id":        data.ID,
			"is_match":  args.IsMatch,
			"is_ban":    args.IsBan,
			"is_white":  args.IsWhite,
		})
		if err != nil {
			err = errors.New("update data, " + err.Error())
			return err
		}
	}
	if args.IsWhite {
		_ = getAllMatchData()
	}
	return nil
}

// ArgsGetAddressByIP 获取IP地理位置参数
type ArgsGetAddressByIP struct {
	//IP地址
	IP string
}

// GetAddressByIP 获取IP地理位置
// TODO: 未完成，寻找合适的外部模块实现即可
func GetAddressByIP(args *ArgsGetAddressByIP) (data string, err error) {
	return
}

// ClearAll 重建数据集合，清理数据
func ClearAll() (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_ipaddr", "true", nil)
	if err != nil {
		return err
	}
	dataByMatchLock.Lock()
	dataByMatch = []FieldsIPAddr{}
	dataByMatchLock.Unlock()
	return nil
}

// ArgsClearIP 清除某个IP参数
type ArgsClearIP struct {
	//IP地址
	IP string `db:"ip"`
}

// ClearIP 清除某个IP
func ClearIP(args *ArgsClearIP) (err error) {
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_ipaddr", "ip", args)
	return
}

// getIP 获取数据
func getIP(ip string) (data FieldsIPAddr, err error) {
	err = checkIP(ip)
	if err != nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, expire_at, ip, is_match, is_ban, is_white FROM core_ipaddr WHERE ip=$1;", ip)
	if err != nil {
		err = errors.New("get ip data, " + err.Error())
	}
	return
}

// getAllMatchData 获取所有正则数据
func getAllMatchData() (err error) {
	dataByMatchLock.Lock()
	defer dataByMatchLock.Unlock()
	var data []FieldsIPAddr
	err = Router2SystemConfig.MainDB.Select(&data, "SELECT id, create_at, update_at, expire_at, ip, is_match, is_ban, is_white FROM core_ipaddr WHERE is_match=true;")
	if err != nil {
		return
	}
	dataByMatch = data
	return
}

// checkMatchData 检查数据是否符合正则表达式
func checkMatchData(ip string) []FieldsIPAddr {
	var checkResult []FieldsIPAddr
	for _, v := range dataByMatch {
		if CoreFilter.MatchStr(v.IP, ip) {
			checkResult = append(checkResult, v)
		}
	}
	return checkResult
}

// checkIP 检查IP是否合法
// 支持ip v4/v6
func checkIP(ip string) error {
	if b := CoreFilter.CheckIP(ip); !b {
		return errors.New("ip type is error")
	}
	return nil
}
