package vpc

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfSubnetInDifferentZone(t *testing.T) {
	type args struct {
		checkConfig  commons.CheckConfig
		vpcToSubnets []VPCToSubnet
		testName     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfSubnetInDifferentZone",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				vpcToSubnets: []VPCToSubnet{
					{
						VpcID: "test",
						Subnets: []types.Subnet{
							{
								SubnetId:         aws.String("test"),
								AvailabilityZone: aws.String("test"),
							},
							{
								SubnetId:         aws.String("test"),
								AvailabilityZone: aws.String("test2"),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfSubnetInDifferentZone(tt.args.checkConfig, tt.args.vpcToSubnets, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfSubnetInDifferentZone() = %v", t)
					}
					tt.args.checkConfig.Wg.Done()
				}

			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfSubnetInDifferentZoneFail(t *testing.T) {
	type args struct {
		checkConfig  commons.CheckConfig
		vpcToSubnets []VPCToSubnet
		testName     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfSubnetInDifferentZone",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				vpcToSubnets: []VPCToSubnet{
					{
						VpcID: "test",
						Subnets: []types.Subnet{
							{
								SubnetId:         aws.String("test"),
								AvailabilityZone: aws.String("test"),
							},
							{
								SubnetId:         aws.String("test"),
								AvailabilityZone: aws.String("test"),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfSubnetInDifferentZone(tt.args.checkConfig, tt.args.vpcToSubnets, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfSubnetInDifferentZone() = %v", t)
					}
					tt.args.checkConfig.Wg.Done()
				}

			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
