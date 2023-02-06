package vpc

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func checkIfOnlyOneVPC(checkConfig commons.CheckConfig, vpcs []types.Vpc, testName string) {
	var check commons.Check
	check.InitCheck("VPC can't be in the same account", "Check if VPC has only one VPC", testName, []string{"Security", "Good Practice"})
	for _, vpc := range vpcs {
		if len(vpcs) > 1 {
			Message := "VPC Id:" + *vpc.VpcId
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *vpc.VpcId}
			check.AddResult(result)
		} else {
			Message := "VPC Id:" + *vpc.VpcId
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *vpc.VpcId}
			check.AddResult(result)
		}
	}

	checkConfig.Queue <- check
}
