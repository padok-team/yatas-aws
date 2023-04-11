package ec2

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func Ec2MonitoringEnabledCondition(resource interface{}) bool {
	instance, ok := resource.(*types.Instance)
	if !ok {
		return false
	}
	return instance.Monitoring.State == types.MonitoringStateEnabled
}

func Ec2PublicIPCondition(resource interface{}) bool {
	instance, ok := resource.(*types.Instance)
	if !ok {
		return false
	}
	return instance.PublicIpAddress == nil
}

func Ec2RunningInVPCCondition(resource interface{}) bool {
	instance, ok := resource.(*types.Instance)
	if !ok {
		return false
	}
	return instance.VpcId != nil && *instance.VpcId != ""
}

// Generate unit test functions for each condition
