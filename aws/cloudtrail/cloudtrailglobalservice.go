package cloudtrail

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfCloudtrailsGlobalServiceEventsEnabled(checkConfig config.CheckConfig, cloudtrails []types.Trail, testName string) {
	var check config.Check
	check.InitCheck("Cloudtrails have Global Service Events Activated", "check if all cloudtrails have global service events enabled", testName)
	for _, cloudtrail := range cloudtrails {
		if !*cloudtrail.IncludeGlobalServiceEvents {
			Message := "Cloudtrail " + *cloudtrail.Name + " has global service events disabled"
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		} else {
			Message := "Cloudtrail " + *cloudtrail.Name + " has global service events enabled"
			result := config.Result{Status: "OK", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
