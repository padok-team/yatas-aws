package ssm

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	ec2Svc "github.com/aws/aws-sdk-go-v2/service/ec2"
	iamSvc "github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/padok-team/yatas-aws/aws/ec2"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {

	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check

	ec2Client := ec2Svc.NewFromConfig(s)
	iamClient := iamSvc.NewFromConfig(s)
	instances := ec2.GetEC2s(ec2Client)
	bastionToIAMPolicies := GetBastionToIAMPolicies(ec2Client, iamClient, instances)
	go commons.CheckTest(checkConfig.Wg, c, "AWS_SSM_001", CheckIfAuditLogsEnabledOnBastionInstance)(checkConfig, bastionToIAMPolicies, "AWS_SSM_001")

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
