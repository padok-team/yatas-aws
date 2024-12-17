package vpc

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfPrivateAndPublicSubnets(t *testing.T) {
	type args struct {
		checkConfig                commons.CheckConfig
		vpcToSubnetWithRouteTables map[string][]SubnetWithRouteTables
		testName                   string
	}
	tests := []struct {
		name       string
		args       args
		expectFail bool
	}{
		{
			name: "VPC with both public and private subnets",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				vpcToSubnetWithRouteTables: map[string][]SubnetWithRouteTables{
					"vpc-12345": {
						{
							SubnetId: "subnet-1",
							RouteTables: []types.RouteTable{
								{
									Routes: []types.Route{
										{GatewayId: aws.String("igw-12345")},
									},
								},
							},
						},
						{
							SubnetId: "subnet-2",
							RouteTables: []types.RouteTable{
								{
									Routes: []types.Route{
										{GatewayId: nil},
									},
								},
							},
						},
					},
				},
				testName: "CheckSubnets",
			},
			expectFail: false,
		},
		{
			name: "VPC with no public subnet",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				vpcToSubnetWithRouteTables: map[string][]SubnetWithRouteTables{
					"vpc-67890": {
						{
							SubnetId: "subnet-1",
							RouteTables: []types.RouteTable{
								{
									Routes: []types.Route{
										{GatewayId: nil},
									},
								},
							},
						},
					},
				},
				testName: "CheckSubnets",
			},
			expectFail: true,
		},
		{
			name: "VPC with no private subnet",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				vpcToSubnetWithRouteTables: map[string][]SubnetWithRouteTables{
					"vpc-54321": {
						{
							SubnetId: "subnet-1",
							RouteTables: []types.RouteTable{
								{
									Routes: []types.Route{
										{GatewayId: aws.String("igw-67890")},
									},
								},
							},
						},
					},
				},
				testName: "CheckSubnets",
			},
			expectFail: true,
		},
		{
			name: "VPC with no subnets",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				vpcToSubnetWithRouteTables: map[string][]SubnetWithRouteTables{
					"vpc-00000": {},
				},
				testName: "CheckSubnets",
			},
			expectFail: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				defer tt.args.checkConfig.Wg.Done()
				CheckIfPrivateAndPublicSubnets(tt.args.checkConfig, tt.args.vpcToSubnetWithRouteTables, tt.args.testName)
			}()

			tt.args.checkConfig.Wg.Wait()
			close(tt.args.checkConfig.Queue)

			for check := range tt.args.checkConfig.Queue {
				failed := false
				for _, result := range check.Results {
					if result.Status == "FAIL" {
						failed = true
					}
				}

				if failed != tt.expectFail {
					t.Errorf("Test %s failed. Expected fail: %v, got fail: %v", tt.name, tt.expectFail, failed)
				}
			}
		})
	}
}
