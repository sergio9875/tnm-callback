package main

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
	"malawi-callback/process"
	"malawi-callback/utils"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
)

func LambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var controller *process.Controller
	var err error
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	// If no name is provided in the HTTP request body, throw an error
	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}
	controller = process.NewController(os.Getenv("SECRET_NAME"))

	controller.PreProcess(utils.StringPtr(uuid.New().String()))
	err = controller.Process(ctx, request)
	if err != nil {
		log.Fatalf("Lambda process failed %s", err.Error())
		return events.APIGatewayProxyResponse{}, nil
	}
	controller.PostProcess()

	return events.APIGatewayProxyResponse{}, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
