package lambda

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/stangirard/yatas/config"
)

func TestCheckIfLambdaNoErrors(t *testing.T) {
	type args struct {
		checkConfig config.CheckConfig
		lambdas     []types.FunctionConfiguration
		testName    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestCheckIfLambdaNoErrors",
			args: args{
				checkConfig: config.CheckConfig{Queue: make(chan config.Check, 1), Wg: &sync.WaitGroup{}},
				lambdas: []types.FunctionConfiguration{
					{
						FunctionName:    aws.String("test"),
						FunctionArn:     aws.String("arn:aws:lambda:eu-west-3:123456789012:function:test"),
						StateReasonCode: types.StateReasonCodeIdle,
					},
				},
				testName: "TestCheckIfLambdaNoErrors",
			},
			want: "OK",
		},
		{
			name: "TestCheckIfLambdaNoErrors",
			args: args{
				checkConfig: config.CheckConfig{Queue: make(chan config.Check, 1), Wg: &sync.WaitGroup{}},
				lambdas: []types.FunctionConfiguration{
					{
						FunctionName:    aws.String("test"),
						FunctionArn:     aws.String("arn:aws:lambda:eu-west-3:123456789012:function:test"),
						StateReasonCode: types.StateReasonCodeEniLimitExceeded,
					},
				},
				testName: "TestCheckIfLambdaNoErrors",
			},
			want: "FAIL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfLambdaNoErrors(tt.args.checkConfig, tt.args.lambdas, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for result := range tt.args.checkConfig.Queue {
					if result.Status != tt.want {
						t.Errorf("CheckIfLambdaNoErrors() = %v, want %v", result.Status, tt.want)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
