package vpc

import (
	"github.com/stangirard/yatas/config"
)

func CheckIfSubnetInDifferentZone(checkConfig config.CheckConfig, vpcToSubnets []VPCToSubnet, testName string) {
	var check config.Check
	check.InitCheck("VPC's Subnets are in different zones", "Check if Subnet are in different zone", testName)
	for _, vpcToSubnet := range vpcToSubnets {
		subnetsAZ := make(map[string]int)
		for _, subnet := range vpcToSubnet.Subnets {
			subnetsAZ[*subnet.AvailabilityZone]++
		}
		if len(subnetsAZ) > 1 {
			Message := "Subnets are in different zone on " + vpcToSubnet.VpcID
			result := config.Result{Status: "OK", Message: Message, ResourceID: vpcToSubnet.VpcID}
			check.AddResult(result)
		} else {
			Message := "Subnets are in same zone on " + vpcToSubnet.VpcID
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: vpcToSubnet.VpcID}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
