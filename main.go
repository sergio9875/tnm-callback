package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"os"
	log "tnm-malawi/connectors/callback/logger"
	"tnm-malawi/connectors/callback/process"
	"tnm-malawi/connectors/callback/utils"

	"github.com/aws/aws-lambda-go/events"
)

var (
	ErrNameNotProvided = errors.New("body was provided in the ApiGateway event")
	res                = events.APIGatewayProxyResponse{}
	controller         *process.Controller
	err                error
)

func LambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//stdout and stderr are sent to AWS CloudWatch Logs
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))
	fmt.Println("request Body:", request.Body)
	fmt.Println("request HTTPMethod:", request.HTTPMethod)
	fmt.Println("request Headers:", request.Headers)

	if (len(request.Body)) < 1 {
		fmt.Println("Request.Body was not provided")
		return events.APIGatewayProxyResponse{}, nil
	}

	controller = process.NewController(os.Getenv("SECRET_NAME"))

	controller.PreProcess(utils.StringPtr(uuid.New().String()))
	res, err = controller.Process(ctx, request)
	log.Println("Processing Lambda request %s\n", res)

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
