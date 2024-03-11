package dto

// schema: https://sandbox.vnpayment.vn/apis/docs/huong-dan-tich-hop/

type VNPCreatePaymentUrl struct {
	BankCode  *string `json:"bankCode" validate:"omitempty"`
	Language  *string `json:"language" validate:"omitempty"`
	ReturnUrl string  `json:"returnUrl" validate:"required"`
}

type VNPReturnQuery struct {
	VnpSecureHash     string `query:"vnp_SecureHash" validate:"required"`
	VnpSecureHashType string `query:"vnp_SecureHashType"`
}

type VNPIpnQuery struct {
	VnpTxnRef         string `query:"vnp_TxnRef" validate:"required"`
	VnpResponseCode   string `query:"vnp_ResponseCode" validate:"required"`
	VnpSecureHash     string `query:"vnp_SecureHash" validate:"required"`
	VnpSecureHashType string `query:"vnp_SecureHashType" validate:"required"`
}

type VNPQuerydr struct {
	OrderId   string `query:"orderId" validate:"required"`
	TransDate string `query:"transDate" validate:"required"`
}

type VNPRefund struct {
	OrderId   string `json:"orderId" validate:"required"`
	TransDate string `json:"transDate" validate:"required"`
	Amount    int64  `json:"amount" validate:"required"`
	TransType string `json:"transType" validate:"required"`
	User      string `json:"user" validate:"required"`
}
