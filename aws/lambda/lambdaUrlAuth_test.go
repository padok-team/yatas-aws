package lambda

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfLambdaUrlAuth(t *testing.T) {
	type args struct {
		checkConfig      commons.CheckConfig
		lambdaUrlConfigs []LambdaUrlConfig
		testName         string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfLambdaUrlAuth",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				lambdaUrlConfigs: []LambdaUrlConfig{
					{
						LambdaName: *aws.String("test"),
						LambdaArn:  *aws.String("arn:aws:lambda:us-east-1:123456789012:function:test"),
						UrlConfigs: []types.FunctionUrlConfig{},
					},
				},
			},
		},
		{
			name: "TestCheckIfLambdaUrlAuth",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				lambdaUrlConfigs: []LambdaUrlConfig{
					{
						LambdaName: *aws.String("test"),
						LambdaArn:  *aws.String("arn:aws:lambda:us-east-1:123456789012:function:test"),
						UrlConfigs: []types.FunctionUrlConfig{
							{AuthType: "AWS_IAM"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfLambdaUrlAuth(tt.args.checkConfig, tt.args.lambdaUrlConfigs, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfLambdaUrlAuth() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfLambdaUrlAuthFail(t *testing.T) {
	type args struct {
		checkConfig      commons.CheckConfig
		lambdaUrlConfigs []LambdaUrlConfig
		testName         string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfLambdaUrlAuth",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				lambdaUrlConfigs: []LambdaUrlConfig{
					{
						LambdaName: *aws.String("test"),
						LambdaArn:  *aws.String("arn:aws:lambda:us-east-1:123456789012:function:test"),
						UrlConfigs: []types.FunctionUrlConfig{
							{AuthType: "NONE"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfLambdaUrlAuth(tt.args.checkConfig, tt.args.lambdaUrlConfigs, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfLambdaUrlAuth() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
