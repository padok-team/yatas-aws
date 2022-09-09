package vpc

import (
	"github.com/stangirard/yatas/plugins/commons"
)

func CheckIfAtLeast2Subnets(checkConfig commons.CheckConfig, vpcToSubnets []VPCToSubnet, testName string) {
	var check commons.Check
	check.InitCheck("VPC have at least 2 subnets", "Check if VPC has at least 2 subnets", testName)
	for _, vpcToSubnet := range vpcToSubnets {

		if len(vpcToSubnet.Subnets) < 2 {
			Message := "VPC " + vpcToSubnet.VpcID + " has less than 2 subnets"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: vpcToSubnet.VpcID}
			check.AddResult(result)
		} else {
			Message := "VPC " + vpcToSubnet.VpcID + " has at least 2 subnets"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: vpcToSubnet.VpcID}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
