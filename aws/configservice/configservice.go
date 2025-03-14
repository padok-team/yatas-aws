package configservice

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/padok-team/yatas-aws/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {
	logger.Logger.Debug("CFG - Checks started")

	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check
	configurationRecorderStatus := GetConfigurationRecorderStatus(s)

	go commons.CheckTest(checkConfig.Wg, c, "AWS_CFG_001", CheckIfConfigServiceIsEnabled)(checkConfig, "AWS_CFG_001", configurationRecorderStatus)

	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)

			checkConfig.Wg.Done()

		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
	logger.Logger.Debug("CFG - Checks done")
}
