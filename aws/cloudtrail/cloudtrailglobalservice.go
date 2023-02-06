package cloudtrail

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfCloudtrailsGlobalServiceEventsEnabled(checkConfig commons.CheckConfig, cloudtrails []types.Trail, testName string) {
	var check commons.Check
	check.InitCheck("Cloudtrails have Global Service Events Activated", "check if all cloudtrails have global service events enabled", testName, []string{"Security", "Good Practice"})
	for _, cloudtrail := range cloudtrails {
		if !*cloudtrail.IncludeGlobalServiceEvents {
			Message := "Cloudtrail " + *cloudtrail.Name + " has global service events disabled"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		} else {
			Message := "Cloudtrail " + *cloudtrail.Name + " has global service events enabled"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
