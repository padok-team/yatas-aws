package elb

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckAlbEnsureHttpsSuccess(t *testing.T) {
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
			name: "Test ALB with HTTPS and HTTP redirect",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				loadBalancers: []LoadBalancerAttributes{
					{
						LoadBalancerType: "application",
						LoadBalancerName: "test-alb",
						LoadBalancerArn:  "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/test-alb/1234567890",
						Listeners: []types.Listener{
							{
								Port: aws.Int32(443),
								DefaultActions: []types.Action{
									{
										Type: types.ActionTypeEnumForward,
									},
								},
							},
							{
								Port: aws.Int32(80),
								DefaultActions: []types.Action{
									{
										Type: types.ActionTypeEnumRedirect,
										RedirectConfig: &types.RedirectActionConfig{
											Protocol:   aws.String("HTTPS"),
											Port:       aws.String("443"),
											Host:       aws.String("#{host}"),
											Path:       aws.String("/#{path}"),
											Query:      aws.String("#{query}"),
											StatusCode: types.RedirectActionStatusCodeEnumHttp301,
										},
									},
								},
							},
						},
					},
				},
				testName: "AWS_ELB_002",
			},
		},
		{
			name: "Test ALB with HTTPS only",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				loadBalancers: []LoadBalancerAttributes{
					{
						LoadBalancerType: "application",
						LoadBalancerName: "test-alb",
						LoadBalancerArn:  "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/test-alb/1234567890",
						Listeners: []types.Listener{
							{
								Port: aws.Int32(443),
								DefaultActions: []types.Action{
									{
										Type: types.ActionTypeEnumForward,
									},
								},
							},
						},
					},
				},
				testName: "AWS_ELB_002",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckAlbEnsureHttps(tt.args.checkConfig, tt.args.loadBalancers, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckAlbEnsureHttps() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckAlbEnsureHttpsFail(t *testing.T) {
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
			name: "Test ALB without HTTPS",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				loadBalancers: []LoadBalancerAttributes{
					{
						LoadBalancerType: "application",
						LoadBalancerName: "test-alb",
						LoadBalancerArn:  "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/test-alb/1234567890",
						Listeners: []types.Listener{
							{
								Port: aws.Int32(80),
								DefaultActions: []types.Action{
									{
										Type: types.ActionTypeEnumForward,
									},
								},
							},
						},
					},
				},
				testName: "AWS_ELB_002",
			},
		},
		{
			name: "Test ALB with HTTPS but no HTTP redirect",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				loadBalancers: []LoadBalancerAttributes{
					{
						LoadBalancerType: "application",
						LoadBalancerName: "test-alb",
						LoadBalancerArn:  "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/test-alb/1234567890",
						Listeners: []types.Listener{
							{
								Port: aws.Int32(443),
								DefaultActions: []types.Action{
									{
										Type: types.ActionTypeEnumForward,
									},
								},
							},
							{
								Port: aws.Int32(80),
								DefaultActions: []types.Action{
									{
										Type: types.ActionTypeEnumForward,
									},
								},
							},
						},
					},
				},
				testName: "AWS_ELB_002",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckAlbEnsureHttps(tt.args.checkConfig, tt.args.loadBalancers, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckAlbEnsureHttps() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
