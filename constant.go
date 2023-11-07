package payletter

import (
	"github.com/whitecubeinc/go-utils"
	"strings"
)

const (
	registerAutoPayUrl    = "https://pgapi.payletter.com/v1.0/payments/request"
	transactionAutoPayUrl = "https://pgapi.payletter.com/v1.0/payments/autopay"
	cancelTransactionUrl  = "https://pgapi.payletter.com/v1.0/payments/cancel"
)

var (
	PgCode = utils.NewStringEnum[pgCode](nil, strings.ToLower)
)
