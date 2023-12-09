package VCodeImageCore

import (
	"errors"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/mojocn/base64Captcha"
	"github.com/mojocn/base64Captcha/store"
	"strings"
	"time"
)

var (
	//验证码验证过期时间
	expiredTime int64 = 120
	//验证码生成时间限制
	// 同一个token不能连续N个时间获取验证码
	intervalTime int64 = 1

	//图形验证码 参数配置
	imageConfig base64Captcha.ConfigCharacter
	//store存储器
	base64CaptchaStore store.Store
	// GCLimitNumber The number of captchas created that triggers garbage collection used by default store.
	// 默认图像验证GC清理的上限个数
	gcLimitNumber int
	// Expiration time of captchas used by default store.
	// 内存保存验证码的时限
	expiration time.Duration
)

// Init 初始化
func Init() (err error) {
	//设定存储器
	gcLimitNumber = 10240
	expiration = 10 * time.Minute
	base64CaptchaStore = store.NewMemoryStore(gcLimitNumber, expiration)
	base64Captcha.SetCustomStore(base64CaptchaStore)
	//初始化图像配置
	imageConfig = base64Captcha.ConfigCharacter{
		Height: 60,
		Width:  240,
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:               base64Captcha.CaptchaModeNumberAlphabet,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsShowHollowLine:   true,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         6,
	}
	return
}

// ArgsGenerate 生成验证码参数
type ArgsGenerate struct {
	//会话
	Token int64 `db:"token"`
}

// Generate 生成验证码
func Generate(args *ArgsGenerate) (imgData base64Captcha.CaptchaInterface, err error) {
	//要反馈的空数据，失败时反馈
	var errorResult base64Captcha.CaptchaInterface
	//获取最近的验证码
	var data FieldsVCodeType
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM core_vcode_image WHERE token=$1", args.Token)
	//如果存在数据，则检查是否没有超出获取间隔时间限制
	if data.Token < 1 {
		if data.CreateAt.Unix()+intervalTime > CoreFilter.GetNowTime().Unix() {
			return errorResult, errors.New("interval time is too short")
		}
		//通过后，删除该数据
		if _, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_vcode_image", "token=:token", args); err != nil {
			//删除失败，记录日志
			//return errorResult, errors.New("cannot delete old data")
		}
	}
	//生成新的验证码数据
	// 注意，该key并不是最终的key，需
	key, captchaInterface := base64Captcha.GenerateCaptcha("", imageConfig)
	value := base64CaptchaStore.Get(key, true)
	if _, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO core_vcode_image(token, value) VALUES(:token, :value)", map[string]interface{}{
		"token": args.Token,
		"value": value,
	}); err != nil {
		return errorResult, errors.New("cannot create new data, " + err.Error())
	}
	return captchaInterface, nil
}

// ArgsCheck 验证验证码参数
type ArgsCheck struct {
	//会话
	Token int64
	//验证的值
	Value string
}

// Check 验证验证码
func Check(args *ArgsCheck) bool {
	//检查配置项开关
	// VerificationCodeImageON
	VerificationCodeImageON, err := BaseConfig.GetDataBool("VerificationCodeImageON")
	if err != nil {
		CoreLog.Error("get config by VerificationCodeImageON, ", err)
		VerificationCodeImageON = true
	}
	if !VerificationCodeImageON {
		return true
	}
	//查询并验证
	var data FieldsVCodeType
	if err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, token, value FROM core_vcode_image WHERE token=$1 ORDER BY id DESC LIMIT 1;", args.Token); err != nil {
		return false
	}
	if data.CreateAt.Unix()+expiredTime > CoreFilter.GetNowTime().Unix() && strings.ToLower(data.Value) == strings.ToLower(args.Value) {
		if _, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_vcode_image", "token=:token", args); err != nil {
			//删除失败，记录日志
			CoreLog.Error("check vcode is error, cannot delete vcode data, " + err.Error())
		}
		return true
	}
	return false
}

// 清除所有验证码
// 将清除所有验证码，该方法可用于系统维护等操作
// 注意，该操作将直接操作原始数据，不是普通的delete删除操作
func Clear() (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_vcode_image", "true", nil)
	return
}
