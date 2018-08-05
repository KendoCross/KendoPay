package kendopay

import (
	"fmt"
	"testing"
	"time"
)

//金账户注册
func TestJzhReg(t *testing.T) {
	jzhRegInfo := JzhRegister{
		CstmNM:   "王道",
		CertifTP: "0",
		CertifID: "612429198812220777",
		MobileNo: "13410100777",         // 手机号码
		CityID:   "5840",                // 开户行地区
		BankID:   "113",                 //银行
		BankNm:   "深圳支行",                //支行名称
		ActNo:    "6230580000144090777", //账户
		Remark:   "胡乱备注",                //备注
	}

	Register(jzhRegInfo)
	fmt.Println("---------------------金账户注册------------------------------------------")
	fmt.Println()
}

//通联快捷支付
func TestTLQuickPay(t *testing.T) {

	// 请求参数
	fastTrx := QuickTradeReqFASTTRX{
		BUSINESS_CODE: "19900",                             // 业务代码
		SUBMIT_TIME:   time.Now().Format("20060102150405"), // 提交时间（YYYYMMDDHHMMSS）
		AGRMNO:        "AIP9549180803000001424",            // 协议号（签约时返回的协议号）
		ACCOUNT_NO:    "6214850219949549",                  // 账号（借记卡或信用卡）
		ACCOUNT_NAME:  "幸福",                                // 账号名（借记卡或信用卡上的所有人姓名）
		AMOUNT:        "12345",                             // 金额(整数，单位分)
		CURRENCY:      "CNY",                               // 货币类型(人民币：CNY, 港元：HKD，美元：USD。不填时，默认为人民币)
		ID_TYPE:       "0",                                 // 开户证件类型（0身份证，1户口簿，2护照，3军官证，4士兵证...）
		ID:            "370613198705308692",                // 证件号
		TEL:           "18689262774",                       // 手机号
		CUST_USERID:   "github.com/ikaiguang",              // 自定义用户号（商户自定义的用户号，开发人员可当作备注字段使用）
		SUMMARY:       "交易附言",                              // 交易附言（填入网银的交易备注）
		REMARK:        "不备注",                               // 备注（供商户填入参考信息）
	}

	result, err := Allinpay.Collect(fastTrx)
	if err != nil {
		fmt.Printf("%#v \n", err)
	}

	fmt.Println(result)

	fmt.Println("---------------------通联快捷支付------------------------------------------")
	fmt.Println()
}

//宝付裸扣
func TestBFBareCollect(t *testing.T) {
	Baofoo.BareCollect()
	fmt.Println("---------------------宝付裸扣------------------------------------------")
	fmt.Println()
}
