package dynamodb

import (
	"github.com/stangirard/yatas/config"
)

func CheckIfDynamodbContinuousBackupsEnabled(checkConfig config.CheckConfig, dynamodbs []TableBackups, testName string) {
	var check config.Check
	check.InitCheck("Dynamodb have continuous backup enabled with PITR", "Check if DynamoDB continuous backups are enabled", testName)
	for _, d := range dynamodbs {
		if d.Backups.ContinuousBackupsStatus != "ENABLED" {
			Message := "Dynamodb continuous backups are not enabled on " + d.TableName
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: d.TableName}
			check.AddResult(result)
		} else {
			Message := "Dynamodb continuous backups are enabled on " + d.TableName
			result := config.Result{Status: "OK", Message: Message, ResourceID: d.TableName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
