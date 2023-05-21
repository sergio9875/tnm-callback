package main

import (
	"context"
	"github.com/google/uuid"
	"os"
	log "tnm-malawi/connectors/callback/logger"
	"tnm-malawi/connectors/callback/process"
	"tnm-malawi/connectors/callback/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	invokeCount = 0
	controller  *process.Controller
)

const DefaultInvokeCount = 15

func Init() {
	controller = process.NewController(os.Getenv("SECRET_NAME"))
	invokeCount = 0
}

func init() {
	// used to init anything special
}

func LambdaHandler(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Debug("ROOT", "version: <GIT_HASH>")
	if invokeCount == 0 {
		Init()
	}

	invokeCount = invokeCount + 1
	if invokeCount > *utils.SafeAtoi(utils.Getenv("MAX_INVOKE", "15"), utils.IntPtr(DefaultInvokeCount)) {
		// reset global variables to nil
		controller.ShutDown()
		Init()
		invokeCount = 1
	}

	requestId := uuid.New().String()
	controller.PreProcess(utils.StringPtr(requestId))
	res, err := controller.Process(ctx, event)
	if err != nil {
		controller.PostProcess()
		return res, err
	}
	controller.PostProcess()
	return res, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
