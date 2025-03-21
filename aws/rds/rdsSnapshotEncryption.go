package rds

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func checkIfRDSSnapshotEncryptionEnabled(checkConfig commons.CheckConfig, snapshots []types.DBSnapshot, testName string) {
	var check commons.Check
	check.InitCheck("RDS snapshots are encrypted", "Check if RDS snapshot encryption is enabled", testName, []string{"Security", "Good Practice"})
	for _, snapshot := range snapshots {
		if !aws.ToBool(snapshot.Encrypted) {
			Message := "RDS snapshot encryption is not enabled on " + *snapshot.DBInstanceIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *snapshot.DBSnapshotArn}
			check.AddResult(result)
		} else {
			Message := "RDS snapshot encryption is enabled on " + *snapshot.DBInstanceIdentifier
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *snapshot.DBSnapshotArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
