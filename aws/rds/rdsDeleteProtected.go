package rds

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfDeleteProtectionEnabled(checkConfig commons.CheckConfig, instances []types.DBInstance, testName string) {
	var check commons.Check
	check.InitCheck("RDS have the deletion protection enabled", "Check if RDS delete protection is enabled", testName, []string{"Security", "Good Practice"})
	for _, instance := range instances {
		if aws.ToBool(instance.DeletionProtection) {
			Message := "RDS delete protection is enabled on " + *instance.DBInstanceIdentifier
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		} else {
			Message := "RDS delete protection is not enabled on " + *instance.DBInstanceIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
