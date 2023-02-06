package ec2

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfMonitoringEnabled(checkConfig commons.CheckConfig, instances []types.Instance, testName string) {
	var check commons.Check
	check.InitCheck("EC2s have the monitoring option enabled", "Check if all instances have monitoring enabled", testName, []string{"Security", "Good Practice"})
	for _, instance := range instances {
		if instance.Monitoring.State != types.MonitoringStateEnabled {
			Message := "EC2 instance " + *instance.InstanceId + " has no monitoring enabled"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *instance.InstanceId}
			check.AddResult(result)
		} else {
			Message := "EC2 instance " + *instance.InstanceId + " has monitoring enabled"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *instance.InstanceId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
