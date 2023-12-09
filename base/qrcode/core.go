package BaseQRCode

import (
	"encoding/base64"
	"github.com/skip2/go-qrcode"
)

//ArgsGetQRCode 生成一般的二维码
// 反馈base64数据
type ArgsGetQRCode struct {
	//二维码内容
	Param string `db:"param" json:"param"`
	//尺寸
	Size int `db:"size" json:"size" check:"intThan0"`
}

//GetQRCode 生成二维码
func GetQRCode(args *ArgsGetQRCode) (result string, err error) {
	if args.Size < 1 {
		args.Size = 1
	}
	if args.Size > 2048 {
		args.Size = 2048
	}
	var resultByte []byte
	resultByte, err = qrcode.Encode(args.Param, qrcode.Medium, args.Size)
	if err != nil {
		return
	}
	result = base64.StdEncoding.EncodeToString(resultByte)
	return
}

//GetQRCodeByte 生成二维码二进制
func GetQRCodeByte(args *ArgsGetQRCode) (result []byte, err error) {
	if args.Size < 1 {
		args.Size = 1
	}
	if args.Size > 2048 {
		args.Size = 2048
	}
	result, err = qrcode.Encode(args.Param, qrcode.Medium, args.Size)
	if err != nil {
		return
	}
	return
}
