package payletter

import (
	"github.com/whitecubeinc/go-utils"
	"strings"
)

type pgCode struct {
	CreditCard string
	EasyBank   string
	NaverPay   string
	NaverCard  string
	NaverPoint string
}

func (o *pgCode) IsNaverCode(code string) bool {
	switch code {
	case o.NaverPoint, o.NaverCard, o.NaverPay:
		return true
	default:
		return false
	}
}

type payletterCardCode struct {
	P001     string `value:"비씨카드"`
	P002     string `value:"KB국민카드"`
	P003     string `value:"하나카드"`
	P004     string `value:"삼성카드"`
	P005     string `value:"신한카드"`
	P006     string `value:"현대카드"`
	P007     string `value:"롯데카드"`
	P008     string `value:"NH카드"`
	P009     string `value:"씨티카드"`
	P010     string `value:"수협카드"`
	P011     string `value:"우리카드"`
	P012     string `value:"신협카드"`
	P013     string `value:"광주카드"`
	P014     string `value:"전북카드"`
	P015     string `value:"제주카드"`
	P016     string `value:"우체국카드"`
	P017     string `value:"MG새마을금고"`
	P018     string `value:"저축은행"`
	P019     string `value:"카카오뱅크카드"`
	P020     string `value:"은련카드"`
	P021     string `value:"해외 VISA 카드"`
	P022     string `value:"해외 MASTER 카드"`
	P023     string `value:"해외 JCB 카드"`
	P024     string `value:"해외 AMX 카드"`
	P025     string `value:"해외 DINERS 카드"`
	ValueMap map[string]string
}

type transactionDateType struct {
	Transaction string
	Settle      string
}

const (
	registerAutoPayUrl                = "https://pgapi.payletter.com/v1.0/payments/request"
	transactionAutoPayUrl             = "https://pgapi.payletter.com/v1.0/payments/autopay"
	cancelTransactionUrl              = "https://pgapi.payletter.com/v1.0/payments/cancel"
	partialCancelTransactionUrl       = "https://pgapi.payletter.com/v1.0/payments/cancel/partial"
	easyPayRegisterUrl                = "https://ppay.payletter.com/api/url/request/register-method"
	easyPayRegisterTestUrl            = "https://testppay.payletter.com/api/url/request/register-method"
	easyPayGetRegisteredMethodUrl     = "https://ppay.payletter.com/api/user/methods"
	easyPayGetRegisteredMethodTestUrl = "https://testppay.payletter.com/api/user/methods"
	easyPayCancelUrl                  = "https://ppay.payletter.com/api/payments/cancel"
	easyPayCancelTestUrl              = "https://testppay.payletter.com/api/payments/cancel"
	easyPayTransactionUrl             = "https://ppay.payletter.com/api/url/request/request-payment"
	easyPayTransactionTestUrl         = "https://testppay.payletter.com/api/url/request/request-payment"
	normalTransactionUrl              = "https://pgapi.payletter.com/v1.0/payments/request"
	getTransactionListUrl             = "https://pgapi.payletter.com/v1.0/payments/transaction/list"
)

var (
	PgCode   = utils.NewStringEnum[pgCode](nil, strings.ToLower)
	CardCode = utils.NewConstantFromTag[payletterCardCode](strings.ToUpper)
	BankCode = map[string]string{
		"003": "IBK기업은행",
		"002": "KDB산업은행",
		"004": "KB국민은행",
		"007": "수협은행",
		"011": "NH농협은행",
		"020": "우리은행",
		"023": "SC제일은행",
		"027": "한국씨티은행",
		"031": "대구은행",
		"032": "부산은행",
		"034": "광주은행",
		"035": "제주은행",
		"037": "전북은행",
		"039": "경남은행",
		"081": "하나은행",
		"088": "신한은행",
		"089": "케이뱅크",
		"090": "카카오뱅크",
		"012": "농협중앙회",
		"045": "새마을금고중앙회",
		"048": "신협중앙회",
		"064": "산림조합중앙회",
		"071": "우체국은행",
	}
	TransactionDateType = utils.NewStringEnum[transactionDateType](nil, strings.ToLower)
)
