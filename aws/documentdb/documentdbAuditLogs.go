package documentdb

import (
	"github.com/aws/aws-sdk-go-v2/service/docdb/types"
	"github.com/padok-team/yatas/plugins/commons"
)

// checkIfDocumentDBAuditLogsEnabled verifies that audit log export to CloudWatch Logs
// is enabled on each DocumentDB cluster. Audit logs record all authenticated operations
// (connections, disconnections, and queries) and are essential for compliance and
// security incident investigations.
func checkIfDocumentDBAuditLogsEnabled(checkConfig commons.CheckConfig, clusters []types.DBCluster, testName string) {
	var check commons.Check
	check.InitCheck("DocumentDB clusters export audit logs to CloudWatch", "Check if DocumentDB audit log export to CloudWatch is enabled", testName, []string{"Security", "Good Practice"})
	for _, cluster := range clusters {
		found := false
		for _, export := range cluster.EnabledCloudwatchLogsExports {
			if export == "audit" {
				found = true
				break
			}
		}
		if found {
			Message := "DocumentDB audit logs are enabled on " + *cluster.DBClusterIdentifier
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *cluster.DBClusterArn}
			check.AddResult(result)
		} else {
			Message := "DocumentDB audit logs are not enabled on " + *cluster.DBClusterIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *cluster.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
