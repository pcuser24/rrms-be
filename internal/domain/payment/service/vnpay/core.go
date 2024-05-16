package vnpay

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/payment/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

var (
	ErrUnauthorizedUser = errors.New("unauthorized user")
	ErrBadStatusPayment = errors.New("bad status payment")
)

func (s *VnPayService) CreatePaymentUrl(ipAddr string, userId uuid.UUID, paymentId int64, data *dto.VNPCreatePaymentUrl) (string, error) {
	// Get payment by id
	payment, err := s.repo.GetPaymentById(context.Background(), paymentId)
	if err != nil {
		return "", err
	}
	if payment.UserID != userId {
		return "", ErrUnauthorizedUser
	}
	if payment.Status == "COMPLETED" {
		return "", ErrBadStatusPayment
	}

	tz, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return "", err
	}
	date := time.Now()
	createDate := date.In(tz).Format("20060102150405") // YYYYMMDDHHMMSS
	orderId := date.In(tz).Format("02150405")          // DDHHmmss

	if data.Language == nil || *data.Language == "" {
		data.Language = new(string)
		*data.Language = "vn"
	}

	currCode := "VND"

	vnpParams := make(map[string]string)
	vnpParams["vnp_Version"] = "2.1.0"
	vnpParams["vnp_Command"] = "pay"
	vnpParams["vnp_TmnCode"] = s.vnpTmnCode
	vnpParams["vnp_Locale"] = *data.Language
	vnpParams["vnp_CurrCode"] = currCode
	vnpParams["vnp_TxnRef"] = orderId
	vnpParams["vnp_OrderInfo"] = fmt.Sprintf("[%d]%s", paymentId, payment.OrderInfo)
	vnpParams["vnp_OrderType"] = "other"
	vnpParams["vnp_Amount"] = strconv.FormatInt(payment.Amount*100, 10)
	vnpParams["vnp_ReturnUrl"] = data.ReturnUrl
	vnpParams["vnp_IpAddr"] = ipAddr
	vnpParams["vnp_CreateDate"] = createDate

	if data.BankCode != nil && *data.BankCode != "" {
		vnpParams["vnp_BankCode"] = *data.BankCode
	}

	vnpParams = sortObject(vnpParams)

	// stringify the params
	signData := stringify(vnpParams)
	log.Println("signData", signData)

	// hash the params
	h := hmac.New(sha512.New, []byte(s.vnpHashSecret))
	h.Write([]byte(signData))
	signed := hex.EncodeToString(h.Sum(nil))
	vnpParams["vnp_SecureHash"] = signed

	vnpUrl := s.vnpUrl + "?" + stringify(vnpParams)

	return vnpUrl, nil
}

var (
	ErrInvalidHash = errors.New("invalid hash")
)

func (s *VnPayService) Return(query map[string]string) error {
	secureHash := query["vnp_SecureHash"]
	delete(query, "vnp_SecureHash")
	delete(query, "vnp_SecureHashType")

	vnpParams := sortObject(query)

	signData := stringify(vnpParams)
	h := hmac.New(sha512.New, []byte(s.vnpHashSecret))
	h.Write([]byte(signData))
	signed := hex.EncodeToString(h.Sum(nil))
	if secureHash != signed {
		return ErrInvalidHash
	}

	// TODO: compare with data in database
	// get paymentId from query["vnp_OrderInfo"]. query["vnp_OrderInfo"] is in this format "[paymentId][paymentType]orderInfo"
	var paymentId int64
	orderInfo := query["vnp_OrderInfo"]
	end := strings.Index(orderInfo, "]")
	if end != -1 {
		id, err := strconv.ParseInt(orderInfo[1:end], 10, 64)
		if err != nil {
			return err
		}
		paymentId = id
	}

	paymentUpdatePayload := dto.UpdatePayment{
		ID:      paymentId,
		OrderId: types.Ptr(query["vnp_TxnRef"]),
		Status: types.Ptr(database.PAYMENTSTATUS(
			utils.Ternary(
				slices.Contains([]string{"00", "07"}, query["vnp_ResponseCode"]),
				database.PAYMENTSTATUSSUCCESS, database.PAYMENTSTATUSFAILED,
			),
		)),
	}

	err := s.HandleReturn(&paymentUpdatePayload, orderInfo[end+1:])
	if err != nil {
		return err
	}

	return s.repo.UpdatePayment(context.Background(), &paymentUpdatePayload)
}

type IpnReturn struct {
	RspCode string `json:"RspCode"`
	Message string `json:"Message"`
}

func (s *VnPayService) Ipn(query map[string]string) IpnReturn {
	secureHash := query["vnp_SecureHash"]

	// orderId := query["vnp_TxnRef"]
	rspCode := query["vnp_ResponseCode"]

	delete(query, "vnp_SecureHash")
	delete(query, "vnp_SecureHashType")

	vnpParams := sortObject(query)
	signData := stringify(vnpParams)
	h := hmac.New(sha512.New, []byte(s.vnpHashSecret))
	h.Write([]byte(signData))
	signed := hex.EncodeToString(h.Sum(nil))

	// Giả sử "0" là trạng thái khởi tạo giao dịch, chưa có IPN. Trạng thái này được lưu khi yêu cầu thanh toán chuyển hướng sang Cổng thanh toán VNPAY tại đầu khởi tạo đơn hàng.
	//paymentStatus := "1"; // Giả sử "1" là trạng thái thành công bạn cập nhật sau IPN được gọi và trả kết quả về nó
	//paymentStatus := "2"; // Giả sử "2" là trạng thái thất bại bạn cập nhật sau IPN được gọi và trả kết quả về nó
	paymentStatus := "0"
	// Mã đơn hàng "giá trị của vnp_TxnRef" VNPAY phản hồi tồn tại trong CSDL của bạn
	checkOrderId := true
	// Kiểm tra số tiền "giá trị của vnp_Amout/100" trùng khớp với số tiền của đơn hàng trong CSDL của bạn
	checkAmount := true

	// verify checksum
	if secureHash == signed {
		if checkOrderId {
			if checkAmount {
				// verify transaction status before updating payment status
				if paymentStatus == "0" {
					if rspCode == "00" {
						// success
						// paymentStatus = "1"
						// TODO: update payment status to SUCCESS in database
						return IpnReturn{
							RspCode: "00",
							Message: "Success",
						}
					} else {
						// fail
						// paymentStatus = "2"
						// TODO: update payment status to FAILED in database
						return IpnReturn{
							RspCode: "00",
							Message: "Success",
						}
					}
				} else {
					return IpnReturn{
						RspCode: "02",
						Message: "This order has been updated to the payment status",
					}
				}
			} else {
				return IpnReturn{
					RspCode: "04",
					Message: "Amount invalid",
				}
			}
		} else {
			return IpnReturn{
				RspCode: "01",
				Message: "Order not found",
			}
		}
	} else {
		return IpnReturn{
			RspCode: "97",
			Message: "Checksum failed",
		}
	}
}

func (s *VnPayService) Querydr(ipAddr string, d *dto.VNPQuerydr) (*http.Response, error) {
	tz, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return nil, err
	}

	date := time.Now()
	vnpRequestId := date.In(tz).Format("150405") // HHmmss
	vnpVersion := "2.1.0"
	vnpCommand := "querydr"
	vnpTxnRef := d.OrderId
	vnpTransactionDate := d.TransDate
	vnpOrderInfo := "Truy van GD ma:" + vnpTxnRef
	vnpIpAddr := ipAddr
	// currCode := "VND"
	vnpCreateDate := date.In(tz).Format("20060102150405") // YYYYMMDDHHMMSS
	vnpTmnCode := s.vnpTmnCode

	data := strings.Join([]string{vnpRequestId, vnpVersion, vnpCommand, vnpTmnCode, vnpTxnRef, vnpTransactionDate, vnpCreateDate, vnpIpAddr, vnpOrderInfo}, "|")

	h := hmac.New(sha512.New, []byte(s.vnpHashSecret))
	h.Write([]byte(data))
	vnpSecureHash := hex.EncodeToString(h.Sum(nil))

	return sendHttpRequest(s.vnpApi, http.MethodPost, map[string]string{
		"vnp_RequestId":       vnpRequestId,
		"vnp_Version":         vnpVersion,
		"vnp_Command":         vnpCommand,
		"vnp_TmnCode":         vnpTmnCode,
		"vnp_TxnRef":          vnpTxnRef,
		"vnp_OrderInfo":       vnpOrderInfo,
		"vnp_TransactionDate": vnpTransactionDate,
		"vnp_CreateDate":      vnpCreateDate,
		"vnp_IpAddr":          vnpIpAddr,
		"vnp_SecureHash":      vnpSecureHash,
	})
}

func (s *VnPayService) Refund(ipAddr string, d *dto.VNPRefund) (*http.Response, error) {
	tz, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return nil, err
	}

	date := time.Now()
	vnpTxnRef := d.OrderId
	vnpTransactionDate := d.TransDate
	vnpAmount := strconv.FormatInt(d.Amount*100, 10)
	vnpTransactionType := d.TransType
	vnpCreateBy := d.User
	// currCode := "vn"
	vnpRequestId := date.In(tz).Format("150405") // HHmmss
	vnpVersion := "2.1.0"
	vnpCommand := "refund"
	vnpOrderInfo := "Hoan tien GD ma:" + vnpTxnRef
	vnpCreateDate := date.In(tz).Format("20060102150405") // YYYYMMDDHHMMSS
	vnpTransactionNo := "0"
	vnpIpAddr := ipAddr
	vnpTmnCode := s.vnpTmnCode

	data := strings.Join([]string{vnpRequestId, vnpVersion, vnpCommand, vnpTmnCode, vnpTransactionType, vnpTxnRef, vnpAmount, vnpTransactionNo, vnpTransactionDate, vnpCreateBy, vnpCreateDate, vnpIpAddr, vnpOrderInfo}, "|")
	h := hmac.New(sha512.New, []byte(s.vnpHashSecret))
	h.Write([]byte(data))
	vnpSecureHash := hex.EncodeToString(h.Sum(nil))

	return sendHttpRequest(s.vnpApi, http.MethodPost, map[string]string{
		"vnp_RequestId":       vnpRequestId,
		"vnp_Version":         vnpVersion,
		"vnp_Command":         vnpCommand,
		"vnp_TmnCode":         vnpTmnCode,
		"vnp_TransactionType": vnpTransactionType,
		"vnp_TxnRef":          vnpTxnRef,
		"vnp_Amount":          vnpAmount,
		"vnp_TransactionNo":   vnpTransactionNo,
		"vnp_CreateBy":        vnpCreateBy,
		"vnp_OrderInfo":       vnpOrderInfo,
		"vnp_TransactionDate": vnpTransactionDate,
		"vnp_CreateDate":      vnpCreateDate,
		"vnp_IpAddr":          vnpIpAddr,
		"vnp_SecureHash":      vnpSecureHash,
	})
}
