package documentdb

import (
	"github.com/aws/aws-sdk-go-v2/service/docdb/types"
	"github.com/padok-team/yatas/plugins/commons"
)

// checkIfDocumentDBMultiAZEnabled verifies that each DocumentDB cluster has Multi-AZ
// enabled. Multi-AZ deploys read replicas across multiple Availability Zones, ensuring
// automatic failover in case of an AZ outage and improving data durability.
func checkIfDocumentDBMultiAZEnabled(checkConfig commons.CheckConfig, clusters []types.DBCluster, testName string) {
	var check commons.Check
	check.InitCheck("DocumentDB clusters are deployed in multiple availability zones", "Check if DocumentDB multi-AZ is enabled for high availability", testName, []string{"Good Practice"})
	for _, cluster := range clusters {
		if cluster.MultiAZ != nil && *cluster.MultiAZ {
			Message := "DocumentDB multi-AZ is enabled on " + *cluster.DBClusterIdentifier
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *cluster.DBClusterArn}
			check.AddResult(result)
		} else {
			Message := "DocumentDB multi-AZ is not enabled on " + *cluster.DBClusterIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *cluster.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
