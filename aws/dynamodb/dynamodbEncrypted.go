package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stangirard/yatas/config"
)

func CheckIfDynamodbEncrypted(checkConfig config.CheckConfig, dynamodbs []*dynamodb.DescribeTableOutput, testName string) {
	var check config.Check
	check.InitCheck("Dynamodbs are encrypted", "Check if DynamoDB encryption is enabled", testName)
	for _, d := range dynamodbs {
		if d.Table != nil && d.Table.SSEDescription != nil && d.Table.SSEDescription.Status == "ENABLED" {
			Message := "Dynamodb encryption is enabled on " + *d.Table.TableName
			result := config.Result{Status: "OK", Message: Message, ResourceID: *d.Table.TableArn}
			check.AddResult(result)

		} else {
			Message := "Dynamodb encryption is not enabled on " + *d.Table.TableName
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *d.Table.TableArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
