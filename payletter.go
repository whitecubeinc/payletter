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
		reqTransactionAutoPay{
			ClientInfo:            o.ClientInfo,
			ReqTransactionAutoPay: req,
		},
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
	cancelData := reqCancelTransaction{
		ClientInfo:           o.ClientInfo,
		ReqCancelTransaction: req,
	}

	apiKey := o.ClientInfo.PaymentAPIKey
	if PgCode.IsNaverCode(cancelData.ReqCancelTransaction.PgCode) { // 네이버페이는 client id 와 api key 가 다름
		apiKey = cancelData.ReqCancelTransaction.NaverAPIKey
		cancelData.ClientInfo.ClientID = cancelData.ReqCancelTransaction.NaverAPIClientId
	}

	payLetterRes := utils.Post[utils.M](
		cancelTransactionUrl,
		cancelData,
		http.Header{
			"Authorization": []string{fmt.Sprintf("PLKEY %s", apiKey)},
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

func (o *PayLetter) PartialCancelTransaction(req ReqPartialCancelTransaction) (res ResPartialCancelTransaction, err error) {
	cancelData := reqPartialCancelTransaction{
		ClientInfo:                  o.ClientInfo,
		ReqPartialCancelTransaction: req,
	}

	apiKey := o.ClientInfo.PaymentAPIKey
	if PgCode.IsNaverCode(cancelData.ReqPartialCancelTransaction.PgCode) { // 네이버페이는 client id 와 api key 가 다름
		apiKey = cancelData.ReqPartialCancelTransaction.NaverAPIKey
		cancelData.ClientInfo.ClientID = cancelData.ReqPartialCancelTransaction.NaverAPIClientId
	}

	payLetterRes := utils.Post[utils.M](
		partialCancelTransactionUrl,
		cancelData,
		http.Header{
			"Authorization": []string{fmt.Sprintf("PLKEY %s", apiKey)},
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

	res = ResPartialCancelTransaction{
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
			method.MethodName = CardCode.ValueMap[method.MethodCode]
		case PgCode.EasyBank:
			method.MethodName = BankCode[method.MethodCode]
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
		easyPayTransactionUrl,
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

func (o *PayLetter) TransactionNormalPay(req ReqTransactionNormalPay) (payLetterRes ResTransactionNormalPay, err error) {
	paymentData := reqPaymentData{
		PgCode:          req.PgCode,
		ServiceName:     req.ServiceName,
		ClientID:        o.ClientID,
		UserID:          int64(req.UserID),
		UserName:        req.UserName,
		OrderNo:         req.OrderNo,
		Amount:          req.Amount,
		ProductName:     req.ProductName,
		EmailFlag:       req.EmailFlag,
		EmailAddr:       req.EmailAddr,
		AutoPayFlag:     "N",
		CustomParameter: strconv.Itoa(req.CustomParameter),
		ReturnUrl:       req.ReturnUrl,
		CallbackUrl:     req.CallbackUrl,
		CancelUrl:       req.CancelUrl,
	}

	apiKey := o.PaymentAPIKey
	if PgCode.IsNaverCode(req.PgCode) { // 네이버페이는 client id 와 api key 가 다름
		apiKey = req.NaverAPIKey
		paymentData.ClientID = req.NaverAPIClientId
	}

	payLetterRes = utils.Post[ResTransactionNormalPay](
		normalTransactionUrl,
		paymentData,
		http.Header{
			"Authorization": []string{fmt.Sprintf("PLKEY %s", apiKey)},
			"Content-Type":  []string{"application/json"},
		},
	)

	if payLetterRes.Code != nil {
		err = errors.New(fmt.Sprintf("[%d]%s", *payLetterRes.Code, payLetterRes.Message))
		return
	}
	return
}

func (o *PayLetter) GetTransactionList(req ReqGetTransactionList) (res ResGetTransactionList, err error) {
	switch req.DateType {
	case TransactionDateType.Transaction, TransactionDateType.Settle:
	default:
		err = errors.New("유효하지 않은 date type")
	}

	reqParam := map[string]string{
		"date":      req.Date,
		"date_type": req.DateType,
		"pgcode":    req.PgCode,
		"client_id": o.ClientID,
	}

	apiKey := o.SearchAPIKey
	if PgCode.IsNaverCode(req.PgCode) { // 네이버페이는 client id 와 api key 가 다름
		apiKey = req.NaverAPISearchKey
		reqParam["client_id"] = req.NaverAPIClientID
	}

	res = utils.Get[ResGetTransactionList](
		getTransactionListUrl,
		reqParam,
		http.Header{
			"Authorization": []string{fmt.Sprintf("PLKEY %s", apiKey)},
			"Content-Type":  []string{"application/json"},
		},
	)

	if res.Code != nil {
		err = errors.New(fmt.Sprintf("[%d]%s", *res.Code, *res.Message))
	}

	return
}
