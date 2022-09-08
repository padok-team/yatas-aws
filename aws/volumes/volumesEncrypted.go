package volumes

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/config"
)

func checkIfEncryptionEnabled(checkConfig config.CheckConfig, volumes []types.Volume, testName string) {
	var check config.Check
	check.InitCheck("EC2's volumes are encrypted", "Check if EC2 encryption is enabled", testName)
	for _, volume := range volumes {
		if volume.Encrypted != nil && *volume.Encrypted {
			Message := "EC2 encryption is enabled on " + *volume.VolumeId
			result := config.Result{Status: "OK", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		} else {
			Message := "EC2 encryption is not enabled on " + *volume.VolumeId
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
