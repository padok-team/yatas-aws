package volumes

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfAllSnapshotsEncrypted(checkConfig commons.CheckConfig, snapshots []types.Snapshot, testName string) {
	var check commons.Check
	check.InitCheck("EC2's Snapshots are encrypted", "Check if all snapshots are encrypted", testName, []string{"Security", "Good Practice", "HDS"})
	for _, snapshot := range snapshots {
		if snapshot.Encrypted == nil || !*snapshot.Encrypted {
			Message := "Snapshot " + *snapshot.SnapshotId + " is not encrypted"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *snapshot.SnapshotId}
			check.AddResult(result)
		} else {
			Message := "Snapshot " + *snapshot.SnapshotId + " is encrypted"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *snapshot.SnapshotId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
