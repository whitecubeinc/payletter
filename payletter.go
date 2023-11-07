package payletter

import (
	"errors"
	"fmt"
	"github.com/whitecubeinc/go-utils"
	"net/http"
)

type PayLetter struct{}

func GetPayLetter() IPayLetter {
	return &PayLetter{}
}

func (o *PayLetter) RegisterAutoPay(req ReqRegisterAutoPay) (res ResRegisterAutoPay, err error) {
	paymentData := reqPaymentData{
		PgCode:          req.PgCode,
		ClientID:        req.ClientID,
		ServiceName:     req.ServiceName,
		UserID:          req.UserID,
		UserName:        req.UserName,
		OrderNo:         req.OrderNo,
		Amount:          req.Amount,
		ProductName:     req.ProductName,
		EmailFlag:       "N",
		AutoPayFlag:     "Y",
		ReceiptFlag:     "N",
		CustomParameter: req.CustomParameter,
		ReturnUrl:       req.ReturnUrl,
		CallbackUrl:     req.CallbackEndpoint,
		CancelUrl:       req.CancelUrl,
	}

	payLetterRes := utils.Post[utils.M](
		registerAutoPayUrl,
		paymentData,
		http.Header{
			"Authorization": []string{fmt.Sprintf("PLKEY %s", req.APIKey)},
			"Content-Type":  []string{"application/json"},
		},
	)

	if v, exists := payLetterRes["error"]; exists { // 500 error
		e := v.(map[string]any)
		err = errors.New(fmt.Sprintf("[%v]%v", e["code"], e["message"]))
		return
	}

	if code, exists := payLetterRes["code"]; exists {
		err = errors.New(fmt.Sprintf("[%v]%v", code, payLetterRes["message"]))
		return
	}

	res = ResRegisterAutoPay{
		MobileUrl: payLetterRes["mobile_url"].(string),
		OnlineUrl: payLetterRes["online_url"].(string),
		Token:     payLetterRes["token"].(string),
	}

	return
}

func (o *PayLetter) TransactionAutoPay(req ReqTransactionAutoPay) (res ResTransactionAutoPay, err error) {
	payLetterRes := utils.Post[utils.M](
		transactionAutoPayUrl,
		req,
		http.Header{
			"Authorization": []string{fmt.Sprintf("PLKEY %s", req.APIKey)},
			"Content-Type":  []string{"application/json"},
		},
	)
	if v, exists := payLetterRes["error"]; exists { // 500 error
		e := v.(map[string]any)
		err = errors.New(fmt.Sprintf("[%v]%v", e["code"], e["message"]))
		return
	}

	if code, exists := payLetterRes["code"]; exists {
		err = errors.New(fmt.Sprintf("[%v]%v", code, payLetterRes["message"]))
		return
	}

	res = ResTransactionAutoPay{
		TID:             payLetterRes["tid"].(string),
		CID:             payLetterRes["cid"].(string),
		Amount:          utils.Any2IntMust(payLetterRes["amount"]),
		BillKey:         payLetterRes["billkey"].(string),
		TransactionDate: payLetterRes["transaction_date"].(string),
	}
	return
}

func (o *PayLetter) CancelTransaction(req ReqCancelTransaction) (res ResCancelTransaction, err error) {
	payLetterRes := utils.Post[utils.M](
		cancelTransactionUrl,
		req,
		http.Header{
			"Authorization": []string{fmt.Sprintf("PLKEY %s", req.APIKey)},
			"Content-Type":  []string{"application/json"},
		},
	)

	if v, exists := payLetterRes["error"]; exists { // 500 error
		e := v.(map[string]any)
		err = errors.New(fmt.Sprintf("[%v]%v", e["code"], e["message"]))
		return
	}

	if code, exists := payLetterRes["code"]; exists {
		// 에러 발생
		err = errors.New(fmt.Sprintf("[%v]%v", code, payLetterRes["message"]))
		return
	}

	res = ResCancelTransaction{
		TID:    payLetterRes["tid"].(string),
		CID:    payLetterRes["cid"].(string),
		Amount: utils.Any2IntMust(payLetterRes["amount"]),
	}

	return
}
