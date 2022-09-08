package rds

import (
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfClusterLoggingEnabled(checkConfig config.CheckConfig, instances []types.DBCluster, testName string) {
	var check config.Check
	check.InitCheck("Aurora RDS logs are exported to cloudwatch", "Check if Aurora RDS logging is enabled", testName)
	for _, instance := range instances {
		if instance.EnabledCloudwatchLogsExports != nil {
			found := false
			for _, export := range instance.EnabledCloudwatchLogsExports {
				if export == "audit" {
					Message := "RDS logging is enabled on " + *instance.DBClusterIdentifier
					result := config.Result{Status: "OK", Message: Message, ResourceID: *instance.DBClusterArn}
					check.AddResult(result)
					found = true

					break

				}
			}
			if !found {
				Message := "RDS logging is not enabled on " + *instance.DBClusterIdentifier
				result := config.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBClusterArn}
				check.AddResult(result)
				continue
			}
		} else {
			Message := "RDS logging is not enabled on " + *instance.DBClusterIdentifier
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
