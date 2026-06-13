package documentdb

import (
	"github.com/padok-team/yatas/plugins/commons"
)

// checkIfDocumentDBTLSEnabled verifies that the "tls" parameter is set to "enabled"
// in the cluster parameter group of each DocumentDB cluster.
// When TLS is enforced, clients must use an SSL/TLS connection, preventing
// unencrypted data from being transmitted over the network.
func checkIfDocumentDBTLSEnabled(checkConfig commons.CheckConfig, clusters []ClusterWithTLSParam, testName string) {
	var check commons.Check
	check.InitCheck("DocumentDB clusters enforce TLS connections", "Check if DocumentDB TLS parameter is enabled in the cluster parameter group", testName, []string{"Security", "Good Practice"})
	for _, c := range clusters {
		if c.TLSValue == "enabled" {
			Message := "DocumentDB TLS is enabled on " + *c.Cluster.DBClusterIdentifier
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *c.Cluster.DBClusterArn}
			check.AddResult(result)
		} else {
			Message := "DocumentDB TLS is not enabled on " + *c.Cluster.DBClusterIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *c.Cluster.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
