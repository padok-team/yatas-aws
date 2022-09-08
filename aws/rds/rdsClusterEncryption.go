package rds

import (
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/config"
)

func checkIfClusterEncryptionEnabled(checkConfig config.CheckConfig, instances []types.DBCluster, testName string) {
	var check config.Check
	check.InitCheck("Aurora RDS are encrypted", "Check if Aurora RDS encryption is enabled", testName)
	for _, instance := range instances {
		if !instance.StorageEncrypted {
			Message := "RDS encryption is not enabled on " + *instance.DBClusterIdentifier
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		} else {
			Message := "RDS encryption is enabled on " + *instance.DBClusterIdentifier
			result := config.Result{Status: "OK", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
