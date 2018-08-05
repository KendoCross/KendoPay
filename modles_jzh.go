package kendopay

type jzhSign struct {
	Signature string `xml:"signature,omitempty"`
}

type jzhPlainReg struct {
	RespCode   string `xml:"resp_code,omitempty"`
	RespDesc   string `xml:"resp_desc,omitempty"`
	MchntTxnSn string `xml:"mchnt_txn_ssn,omitempty"`
}
type jzhRespReg struct {
	jzhSign
	Plain jzhPlainReg `xml:"plain,omitempty"`
}

//金账户注册所需信息
type JzhRegister struct {
	CstmNM   string // 客户姓名
	CertifTP string // 证件类型0:身份证7：其他证件
	CertifID string // 身份证号码/证件
	MobileNo string // 手机号码
	CityID   string // 开户行地区
	BankID   string //银行
	BankNm   string //支行名称
	ActNo    string //账户
	Remark   string //备注
}
