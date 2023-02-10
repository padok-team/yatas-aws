package cognito

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfCognitoAllowsUnauthenticated(t *testing.T) {
	type args struct {
		checkConfig  commons.CheckConfig
		cognitoPools []cognitoidentity.DescribeIdentityPoolOutput
		testName     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfCognitoAllowsUnauthenticated",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				cognitoPools: []cognitoidentity.DescribeIdentityPoolOutput{
					{
						IdentityPoolName:               aws.String("test"),
						IdentityPoolId:                 aws.String("us-east-1:cb21213c-a931-11ed-afa1-0242ac120002"),
						AllowUnauthenticatedIdentities: false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfCognitoAllowsUnauthenticated(tt.args.checkConfig, tt.args.cognitoPools, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfCognitoAllowsUnauthenticated() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfCognitoAllowsUnauthenticatedFail(t *testing.T) {
	type args struct {
		checkConfig  commons.CheckConfig
		cognitoPools []cognitoidentity.DescribeIdentityPoolOutput
		testName     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfCognitoAllowsUnauthenticated",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				cognitoPools: []cognitoidentity.DescribeIdentityPoolOutput{
					{
						IdentityPoolName:               aws.String("test"),
						IdentityPoolId:                 aws.String("us-east-1:cb21213c-a931-11ed-afa1-0242ac120002"),
						AllowUnauthenticatedIdentities: true,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfCognitoAllowsUnauthenticated(tt.args.checkConfig, tt.args.cognitoPools, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfCognitoAllowsUnauthenticated() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
