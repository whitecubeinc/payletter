package payletter

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type IPayLetter interface {
	// RegisterAutoPay 자동 결제 수단 등록
	RegisterAutoPay(req ReqRegisterAutoPay) (res ResRegisterAutoPay, err error)
	// TransactionAutoPay 자동 결제 수단으로 결제
	TransactionAutoPay(req ReqTransactionAutoPay) (res ResTransactionAutoPay, err error)
	// CancelTransaction 결제 취소
	CancelTransaction(req ReqCancelTransaction) (res ResCancelTransaction, err error)
	// RegisterEasyPay 간편결제 결제 수단 등록
	RegisterEasyPay(req ReqRegisterEasyPay) (res ResRegisterEasyPay, err error)
	// GetRegisteredEasyPayMethods 간편결제 등록한 결제 수단 목록 조회
	GetRegisteredEasyPayMethods(req ReqGetRegisteredEasyPayMethod) (res ResPayletterGetEasyPayMethods, err error)
}

type ClientInfo struct {
	APIKey   string `json:"-"`
	ClientID string `json:"client_id"`
	IpAddr   string `json:"ip_addr"`
}

type ReqRegisterAutoPay struct {
	ClientInfo
	PgCode          string
	ServiceName     string
	UserID          int64
	UserName        string
	OrderNo         string
	Amount          int
	ProductName     string
	CustomParameter string
	ReturnUrl       string // POST 결제 성공 response = ResPaymentData
	CancelUrl       string // GET 결제 중간에 취소
	CallbackUrl     string // POST
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

type ReqRegisterEasyPay struct {
	ClientInfo
	UserID        int    `json:"user_id"`
	ServiceName   string `json:"service_name"`
	PaymentMethod string `json:"payment_method"`
	ReturnUrl     string `json:"return_url"`
	CancelUrl     string `json:"cancel_url"`
	ReqDate       string `json:"req_date"`
	HashData      string `json:"hash_data"`
}

type ResRegisterEasyPay struct {
	Token       *string `json:"token"`
	RedirectUrl *string `json:"redirect_url"`
	Code        *int    `json:"code"`
	Message     string  `json:"message"`
}

type ReqGetRegisteredEasyPayMethod struct {
	ClientInfo
	UserID   int    `json:"user_id"`
	ReqDate  string `json:"req_date"`
	HashData string `json:"hash_data"`
}

type ResPayletterGetEasyPayMethods struct {
	TotalCount       int                    `json:"total_count"`
	JoinDate         string                 `json:"join_date"`
	MethodCount      []PayletterMethodCount `json:"method_count"`
	MethodList       []PayletterMethod      `json:"method_list"`
	PasswordSkipFlag string                 `json:"password_skip_flag"`
	Code             *string                `json:",omitempty"`
	Message          string                 `json:",omitempty"`
}

type PayletterMethodCount struct {
	PaymentMethod string
	Count         int
}

type PayletterMethod struct {
	// payletter response
	PaymentMethod         string `json:"payment_method"`
	Billkey               string `json:"billkey"`
	AliasName             string `json:"alias_name"`
	FavoriteFlag          string `json:"favorite_flag"`
	MethodRegDate         string `json:"method_reg_date"`
	MethodInfo            string `json:"method_info"`
	MethodCode            string `json:"method_code"`
	MethodImgUrl          string `json:"method_img_url"`
	CardTypeCode          string `json:"card_type_code"`
	InstallmentUseFlag    string `json:"installment_use_flag"`
	MinInstallmentAmount  int    `json:"min_installment_amount"`
	InstallmentMonths     string `json:"installment_months"`
	FreeInstallmentMonths string `json:"free_installment_months"`
	ProductCode           string `json:"product_code"`
	ProductName           string `json:"product_name"`
	LastTranDate          string `json:"last_tran_date"`

	// method code에 따른 method name
	MethodName string
}
