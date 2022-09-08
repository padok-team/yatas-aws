package cloudtrail

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfCloudtrailsMultiRegion(checkConfig config.CheckConfig, cloudtrails []types.Trail, testName string) {
	var check config.Check
	check.InitCheck("Cloudtrails are in multiple regions", "check if all cloudtrails are multi region", testName)
	for _, cloudtrail := range cloudtrails {
		if !*cloudtrail.IsMultiRegionTrail {
			Message := "Cloudtrail " + *cloudtrail.Name + " is not multi region"
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		} else {
			Message := "Cloudtrail " + *cloudtrail.Name + " is multi region"
			result := config.Result{Status: "OK", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
