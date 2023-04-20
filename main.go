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
	//return events.APIGatewayProxyResponse{
	//	Body:       "Hello " + request.Body,
	//	StatusCode: 200,
	//}, nil

}

func main() {
	lambda.Start(LambdaHandler)
}

//
//import (
//	"errors"
//	"github.com/aws/aws-lambda-go/events"
//	"github.com/aws/aws-lambda-go/lambda"
//	"github.com/google/uuid"
//	log "malawi-callback/logger"
//	"malawi-callback/process"
//	"malawi-callback/utils"
//)
//
//var (
//	// ErrNameNotProvided is thrown when a name is not provided
//	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
//)
//
//// var invokeCount = 0
//var controller *process.Controller
//
//// func Init() {
//// \\	controller = process.NewController(os.Getenv("SECRET_NAME"))
////
////		invokeCount = 0
////	}
////
////	func init() {
////		// used to init anything special
////	}
//func LambdaHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//	log.Debug("ROOT", "version: <GIT_HASH>")
//	// stdout and stderr are sent to AWS CloudWatch Logs
//	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)
//
//	// If no name is provided in the HTTP request body, throw an error
//	if len(request.Body) < 1 {
//		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
//	}
//
//	controller.PreProcess(utils.StringPtr(uuid.New().String()))
//	err := controller.Process(request)
//	if err != nil {
//		log.Fatalf("Lambda process failed %s", err.Error())
//		return events.APIGatewayProxyResponse{}, nil
//	}
//	controller.PostProcess()
//
//	//return events.APIGatewayProxyResponse{
//	//	Body:       "Hello " + request.Body,
//	//	StatusCode: 200,
//	//}, nil
//
//	return events.APIGatewayProxyResponse{}, nil
//}
//
//// LambdaHandler - Listen to S3 events and start processing
////func LambdaHandler1(ctx context.Context, sqsEvent events.SQSEvent) error {
////	log.Debug("ROOT", "version: <GIT_HASH>")
////
////	if invokeCount == 0 {
////		Init()
////	}
////
////	invokeCount = invokeCount + 1
////	if invokeCount > utils.SafeAtoi(utils.Getenv("MAX_INVOKE", "15"), 15) {
////		// reset global variables to nil
////		controller.ShutDown()
////
////	}
////
////	for _, record := range sqsEvent.Records {
////		controller.PreProcess(utils.StringPtr(uuid.New().String()))
////		err := controller.Process(ctx, record)
////		if err != nil {
////			log.Fatalf("Lambda process failed %s", err.Error())
////			return err
////		}
////		controller.PostProcess()
////	}
////	return nil
////}
//
//func main() {
//	lambda.Start(LambdaHandler)
//}
