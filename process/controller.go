package process

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"os"
	"strconv"
	"strings"
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

func (c *Controller) Process(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	c.sendSumoMessages(ctx, "start tnm-malawi get callback process", request)

	var err error
	var res *events.APIGatewayProxyResponse
	msgBody := new(models.IncomingRequest)
	pgwResponse := new(models.PaymentGatewayResponse)

	if err = c.handleGetMessage(request.Body, &msgBody); err != nil {
		c.sendSumoMessages(ctx, utils.JsonIt(err), nil)
		res = &events.APIGatewayProxyResponse{
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            err.Error(),
			StatusCode:      400,
			IsBase64Encoded: false,
		}
		return res, err
	}
	url := c.config.DpoPygwUrl

	if strings.Trim(os.Getenv("PGW_URL"), "") != "sm" {
		url = os.Getenv("PGW_URL")
	}

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
		pgwRequest, "to:", url)

	if err = (*c.httpClient).PostWithJsonResponse(url, headers, pgwRequest, pgwResponse); err != nil {
		c.sendSumoMessages(ctx, utils.JsonIt(err), nil)
		res = &events.APIGatewayProxyResponse{
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            err.Error(),
			StatusCode:      400,
			IsBase64Encoded: false,
		}
		return res, err
	}

	log.Infof(*c.requestId, "successfully retrieved payment gateway response %v", pgwResponse)
	code, err := strconv.Atoi(pgwResponse.Code)

	if err != nil {
		fmt.Println("Error during conversion")
		return nil, nil
	}
	res = &events.APIGatewayProxyResponse{
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            enums.PGW_RRESPONSE,
		StatusCode:      code,
		IsBase64Encoded: false,
	}
	return res, nil
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

	log.Info(*c.requestId, "Message Successfully pushed")
}
