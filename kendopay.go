package kendopay

//签代扣协约资所需信息
type SignInfo struct {
	BankCode    string
	AccountType int
	AccountNo   string
	AccountName string
	AccountProp int
	IDType      int
	IDCardNo    string
	PhoneNo     string
	Reamrk      string
	Merrem      string
}

// 支付协议的相关
type Protocol interface {
	//签约申请，获取验证码
	SignApply(info SignInfo)
	//重发验证码
	ReSendSMS()
	//签约（验证码）确认
	SignCfrm()
}

type Defray interface {
}
