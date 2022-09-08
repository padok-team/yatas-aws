package vpc

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/config"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *config.Config, queue chan []config.Check) {

	var checkConfig config.CheckConfig
	checkConfig.Init(s, c)
	var checks []config.Check
	vpcs := GetListVPC(s)
	subnetsforvpcs := GetSubnetForVPCS(s, vpcs)
	internetGatewaysForVpc := GetInternetGatewaysForVpc(s, vpcs)
	vpcFlowLogs := GetFlowLogsForVpc(s, vpcs)

	go config.CheckTest(checkConfig.Wg, c, "AWS_VPC_001", checkCIDR20)(checkConfig, vpcs, "AWS_VPC_001")
	go config.CheckTest(checkConfig.Wg, c, "AWS_VPC_002", checkIfOnlyOneVPC)(checkConfig, vpcs, "AWS_VPC_002")
	go config.CheckTest(checkConfig.Wg, c, "AWS_VPC_003", checkIfOnlyOneGateway)(checkConfig, internetGatewaysForVpc, "AWS_VPC_003")
	go config.CheckTest(checkConfig.Wg, c, "AWS_VPC_004", checkIfVPCFLowLogsEnabled)(checkConfig, vpcFlowLogs, "AWS_VPC_004")
	go config.CheckTest(checkConfig.Wg, c, "AWS_VPC_005", CheckIfAtLeast2Subnets)(checkConfig, subnetsforvpcs, "AWS_VPC_005")
	go config.CheckTest(checkConfig.Wg, c, "AWS_VPC_006", CheckIfSubnetInDifferentZone)(checkConfig, subnetsforvpcs, "AWS_VPC_006")
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
