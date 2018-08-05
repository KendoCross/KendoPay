package utility

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"sort"
)

//金账户测试环境
const (
	JzhTestAdrr   = "https://jzh-test.fuiou.com/jzh/"
	JzhMchntCd    = "0006510F0106121"
	JzhAPISuccess = "0000"
)

var jzhSignaKey *rsa.PrivateKey
var jzhVerfyKey *rsa.PublicKey

func init() {
	carFile, err := ioutil.ReadFile("Assets/prkey.key")
	if err != nil {
		return
	}
	pemBlock, _ := pem.Decode(carFile)
	if pemBlock == nil {
		return
	}
	parsedKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)

	if err != nil {
		return
	}
	jzhSignaKey = parsedKey

	verify, err := ioutil.ReadFile("Assets/pbkey.key")
	if err != nil {
		return
	}

	block, _ := pem.Decode(verify)
	if block == nil {
		return
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}

	jzhVerfyKey = pub.(*rsa.PublicKey)

}

func rsaSign(origData string) (string, error) {

	h := sha1.New()
	h.Write([]byte(origData))
	digest := h.Sum(nil)

	s, err := rsa.SignPKCS1v15(nil, jzhSignaKey, crypto.SHA1, digest)
	if err != nil {
		fmt.Errorf("rsaSign SignPKCS1v15 error")
		return "", err
	}
	data := base64.StdEncoding.EncodeToString(s)
	return string(data), nil
}

//金账户签名
func JZHSigna(dic map[string]string) url.Values {
	//STEP 1, 对key进行升序排序.
	sortedkeys := make([]string, 0)
	dicParms := make(url.Values)
	for k := range dic {
		sortedkeys = append(sortedkeys, k)
	}
	sort.Strings(sortedkeys)

	//STEP2, 对key=value的键值对用|连接起来，略过空值
	var signStrings string
	for _, k := range sortedkeys {
		signStrings += k + "=" + dic[k] + "|"
		dicParms[k] = []string{dic[k]}
	}
	signStrings = signStrings[:len(signStrings)-2]

	signStr, err := rsaSign(signStrings)
	if err != nil {

	}
	dicParms["signature"] = []string{signStr}

	return dicParms
}

func reflectSignValue(res interface{}) (sign string, err error) {
	// reflect value
	reflectValue := reflect.ValueOf(res)
	// pointer
	if reflectValue.Kind() != reflect.Ptr {
		err = errors.New("result interface not a pointer type : reflect.Ptr")
		return sign, err
	}

	signValue := reflectValue.Elem().FieldByName("Signature")
	if !signValue.IsValid() {
		err = errors.New("error : Signature flied not exist")
		return sign, err
	}
	// string
	if signValue.Kind() != reflect.String {
		err = errors.New("signValue not a string type : reflect.String")
		return sign, err
	}
	// get sign
	sign = signValue.String()
	// assignment value
	signValue.SetString("")
	// return
	return sign, err
}

//金账户验签
func JzhVerify(respXml string, respModel interface{}) error {

	if err := xml.Unmarshal([]byte(respXml), respModel); err != nil {
		return errors.New("反序列化失败！")
	}
	reg, err := regexp.Compile(`<plain>.*</plain>`)
	if err != nil {
		return errors.New("无<plain>节点")
	}
	plainTxt := reg.FindAllString(respXml, 1)
	if len(plainTxt) < 1 {
		return errors.New("无<plain>节点")
	}
	hash := crypto.SHA1
	h := hash.New()
	if _, err := h.Write([]byte(plainTxt[0])); err != nil {
		return errors.New("SHA1失败")
	}
	hashed := h.Sum(nil)

	signStr, err := reflectSignValue(respModel)
	if err != nil {
		return errors.New("反射获取签名信息失败！")
	}

	signData, err := base64.StdEncoding.DecodeString(signStr)
	if err != nil {
		return errors.New("反序列化失败！")
	}

	if err := rsa.VerifyPKCS1v15(jzhVerfyKey, hash, hashed, signData); err != nil {
		return err
	}
	return nil
}

//金账户Post请求
func JzhPost(action string, body io.Reader) (resp *http.Response, err error) {
	httpClient := &http.Client{}
	resp, err = httpClient.Post(JzhTestAdrr+action, "application/x-www-form-urlencoded", body)
	return
}
