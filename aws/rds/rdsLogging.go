package rds

import (
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfLoggingEnabled(checkConfig commons.CheckConfig, instances []types.DBInstance, testName string) {
	var check commons.Check
	check.InitCheck("RDS logs are exported to cloudwatch", "Check if RDS logging is enabled", testName, []string{"Security", "Good Practice"})
	for _, instance := range instances {
		if instance.EnabledCloudwatchLogsExports != nil {
			found := false
			for _, export := range instance.EnabledCloudwatchLogsExports {
				if export == "audit" {
					Message := "RDS logging is enabled on " + *instance.DBInstanceIdentifier
					result := commons.Result{Status: "OK", Message: Message, ResourceID: *instance.DBInstanceArn}
					check.AddResult(result)
					found = true

					break

				}
			}
			if !found {
				Message := "RDS logging is not enabled on " + *instance.DBInstanceIdentifier
				result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
				check.AddResult(result)
				continue
			}
		} else {
			Message := "RDS logging is not enabled on " + *instance.DBInstanceIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
