package rds

import (
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfDeleteProtectionEnabled(checkConfig config.CheckConfig, instances []types.DBInstance, testName string) {
	var check config.Check
	check.InitCheck("RDS have the deletion protection enabled", "Check if RDS delete protection is enabled", testName)
	for _, instance := range instances {
		if instance.DeletionProtection {
			Message := "RDS delete protection is enabled on " + *instance.DBInstanceIdentifier
			result := config.Result{Status: "OK", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		} else {
			Message := "RDS delete protection is not enabled on " + *instance.DBInstanceIdentifier
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
