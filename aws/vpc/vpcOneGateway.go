package vpc

import (
	"github.com/stangirard/yatas/config"
)

func checkIfOnlyOneGateway(checkConfig config.CheckConfig, vpcInternetGateways []VpcToInternetGateway, testName string) {
	var check config.Check
	check.InitCheck("VPC only have one Gateway", "Check if VPC has only one gateway", testName)
	for _, vpcInternetGateway := range vpcInternetGateways {
		if len(vpcInternetGateway.InternetGateways) > 1 {
			Message := "VPC has more than one gateway on " + vpcInternetGateway.VpcID
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: vpcInternetGateway.VpcID}
			check.AddResult(result)
		} else {
			Message := "VPC has only one gateway on " + vpcInternetGateway.VpcID
			result := config.Result{Status: "OK", Message: Message, ResourceID: vpcInternetGateway.VpcID}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
