package process

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"net/http"
	"tnm-malawi/connectors/callback/enums"
	log "tnm-malawi/connectors/callback/logger"
	"tnm-malawi/connectors/callback/models"
	repo "tnm-malawi/connectors/callback/repository"
	"tnm-malawi/connectors/callback/request"
	"tnm-malawi/connectors/callback/utils"
)

// Controller container
type Controller struct {
	secretHolder *SecretIDHolder
	sumoProducer *SQSProducer
	config       *models.SecretModel
	repository   *repo.Repository
	httpClient   *request.IRequest
	requestId    *string
}

func NewController(secret string) *Controller {
	controller := Controller{
		requestId: utils.StringPtr("ROOT"),
	}
	controller.initSecret(secret)
	controller.initRepository()
	controller.initClient()
	controller.initSumoProducer()
	return &controller
}

func (c *Controller) ShutDown() {
	c.config = nil
	c.sumoProducer = nil
	c.secretHolder = nil
	c.httpClient = nil
	c.repository = nil
}

func (c *Controller) PreProcess(pid *string) {
	c.requestId = pid
}

func (c *Controller) PostProcess() {
	c.requestId = utils.StringPtr("ROOT")
}

func (c *Controller) Process(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	c.sendSumoMessages(ctx, "start tnm-malawi get callback process", request)

	var err error
	msgBody := new(models.IncomingRequest)
	pgwResponse := new(models.PaymentGatewayResponse)
	err = json.Unmarshal([]byte(request.Body), &msgBody)

	if err != nil {
		res := models.Res{
			StatusCode:        enums.RESULT_CODE_SUCCESS,
			StatusDescription: err.Error(),
			IsBase64Encoded:   false,
			Headers:           make(map[string]string),
		}
		c.sendSumoMessages(ctx, enums.ERROR_MSG_UNMARSHL, res)
		log.Fatalf(*c.requestId, enums.ERROR_MSG_UNMARSHL+"%s", err.Error())
	}

	url := "http://sergeyk-3g.dev.directpay.online/PaymentGateway/paymentGateway.php"
	//url := c.config.DpoPygwUrl
	//if strings.Trim(os.Getenv("PGW_URL"), "") != "sm" {
	//	url = os.Getenv("PGW_URL")
	//}
	log.Infof(*c.requestId, "message body", msgBody)
	log.Infof(*c.requestId, "pgw url", url)

	var statusCode int
	if msgBody.ResultCode == enums.RESULT_CODE_SUCCESS {
		statusCode = enums.PGW_STATUS_SUCCESS
	} else {
		statusCode = enums.PGW_STATUS_FAILED
	}

	log.Infof(*c.requestId, "Status code assigned:", statusCode)

	pgwRequest := c.mapPaymentGatewayRequest(msgBody, statusCode)

	c.sendSumoMessages(ctx, "payment gateway request", pgwRequest)

	log.Infof(*c.requestId, "trying to send request to payment gateway",
		pgwRequest, "to:", url)
	headers := make(map[string]string, 0)

	if err := (*c.httpClient).PostWithJsonResponse(url, headers, pgwRequest, pgwResponse); err != nil {
		c.sendSumoMessages(ctx, utils.JsonIt(err), nil)
		return &events.APIGatewayProxyResponse{
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            err.Error(),
			StatusCode:      400,
			IsBase64Encoded: false,
		}, nil

	}

	log.Infof(*c.requestId, "successfully retrieved payment gateway response %v", pgwResponse)
	var resp = models.Res{}

	if pgwResponse.Code == enums.PGW_FAILED {
		resp = models.Res{
			StatusCode:        msgBody.ResultCode,
			StatusDescription: msgBody.ResultDesc,
			Headers:           map[string]string{"Content-Type": "application/json"},
			ExternalRef:       msgBody.ExternalRef,
			PgwDescription:    enums.PGW_FAILED_BODY,
			PgwStatusCode:     pgwResponse.Code,
		}
	} else {
		resp = models.Res{
			StatusCode:        msgBody.ResultCode,
			StatusDescription: msgBody.ResultDesc,
			Headers:           map[string]string{"Content-Type": "application/json"},
			ExternalRef:       msgBody.ExternalRef,
			PgwDescription:    enums.PGW_SUCCESS_BODY,
			PgwStatusCode:     pgwResponse.Code,
		}
	}

	return &events.APIGatewayProxyResponse{
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            utils.JsonIt(&resp),
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
	}, nil
}

func (c *Controller) sendSumoMessages(ctx context.Context, message string, params interface{}) {
	if params != nil {
		params = fmt.Sprintf("%+v", params)
	}

	sumo := &models.SumoPusherMessage{
		Category: "malawi",
		SumoPayload: models.SumoPayload{
			Stack:   *c.requestId,
			Message: "[tnm-tnm-malawi/connectors/callback-status-check] " + message,
			Params:  params,
		},
	}

	messageBody, err := json.Marshal(sumo)
	if err != nil {
		log.Error(*c.requestId, "Error Create Message Body For API Gateway: ", err.Error())
		return
	}

	sqsMessage := &sqs.SendMessageInput{
		MessageBody: aws.String(string(messageBody)),
	}

	_, err = c.sumoProducer.SendMsg(ctx, sqsMessage)

	if err != nil {
		log.Error(*c.requestId, "Error while pushing to Api Gateway: ", err.Error())
		return
	}
}
