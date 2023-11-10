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
	RegisterEasyPay(req ReqRegisterEasyPay) (res ResEasyPayUI, err error)
	// GetRegisteredEasyPayMethods 간편결제 등록한 결제 수단 목록 조회
	GetRegisteredEasyPayMethods(req ReqGetRegisteredEasyPayMethod) (res ResPayLetterGetEasyPayMethods, err error)
	// CancelEasyPay 간편결제 취소
	CancelEasyPay(req ReqCancelEasyPay) (res ResCancelEasyPay, err error)
	// TransactionEasyPay 간편결제 수단으로 결제
	TransactionEasyPay(req ReqTransactionEasyPay) (res ResEasyPayUI, err error)
	// TransactionNormalPay 일반 페이레터 결제
	TransactionNormalPay(req ReqTransactionNormalPay) (res ResTransactionNormalPay, err error)
}

type ClientInfo struct {
	PaymentAPIKey string `json:"-"` // PAYMENT KEY
	SearchAPIKey  string `json:"-"` // SEARCH KEY
	ClientID      string `json:"client_id"`
	IpAddr        string `json:"ip_addr"`
}

type ReqRegisterAutoPay struct {
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
	//Token     string // 결제 인증 토큰
}

type ReqTransactionAutoPay struct {
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

type CommonTransactionData struct {
	PgCode          string
	UserID          int
	UserName        string
	ServiceName     string
	OrderNo         string
	Amount          int
	ProductName     string
	EmailFlag       string
	EmailAddr       string
	CustomParameter int
	ReturnUrl       string
	CallbackUrl     string
	CancelUrl       string
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
	EmailAddr       string `json:"email_addr"`
	AutoPayFlag     string `json:"autopay_flag"`
	ReceiptFlag     string `json:"receipt_flag"`
	CustomParameter string `json:"custom_parameter"`
	ReturnUrl       string `json:"return_url"`
	CallbackUrl     string `json:"callback_url"`
	CancelUrl       string `json:"cancel_url"`
	ReqDate         string `json:"req_date,omitempty"`
	HashData        string `json:"hash_data,omitempty"`
	BillKey         string `json:"billkey,omitempty"`
	ReceiptType     string `json:"receipt_type,omitempty"`
	ReceiptInfo     string `json:"receipt_info,omitempty"`
	InstallMonth    string `json:"install_month,omitempty"`
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
	InstallMonth         string `json:"install_month"`
	CardCode             string `json:"card_code"`
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

func (o *ResPaymentData) ReplacePayInfo() {
	if o.PgCode == PgCode.EasyBank {
		o.PayInfo = BankCode[o.CardCode]
	} else {
		o.PayInfo = CardCode.ValueMap[o.CardCode]
	}
}

type ReqCancelTransaction struct {
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
	ClientID      string `json:"client_id"`
	UserID        int    `json:"user_id"`
	ServiceName   string `json:"service_name"`
	PaymentMethod string `json:"payment_method"`
	ReturnUrl     string `json:"return_url"`
	CancelUrl     string `json:"cancel_url"`
	ReqDate       string `json:"req_date"`
	HashData      string `json:"hash_data"`
}

func (o *ReqRegisterEasyPay) setClientID(clientID string) {
	o.ClientID = clientID
}

func (o *ReqRegisterEasyPay) setHashData(apiKey string, clientId string) {
	originHashString := fmt.Sprintf("%s%d%s%s", clientId, o.UserID, o.ReqDate, apiKey)
	h := sha256.Sum256([]byte(originHashString))

	o.HashData = hex.EncodeToString(h[:])
}

type ResEasyPayUI struct {
	Token       *string `json:"token"`
	RedirectUrl *string `json:"redirect_url"`
	Code        *int    `json:"code,omitempty"`
	Message     string  `json:"message,omitempty"`
	OrderNo     string  `json:"order_no,omitempty"` // 간편결제에서만 사용하는 field
}

type ReqGetRegisteredEasyPayMethod struct {
	UserID  int    `json:"user_id"`
	ReqDate string `json:"req_date"`
}

func (o *ReqGetRegisteredEasyPayMethod) createHashData(apiKey string, clientId string) string {
	originHashString := fmt.Sprintf("%s%d%s%s", clientId, o.UserID, o.ReqDate, apiKey)
	h := sha256.Sum256([]byte(originHashString))

	ownHash := hex.EncodeToString(h[:])

	return ownHash
}

type ResPayLetterGetEasyPayMethods struct {
	TotalCount       int                  `json:"total_count"`
	JoinDate         string               `json:"join_date"`
	MethodCount      []EasyPayMethodCount `json:"method_count"`
	MethodList       []EasyPayMethod      `json:"method_list"`
	PasswordSkipFlag string               `json:"password_skip_flag"`
	Code             *int                 `json:"code,omitempty"`
	Message          string               `json:"message,omitempty"`
}

type EasyPayMethodCount struct {
	PaymentMethod string `json:"paymentMethod"`
	Count         int    `json:"count"`
}

type EasyPayMethod struct {
	// payletter response
	PaymentMethod         string `json:"payment_method"`
	BillKey               string `json:"billkey"`
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

	// method code 에 따른 method name
	MethodName string `json:"method_name"`
}

type ReqCancelEasyPay struct {
	ClientID string `json:"client_id"`
	UserID   int    `json:"user_id"`
	Tid      string `json:"tid"`
	Amount   int    `json:"amount"`
	ReqDate  string `json:"req_date"`
	HashData string `json:"hash_data"`
	IpAddr   string `json:"ip_addr"`
}

func (o *ReqCancelEasyPay) setClientID(clientID string) {
	o.ClientID = clientID
}

func (o *ReqCancelEasyPay) setIPAddress(ipAddr string) {
	o.IpAddr = ipAddr
}

func (o *ReqCancelEasyPay) setHashData(clientId, apiKey string) {
	originHashString := fmt.Sprintf("%s%s%d%s%s", clientId, o.Tid, o.Amount, o.ReqDate, apiKey)
	h := sha256.Sum256([]byte(originHashString))

	o.HashData = hex.EncodeToString(h[:])
}

type ResCancelEasyPay struct {
	Tid        string `json:"tid"`
	Cid        string `json:"cid"`
	Amount     int    `json:"amount"`
	CancelDate string `json:"cancel_date"`
	Code       *int   `json:"code"`
	Message    string `json:"message"`
}

type ReqTransactionEasyPay struct {
	CommonTransactionData
	ReqDate      string
	BillKey      string
	ReceiptFlag  string
	ReceiptType  string
	ReceiptInfo  string
	InstallMonth int
}

func (o *ReqTransactionEasyPay) createHashData(clientId, apiKey string) string {
	originHashString := fmt.Sprintf("%s%d%s%s", clientId, o.UserID, o.ReqDate, apiKey)
	h := sha256.Sum256([]byte(originHashString))

	return hex.EncodeToString(h[:])
}

type ReqTransactionNormalPay struct {
	CommonTransactionData
	NaverAPIClientId string
	NaverAPIKey      string
}

type ResTransactionNormalPay struct {
	MobileUrl string `json:"mobile_url"`
	OnlineUrl string `json:"online_url"`
	OrderNo   string `json:"order_no,omitempty"`
	Token     int64  `json:"token"`
	Code      *int   `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
}
