package vpc

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/plugins/commons"
)

func TestCheckInboudPortSucess(t *testing.T) {
	type args struct {
		checkConfig commons.CheckConfig
		sg2SGRs     []SGToSecurityGroupRules
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test inbound 0.0.0.0/0 80 & 443",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				sg2SGRs: []SGToSecurityGroupRules{
					{
						vpcID: "vpc1",
						securityGroup: types.SecurityGroup{
							Description: aws.String("test"),
							GroupId:     aws.String("1"),
							GroupName:   aws.String("group-1"),
						},
						securityGroupRules: []types.SecurityGroupRule{
							{
								CidrIpv4:            aws.String("0.0.0.0/0"),
								FromPort:            aws.Int32(80),
								GroupId:             aws.String("1"),
								IsEgress:            aws.Bool(false),
								SecurityGroupRuleId: aws.String("rule-1"),
								ToPort:              aws.Int32(80),
							},
							{
								CidrIpv4:            aws.String("0.0.0.0/0"),
								FromPort:            aws.Int32(443),
								GroupId:             aws.String("1"),
								IsEgress:            aws.Bool(false),
								SecurityGroupRuleId: aws.String("rule-2"),
								ToPort:              aws.Int32(443),
							},
						},
					},
				},
				testName: "test success: basic",
			},
		},
		{
			name: "Test not inbound but 0.0.0.0/0",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				sg2SGRs: []SGToSecurityGroupRules{
					{
						vpcID: "vpc2",
						securityGroup: types.SecurityGroup{
							Description: aws.String("test"),
							GroupId:     aws.String("1"),
							GroupName:   aws.String("group-1"),
						},
						securityGroupRules: []types.SecurityGroupRule{
							{
								CidrIpv4:            aws.String("0.0.0.1/0"),
								FromPort:            aws.Int32(79),
								GroupId:             aws.String("1"),
								IsEgress:            aws.Bool(false),
								SecurityGroupRuleId: aws.String("rule-1"),
								ToPort:              aws.Int32(444),
							},
							{
								CidrIpv4:            aws.String("0.0.0.0/0"),
								FromPort:            aws.Int32(79),
								GroupId:             aws.String("1"),
								IsEgress:            aws.Bool(true),
								SecurityGroupRuleId: aws.String("rule-2"),
								ToPort:              aws.Int32(444),
							},
						},
					},
				},
				testName: "test success: not concerned",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckInboudPort(tt.args.checkConfig, tt.args.sg2SGRs, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckInboudPort() = %v", t)
					}
					tt.args.checkConfig.Wg.Done()
				}

			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckInboudPortFail(t *testing.T) {
	type args struct {
		checkConfig commons.CheckConfig
		sg2SGRs     []SGToSecurityGroupRules
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test inbound 0.0.0.0/0 !=  (80 || 443)",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				sg2SGRs: []SGToSecurityGroupRules{
					{
						vpcID: "vpc3",
						securityGroup: types.SecurityGroup{
							Description: aws.String("test"),
							GroupId:     aws.String("1"),
							GroupName:   aws.String("group-1"),
						},
						securityGroupRules: []types.SecurityGroupRule{
							{
								CidrIpv4:            aws.String("0.0.0.0/0"),
								FromPort:            aws.Int32(80),
								GroupId:             aws.String("1"),
								IsEgress:            aws.Bool(false),
								SecurityGroupRuleId: aws.String("rule-1"),
								ToPort:              aws.Int32(443),
							},
						},
					},
				},
				testName: "test failed: 80-443",
			},
		},
		{
			name: "Test inbound 0.0.0.0/0 !=  (80 || 443)",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				sg2SGRs: []SGToSecurityGroupRules{
					{
						vpcID: "vpc4",
						securityGroup: types.SecurityGroup{
							Description: aws.String("test"),
							GroupId:     aws.String("1"),
							GroupName:   aws.String("group-1"),
						},
						securityGroupRules: []types.SecurityGroupRule{
							{
								CidrIpv4:            aws.String("0.0.0.0/0"),
								FromPort:            aws.Int32(79),
								GroupId:             aws.String("1"),
								IsEgress:            aws.Bool(false),
								SecurityGroupRuleId: aws.String("rule-1"),
								ToPort:              aws.Int32(79),
							},
						},
					},
				},
				testName: "test failed: 79",
			},
		},
		{
			name: "Test inbound 0.0.0.0/0 !=  (80 || 443)",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				sg2SGRs: []SGToSecurityGroupRules{
					{
						vpcID: "vpc5",
						securityGroup: types.SecurityGroup{
							Description: aws.String("test"),
							GroupId:     aws.String("1"),
							GroupName:   aws.String("group-1"),
						},
						securityGroupRules: []types.SecurityGroupRule{
							{
								CidrIpv4:            aws.String("0.0.0.0/0"),
								FromPort:            aws.Int32(79),
								GroupId:             aws.String("1"),
								IsEgress:            aws.Bool(false),
								SecurityGroupRuleId: aws.String("rule-1"),
								ToPort:              aws.Int32(80),
							},
						},
					},
				},
				testName: "test failed: 79-80",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckInboudPort(tt.args.checkConfig, tt.args.sg2SGRs, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckInboudPort() = %v", t)
					}
					tt.args.checkConfig.Wg.Done()
				}

			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
