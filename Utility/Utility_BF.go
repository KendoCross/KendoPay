package utility

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"io/ioutil"

	"golang.org/x/crypto/pkcs12"
)

const BFTerminalId = "100000990"
const BFMemberId = "100000276"
const BFBaseAdrr = "https://vgw.baofoo.com/cutpayment/api/backTransRequest"

var BFSignCert *x509.Certificate
var BFSignPriKey *rsa.PrivateKey
var BFVerifyKey *rsa.PublicKey

//宝付密钥初始化
func init() {
	carFile, err := ioutil.ReadFile("Assets/bfkey.pfx")
	if err != nil {
		return
	}

	priKey, x509cert, err := pkcs12.Decode(carFile, `123456`)
	if err != nil {
		return
	}

	BFSignPriKey = priKey.(*rsa.PrivateKey)
	BFSignCert = x509cert

	verify, err := ioutil.ReadFile("Assets/bfkey.cer")
	if err != nil {
		return
	}

	block, _ := pem.Decode(verify)
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	BFVerifyKey = cert.PublicKey.(*rsa.PublicKey)
}

//宝付私钥加密
func BFEncrypt(srcStr string) (string, error) {
	base64Slice := make([]byte, base64.StdEncoding.EncodedLen(len(srcStr)))
	base64.StdEncoding.Encode(base64Slice, []byte(srcStr))
	// h := sha1.New()
	// h.Write(base64Slice)
	// digest := h.Sum(nil)
	// var _ = digest
	dstArr, err := priKeyByte(BFSignPriKey, base64Slice, true)
	//dstArr, err := rsa.SignPKCS1v15(nil, BFSignPriKey, crypto.SHA1, digest)
	//dstArr, err := rsa.EncryptPKCS1v15(nil, TLSignCert.PublicKey.(*rsa.PublicKey), base64Slice)
	if err != nil {
		return "", err
	}
	data := hex.EncodeToString(dstArr)
	return data, nil
}

//宝付公钥解密
func BFDecrypt(sign string) (string, error) {
	signData, _ := hex.DecodeString(sign)
	plainText, err := pubKeyByte(BFVerifyKey, signData, false)
	if err != nil {
		return "", err
	}
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(plainText)))
	_, err = base64.StdEncoding.Decode(dst, plainText)
	if err != nil {
		return "", err
	}
	return string(dst), nil
}
