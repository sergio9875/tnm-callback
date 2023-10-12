package process

import (
	"tnm-malawi/connectors/callback/models"
)

func (c *Controller) mapPaymentGatewayRequest(msgBody *models.IncomingRequest, statusCode int) *models.PaymentGatewayRequest {

	paymentDetails := models.PaymentDetails{
		Code:              statusCode,
		Explanation:       msgBody.ResultDescription,
		Paymentreference:  msgBody.ReceiptNumber,
		Sequenceid:        msgBody.TransactionId,
		Success:           msgBody.Success,
		Terminalmno:       "TnmMalawiPayment",
		Approvalreference: msgBody.ReceiptNumber,
		Timestamp:         msgBody.ResultTime,
	}
	return &models.PaymentGatewayRequest{
		Paymentmethod:  "mobile",
		Paymenttype:    "Process",
		Paymentdetails: paymentDetails,
	}

}
