package cloudtrail

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfCloudtrailsEncrypted(checkConfig config.CheckConfig, cloudtrails []types.Trail, testName string) {

	var check config.Check
	check.InitCheck("Cloudtrails are encrypted", "check if all cloudtrails are encrypted", testName)
	for _, cloudtrail := range cloudtrails {
		if cloudtrail.KmsKeyId == nil || *cloudtrail.KmsKeyId == "" {
			Message := "Cloudtrail " + *cloudtrail.Name + " is not encrypted"
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		} else {
			Message := "Cloudtrail " + *cloudtrail.Name + " is encrypted"
			result := config.Result{Status: "OK", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
