package process

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"malawi-callback/enums"
	log "malawi-callback/logger"
	"malawi-callback/models"
	repo "malawi-callback/repository"
	"malawi-callback/request"
	"malawi-callback/utils"
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
	controller.initSumoProducer()
	controller.initRepository()
	controller.initClient()
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

func (c *Controller) Process(ctx context.Context, request events.APIGatewayProxyRequest) error {
	c.sendSumoMessages(ctx, "start tnm-malawi get callback process", request)

	var err error
	msgBody := new(models.IncomingRequest)
	pgwResponse := new(models.PaymentGatewayResponse)

	if err = c.handleGetMessage(request.Body, &msgBody); err != nil {
		c.sendSumoMessages(ctx, err.Error(), nil)
		return err
	}

	//pgwUrl := c.config.DpoPygwUrl
	pgwUrl := "http://sergeyk-3g.dev.directpay.online/PaymentGateway/paymentGateway.php"
	headers := make(map[string]string, 0)

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
		pgwRequest, "to:", pgwUrl)

	if err := (*c.httpClient).PostWithJsonResponse(pgwUrl, headers, pgwRequest, pgwResponse); err != nil {
		return err
	}

	log.Infof(*c.requestId, "successfully retrieved payment gateway response %v", pgwResponse)

	return nil

}

func (c *Controller) sendSumoMessages(ctx context.Context, message string, params interface{}) {

	if params != nil {
		params = fmt.Sprintf("%+v", params)
	}

	sumo := &models.SumoPusherMessage{
		Category: "tnm",
		SumoPayload: models.SumoPayload{
			Stack:   *c.requestId,
			Message: "[tnm-malawi-callback-status-check] " + message,
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
		log.Error(*c.requestId, "Error while pushing to sqs producer: ", err.Error())
		return
	}

	log.Info(*c.requestId, "Message Successfully pushed")
}
