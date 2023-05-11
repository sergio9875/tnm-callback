package process

import (
	"malawi-callback/models"
	"time"
)

func (c *Controller) mapPaymentGatewayRequest(msgBody *models.IncomingRequest, statusCode int) *models.PaymentGatewayRequest {

	tm := time.Now().Format("02-Jan-2006 15:04:05")

	paymentDetails := models.PaymentDetails{
		Code:              statusCode,
		StatusCode:        msgBody.ResultCode,
		Explanation:       msgBody.ResultDesc,
		Paymentreference:  msgBody.TransactionId,
		Sequenceid:        msgBody.ExternalRef,
		Terminalmno:       "TnmMalawiPayment",
		Terminalsettings:  nil,
		Amount:            "",
		Currency:          "",
		Msisdn:            "",
		Approvalreference: msgBody.TransactionId,
		Timestamp:         tm,
		Customerdetails:   nil,
		Mbtid:             msgBody.ExternalRef,
	}
	return &models.PaymentGatewayRequest{
		Paymentmethod:  "mobile",
		Transunq:       nil,
		Paymenttype:    "Process",
		Paymentdetails: paymentDetails,
	}

}
