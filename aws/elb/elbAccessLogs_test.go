package elb

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfAccessLogsEnabled(t *testing.T) {
	type args struct {
		checkConfig   commons.CheckConfig
		loadBalancers []LoadBalancerAttributes
		testName      string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfAccessLogsEnabled",
			args: args{
				checkConfig: commons.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan commons.Check, 1),
				},
				loadBalancers: []LoadBalancerAttributes{
					{
						LoadBalancerName: "test",
						LoadBalancerArn:  "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test/1a2b3c4d5e6f",
						Output: &elasticloadbalancingv2.DescribeLoadBalancerAttributesOutput{
							Attributes: []types.LoadBalancerAttribute{
								{
									Key:   aws.String("access_logs.s3.enabled"),
									Value: aws.String("true"),
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfAccessLogsEnabled(tt.args.checkConfig, tt.args.loadBalancers, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckifAccessLogsEnabled() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfAccessLogsEnabledFail(t *testing.T) {
	type args struct {
		checkConfig   commons.CheckConfig
		loadBalancers []LoadBalancerAttributes
		testName      string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfAccessLogsEnabled",
			args: args{
				checkConfig: commons.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan commons.Check, 1),
				},
				loadBalancers: []LoadBalancerAttributes{
					{
						LoadBalancerName: "test",
						LoadBalancerArn:  "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test/1a2b3c4d5e6f",
						Output: &elasticloadbalancingv2.DescribeLoadBalancerAttributesOutput{
							Attributes: []types.LoadBalancerAttribute{
								{
									Key:   aws.String("access_logs.s3.enabled"),
									Value: aws.String("false"),
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfAccessLogsEnabled(tt.args.checkConfig, tt.args.loadBalancers, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckifAccessLogsEnabled() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
