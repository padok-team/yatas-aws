package rds

import (
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfClusterDeleteProtectionEnabled(checkConfig config.CheckConfig, instances []types.DBCluster, testName string) {
	var check config.Check
	check.InitCheck("Aurora RDS have the deletion protection enabled", "Check if Aurora RDS delete protection is enabled", testName)
	for _, instance := range instances {
		if instance.DeletionProtection != nil && *instance.DeletionProtection {
			Message := "RDS delete protection is enabled on " + *instance.DBClusterIdentifier
			result := config.Result{Status: "OK", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		} else {
			Message := "RDS delete protection is not enabled on " + *instance.DBClusterIdentifier
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
