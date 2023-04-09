package cognito

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check
	cognitoPools := GetCognitoPools(s)
	cognitoPoolsDetailed := GetDetailedCognitoPool(s, cognitoPools)

	go commons.CheckTest(checkConfig.Wg, c, "AWS_COG_001", CheckIfCognitoAllowsUnauthenticated)(checkConfig, cognitoPoolsDetailed, "AWS_COG_001")

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
