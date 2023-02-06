package cloudtrail

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfCloudtrailsEncrypted(checkConfig commons.CheckConfig, cloudtrails []types.Trail, testName string) {

	var check commons.Check
	check.InitCheck("Cloudtrails are encrypted", "check if all cloudtrails are encrypted", testName, []string{"Security", "Good Practice"})
	for _, cloudtrail := range cloudtrails {
		if cloudtrail.KmsKeyId == nil || *cloudtrail.KmsKeyId == "" {
			Message := "Cloudtrail " + *cloudtrail.Name + " is not encrypted"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		} else {
			Message := "Cloudtrail " + *cloudtrail.Name + " is encrypted"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
