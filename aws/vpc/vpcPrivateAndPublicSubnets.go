package vpc

import (
	"strings"

	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfPrivateAndPublicSubnets(checkConfig commons.CheckConfig, vpcToSubnetWithRouteTables map[string][]SubnetWithRouteTables, testName string) {
	var check commons.Check
	check.InitCheck("VPCs have private and public subnets", "Check if VPCs have private and public subnets", testName, []string{"Security", "Good Practice", "HDS"})

	for vpcId, subnetsWithRouteTables := range vpcToSubnetWithRouteTables {
		hasPublicSubnet := false
		hasPrivateSubnet := false

		for _, subnetWithRouteTable := range subnetsWithRouteTables {
			isPublic := false

			for _, routeTable := range subnetWithRouteTable.RouteTables {
				for _, route := range routeTable.Routes {
					if route.GatewayId != nil && strings.HasPrefix(*route.GatewayId, "igw-") {
						isPublic = true
						break
					}
				}
				if isPublic {
					break
				}
			}

			if isPublic {
				hasPublicSubnet = true
			} else {
				hasPrivateSubnet = true
			}
		}

		if !hasPublicSubnet {
			Message := "VPC " + vpcId + " has no Public subnet"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: vpcId}
			check.AddResult(result)
		} else if !hasPrivateSubnet {
			Message := "VPC " + vpcId + " has no Private subnet"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: vpcId}
			check.AddResult(result)
		} else {
			Message := "VPC " + vpcId + " has at least on public and one private subnets"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: vpcId}
			check.AddResult(result)
		}
	}

	checkConfig.Queue <- check
}
