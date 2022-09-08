package vpc

import (
	"github.com/stangirard/yatas/config"
)

func CheckIfAtLeast2Subnets(checkConfig config.CheckConfig, vpcToSubnets []VPCToSubnet, testName string) {
	var check config.Check
	check.InitCheck("VPC have at least 2 subnets", "Check if VPC has at least 2 subnets", testName)
	for _, vpcToSubnet := range vpcToSubnets {

		if len(vpcToSubnet.Subnets) < 2 {
			Message := "VPC " + vpcToSubnet.VpcID + " has less than 2 subnets"
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: vpcToSubnet.VpcID}
			check.AddResult(result)
		} else {
			Message := "VPC " + vpcToSubnet.VpcID + " has at least 2 subnets"
			result := config.Result{Status: "OK", Message: Message, ResourceID: vpcToSubnet.VpcID}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
