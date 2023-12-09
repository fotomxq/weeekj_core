package BaseWeixinApp

//微信APP相关授权等操作接口封装

var (
	//微信授权接口地址
	userLoginAuthURL = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=$1&secret=$2&code=$3&grant_type=authorization_code"
	//授权获取用户基本信息接口
	userLoginAuthInfoURL = "https://api.weixin.qq.com/sns/userinfo?access_token=$1&openid=$2"
)
