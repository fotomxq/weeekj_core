package WeixinPayV3Key

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
)

var (
	aesKey          string
	mchId           string
	privateSerialNo string
	appId           string
)

func getRandomString(i int) string    { return "" }
func hasha256(s string) []byte        { return []byte{} }
func base64DecodeStr(s string) string { return "" }

// 对消息的散列值进行数字签名
func signPKCS1v15(msg, privateKey []byte, hashType crypto.Hash) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key decode error")
	}
	pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("parse private key error")
	}
	key, ok := pri.(*rsa.PrivateKey)
	if ok == false {
		return nil, errors.New("private key format error")
	}
	sign, err := rsa.SignPKCS1v15(cryptoRand.Reader, key, hashType, msg)
	if err != nil {
		return nil, errors.New("sign error")
	}
	return sign, nil
}

// base编码
func base64EncodeStr(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// 生成身份认证信息
func authorization(method string, paramMap map[string]interface{}, rawUrl string) (token string, err error) {
	var body string
	if len(paramMap) != 0 {
		paramJsonBytes, err := json.Marshal(paramMap)
		if err != nil {
			return token, err
		}
		body = string(paramJsonBytes)
	}
	urlPart, err := url.Parse(rawUrl)
	if err != nil {
		return token, err
	}
	canonicalUrl := urlPart.RequestURI()
	timestamp := CoreFilter.GetNowTime().Unix()
	nonce := getRandomString(32)
	message := fmt.Sprintf("%s\n%s\n%d\n%s\n%s\n", method, canonicalUrl, timestamp, nonce, body)
	open, err := os.Open("./private.pem") // 商户私有证书路径或者从数据库读取
	if err != nil {
		return token, err
	}
	defer open.Close()
	privateKey, err := ioutil.ReadAll(open)
	if err != nil {
		return token, err
	}
	signBytes, err := signPKCS1v15(hasha256(message), privateKey, crypto.SHA256)
	if err != nil {
		return token, err
	}
	sign := base64EncodeStr(signBytes)
	token = fmt.Sprintf("mchid=\"%s\",nonce_str=\"%s\",timestamp=\"%d\",serial_no=\"%s\",signature=\"%s\"",
		mchId, nonce, timestamp, privateSerialNo, sign)
	return token, nil
}

// 报文解密
func decryptGCM(aesKey, nonceV, ciphertextV, additionalDataV string) ([]byte, error) {
	key := []byte(aesKey)
	nonce := []byte(nonceV)
	additionalData := []byte(additionalDataV)
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextV)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, additionalData)
	if err != nil {
		return nil, err
	}
	return plaintext, err
}

// 获取公钥
const publicKeyUrl = "https://api.mch.weixin.qq.com/v3/certificates"

type TokenResponse struct {
	Data []TokenResponseData `json:"data"`
}
type TokenResponseData struct {
	EffectiveTime      string             `json:"effective_time"`
	EncryptCertificate EncryptCertificate `json:"encrypt_certificate"`
	ExpireTime         string             `json:"expire_time"`
	SerialNo           string             `json:"serial_no"`
}
type EncryptCertificate struct {
	Algorithm      string `json:"algorithm"`
	AssociatedData string `json:"associated_data"`
	Ciphertext     string `json:"ciphertext"`
	Nonce          string `json:"nonce"`
}

var publicSyncMap sync.Map

// 获取公钥
func getPublicKey() (key string, err error) {
	var prepareTime int64 = 24 * 3600 * 3 // 证书提前三天过期旧证书，获取新证书
	nowTime := CoreFilter.GetNowTime().Unix()
	// 读取公钥缓存数据
	cacheValueKey := fmt.Sprintf("app_id:%s:public_key:value", appId)
	cacheExpireTimeKey := fmt.Sprintf("app_id:%s:public_key:expire_time", appId)
	cacheValue, keyValueOk := publicSyncMap.Load(cacheValueKey)
	cacheExpireTime, expireTimeOk := publicSyncMap.Load(cacheExpireTimeKey)
	if keyValueOk && expireTimeOk {
		// 格式化时间
		local, _ := time.LoadLocation("Local")
		location, _ := time.ParseInLocation(time.RFC3339, cacheExpireTime.(string), local)
		// 判断是否过期，证书没有过期直接返回
		if location.Unix()-prepareTime > nowTime {
			return cacheValue.(string), nil
		}
	}
	token, err := authorization(http.MethodGet, nil, publicKeyUrl)
	if err != nil {
		return key, err
	}
	request, err := http.NewRequest(http.MethodGet, publicKeyUrl, nil)
	if err != nil {
		return key, err
	}
	request.Header.Add("Authorization", "WECHATPAY2-SHA256-RSA2048 "+token)
	request.Header.Add("User-Agent", "用户代理(https://zh.wikipedia.org/wiki/User_agent)")
	request.Header.Add("Content-type", "application/json;charset='utf-8'")
	request.Header.Add("Accept", "application/json")
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return key, err
	}
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return key, err
	}
	//fmt.Println(string(bodyBytes))
	var tokenResponse TokenResponse
	if err = json.Unmarshal(bodyBytes, &tokenResponse); err != nil {
		return key, err
	}
	for _, encryptCertificate := range tokenResponse.Data {
		// 格式化时间
		local, _ := time.LoadLocation("Local")
		location, err := time.ParseInLocation(time.RFC3339, encryptCertificate.ExpireTime, local)
		if err != nil {
			return key, err
		}
		// 判断是否过期，证书没有过期直接返回
		if location.Unix()-prepareTime > nowTime {
			decryptBytes, err := decryptGCM(aesKey, encryptCertificate.EncryptCertificate.Nonce, encryptCertificate.EncryptCertificate.Ciphertext,
				encryptCertificate.EncryptCertificate.AssociatedData)
			if err != nil {
				return key, err
			}
			key = string(decryptBytes)
			publicSyncMap.Store(cacheValueKey, key)
			publicSyncMap.Store(cacheExpireTimeKey, encryptCertificate.ExpireTime)
			return key, nil
		}
	}
	return key, errors.New("get public key error")
}

// VerifyRsaSign 验证数字签名
func VerifyRsaSign(msg []byte, sign []byte, publicStr []byte, hashType crypto.Hash) bool {
	//pem解码
	block, _ := pem.Decode(publicStr)
	//x509解码
	publicKeyInterface, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}
	publicKey := publicKeyInterface.PublicKey.(*rsa.PublicKey)
	//验证数字签名
	err = rsa.VerifyPKCS1v15(publicKey, hashType, msg, sign) //crypto.SHA1
	return err == nil
}

// signatureValidate 验证签名
func signatureValidate(timeStamp, rawPost, nonce, signature string) (bool, error) {
	signature = base64DecodeStr(signature)
	message := fmt.Sprintf("%s\n%s\n%s\n", timeStamp, nonce, rawPost)
	publicKey, err := getPublicKey()
	if err != nil {
		return false, err
	}

	return VerifyRsaSign(hasha256(message), []byte(signature), []byte(publicKey), crypto.SHA256), nil
}

type NotifyResponse struct {
	CreateTime string         `json:"create_time"`
	Resource   NotifyResource `json:"resource"`
}
type NotifyResource struct {
	Ciphertext     string `json:"ciphertext"`
	AssociatedData string `json:"associated_data"`
	Nonce          string `json:"nonce"`
}

func notifyDecrypt(rawPost string) (decrypt string, err error) {
	var notifyResponse NotifyResponse
	if err = json.Unmarshal([]byte(rawPost), &notifyResponse); err != nil {
		return decrypt, err
	}
	decryptBytes, err := decryptGCM(aesKey, notifyResponse.Resource.Nonce, notifyResponse.Resource.Ciphertext,
		notifyResponse.Resource.AssociatedData)
	if err != nil {
		return decrypt, err
	}
	decrypt = string(decryptBytes)
	return decrypt, nil
}
