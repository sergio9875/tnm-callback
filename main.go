package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"os"
	"tnm-malawi/connectors/callback/enums"
	log "tnm-malawi/connectors/callback/logger"
	"tnm-malawi/connectors/callback/models"
	"tnm-malawi/connectors/callback/process"
	"tnm-malawi/connectors/callback/utils"

	"github.com/aws/aws-lambda-go/events"
)

var (
	res        = models.Response{}
	controller *process.Controller
	err        error
)

func LambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (models.Response, error) {
	//stdout and stderr are sent to AWS CloudWatch Logs
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))
	fmt.Println("request Body:", request.Body)
	fmt.Println("request HTTPMethod:", request.HTTPMethod)
	fmt.Println("request Headers:", request.Headers)

	if (len(request.Body)) < 1 {
		fmt.Println("Request.Body was not provided")
		return models.Response{
			Body: enums.ERROR_MSG_EMPTY_REQ,
		}, nil
	}

	controller = process.NewController(os.Getenv("SECRET_NAME"))

	controller.PreProcess(utils.StringPtr(uuid.New().String()))

	res, err = controller.Process(ctx, request)
	log.Println("Processing Lambda request %s\n", res)

	if err != nil {
		log.Fatalf("Lambda process failed %s", err.Error())
		return models.Response{
			Body: "Lambda process failed",
		}, nil
	}
	controller.PostProcess()

	res2 := models.Response{
		Body:          res.Body,
		StatusCode:    res.StatusCode,
		TransactionId: res.TransactionId,
	}
	log.Println("Final Result", res2)

	return res2, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
