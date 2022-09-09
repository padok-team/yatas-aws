package volumes

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/plugins/commons"
)

func CheckIfVolumeIsUsed(checkConfig commons.CheckConfig, volumes []types.Volume, testName string) {
	var check commons.Check
	check.InitCheck("EC2's volumes are unused", "Check if EC2 volumes are unused", testName)
	for _, volume := range volumes {
		if volume.State != types.VolumeStateInUse && volume.State != types.VolumeStateDeleted {
			Message := "EC2 volume is unused " + *volume.VolumeId
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		} else if volume.State == types.VolumeStateDeleted {
			continue
		} else {
			Message := "EC2 volume is in use " + *volume.VolumeId
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
