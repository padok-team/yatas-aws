package autoscaling

import (
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/stangirard/yatas/plugins/commons"
)

func CheckIfInTwoAvailibilityZones(checkConfig commons.CheckConfig, groups []types.AutoScalingGroup, testName string) {
	var check commons.Check
	check.InitCheck("Autoscaling group are in two availability zones", "Check if all autoscaling groups have at least two availability zones", testName)
	for _, group := range groups {
		if len(group.AvailabilityZones) < 2 {
			Message := "Autoscaling group " + *group.AutoScalingGroupName + " has less than two availability zones"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *group.AutoScalingGroupName}
			check.AddResult(result)
		} else {
			Message := "Autoscaling group " + *group.AutoScalingGroupName + " has two availability zones"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *group.AutoScalingGroupName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
