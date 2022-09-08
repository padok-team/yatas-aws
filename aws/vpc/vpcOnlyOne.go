package vpc

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/config"
)

func checkIfOnlyOneVPC(checkConfig config.CheckConfig, vpcs []types.Vpc, testName string) {
	var check config.Check
	check.InitCheck("VPC can't be in the same account", "Check if VPC has only one VPC", testName)
	for _, vpc := range vpcs {
		if len(vpcs) > 1 {
			Message := "VPC Id:" + *vpc.VpcId
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *vpc.VpcId}
			check.AddResult(result)
		} else {
			Message := "VPC Id:" + *vpc.VpcId
			result := config.Result{Status: "OK", Message: Message, ResourceID: *vpc.VpcId}
			check.AddResult(result)
		}
	}

	checkConfig.Queue <- check
}
