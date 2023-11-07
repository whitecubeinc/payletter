package payletter

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type pgCode struct {
	CreditCard string
}

type IPayLetter interface {
	// RegisterAutoPay 자동 결제 수단 등록
	RegisterAutoPay(req ReqRegisterAutoPay) (res ResRegisterAutoPay, err error)
	// TransactionAutoPay 자동 결제 수단으로 결제
	TransactionAutoPay(req ReqTransactionAutoPay) (res ResTransactionAutoPay, err error)
	// CancelTransaction 결제 취소
	CancelTransaction(req ReqCancelTransaction) (res ResCancelTransaction, err error)
}

type ClientInfo struct {
	APIKey   string `json:"-"`
	ClientID string `json:"client_id"`
	IpAddr   string `json:"ip_addr"`
}

type ReqRegisterAutoPay struct {
	ClientInfo
	PgCode           string
	ServiceName      string
	UserID           int64
	UserName         string
	OrderNo          string
	Amount           int
	ProductName      string
	CustomParameter  string
	ReturnUrl        string // POST 결제 성공 response = ResPaymentData
	CancelUrl        string // GET 결제 중간에 취소
	CallbackEndpoint string // POST
}

type ResRegisterAutoPay struct {
	OnlineUrl string // PC 환경 결제 창 호출 URL
	MobileUrl string // 모바일 환경 결제 창 호출 URL
	Token     string // 결제 인증 토큰
}

type ReqTransactionAutoPay struct {
	ClientInfo
	PgCode      string `json:"pgcode"`
	ServiceName string `json:"service_name"`
	UserID      int64  `json:"user_id"`
	UserName    string `json:"user_name"`
	OrderNo     string `json:"order_no"`
	Amount      int    `json:"amount"`
	ProductName string `json:"product_name"`
	BillKey     string `json:"billkey"`
}

type ResTransactionAutoPay struct {
	TID             string `json:"tid"`
	CID             string `json:"cid"`
	Amount          int    `json:"amount"`
	BillKey         string `json:"billkey"` // 자동결제 재결제용 키
	TransactionDate string `json:"transaction_date"`
}

type reqPaymentData struct {
	PgCode          string `json:"pgcode"`
	ClientID        string `json:"client_id"`
	ServiceName     string `json:"service_name"`
	UserID          int64  `json:"user_id"`
	UserName        string `json:"user_name"`
	OrderNo         string `json:"order_no"`
	Amount          int    `json:"amount"`
	ProductName     string `json:"product_name"`
	EmailFlag       string `json:"email_flag"`
	AutoPayFlag     string `json:"autopay_flag"`
	ReceiptFlag     string `json:"receipt_flag"`
	CustomParameter string `json:"custom_parameter"`
	ReturnUrl       string `json:"return_url"`
	CallbackUrl     string `json:"callback_url"`
	CancelUrl       string `json:"cancel_url"`
}

type ResPaymentData struct {
	Code                 string `json:"code" form:"code"`
	Message              string `json:"message" form:"message"`
	UserID               string `json:"user_id" form:"user_id"`
	UserName             string `json:"user_name" form:"user_name"`
	Amount               int    `json:"amount" form:"amount"`
	TaxAmount            int    `json:"tax_amount" form:"tax_amount"`
	TaxFreeAmount        int    `json:"taxfree_amount" form:"taxfree_amount"`
	Tid                  string `json:"tid" form:"tid"`
	Cid                  string `json:"cid" form:"cid"`
	OrderNo              string `json:"order_no" form:"order_no"`
	ServiceName          string `json:"service_name" form:"service_name"`
	ProductName          string `json:"product_name" form:"product_name"`
	CustomParameter      string `json:"custom_parameter" form:"custom_parameter"`
	TransactionDate      string `json:"transaction_date" form:"transaction_date"`
	PayInfo              string `json:"pay_info" form:"pay_info"`
	PgCode               string `json:"pgcode" form:"pgcode"`
	DomesticFlag         string `json:"domestic_flag" form:"domestic_flag"`
	BillKey              string `json:"billkey" form:"billkey"`
	CardInfo             string `json:"card_info" form:"card_info"`
	PayHash              string `json:"payhash" form:"payhash"`
	DisposableCupDeposit int    `json:"disposable_cup_deposit" form:"disposable_cup_deposit"`
	CashReceipt          struct {
		Code      string `json:"code" form:"code"`
		Message   string `json:"message" form:"message"`
		Cid       string `json:"cid" form:"cid"`
		DealNo    string `json:"deal_no" form:"deal_no"`
		IssueType string `json:"issue_type" form:"issue_type"`
		PayerSid  string `json:"payer_sid" form:"payer_sid"`
		Type      string `json:"type" form:"type"`
	} `json:"cash_receipt" form:"cash_receipt"`
}

func (o *ResPaymentData) Validate(paymentAPIKey string) (err error) {
	pgHashText := fmt.Sprintf("%s%d%s%s", o.UserID, o.Amount, o.Tid, paymentAPIKey)
	h := sha256.Sum256([]byte(pgHashText))
	pgHash := strings.ToUpper(hex.EncodeToString(h[:]))

	if pgHash != o.PayHash {
		err = errors.New("pgHash 검증 실패")
	}

	return
}

type ReqCancelTransaction struct {
	ClientInfo
	PgCode string `json:"pgcode"`
	UserID int64  `json:"user_id"`
	TID    string `json:"tid"`
}

type ResCancelTransaction struct {
	TID    string `json:"tid"`
	CID    string `json:"cid"`
	Amount int    `json:"amount"`
}
