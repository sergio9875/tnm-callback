package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"os"
	"testing"
)

func TestLambdaHandler(t *testing.T) {
	err := os.Setenv("AWS_REGION", "eu-west-1")
	if err != nil {
		return
	}
	err = os.Setenv("SECRET_NAME", "configuration/dpo/callback/config")
	if err != nil {
		return
	}

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
	"transaction_id": "R3T5T9G9",
	"external_ref":"1348800",
	"response_time": "2023-05-17 18:12:55"
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
				var err error
				if (err != nil) != tt.wantErr {
					t.Errorf("LambdaHandler() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			})
		}
	}
}
