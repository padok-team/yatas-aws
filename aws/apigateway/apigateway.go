package apigateway

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/stangirard/yatas/config"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *config.Config, queue chan []config.Check) {
	var checkConfig config.CheckConfig
	checkConfig.Init(s, c)
	var checks []config.Check
	svc := apigateway.NewFromConfig(s)
	apis := GetApiGateways(svc)
	stages := GetAllStagesApiGateway(svc, apis)
	go config.CheckTest(checkConfig.Wg, c, "AWS_APG_001", CheckIfStagesCloudwatchLogsExist)(checkConfig, stages, "AWS_APG_001")
	go config.CheckTest(checkConfig.Wg, c, "AWS_APG_002", CheckIfStagesProtectedByAcl)(checkConfig, stages, "AWS_APG_002")
	go config.CheckTest(checkConfig.Wg, c, "AWS_APG_003", CheckIfTracingEnabled)(checkConfig, stages, "AWS_APG_003")

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
