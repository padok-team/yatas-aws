package eks

import (
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfLoggingIsEnabled(checkConfig config.CheckConfig, clusters []types.Cluster, testName string) {
	var check config.Check
	check.InitCheck("EKS clusters have logging enabled", "Check if logging is enabled", testName)
	for _, cluster := range clusters {
		if cluster.Logging != nil && len(cluster.Logging.ClusterLogging) > 0 {
			Message := "Logging is enabled for cluster " + *cluster.Name
			result := config.Result{Status: "OK", Message: Message, ResourceID: *cluster.Name}
			check.AddResult(result)
		} else {
			Message := "Logging is not enabled for cluster " + *cluster.Name
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *cluster.Name}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
