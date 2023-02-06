package ec2

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfEC2PublicIP(checkConfig commons.CheckConfig, instances []types.Instance, testName string) {
	var check commons.Check
	check.InitCheck("EC2s don't have a public IP", "Check if all instances have a public IP", testName, []string{"Security", "Good Practice"})
	for _, instance := range instances {
		if instance.PublicIpAddress != nil {
			Message := "EC2 instance " + *instance.InstanceId + " has a public IP" + *instance.PublicIpAddress
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *instance.InstanceId}
			check.AddResult(result)
		} else {
			Message := "EC2 instance " + *instance.InstanceId + " has no public IP "
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *instance.InstanceId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
