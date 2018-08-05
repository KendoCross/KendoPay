package kendopay

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	utility "./Utility"
)

func Register(regInfo interface{}) error {

	var regModel = regInfo.(JzhRegister)
	var reqSN = fmt.Sprintf("KendoPay_%4d%d", rand.Intn(10000), time.Now().Unix())
	dicInfo := map[string]string{
		"ver":            "0.44",
		"mchnt_cd":       utility.JzhMchntCd,
		"mchnt_txn_ssn":  reqSN,
		"cust_nm":        regModel.CstmNM,
		"certif_id":      regModel.CertifID,
		"mobile_no":      regModel.MobileNo,
		"city_id":        regModel.CityID,
		"parent_bank_id": regModel.BankID,
		"bank_nm":        regModel.BankNm,
		"capAcntNo":      regModel.ActNo,
		"rem":            regModel.Remark,
	}

	ulrDic := utility.JZHSigna(dicInfo)
	//接口兼容问题，此字段无需签名
	ulrDic["certif_tp"] = []string{regModel.CertifTP}

	reqParms := ulrDic.Encode()
	reqBody := bytes.NewBuffer([]byte(reqParms))

	resp, err := utility.JzhPost("reg.action", reqBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	plainStr := string(body)
	respM := &jzhRespReg{}
	if err := utility.JzhVerify(plainStr, respM); err != nil {
		errStr := "金账户验签不通过" + err.Error()
		fmt.Println(errStr)
		return errors.New(errStr)
	}
	if respM.Plain.RespCode != utility.JzhAPISuccess {
		errStr := "金账户注册失败：" + respM.Plain.RespDesc
		fmt.Println(respM.Plain.RespDesc)
		return errors.New(errStr)
	}
	return nil
}
