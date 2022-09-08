package cloudtrail

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/config"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *config.Config, queue chan []config.Check) {
	var checkConfig config.CheckConfig
	checkConfig.Init(s, c)
	var checks []config.Check
	cloudtrails := GetCloudtrails(s)

	go config.CheckTest(checkConfig.Wg, c, "AWS_CLD_001", CheckIfCloudtrailsEncrypted)(checkConfig, cloudtrails, "AWS_CLD_001")
	go config.CheckTest(checkConfig.Wg, c, "AWS_CLD_002", CheckIfCloudtrailsGlobalServiceEventsEnabled)(checkConfig, cloudtrails, "AWS_CLD_002")
	go config.CheckTest(checkConfig.Wg, c, "AWS_CLD_003", CheckIfCloudtrailsMultiRegion)(checkConfig, cloudtrails, "AWS_CLD_003")

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
