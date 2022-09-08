package rds

import (
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/config"
)

func checkIfClusterRDSPrivateEnabled(checkConfig config.CheckConfig, instances []types.DBCluster, testName string) {
	var check config.Check
	check.InitCheck("Aurora RDS aren't publicly accessible", "Check if Aurora RDS private is enabled", testName)
	for _, instance := range instances {
		if instance.PubliclyAccessible != nil && *instance.PubliclyAccessible {
			Message := "RDS private is not enabled on " + *instance.DBClusterIdentifier
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)

		} else {

			Message := "RDS private is enabled on " + *instance.DBClusterIdentifier
			result := config.Result{Status: "OK", Message: Message, ResourceID: *instance.DBClusterArn}
			check.AddResult(result)

		}
	}
	checkConfig.Queue <- check
}
