package cloudtrail

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfCloudtrailIsEnabled(checkConfig commons.CheckConfig, trailStatus []cloudtrail.GetTrailStatusOutput, testName string) {
	var check commons.Check
	check.InitCheck("Cloudtrail is enabled", "Check if Cloudtrail is enabled", testName, []string{"Security", "Good Practice"})
	result := commons.Result{Status: "FAIL", Message: "Cloudtrail is not enabled"}

	if len(trailStatus) == 0 {
		check.AddResult(result)
		checkConfig.Queue <- check
		return
	}

	for _, trail := range trailStatus {
		if aws.ToBool(trail.IsLogging) {
			result = commons.Result{Status: "OK", Message: "Cloudtrail is enabled"}
			break
		}
	}

	check.AddResult(result)
	checkConfig.Queue <- check
}
