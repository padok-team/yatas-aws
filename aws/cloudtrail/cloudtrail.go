package cloudtrail

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/padok-team/yatas-aws/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {
	logger.Logger.Debug("CLD - Checks started")
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check
	cloudtrails := GetCloudtrails(s)
	eventSelectors := GetEventSelectorsForIsLoggingTrail(s, cloudtrails)

	go commons.CheckTest(checkConfig.Wg, c, "AWS_CLD_001", CheckIfCloudtrailsEncrypted)(checkConfig, cloudtrails, "AWS_CLD_001")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_CLD_002", CheckIfCloudtrailsGlobalServiceEventsEnabled)(checkConfig, cloudtrails, "AWS_CLD_002")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_CLD_003", CheckIfCloudtrailsMultiRegion)(checkConfig, cloudtrails, "AWS_CLD_003")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_CLD_004", CheckIfCloudtrailIsEnabled)(checkConfig, eventSelectors, "AWS_CLD_004")

	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)

			checkConfig.Wg.Done()

		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
	logger.Logger.Debug("CLD - Checks done")
}
