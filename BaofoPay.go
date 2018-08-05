package kendopay

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	utility "./Utility"
)

var Baofoo *BaofooPay

type BaofooPay struct {
	BFTerminalId string
	BFMemberId   string
	BFBaseAdrr   string
}

//宝付密钥初始化
func init() {
	Baofoo = &BaofooPay{utility.BFTerminalId, utility.BFMemberId, utility.BFBaseAdrr}
}

//宝付裸扣
func (bf BaofooPay) BareCollect() {
	var reqSN = fmt.Sprintf("KendoPay%4d%d", rand.Intn(10000), time.Now().Unix())
	var transDate = time.Now().Format("20060102150405")
	var serialNo = fmt.Sprintf("%3d%d", rand.Intn(1000), time.Now().Unix())
	dic := make(map[string]interface{})
	dic["txn_sub_type"] = "13"
	dic["txn_sub_type"] = "13"
	dic["biz_type"] = "0000"
	dic["terminal_id"] = utility.BFTerminalId
	dic["member_id"] = utility.BFMemberId
	dic["pay_code"] = "ICBC"
	dic["pay_cm"] = "1"
	dic["acc_no"] = "6222020111122220000"
	dic["id_card_type"] = "01"
	dic["id_card"] = "340101198108119852"
	dic["id_holder"] = "张宝"
	dic["mobile"] = "18689262776"
	dic["trans_id"] = reqSN
	dic["txn_amt"] = "67800"
	dic["trade_date"] = transDate
	dic["additional_info"] = "不备注，0731"
	dic["req_reserved"] = "0731，备注不"
	dic["trans_serial_no"] = serialNo

	contentJson, err := json.Marshal(dic)

	if err != nil {
		return
	}
	signaContent, err := utility.BFEncrypt(string(contentJson))

	if err != nil {
		return
	}

	dicParms := url.Values{
		"version":      {"4.0.0.0"},
		"terminal_id":  {utility.BFTerminalId},
		"txn_type":     {"0431"},
		"txn_sub_type": {"13"},
		"member_id":    {utility.BFMemberId},
		"data_type":    {"json"},
		"data_content": {signaContent},
	}

	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cookieJar, _ := cookiejar.New(nil)

	httpClient := &http.Client{
		Jar:       cookieJar,
		Transport: tr,
	}

	reqParms := dicParms.Encode()
	reqBody := bytes.NewBuffer([]byte(reqParms))
	resp, err := httpClient.Post(utility.BFBaseAdrr, "application/x-www-form-urlencoded", reqBody)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	plainStr, err := utility.BFDecrypt(string(body))
	if err != nil {
		return
	}

	fmt.Println(plainStr)
}
