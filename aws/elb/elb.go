package elb

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {

	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check
	loadBalancers := GetElasticLoadBalancers(s)
	la := GetLoadBalancersAttributes(s, loadBalancers)
	go commons.CheckTest(checkConfig.Wg, c, "AWS_ELB_001", CheckIfAccessLogsEnabled)(checkConfig, la, "AWS_ELB_001")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_ELB_002", CheckAlbEnsureHttps)(checkConfig, la, "AWS_ELB_002")
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
