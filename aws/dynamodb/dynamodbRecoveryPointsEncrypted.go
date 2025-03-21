package dynamodb

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfEncryptionDynamodbRecoveryPointsEnabled(checkConfig commons.CheckConfig, tableRecoveryPoints []TableRecoveryPoints, testName string) {
	var check commons.Check
	check.InitCheck("Dynamodb recovery points are encrypted", "Check if DynamoDB recovery point encryption is enabled", testName, []string{"Security", "Good Practice"})
	for _, t := range tableRecoveryPoints {
		for _, r := range t.RecoveryPoints {
			if *r.EncryptionKeyArn != "" {
				Message := "Dynamodb recovery point encryption is enabled on table " + t.TableName
				result := commons.Result{Status: "OK", Message: Message, ResourceID: *r.RecoveryPointArn}
				check.AddResult(result)

			} else {
				Message := "Dynamodb recovery point encryption is not enabled on table " + t.TableName
				result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *r.RecoveryPointArn}
				check.AddResult(result)
			}
		}
	}
	checkConfig.Queue <- check
}
