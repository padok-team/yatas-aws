package ssm

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfAuditLogsEnabledOnBastionInstanceSuccess(t *testing.T) {
	type args struct {
		checkConfig      commons.CheckConfig
		ec2ToIAMPolicies []BastionToIAMPolicies
		testName         string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test bastion with correct audit logs permissions",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				ec2ToIAMPolicies: []BastionToIAMPolicies{
					{
						Instance: types.Instance{
							InstanceId: aws.String("i-1234567890abcdef0"),
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("bastion-prod"),
								},
							},
						},
						Policies: []string{
							`{
								"Version": "2012-10-17",
								"Statement": [{
									"Sid": "SSMLogging",
									"Effect": "Allow",
									"Action": ["s3:PutObject", "s3:PutObjectAcl"],
									"Resource": ["arn:aws:s3:::ssm-logging-bucket/*"]
								}]
							}`,
						},
					},
				},
				testName: "AWS_SSM_001",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfAuditLogsEnabledOnBastionInstance(tt.args.checkConfig, tt.args.ec2ToIAMPolicies, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfAuditLogsEnabledOnBastionInstance() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfAuditLogsEnabledOnBastionInstanceFail(t *testing.T) {
	type args struct {
		checkConfig      commons.CheckConfig
		ec2ToIAMPolicies []BastionToIAMPolicies
		testName         string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test bastion with missing permissions",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				ec2ToIAMPolicies: []BastionToIAMPolicies{
					{
						Instance: types.Instance{
							InstanceId: aws.String("i-1234567890abcdef0"),
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("bastion-prod"),
								},
							},
						},
						Policies: []string{
							`{
								"Version": "2012-10-17",
								"Statement": [{
									"Effect": "Allow",
									"Action": "s3:PutObject",
									"Resource": "arn:aws:s3:::ssm-logging-bucket/*"
								}]
							}`,
						},
					},
				},
				testName: "AWS_SSM_001",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfAuditLogsEnabledOnBastionInstance(tt.args.checkConfig, tt.args.ec2ToIAMPolicies, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfAuditLogsEnabledOnBastionInstance() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
