package dynamodb

import (
	"github.com/stangirard/yatas/plugins/commons"
)

func CheckIfDynamodbContinuousBackupsEnabled(checkConfig commons.CheckConfig, dynamodbs []TableBackups, testName string) {
	var check commons.Check
	check.InitCheck("Dynamodb have continuous backup enabled with PITR", "Check if DynamoDB continuous backups are enabled", testName)
	for _, d := range dynamodbs {
		if d.Backups.ContinuousBackupsStatus != "ENABLED" {
			Message := "Dynamodb continuous backups are not enabled on " + d.TableName
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: d.TableName}
			check.AddResult(result)
		} else {
			Message := "Dynamodb continuous backups are enabled on " + d.TableName
			result := commons.Result{Status: "OK", Message: Message, ResourceID: d.TableName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
