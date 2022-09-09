package vpc

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {

	var checkConfig commons.CheckConfig
	checkConfig.Init(s, c)
	var checks []commons.Check
	vpcs := GetListVPC(s)
	subnetsforvpcs := GetSubnetForVPCS(s, vpcs)
	internetGatewaysForVpc := GetInternetGatewaysForVpc(s, vpcs)
	vpcFlowLogs := GetFlowLogsForVpc(s, vpcs)

	go commons.CheckTest(checkConfig.Wg, c, "AWS_VPC_001", checkCIDR20)(checkConfig, vpcs, "AWS_VPC_001")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_VPC_002", checkIfOnlyOneVPC)(checkConfig, vpcs, "AWS_VPC_002")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_VPC_003", checkIfOnlyOneGateway)(checkConfig, internetGatewaysForVpc, "AWS_VPC_003")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_VPC_004", checkIfVPCFLowLogsEnabled)(checkConfig, vpcFlowLogs, "AWS_VPC_004")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_VPC_005", CheckIfAtLeast2Subnets)(checkConfig, subnetsforvpcs, "AWS_VPC_005")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_VPC_006", CheckIfSubnetInDifferentZone)(checkConfig, subnetsforvpcs, "AWS_VPC_006")
	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)

			checkConfig.Wg.Done()

		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
