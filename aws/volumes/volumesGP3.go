package volumes

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/plugins/commons"
)

func CheckIfVolumesTypeGP3(checkConfig commons.CheckConfig, volumes []types.Volume, testName string) {
	var check commons.Check
	check.InitCheck("EC2 are using GP3", "Check if all volumes are of type gp3", testName)
	for _, volume := range volumes {
		if volume.VolumeType != "gp3" {
			Message := "Volume " + *volume.VolumeId + " is not of type gp3"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		} else {
			Message := "Volume " + *volume.VolumeId + " is of type gp3"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
