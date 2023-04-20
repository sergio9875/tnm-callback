package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"os"
	"testing"
)

func TestLambdaHandler(t *testing.T) {
	os.Setenv("AWS_REGION", "eu-west-1")
	os.Setenv("SECRET_NAME", "configuration/dpo/callback/config")

	type args struct {
		ctx        context.Context
		apiGateway events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Integration Test",
			args: args{
				ctx: context.TODO(),
				apiGateway: events.APIGatewayProxyRequest{Body: `{
	"conversation_id": "AG_20200714_00006368d1caf5faf759",
	"result_code": "200",
	"result_desc": "Process service request successfully.",
	"transaction_id": "1348437",
	"external_ref":"F509897",
	"response_time": "2020-07-14 13:12:55"
	}`},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if os.Getenv("SKIP_MAIN") == "" {
			t.Run(tt.name, func(t *testing.T) {
				_, _ = LambdaHandler(tt.args.ctx, tt.args.apiGateway)
				//if (err != ) != tt.wantErr {
				//	t.Errorf("LambdaHandler() error = %v, wantErr %v", err, tt.wantErr)
				//	return
				//}
			})
		}
	}

	//for _, test := range tests {
	//	response, err := main.Handler(test.request)
	//	assert.IsType(t, test.err, err)
	//	assert.Equal(t, test.expect, response.Body)
	//}

	//	tests := []struct {
	//		name    string
	//		args    args
	//		want    bool
	//		wantErr bool
	//	}{
	//		{
	//			name: "Integration Test",
	//			args: args{
	//				ctx: context.TODO(),
	//				apiGateway:
	//				events.APIGatewayProxyRequest {Body: `{
	//"conversation_id": "AG_20200714_00006368d1caf5faf759",
	//"result_code": "500",
	//"result_desc": "Process service request successfully.",
	//"transaction_id": "1348435",
	//"external_ref":"F509897",
	//"response_time": "2020-07-14 13:12:55"
	//}`},
	//						ReceiptHandle: "ErrNameNotProvided",
	//					},
	//				}},
	//			},
	//			want:    true,
	//			wantErr: false,
	//		},
	//	}
	//	for _, tt := range tests {
	//		if os.Getenv("SKIP_MAIN") == "" {
	//			t.Run(tt.name, func(t *testing.T) {
	//				err := Handler(tt.args.request)
	//				if (err != nil) != tt.wantErr {
	//					t.Errorf("LambdaHandler() error = %v, wantErr %v", err, tt.wantErr)
	//					return
	//				}
	//			})
	//		}
	//	}
}
