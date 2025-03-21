package rds

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func checkIfClusterSnapshotEncryptionEnabled(checkConfig commons.CheckConfig, snapshots []types.DBClusterSnapshot, testName string) {
	var check commons.Check
	check.InitCheck("Aurora snapshots are encrypted", "Check if Aurora snapshot encryption is enabled", testName, []string{"Security", "Good Practice"})
	for _, snapshot := range snapshots {
		if !aws.ToBool(snapshot.StorageEncrypted) {
			Message := "Aurora snapshot encryption is not enabled on " + *snapshot.DBClusterIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *snapshot.DBClusterSnapshotArn}
			check.AddResult(result)
		} else {
			Message := "Aurora snapshot encryption is enabled on " + *snapshot.DBClusterIdentifier
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *snapshot.DBClusterSnapshotArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
