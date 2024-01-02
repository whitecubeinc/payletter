package payletter

import (
	"errors"
	"fmt"
	"github.com/whitecubeinc/go-utils"
	"net/http"
	"strconv"
	"time"
)

type MockPayLetter struct {
	ClientInfo
	Success bool
}

// GetSuccessMockPayLetter 무조건 결제 성공하는 Mock pay letter
func GetSuccessMockPayLetter(c ClientInfo) IPayLetter {
	return &MockPayLetter{
		ClientInfo: c,
		Success:    true,
	}
}

// GetFailMockPayLetter 무조건 결제 실패하는 Mock pay letter
func GetFailMockPayLetter(c ClientInfo) IPayLetter {
	return &MockPayLetter{
		ClientInfo: c,
		Success:    false,
	}
}

func (o *MockPayLetter) RegisterAutoPay(req ReqRegisterAutoPay) (res ResRegisterAutoPay, err error) {
	paymentData := reqPaymentData{
		PgCode:          req.PgCode,
		ClientID:        o.ClientID,
		ServiceName:     req.ServiceName,
		UserID:          req.UserID,
		UserName:        req.UserName,
		OrderNo:         req.OrderNo,
		Amount:          0,
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

func (o *MockPayLetter) TransactionAutoPay(req ReqTransactionAutoPay) (res ResTransactionAutoPay, err error) {
	if o.Success {
		res = ResTransactionAutoPay{
			TID:             "tid",
			CID:             "cid",
			Amount:          req.Amount,
			BillKey:         req.BillKey,
			TransactionDate: time.Now().String(),
		}
	} else {
		err = errors.New("fake mock pay letter")
	}

	return
}

func (o *MockPayLetter) CancelTransaction(req ReqCancelTransaction) (res ResCancelTransaction, err error) {
	if o.Success {
		res.TID = req.TID
	} else {
		err = errors.New("fake mock pay letter")
	}
	return
}

func (o *MockPayLetter) PartialCancelTransaction(req ReqPartialCancelTransaction) (res ResPartialCancelTransaction, err error) {
	if o.Success {
		res.TID = req.TID
	} else {
		err = errors.New("fake mock pay letter")
	}
	return
}

func (o *MockPayLetter) RegisterEasyPay(req ReqRegisterEasyPay) (res ResEasyPayUI, err error) {
	req.setClientID(o.ClientID)
	req.setHashData(o.PaymentAPIKey, o.ClientID)

	payletterRes := utils.Post[ResEasyPayUI](
		easyPayRegisterTestUrl,
		req,
		http.Header{
			"Authorization": []string{fmt.Sprintf("PLKEY %s", o.PaymentAPIKey)},
			"Content-Type":  []string{"application/json"},
		},
	)

	if payletterRes.Code != nil {
		// 에러 발생
		err = errors.New(fmt.Sprintf("[%d]%s", *payletterRes.Code, payletterRes.Message))
		return
	}

	res = payletterRes
	return
}

func (o *MockPayLetter) GetRegisteredEasyPayMethods(req ReqGetRegisteredEasyPayMethod) (res ResPayLetterGetEasyPayMethods, err error) {
	params := map[string]string{
		"client_id": o.ClientID,
		"user_id":   strconv.Itoa(req.UserID),
		"req_date":  req.ReqDate,
		"hash_data": req.createHashData(o.PaymentAPIKey, o.ClientID),
	}

	payletterRes := utils.Get[ResPayLetterGetEasyPayMethods](
		easyPayGetRegisteredMethodTestUrl,
		params,
		http.Header{
			"Authorization": []string{fmt.Sprintf("PLKEY %s", o.SearchAPIKey)},
			"Content-Type":  []string{"application/json"},
		},
	)
	if payletterRes.Code != nil {
		err = errors.New(fmt.Sprintf("[%d]%s", *payletterRes.Code, payletterRes.Message))
		return
	}

	if payletterRes.MethodList == nil {
		payletterRes.MethodList = make([]EasyPayMethod, 0)
	}

	if payletterRes.MethodCount == nil {
		payletterRes.MethodCount = make([]EasyPayMethodCount, 0)
	}

	for idx, method := range payletterRes.MethodList {
		switch method.PaymentMethod {
		case PgCode.CreditCard:
			method.MethodName = CardCode.ValueMap[method.MethodCode]
		case PgCode.EasyBank:
			method.MethodName = BankCode[method.MethodCode]
		}
		payletterRes.MethodList[idx] = method
	}
	res = payletterRes

	return
}

func (o *MockPayLetter) CancelEasyPay(req ReqCancelEasyPay) (payLetterRes ResCancelEasyPay, err error) {
	req.setClientID(o.ClientID)
	req.setIPAddress(o.IpAddr)
	req.setHashData(o.ClientID, o.PaymentAPIKey)

	payLetterRes = utils.Post[ResCancelEasyPay](
		easyPayCancelTestUrl,
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

func (o *MockPayLetter) TransactionEasyPay(req ReqTransactionEasyPay) (payLetterRes ResEasyPayUI, err error) {
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
		easyPayTransactionTestUrl,
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

func (o *MockPayLetter) TransactionNormalPay(req ReqTransactionNormalPay) (payLetterRes ResTransactionNormalPay, err error) {
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
		CustomParameter: strconv.Itoa(req.CustomParameter),
		ReturnUrl:       req.ReturnUrl,
		CallbackUrl:     req.CallbackUrl,
		CancelUrl:       req.CancelUrl,
	}

	apiKey := o.PaymentAPIKey
	if req.PgCode == PgCode.NaverPay {
		paymentData.ClientID = req.NaverAPIClientId
		apiKey = req.NaverAPIKey
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
