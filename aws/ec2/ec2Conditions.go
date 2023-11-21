package ec2

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/padok-team/yatas-aws/aws/awschecks"
)

func Ec2MonitoringEnabledCondition(resource awschecks.Resource) bool {
	instanceResource, ok := resource.(*EC2Instance)
	if !ok {
		return false
	}
	return instanceResource.Instance.Monitoring.State == types.MonitoringStateEnabled
}

func Ec2PublicIPCondition(resource awschecks.Resource) bool {
	instanceResource, ok := resource.(*EC2Instance)
	if !ok {
		return false
	}
	return instanceResource.Instance.PublicIpAddress == nil
}

func Ec2RunningInVPCCondition(resource awschecks.Resource) bool {
	instanceResource, ok := resource.(*EC2Instance)
	if !ok {
		return false
	}
	return instanceResource.Instance.VpcId != nil && *instanceResource.Instance.VpcId != ""
}
