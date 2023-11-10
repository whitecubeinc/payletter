package payletter

import (
	"errors"
	"fmt"
	"github.com/whitecubeinc/go-utils"
	"net/http"
	"strconv"
)

type PayLetter struct {
	ClientInfo
}

func GetPayLetter(c ClientInfo) IPayLetter {
	return &PayLetter{
		ClientInfo: c,
	}
}

func (o *PayLetter) RegisterAutoPay(req ReqRegisterAutoPay) (res ResRegisterAutoPay, err error) {
	paymentData := reqPaymentData{
		PgCode:          req.PgCode,
		ClientID:        o.ClientID,
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
		CallbackUrl:     req.CallbackUrl,
		CancelUrl:       req.CancelUrl,
	}

	payLetterRes := utils.Post[utils.M](
		registerAutoPayUrl,
		paymentData,
		http.Header{
			"Authorization": []string{fmt.Sprintf("PLKEY %s", o.PaymentAPIKey)},
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
	}

	return
}

func (o *PayLetter) TransactionAutoPay(req ReqTransactionAutoPay) (res ResTransactionAutoPay, err error) {
	payLetterRes := utils.Post[utils.M](
		transactionAutoPayUrl,
		req,
		http.Header{
			"Authorization": []string{fmt.Sprintf("PLKEY %s", o.PaymentAPIKey)},
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
			"Authorization": []string{fmt.Sprintf("PLKEY %s", o.PaymentAPIKey)},
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

func (o *PayLetter) RegisterEasyPay(req ReqRegisterEasyPay) (payLetterRes ResEasyPayUI, err error) {
	req.setClientID(o.ClientID)
	req.setHashData(o.PaymentAPIKey, o.ClientID)

	payLetterRes = utils.Post[ResEasyPayUI](
		easyPayRegisterUrl,
		req,
		http.Header{
			"Authorization": []string{fmt.Sprintf("PLKEY %s", o.PaymentAPIKey)},
			"Content-Type":  []string{"application/json"},
		},
	)

	if payLetterRes.Code != nil {
		// 에러 발생
		err = errors.New(fmt.Sprintf("[%d]%s", *payLetterRes.Code, payLetterRes.Message))
		return
	}
	return
}

func (o *PayLetter) GetRegisteredEasyPayMethods(req ReqGetRegisteredEasyPayMethod) (payLetterRes ResPayLetterGetEasyPayMethods, err error) {
	params := map[string]string{
		"client_id": o.ClientID,
		"user_id":   strconv.Itoa(req.UserID),
		"req_date":  req.ReqDate,
		"hash_data": req.createHashData(o.PaymentAPIKey, o.ClientID),
	}

	payLetterRes = utils.Get[ResPayLetterGetEasyPayMethods](
		easyPayGetRegisteredMethodUrl,
		params,
		http.Header{
			"Authorization": []string{fmt.Sprintf("PLKEY %s", o.SearchAPIKey)},
			"Content-Type":  []string{"application/json"},
		},
	)
	if payLetterRes.Code != nil {
		err = errors.New(fmt.Sprintf("[%d]%s", *payLetterRes.Code, payLetterRes.Message))
		return
	}

	if payLetterRes.MethodList == nil {
		payLetterRes.MethodList = make([]EasyPayMethod, 0)
	}

	if payLetterRes.MethodCount == nil {
		payLetterRes.MethodCount = make([]EasyPayMethodCount, 0)
	}

	for idx, method := range payLetterRes.MethodList {
		switch method.PaymentMethod {
		case PgCode.CreditCard:
			method.MethodName = PayletterCardCode.ValueMap[method.MethodCode]
		case PgCode.Easybank:
			method.MethodName = PayletterBankCode[method.MethodCode]
		}
		payLetterRes.MethodList[idx] = method
	}

	return
}

func (o *PayLetter) CancelEasyPay(req ReqCancelEasyPay) (payLetterRes ResCancelEasyPay, err error) {
	req.setClientID(o.ClientID)
	req.setIPAddress(o.IpAddr)
	req.setHashData(o.ClientID, o.PaymentAPIKey)

	payLetterRes = utils.Post[ResCancelEasyPay](
		easyPayCancelUrl,
		req,
		http.Header{
			"Authorization": []string{fmt.Sprintf("PLKEY %s", o.PaymentAPIKey)},
			"Content-Type":  []string{"application/json"},
		},
	)
	if payLetterRes.Code != nil {
		err = errors.New(fmt.Sprintf("[%d]%s", *payLetterRes.Code, payLetterRes.Message))
		return
	}

	return
}

func (o *PayLetter) TransactionEasyPay(req ReqTransactionEasyPay) (payLetterRes ResEasyPayUI, err error) {
	paymentData := reqPaymentData{
		PgCode:          req.PgCode,
		ClientID:        o.ClientID,
		ServiceName:     req.ServiceName,
		UserID:          int64(req.UserID),
		UserName:        req.UserName,
		OrderNo:         req.OrderNo,
		Amount:          req.Amount,
		ProductName:     req.ProductName,
		EmailFlag:       req.EmailFlag,
		EmailAddr:       req.EmailAddr,
		AutoPayFlag:     "N",
		ReceiptFlag:     req.ReceiptFlag,
		CustomParameter: strconv.Itoa(req.CustomParameter),
		ReturnUrl:       req.ReturnUrl,
		CallbackUrl:     req.CallbackUrl,
		CancelUrl:       req.CancelUrl,
		ReqDate:         req.ReqDate,
		HashData:        req.createHashData(o.ClientID, o.PaymentAPIKey),
		BillKey:         req.BillKey,
		ReceiptType:     req.ReceiptType,
		ReceiptInfo:     req.ReceiptInfo,
		InstallMonth:    fmt.Sprintf("%02d", req.InstallMonth),
	}

	payLetterRes = utils.Post[ResEasyPayUI](
		registerAutoPayUrl,
		paymentData,
		http.Header{
			"Authorization": []string{fmt.Sprintf("PLKEY %s", o.PaymentAPIKey)},
			"Content-Type":  []string{"application/json"},
		},
	)

	if payLetterRes.Code != nil {
		err = errors.New(fmt.Sprintf("[%d]%s", *payLetterRes.Code, payLetterRes.Message))
		return
	}

	return
}
