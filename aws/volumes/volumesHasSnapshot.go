package volumes

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfAllVolumesHaveSnapshots(checkConfig commons.CheckConfig, snapshot2Volumes couple, testName string) {
	var check commons.Check
	check.InitCheck("EC2 have snapshots", "Check if all volumes have snapshots", testName, []string{"Security", "Good Practice"})
	for _, volume := range snapshot2Volumes.Volume {
		ok := false
		for _, snapshot := range snapshot2Volumes.Snapshot {
			if *snapshot.VolumeId == *volume.VolumeId {
				Message := "Volume " + *volume.VolumeId + " has snapshot " + *snapshot.SnapshotId
				result := commons.Result{Status: "OK", Message: Message, ResourceID: *volume.VolumeId}
				check.AddResult(result)
				ok = true
				break
			}
		}
		if !ok {
			Message := "Volume " + *volume.VolumeId + " has no snapshot"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
