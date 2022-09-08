package cloudfront

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/stangirard/yatas/config"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *config.Config, queue chan []config.Check) {

	var checkConfig config.CheckConfig
	checkConfig.Init(s, c)
	var checks []config.Check
	svc := cloudfront.NewFromConfig(s)
	d := GetAllCloudfront(svc)
	s2c := GetAllDistributionConfig(svc, d)
	go config.CheckTest(checkConfig.Wg, c, "AWS_CFT_001", CheckIfCloudfrontTLS1_2Minimum)(checkConfig, d, "AWS_CFT_001")
	go config.CheckTest(checkConfig.Wg, c, "AWS_CFT_002", CheckIfHTTPSOnly)(checkConfig, d, "AWS_CFT_002")
	go config.CheckTest(checkConfig.Wg, c, "AWS_CFT_003", CheckIfStandardLogginEnabled)(checkConfig, s2c, "AWS_CFT_003")
	go config.CheckTest(checkConfig.Wg, c, "AWS_CFT_004", CheckIfCookieLogginEnabled)(checkConfig, s2c, "AWS_CFT_004")
	go config.CheckTest(checkConfig.Wg, c, "AWS_CFT_005", CheckIfACLUsed)(checkConfig, s2c, "AWS_CFT_005")

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
