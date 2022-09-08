package rds

import (
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/config"
)

func checkIfBackupEnabled(checkConfig config.CheckConfig, instances []types.DBInstance, testName string) {
	var check config.Check
	check.InitCheck("RDS are backedup automatically with PITR", "Check if RDS backup is enabled", testName)
	for _, instance := range instances {
		if instance.BackupRetentionPeriod == 0 {
			Message := "RDS backup is not enabled on " + *instance.DBInstanceIdentifier
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		} else {
			Message := "RDS backup is enabled on " + *instance.DBInstanceIdentifier
			result := config.Result{Status: "OK", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
