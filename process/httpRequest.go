package process

import (
	"fmt"
	"malawi-callback/models"
	"time"
)

func (c *Controller) mapPaymentGatewayRequest(msgBody *models.IncomingRequest, statusCode int) *models.PaymentGatewayRequest {

	mbtPerId, err := c.getMbtPerId(msgBody.TransactionId)

	if err != nil {
		fmt.Println("Error with mbtPerId")
	}

	tm := time.Now().Format("02-Jan-2006 15:04:05")

	paymentDetails := models.PaymentDetails{
		Code:              statusCode,
		Explanation:       msgBody.ResultDesc,
		Paymentreference:  mbtPerId.Mbt,
		Sequenceid:        msgBody.TransactionId,
		Terminalmno:       "TnmMalawiPayment",
		Terminalsettings:  nil,
		Amount:            "",
		Currency:          "",
		Msisdn:            "",
		Approvalreference: msgBody.ConversationId,
		Timestamp:         tm,
		Customerdetails:   nil,
		Mbtid:             msgBody.TransactionId,
	}
	return &models.PaymentGatewayRequest{
		Paymentmethod:  "mobile",
		Transunq:       nil,
		Paymenttype:    "Process",
		Paymentdetails: paymentDetails,
	}

}
