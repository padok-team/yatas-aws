package vpc

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func GetListVPC(s aws.Config) []types.Vpc {
	svc := ec2.NewFromConfig(s)
	var vpcs []types.Vpc
	input := &ec2.DescribeVpcsInput{}
	result, err := svc.DescribeVpcs(context.TODO(), input)
	if err != nil {
		fmt.Println(err)
	}
	vpcs = append(vpcs, result.Vpcs...)
	for {
		if result.NextToken == nil {
			break
		}
		input.NextToken = result.NextToken
		result, err = svc.DescribeVpcs(context.TODO(), input)
		if err != nil {
			fmt.Println(err)
		}
		vpcs = append(vpcs, result.Vpcs...)
	}
	return vpcs
}

type VPCToSubnet struct {
	VpcID   string
	Subnets []types.Subnet
}

func GetSubnetForVPCS(s aws.Config, vpcs []types.Vpc) []VPCToSubnet {
	svc := ec2.NewFromConfig(s)
	var vpcSubnets []VPCToSubnet
	for _, vpc := range vpcs {
		input := &ec2.DescribeSubnetsInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []string{*vpc.VpcId},
				},
			},
		}
		result, err := svc.DescribeSubnets(context.TODO(), input)
		if err != nil {
			fmt.Println(err)
		}
		vpcSubnets = append(vpcSubnets, VPCToSubnet{
			VpcID:   *vpc.VpcId,
			Subnets: result.Subnets,
		})
		for {
			if result.NextToken == nil {
				break
			}
			input.NextToken = result.NextToken
			result, err = svc.DescribeSubnets(context.TODO(), input)
			if err != nil {
				fmt.Println(err)
			}
			vpcSubnets = append(vpcSubnets, VPCToSubnet{
				VpcID:   *vpc.VpcId,
				Subnets: result.Subnets,
			})
		}
	}
	return vpcSubnets
}

type VpcToInternetGateway struct {
	VpcID            string
	InternetGateways []types.InternetGateway
}

func GetInternetGatewaysForVpc(s aws.Config, vpcs []types.Vpc) []VpcToInternetGateway {
	svc := ec2.NewFromConfig(s)
	var vpcInternetGateways []VpcToInternetGateway
	for _, vpc := range vpcs {
		input := &ec2.DescribeInternetGatewaysInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("attachment.vpc-id"),
					Values: []string{*vpc.VpcId},
				},
			},
		}
		result, err := svc.DescribeInternetGateways(context.TODO(), input)
		if err != nil {
			fmt.Println(err)
		}
		vpcInternetGateways = append(vpcInternetGateways, VpcToInternetGateway{
			VpcID:            *vpc.VpcId,
			InternetGateways: result.InternetGateways,
		})
		for {
			if result.NextToken == nil {
				break
			}
			input.NextToken = result.NextToken
			result, err = svc.DescribeInternetGateways(context.TODO(), input)
			if err != nil {
				fmt.Println(err)
			}
			vpcInternetGateways = append(vpcInternetGateways, VpcToInternetGateway{
				VpcID:            *vpc.VpcId,
				InternetGateways: result.InternetGateways,
			})
		}
	}
	return vpcInternetGateways
}

type VpcToFlowLogs struct {
	VpcID    string
	FlowLogs []types.FlowLog
}

func GetFlowLogsForVpc(s aws.Config, vpcs []types.Vpc) []VpcToFlowLogs {
	svc := ec2.NewFromConfig(s)
	var vpcFlowLogs []VpcToFlowLogs
	for _, vpc := range vpcs {
		input := &ec2.DescribeFlowLogsInput{
			Filter: []types.Filter{
				{
					Name:   aws.String("resource-id"),
					Values: []string{*vpc.VpcId},
				},
			},
		}
		result, err := svc.DescribeFlowLogs(context.TODO(), input)
		if err != nil {
			fmt.Println(err)
		}
		vpcFlowLogs = append(vpcFlowLogs, VpcToFlowLogs{
			VpcID:    *vpc.VpcId,
			FlowLogs: result.FlowLogs,
		})
		for {
			if result.NextToken == nil {
				break
			}
			input.NextToken = result.NextToken
			result, err = svc.DescribeFlowLogs(context.TODO(), input)
			if err != nil {
				fmt.Println(err)
			}
			vpcFlowLogs = append(vpcFlowLogs, VpcToFlowLogs{
				VpcID:    *vpc.VpcId,
				FlowLogs: result.FlowLogs,
			})
		}
	}
	return vpcFlowLogs
}

type VpcToSecurityGroups struct {
	vpcID          string
	securityGroups []types.SecurityGroup
}

func GetSecurityGroupForVpc(s aws.Config, vpcs []types.Vpc) []VpcToSecurityGroups {
	svc := ec2.NewFromConfig(s)
	var vpcSecurityGroups []VpcToSecurityGroups
	for _, vpc := range vpcs {
		input := &ec2.DescribeSecurityGroupsInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []string{*vpc.VpcId},
				},
			},
		}
		result, err := svc.DescribeSecurityGroups(context.TODO(), input)
		if err != nil {
			fmt.Println(err)
		}
		vpcSecurityGroups = append(vpcSecurityGroups, VpcToSecurityGroups{
			vpcID:          *vpc.VpcId,
			securityGroups: result.SecurityGroups,
		})
		for {
			if result.NextToken == nil {
				break
			}
			input.NextToken = result.NextToken
			result, err = svc.DescribeSecurityGroups(context.TODO(), input)
			if err != nil {
				fmt.Println(err)
			}
			vpcSecurityGroups = append(vpcSecurityGroups, VpcToSecurityGroups{
				vpcID:          *vpc.VpcId,
				securityGroups: result.SecurityGroups,
			})
		}
	}
	return vpcSecurityGroups
}

type SGToSecurityGroupRules struct {
	vpcID              string
	securityGroup      types.SecurityGroup
	securityGroupRules []types.SecurityGroupRule
}

func GetSecurityGroupRulesForSg(s aws.Config, vpc2SGs []VpcToSecurityGroups) []SGToSecurityGroupRules {
	svc := ec2.NewFromConfig(s)
	var sg2SGRules []SGToSecurityGroupRules
	// For each vpc-id
	for _, vpc2sg := range vpc2SGs {
		// For each Security group
		for _, sg := range vpc2sg.securityGroups {
			input := &ec2.DescribeSecurityGroupRulesInput{
				Filters: []types.Filter{
					{
						Name:   aws.String("group-id"),
						Values: []string{*sg.GroupId},
					},
				},
			}
			result, err := svc.DescribeSecurityGroupRules(context.TODO(), input)
			if err != nil {
				fmt.Println(err)
			}
			sg2SGRules = append(sg2SGRules, SGToSecurityGroupRules{
				vpcID:              vpc2sg.vpcID,
				securityGroup:      sg,
				securityGroupRules: result.SecurityGroupRules,
			})
			for {
				if result.NextToken == nil {
					break
				}
				input.NextToken = result.NextToken
				result, err := svc.DescribeSecurityGroupRules(context.TODO(), input)
				if err != nil {
					fmt.Println(err)
				}
				sg2SGRules = append(sg2SGRules, SGToSecurityGroupRules{
					vpcID:              vpc2sg.vpcID,
					securityGroup:      sg,
					securityGroupRules: result.SecurityGroupRules,
				})
			}
		}
	}
	return sg2SGRules
}
