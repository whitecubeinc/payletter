package payletter

import (
	"github.com/whitecubeinc/go-utils"
	"strings"
)

const (
	registerAutoPayUrl     = "https://pgapi.payletter.com/v1.0/payments/request"
	transactionAutoPayUrl  = "https://pgapi.payletter.com/v1.0/payments/autopay"
	cancelTransactionUrl   = "https://pgapi.payletter.com/v1.0/payments/cancel"
	easyPayRegisterUrl     = "https://pgapi.payletter.com/api/url/request/register-method"
	easyPayRegisterTestUrl = "https://testppay.payletter.com/api/url/request/register-method"
)

var (
	PgCode = utils.NewStringEnum[pgCode](nil, strings.ToLower)
)
