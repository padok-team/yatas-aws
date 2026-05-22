package documentdb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/docdb/types"
	"github.com/padok-team/yatas/plugins/commons"
)

// checkIfDocumentDBEncryptionEnabled verifies that storage encryption (AES-256 via KMS)
// is enabled on each DocumentDB cluster. Encryption at rest protects data stored on disk,
// including automated backups, read replicas, and snapshots.
func checkIfDocumentDBEncryptionEnabled(checkConfig commons.CheckConfig, clusters []types.DBCluster, testName string) {
	var check commons.Check
	check.InitCheck("DocumentDB clusters are encrypted at rest", "Check if DocumentDB storage encryption is enabled", testName, []string{"Security", "Good Practice"})
	for _, cluster := range clusters {
		if !aws.ToBool(cluster.StorageEncrypted) {
			Message := "DocumentDB encryption is not enabled on " + *cluster.DBClusterIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *cluster.DBClusterArn}
			check.AddResult(result)
		} else {
			Message := "DocumentDB encryption is enabled on " + *cluster.DBClusterIdentifier
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *cluster.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
