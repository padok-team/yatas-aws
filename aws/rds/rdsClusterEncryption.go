package rds

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func checkIfClusterEncryptionEnabled(checkConfig commons.CheckConfig, instances []types.DBCluster, testName string) {
	var check commons.Check
	check.InitCheck("Aurora RDS are encrypted", "Check if Aurora RDS encryption is enabled", testName, []string{"Security", "Good Practice", "HDS"})
	for _, instance := range instances {
		if !aws.ToBool(instance.StorageEncrypted) {
			Message := "RDS encryption is not enabled on " + *instance.DBClusterIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		} else {
			Message := "RDS encryption is enabled on " + *instance.DBClusterIdentifier
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
