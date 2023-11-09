package payletter

import (
	"github.com/whitecubeinc/go-utils"
	"strings"
)

const (
	registerAutoPayUrl                = "https://pgapi.payletter.com/v1.0/payments/request"
	transactionAutoPayUrl             = "https://pgapi.payletter.com/v1.0/payments/autopay"
	cancelTransactionUrl              = "https://pgapi.payletter.com/v1.0/payments/cancel"
	easyPayRegisterUrl                = "https://pgapi.payletter.com/api/url/request/register-method"
	easyPayRegisterTestUrl            = "https://testppay.payletter.com/api/url/request/register-method"
	easyPayGetRegisteredMethodUrl     = "https://pgapi.payletter.com/api/user/methods"
	easyPayGetRegisteredMethodTestUrl = "https://testppay.payletter.com/api/user/methods"
)

var (
	PgCode            = utils.NewStringEnum[pgCode](nil, strings.ToLower)
	PayletterCardCode = map[string]string{
		"P001": "비씨카드",
		"P002": "KB국민카드",
		"P003": "하나카드",
		"P004": "삼성카드",
		"P005": "신한카드",
		"P006": "현대카드",
		"P007": "롯데카드",
		"P008": "NH카드",
		"P009": "씨티카드",
		"P010": "수협카드",
		"P011": "우리카드",
		"P012": "신협카드",
		"P013": "광주카드",
		"P014": "전북카드",
		"P015": "제주카드",
		"P016": "우체국카드",
		"P017": "MG새마을금고",
		"P018": "저축은행",
		"P019": "카카오뱅크카드",
		"P020": "은련카드",
		"P021": "해외 VISA 카드",
		"P022": "해외 MASTER 카드",
		"P023": "해외 JCB 카드",
		"P024": "해외 AMX 카드",
		"P025": "해외 DINERS 카드",
	}
	PayletterBankCode = map[string]string{
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
)
