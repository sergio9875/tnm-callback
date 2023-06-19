package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"os"
	log "tnm-malawi/connectors/callback/logger"
	"tnm-malawi/connectors/callback/process"
	"tnm-malawi/connectors/callback/utils"

	"github.com/aws/aws-lambda-go/events"
)

var invokeCount = 0
var controller *process.Controller

const DefaultInvokeCount = 15

func Init() {
	controller = process.NewController(os.Getenv("SECRET_NAME"))
	invokeCount = 0
}

func init() {
	// used to init anything special
}

func LambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	//stdout and stderr are sent to AWS CloudWatch Logs
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))
	fmt.Println("request Body:", request.Body)
	fmt.Println("request HTTPMethod:", request.HTTPMethod)
	fmt.Println("request Headers:", request.Headers)

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
	res, err := controller.Process(ctx, request)
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
