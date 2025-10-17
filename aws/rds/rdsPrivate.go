package rds

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func checkIfRDSPrivateEnabled(checkConfig commons.CheckConfig, instances []types.DBInstance, testName string) {
	var check commons.Check
	check.InitCheck("RDS aren't publicly accessible", "Check if RDS private is enabled", testName, []string{"Security", "Good Practice", "HDS"})
	for _, instance := range instances {
		if aws.ToBool(instance.PubliclyAccessible) {
			Message := "RDS private is not enabled on " + *instance.DBInstanceIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		} else {
			Message := "RDS private is enabled on " + *instance.DBInstanceIdentifier
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
