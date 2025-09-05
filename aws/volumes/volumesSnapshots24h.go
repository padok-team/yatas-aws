package volumes

import (
	"time"

	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfSnapshotYoungerthan24h(checkConfig commons.CheckConfig, vs couple, testName string) {
	var check commons.Check
	check.InitCheck("EC2's snapshots are younger than a day old", "Check if all snapshots are younger than 24h", testName, []string{"Security", "Good Practice", "HDS"})
	for _, volume := range vs.Volume {
		snapshotYoungerThan24h := false
		for _, snapshot := range vs.Snapshot {
			if *snapshot.VolumeId == *volume.VolumeId {
				creationTime := *snapshot.StartTime
				if creationTime.After(time.Now().Add(-24 * time.Hour)) {
					snapshotYoungerThan24h = true
					break
				}
			}
		}
		if !snapshotYoungerThan24h {
			Message := "Volume " + *volume.VolumeId + " has no snapshot younger than 24h"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		} else {
			Message := "Volume " + *volume.VolumeId + " has snapshot younger than 24h"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
