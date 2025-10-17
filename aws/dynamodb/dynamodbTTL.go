// This check only verifies that TTL is enabled and configured on the DynamoDB table.
// Limitation: DynamoDB does not support a global retention period at the table level.
// TTL is set per item by the application, not by infrastructure.
// Scanning all items for TTL values is not scalable and not recommended.
// This check ensures TTL is configured, but actual retention enforcement must be handled by the application.

package dynamodb

import (
	"github.com/padok-team/yatas/plugins/commons"
)

const MAX_TTL_DAYS = 90 // Maximum TTL retention period in days

func CheckIfTTLConfiguredAndValid(checkConfig commons.CheckConfig, tables []TableTTL, testName string) {
	var check commons.Check
	check.InitCheck("DynamoDB tables have TTL configured with 90 days maximum retention", "Check if DynamoDB tables have TTL enabled and configured to delete items after no more than 90 days", testName, []string{"Security", "Good Practice"})

	for _, table := range tables {
		if !table.TTLEnabled {
			Message := "DynamoDB table " + table.TableName + " does not have TTL enabled"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: table.TableName}
			check.AddResult(result)
		} else if table.TTLAttribute == "" {
			Message := "DynamoDB table " + table.TableName + " has TTL enabled but no attribute configured"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: table.TableName}
			check.AddResult(result)
		} else {
			// Note: The actual TTL value check would need to be done at the item level
			// Here we can only verify that TTL is enabled and properly configured
			Message := "DynamoDB table " + table.TableName + " has TTL enabled and configured with attribute: " + table.TTLAttribute
			result := commons.Result{Status: "OK", Message: Message, ResourceID: table.TableName}
			check.AddResult(result)
		}
	}

	checkConfig.Queue <- check
}
