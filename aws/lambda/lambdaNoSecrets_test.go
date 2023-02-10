package lambda

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfLambdaNoSecrets(t *testing.T) {
	type args struct {
		checkConfig commons.CheckConfig
		lambdas     []types.FunctionConfiguration
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfLambdaNoSecrets",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				lambdas: []types.FunctionConfiguration{
					{
						FunctionName: aws.String("test"),
						FunctionArn:  aws.String("arn:aws:lambda:us-east-1:123456789012:function:test"),
						Environment: &types.EnvironmentResponse{
							Variables: map[string]string{},
						},
					},
				},
			},
		},
		{
			name: "TestCheckIfLambdaNoSecrets",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				lambdas: []types.FunctionConfiguration{
					{
						FunctionName: aws.String("test"),
						FunctionArn:  aws.String("arn:aws:lambda:us-east-1:123456789012:function:test"),
						Environment: &types.EnvironmentResponse{
							Variables: map[string]string{
								"my_variable": "test",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfLambdaNoSecrets(tt.args.checkConfig, tt.args.lambdas, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfLambdaNoSecrets() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfLambdaNoSecretsFail(t *testing.T) {
	type args struct {
		checkConfig commons.CheckConfig
		lambdas     []types.FunctionConfiguration
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfLambdaNoSecrets",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				lambdas: []types.FunctionConfiguration{
					{
						FunctionName: aws.String("test"),
						FunctionArn:  aws.String("arn:aws:lambda:us-east-1:123456789012:function:test"),
						Environment: &types.EnvironmentResponse{
							Variables: map[string]string{
								"aws_access_key": "ASIAS6VZTAEWPKBMXQOL",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfLambdaNoSecrets(tt.args.checkConfig, tt.args.lambdas, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfLambdaNoSecrets() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
