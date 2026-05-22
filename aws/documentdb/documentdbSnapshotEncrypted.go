package documentdb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/docdb/types"
	"github.com/padok-team/yatas/plugins/commons"
)

// checkIfDocumentDBSnapshotEncrypted verifies that all manual DocumentDB snapshots
// are encrypted at rest. Unencrypted snapshots expose data at rest and may violate
// compliance requirements (PCI DSS, HDS, SOC 2).
func checkIfDocumentDBSnapshotEncrypted(checkConfig commons.CheckConfig, snapshots []types.DBClusterSnapshot, testName string) {
	var check commons.Check
	check.InitCheck("DocumentDB manual snapshots are encrypted", "Check if DocumentDB manual snapshot encryption is enabled", testName, []string{"Security", "Good Practice"})
	for _, snapshot := range snapshots {
		if !aws.ToBool(snapshot.StorageEncrypted) {
			Message := "DocumentDB snapshot encryption is not enabled on " + *snapshot.DBClusterSnapshotIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *snapshot.DBClusterSnapshotArn}
			check.AddResult(result)
		} else {
			Message := "DocumentDB snapshot encryption is enabled on " + *snapshot.DBClusterSnapshotIdentifier
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *snapshot.DBClusterSnapshotArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
