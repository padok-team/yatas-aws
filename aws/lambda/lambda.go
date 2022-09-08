package lambda

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/config"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *config.Config, queue chan []config.Check) {

	var checkConfig config.CheckConfig
	checkConfig.Init(s, c)
	var checks []config.Check
	lambdas := GetLambdas(s)

	go config.CheckTest(checkConfig.Wg, c, "AWS_LMD_001", CheckIfLambdaPrivate)(checkConfig, lambdas, "AWS_LMD_001")
	go config.CheckTest(checkConfig.Wg, c, "AWS_LMD_002", CheckIfLambdaInSecurityGroup)(checkConfig, lambdas, "AWS_LMD_002")
	go config.CheckTest(checkConfig.Wg, c, "AWS_LMD_003", CheckIfLambdaNoErrors)(checkConfig, lambdas, "AWS_LMD_003")
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
