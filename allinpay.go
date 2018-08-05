package kendopay

import (
	"fmt"
	"math/rand"
	"time"

	utility "./Utility"
)

var Allinpay *TLPay

type TLPay struct {
	TLQuickpayAdrr string
	TLUserDev      string
	MerchantIdDev  string
}

func init() {
	//allinpay = new(TLPay)
	Allinpay = &TLPay{utility.TLQuickpayAdrr, utility.TLUserDev, utility.MerchantIdDev}
}

func InitTLReqHeader(reqHead *TLReqHeader) *TLReqINFO {
	return &TLReqINFO{
		TRX_CODE:    reqHead.TRX_CODE,      // 交易代码
		VERSION:     "04",                  // 版本（04）
		DATA_TYPE:   "2",                   // 数据格式（2：xml格式）
		LEVEL:       reqHead.LEVEL,         // 处理级别（0-9  0优先级最低，默认为5）
		MERCHANT_ID: utility.MerchantIdDev, // 商户代码
		USER_NAME:   utility.TLUserDev,     // 用户名
		USER_PASS:   utility.TLPassWordDev, // 用户密码
		REQ_SN:      reqHead.REQ_SN,        // 交易流水号（必须全局唯一）
		SIGNED_MSG:  reqHead.SIGNED_MSG,    // 签名信息
	}
}

//通联快捷代扣
func (tl TLPay) Collect(tradeReq QuickTradeReqFASTTRX) (result string, err error) {
	var reqSN = fmt.Sprintf("XxxX_%4d%d", rand.Intn(10000), time.Now().Unix())
	// 请求头
	headerReq := &TLReqHeader{
		TRX_CODE:   "310011", // 交易代码
		LEVEL:      "6",      // 处理级别（0-9  0优先级最低，默认为5）
		REQ_SN:     reqSN,    // 交易流水号（必须全局唯一）
		SIGNED_MSG: "",       // 签名信息
	}

	tradeReq.MERCHANT_ID = utility.MerchantIdDev
	payReq := &QuickTradeReq{
		INFO: *InitTLReqHeader(headerReq),
	}
	// xml
	tlReqXML, err := utility.ToTLRequestXmlByte(payReq)
	if err != nil {
		fmt.Println("转化为通联请求XML出错 : ", err)
		return "", err
	}
	bodyByte, err := utility.PostAllinPayXml(tl.TLQuickpayAdrr, tlReqXML)
	if err != nil {
		fmt.Println("通联服务请求异常 : ", err)
		return
	}
	payResp := &QuickTradeResp{}
	if err = utility.VerifyTLResponse(bodyByte, payResp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", payResp)
	return
}
