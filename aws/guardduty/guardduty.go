package guardduty

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {

	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check
	guarddutyDetectors := GetDetectors(s)
	guarddutyFindings := GetHighFindings(s)
	go commons.CheckTest(checkConfig.Wg, c, "AWS_GDT_001", CheckIfGuarddutyEnabled)(checkConfig, "AWS_GDT_001", guarddutyDetectors)
	go commons.CheckTest(checkConfig.Wg, c, "AWS_GDT_002", CheckIfGuarddutyNoHighFindings)(checkConfig, "AWS_GDT_002", guarddutyFindings)
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
