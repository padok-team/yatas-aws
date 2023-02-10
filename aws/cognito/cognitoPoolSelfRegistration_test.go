package cognito

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfCognitoSelfRegistration(t *testing.T) {
	type args struct {
		checkConfig      commons.CheckConfig
		cognitoUserPools []cognitoidentityprovider.DescribeUserPoolOutput
		testName         string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfCognitoSelfRegistration",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				cognitoUserPools: []cognitoidentityprovider.DescribeUserPoolOutput{
					{
						UserPool: &types.UserPoolType{
							Name:                  aws.String("test"),
							Arn:                   aws.String("arn:aws:cognito-idp:us-east-1:123456789012:userpool/us-east-1_test"),
							AdminCreateUserConfig: &types.AdminCreateUserConfigType{AllowAdminCreateUserOnly: true},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfCognitoSelfRegistration(tt.args.checkConfig, tt.args.cognitoUserPools, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfCognitoSelfRegistration() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfCognitoSelfRegistrationFail(t *testing.T) {
	type args struct {
		checkConfig      commons.CheckConfig
		cognitoUserPools []cognitoidentityprovider.DescribeUserPoolOutput
		testName         string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfCognitoSelfRegistration",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				cognitoUserPools: []cognitoidentityprovider.DescribeUserPoolOutput{
					{
						UserPool: &types.UserPoolType{
							Name:                  aws.String("test"),
							Arn:                   aws.String("arn:aws:cognito-idp:us-east-1:123456789012:userpool/us-east-1_test"),
							AdminCreateUserConfig: &types.AdminCreateUserConfigType{AllowAdminCreateUserOnly: false},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfCognitoSelfRegistration(tt.args.checkConfig, tt.args.cognitoUserPools, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfCognitoSelfRegistration() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
