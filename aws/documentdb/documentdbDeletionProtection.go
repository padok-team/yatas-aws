package documentdb

import (
	"github.com/aws/aws-sdk-go-v2/service/docdb/types"
	"github.com/padok-team/yatas/plugins/commons"
)

// checkIfDocumentDBDeletionProtectionEnabled verifies that deletion protection is
// enabled on each DocumentDB cluster. When enabled, the cluster cannot be deleted
// accidentally — it must be explicitly disabled before a deletion can succeed.
func checkIfDocumentDBDeletionProtectionEnabled(checkConfig commons.CheckConfig, clusters []types.DBCluster, testName string) {
	var check commons.Check
	check.InitCheck("DocumentDB clusters have deletion protection enabled", "Check if DocumentDB deletion protection is enabled", testName, []string{"Security", "Good Practice"})
	for _, cluster := range clusters {
		if cluster.DeletionProtection != nil && *cluster.DeletionProtection {
			Message := "DocumentDB deletion protection is enabled on " + *cluster.DBClusterIdentifier
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *cluster.DBClusterArn}
			check.AddResult(result)
		} else {
			Message := "DocumentDB deletion protection is not enabled on " + *cluster.DBClusterIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *cluster.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
