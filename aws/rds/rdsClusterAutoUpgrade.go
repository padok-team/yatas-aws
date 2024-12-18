package rds

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func checkIfClusterAutoUpgradeEnabled(checkConfig commons.CheckConfig, instances []types.DBCluster, testName string) {
	var check commons.Check
	check.InitCheck("Aurora Clusters have minor versions automatically updated", "Check if Aurora RDS minor auto upgrade is enabled", testName, []string{"Security", "Good Practice"})
	for _, instance := range instances {
		if !aws.ToBool(instance.AutoMinorVersionUpgrade) {
			Message := "RDS auto upgrade is not enabled on " + *instance.DBClusterIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		} else {
			Message := "RDS auto upgrade is enabled on " + *instance.DBClusterIdentifier
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
