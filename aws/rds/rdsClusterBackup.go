package rds

import (
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/plugins/commons"
)

func checkIfClusterBackupEnabled(checkConfig commons.CheckConfig, instances []types.DBCluster, testName string) {
	var check commons.Check
	check.InitCheck("Aurora RDS are backedup automatically with PITR", "Check if Aurora RDS backup is enabled", testName)
	for _, instance := range instances {
		if instance.BackupRetentionPeriod == nil || *instance.BackupRetentionPeriod == 0 {
			Message := "RDS backup is not enabled on " + *instance.DBClusterIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		} else {
			Message := "RDS backup is enabled on " + *instance.DBClusterIdentifier
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
