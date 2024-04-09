package RouterParams

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"reflect"
)

// filterParams 过滤参数
// 该方法集合会对参数进行验证、过滤等一系列操作，具体请将识别信息写入json附加内容内，本函数将通过反向映射完成验证操作
// params请使用关联方式进行，否则无法实现参数过滤的操作
// 支持：
//
//	check="xxx" 检查变量是否满足条件
//	filter="xxx" 将变量直接进行参数过滤
//	empty=true 允许变量为空
func filterParams(c *gin.Context, params interface{}) (errField string, errCode string, b bool) {
	paramsType := reflect.TypeOf(params).Elem()
	valueType := reflect.ValueOf(params).Elem()
	step := 0
	//遍历结构，对内容进行解析处理
	for step < paramsType.NumField()-1 {
		//捕捉结构
		vField := paramsType.Field(step)
		vValueType := valueType.Field(step)
		//下一步
		step += 1
		//是否为object子结构
		if vField.Anonymous {
			errField, errCode, b = filterParams(c, params)
			if !b {
				errField = fmt.Sprint(vField.Name, ".", errField)
				return
			}
		} else {
			//不是对象，则说明书单独的变量结构
			//检查该结构过滤器
			// 检查参数的正确性
			checkMark := vField.Tag.Get("check")
			// 按照规则过滤参数，但不绝对拒绝参数
			filterMark := vField.Tag.Get("filter")
			// 参数值的范围，可用于文本长度判断、数字大小范围判断
			minStr := vField.Tag.Get("min")
			maxStr := vField.Tag.Get("max")
			var minVal, maxVal int
			if minStr != "" {
				var err error
				minVal, err = CoreFilter.GetIntByString(minStr)
				if err != nil {
					errField = vField.Name
					errCode = "check_min"
					return
				}
			}
			if maxStr != "" {
				var err error
				maxVal, err = CoreFilter.GetIntByString(maxStr)
				if err != nil {
					errField = vField.Name
					errCode = "check_max"
					return
				}
			}
			isEmpty := vField.Tag.Get("empty") == "true"
			//如果检查和验证全部为空，则跳过检查
			if checkMark == "" && filterMark == "" {
				continue
			}
			//继续检查
			errCode, b = filterParamChild(vValueType.Interface(), checkMark, filterMark, isEmpty, minVal, maxVal)
			if !b {
				errField = vField.Name
				return
			}
		}
	}
	//全部完成反馈成功
	b = true
	return
}

// filterParamChild 根据识别代码，对内容进行检查或过滤操作
func filterParamChild(data interface{}, checkMark string, filterMark string, isEmpty bool, min, max int) (errCode string, b bool) {
	//检查
	switch checkMark {
	case "id":
		//一个ID
		val, isOK := data.(int64)
		if !isOK {
			break
		}
		if isEmpty && val < 1 {
			b = true
		} else {
			b = val > 0
		}
	case "ids":
		//一组ID
		vals, isOK := data.([]int64)
		if !isOK {
			vals, isOK = data.(pq.Int64Array)
			if !isOK {
				break
			}
		}
		if isEmpty && len(vals) < 1 {
			b = true
		} else {
			for _, v := range vals {
				b = v > 0
				if !b {
					break
				}
			}
			if !b {
				break
			}
		}
	case "sn":
		//SN累计数
		val, isOK := data.(int64)
		if !isOK {
			break
		}
		if isEmpty && val == 0 {
			b = true
		} else {
			b = CoreFilter.CheckSN(val)
		}
	case "sns":
		//一组SN累计数
		vals, isOK := data.([]int64)
		if !isOK {
			vals, isOK = data.(pq.Int64Array)
			if !isOK {
				break
			}
		}
		if isEmpty && len(vals) < 1 {
			b = true
		} else {
			for _, v := range vals {
				b = CoreFilter.CheckSN(v)
				if !b {
					break
				}
			}
			if !b {
				break
			}
		}
	case "mark":
		//一个标识码
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckMark(val)
		}
	case "marks":
		//一组标识码
		vals, isOK := data.([]string)
		if !isOK {
			vals, isOK = data.(pq.StringArray)
			if !isOK {
				break
			}
		}
		if isEmpty && len(vals) < 1 {
			b = true
		} else {
			for _, v := range vals {
				b = CoreFilter.CheckMark(v)
				if !b {
					break
				}
			}
			if !b {
				break
			}
		}
	case "mark_page":
		//分页标识码
		// /index.html
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckMarkPage(val)
		}
	case "nationCode":
		//电话国家代码
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckNationCode(val)
		}
	case "phone":
		//电话联系方式
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckPhone(val)
		}
	case "vcode":
		//验证码
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckVcode(val)
		}
	case "page":
		//页数
		val, isOK := data.(int64)
		if !isOK {
			break
		}
		if isEmpty && val == 0 {
			b = true
		} else {
			b = CoreFilter.CheckPage(val)
		}
	case "max":
		//页长
		val, isOK := data.(int64)
		if !isOK {
			break
		}
		if isEmpty && val == 0 {
			b = true
		} else {
			b = CoreFilter.CheckMax(val)
		}
	case "sort":
		//排序规则
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckSort(val)
		}
	case "desc":
		//是否倒叙
		b = true
	case "username":
		//用户名
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckUsername(val)
		}
	case "password":
		//密码
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckPassword(val)
		}
	case "search":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckSearch(val)
		}
	case "filename":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckFileName(val)
		}
	case "expireTime":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckExpireTime(val)
		}
	case "isoTime":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			_, err := CoreFilter.GetTimeByISO(val)
			if err != nil {
				b = false
			} else {
				b = true
			}
		}
	case "defaultTime":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			_, err := CoreFilter.GetTimeByDefault(val)
			if err != nil {
				b = false
			} else {
				b = true
			}
		}
	case "name":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckNiceName(val)
		}
	case "email":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckEmail(val)
		}
	case "des":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckDes(val, min, max)
		}
	case "title":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			if min < 1 {
				min = 1
			}
			if max < 1 {
				max = 150
			}
			b = CoreFilter.CheckDes(val, min, max)
		}
	case "content":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckContent(val, min, max)
		}
	case "timeType":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckTimeType(val)
		}
	case "bool":
		b = true
	case "intThan0":
		val, isOK := data.(int)
		if !isOK {
			break
		}
		if isEmpty && val < 1 {
			b = true
		} else {
			b = val > 0
		}
	case "int64Than0":
		val, isOK := data.(int64)
		if !isOK {
			break
		}
		if isEmpty && val < 1 {
			b = true
		} else {
			b = val > 0
		}
	case "floatThan0":
		val, isOK := data.(float64)
		if !isOK {
			break
		}
		if isEmpty && val < 1 {
			b = true
		} else {
			b = val > 0
		}
	case "host":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckHost(val)
		}
	case "port":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckPort(val)
		}
	case "ip":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckIP(val)
		}
	case "currency":
		val, isOK := data.(int)
		if !isOK {
			break
		}
		if isEmpty && val < 1 {
			b = true
		} else {
			b = CoreFilter.CheckCountry(val)
		}
	case "price":
		val, isOK := data.(int64)
		if !isOK {
			break
		}
		if isEmpty && val < 1 {
			b = true
		} else {
			if val > 0 {
				b = true
			}
		}
	case "country":
		val, isOK := data.(int)
		if !isOK {
			break
		}
		if isEmpty && val < 1 {
			b = true
		} else {
			b = CoreFilter.CheckCountry(val)
		}
	case "province":
		val, isOK := data.(int)
		if !isOK {
			break
		}
		if isEmpty && val < 1 {
			b = true
		} else {
			b = CoreFilter.CheckProvince(val)
		}
	case "cityCode":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckCityCode(val)
		}
	case "address":
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckAddress(val)
		}
	case "address_data":
		val, isOK := data.(CoreSQLAddress.FieldsAddress)
		if !isOK {
			break
		}
		if isEmpty {
			b = true
		} else {
			b = CoreFilter.CheckAddress(val.Address)
			if !b {
				errCode = checkMark
				return
			}
			b = CoreFilter.CheckCountry(val.Country)
			if !b {
				errCode = checkMark
				return
			}
			b = CoreFilter.CheckProvince(val.Province)
			if !b {
				errCode = checkMark
				return
			}
			b = CoreFilter.CheckCity(val.City)
			if !b {
				errCode = checkMark
				return
			}
			b = CoreFilter.CheckNationCode(val.NationCode)
			if !b {
				errCode = checkMark
				return
			}
			b = CoreFilter.CheckPhone(val.Phone)
			if !b {
				errCode = checkMark
				return
			}
			b = CoreFilter.CheckNiceName(val.Name)
			if !b {
				errCode = checkMark
				return
			}
		}
	case "city":
		val, isOK := data.(int)
		if !isOK {
			break
		}
		if isEmpty && val < 1 {
			b = true
		} else {
			b = CoreFilter.CheckCity(val)
		}
	case "mapType":
		val, isOK := data.(int)
		if !isOK {
			break
		}
		if isEmpty && val < 1 {
			b = true
		} else {
			b = CoreFilter.CheckMapType(val)
		}
	case "gps":
		val, isOK := data.(float64)
		if !isOK {
			break
		}
		if isEmpty && val < 1 {
			b = true
		} else {
			b = CoreFilter.CheckGPS(val)
		}
	case "status":
		val, isOK := data.(int)
		if !isOK {
			break
		}
		if isEmpty && val > -1 && val < 99999 {
			b = true
		} else {
			if val > -1 && val < 99999 {
				b = true
			}
		}
	case "params":
		val, isOK := data.(CoreSQLConfig.FieldsConfigsType)
		if !isOK {
			break
		}
		if isEmpty && len(val) < 1 {
			b = true
		} else {
			//参数不能少于1
			if len(val) < 1 {
				errCode = checkMark
				return
			}
			//参数不能超过30条
			if len(val) > 30 {
				errCode = checkMark
				return
			}
			//依次核对数据
			for _, v := range val {
				b = CoreFilter.CheckMark(v.Mark)
				if !b {
					errCode = checkMark
					return
				}
				newStr := CoreFilter.CheckFilterStr(v.Val, 0, 500)
				if v.Val != newStr {
					errCode = checkMark
					return
				}
			}
		}
	case "createInfo":
		val, isOK := data.(CoreSQLFrom.FieldsFrom)
		if !isOK {
			break
		}
		if isEmpty && val.System == "" && val.ID < 1 && val.Name == "" && val.Mark == "" {
			b = true
		} else {
			if val.System != "" {
				b = CoreFilter.CheckMark(val.System)
				if !b {
					errCode = checkMark
					return
				}
			}
			if val.ID > 0 {
				//不验证
			}
			if val.Mark != "" {
				b = CoreFilter.CheckMark(val.Mark)
				if !b {
					errCode = checkMark
					return
				}
			}
			if val.Name != "" {
				b = CoreFilter.CheckNiceName(val.Name)
				if !b {
					errCode = checkMark
					return
				}
			}
		}
	case "sha1":
		//hash哈希值
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckHexSha1(val)
		}
	case "gender":
		//性别
		// 0 男 1 女 2 未知
		val, isOK := data.(int)
		if !isOK {
			break
		}
		if isEmpty && val < 0 {
			b = true
		} else {
			if val == 0 || val == 1 || val == 2 {
				b = true
			}
		}
	case "color":
		//16进制颜色
		val, isOK := data.(string)
		if !isOK {
			break
		}
		if isEmpty && val == "" {
			b = true
		} else {
			b = CoreFilter.CheckColor(val)
		}
	default:
		//无法识别的，失败跳出
		break
	}
	if !b {
		errCode = checkMark
		return
	}
	return
}
