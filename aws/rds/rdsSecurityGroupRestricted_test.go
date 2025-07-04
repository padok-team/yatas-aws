package rds

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	rdsTypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfRDSRestrictedSecurityGroupsRestricted(t *testing.T) {
	checkConfig := commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}}
	instances := []InstanceWithSGs{
		{
			Instance: rdsTypes.DBInstance{
				DBInstanceIdentifier: aws.String("db-restricted"),
				DBInstanceArn:        aws.String("arn:db-restricted"),
			},
			SecurityGroups: []ec2Types.SecurityGroup{
				{
					GroupId: aws.String("sg-restricted"),
					IpPermissions: []ec2Types.IpPermission{
						{
							FromPort: aws.Int32(5432),
							ToPort:   aws.Int32(5432),
							IpRanges: []ec2Types.IpRange{{CidrIp: aws.String("10.0.0.0/16")}},
						},
					},
				},
			},
		},
	}

	checkConfig.Wg.Add(1)
	go checkIfRDSRestrictedSecurityGroups(checkConfig, instances, "AWS_RDS_016")
	go func() {
		for check := range checkConfig.Queue {
			for _, r := range check.Results {
				if r.Status != "OK" {
					t.Errorf("Expected OK, got %s", r.Status)
				}
			}
			checkConfig.Wg.Done()
		}
	}()
	checkConfig.Wg.Wait()
}

func TestCheckIfRDSRestrictedSecurityGroupsOpenIP(t *testing.T) {
	checkConfig := commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}}
	instances := []InstanceWithSGs{
		{
			Instance: rdsTypes.DBInstance{
				DBInstanceIdentifier: aws.String("db-open-ip"),
				DBInstanceArn:        aws.String("arn:db-open-ip"),
			},
			SecurityGroups: []ec2Types.SecurityGroup{
				{
					GroupId: aws.String("sg-open-ip"),
					IpPermissions: []ec2Types.IpPermission{
						{
							FromPort: aws.Int32(3306),
							ToPort:   aws.Int32(3306),
							IpRanges: []ec2Types.IpRange{{CidrIp: aws.String("0.0.0.0/0")}},
						},
					},
				},
			},
		},
	}

	checkConfig.Wg.Add(1)
	go checkIfRDSRestrictedSecurityGroups(checkConfig, instances, "AWS_RDS_016")
	go func() {
		for check := range checkConfig.Queue {
			for _, r := range check.Results {
				if r.Status != "FAIL" || r.Message != "RDS db-open-ip has SG opened to 0.0.0.0/0" {
					t.Errorf("Expected FAIL with open IP, got %s (%s)", r.Status, r.Message)
				}
			}
			checkConfig.Wg.Done()
		}
	}()
	checkConfig.Wg.Wait()
}

func TestCheckIfRDSRestrictedSecurityGroupsOpenPorts(t *testing.T) {
	checkConfig := commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}}
	instances := []InstanceWithSGs{
		{
			Instance: rdsTypes.DBInstance{
				DBInstanceIdentifier: aws.String("db-open-ports"),
				DBInstanceArn:        aws.String("arn:db-open-ports"),
			},
			SecurityGroups: []ec2Types.SecurityGroup{
				{
					GroupId: aws.String("sg-open-ports"),
					IpPermissions: []ec2Types.IpPermission{
						{
							FromPort: aws.Int32(0),
							ToPort:   aws.Int32(65535),
							IpRanges: []ec2Types.IpRange{{CidrIp: aws.String("192.168.0.0/16")}},
						},
					},
				},
			},
		},
	}

	checkConfig.Wg.Add(1)
	go checkIfRDSRestrictedSecurityGroups(checkConfig, instances, "AWS_RDS_016")
	go func() {
		for check := range checkConfig.Queue {
			for _, r := range check.Results {
				if r.Status != "FAIL" || r.Message != "RDS db-open-ports has SG with all ports opened" {
					t.Errorf("Expected FAIL with all ports open, got %s (%s)", r.Status, r.Message)
				}
			}
			checkConfig.Wg.Done()
		}
	}()
	checkConfig.Wg.Wait()
}

func TestCheckIfRDSRestrictedSecurityGroupsOpenBoth(t *testing.T) {
	checkConfig := commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}}
	instances := []InstanceWithSGs{
		{
			Instance: rdsTypes.DBInstance{
				DBInstanceIdentifier: aws.String("db-open-both"),
				DBInstanceArn:        aws.String("arn:db-open-both"),
			},
			SecurityGroups: []ec2Types.SecurityGroup{
				{
					GroupId: aws.String("sg-open-both"),
					IpPermissions: []ec2Types.IpPermission{
						{
							FromPort: aws.Int32(0),
							ToPort:   aws.Int32(65535),
							IpRanges: []ec2Types.IpRange{{CidrIp: aws.String("0.0.0.0/0")}},
						},
					},
				},
			},
		},
	}

	checkConfig.Wg.Add(1)
	go checkIfRDSRestrictedSecurityGroups(checkConfig, instances, "AWS_RDS_016")
	go func() {
		for check := range checkConfig.Queue {
			for _, r := range check.Results {
				expected := "RDS db-open-both has SG opened to 0.0.0.0/0 and all ports opened"
				if r.Status != "FAIL" || r.Message != expected {
					t.Errorf("Expected FAIL with both issues, got %s (%s)", r.Status, r.Message)
				}
			}
			checkConfig.Wg.Done()
		}
	}()
	checkConfig.Wg.Wait()
}
