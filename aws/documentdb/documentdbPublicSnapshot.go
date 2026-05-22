package documentdb

import (
	"github.com/padok-team/yatas/plugins/commons"
)

// checkIfDocumentDBSnapshotNotPublic verifies that no manual DocumentDB snapshot is
// publicly accessible. A public snapshot can be restored by any AWS account,
// leading to a potential full data exfiltration. This is a critical security risk.
func checkIfDocumentDBSnapshotNotPublic(checkConfig commons.CheckConfig, snapshots []SnapshotWithPublicFlag, testName string) {
	var check commons.Check
	check.InitCheck("DocumentDB manual snapshots are not publicly accessible", "Check if DocumentDB manual snapshots are private", testName, []string{"Security"})
	for _, s := range snapshots {
		if s.IsPublic {
			Message := "DocumentDB snapshot is publicly accessible: " + *s.Snapshot.DBClusterSnapshotIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *s.Snapshot.DBClusterSnapshotArn}
			check.AddResult(result)
		} else {
			Message := "DocumentDB snapshot is not publicly accessible: " + *s.Snapshot.DBClusterSnapshotIdentifier
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *s.Snapshot.DBClusterSnapshotArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
