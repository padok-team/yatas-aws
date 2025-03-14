package dynamodb

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/padok-team/yatas-aws/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {
	logger.Logger.Debug("DYN - Checks started")
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check
	dynamodbs := GetDynamodbs(s)
	gt := GetTables(s, dynamodbs)
	gb := GetContinuousBackups(s, dynamodbs)
	go commons.CheckTest(checkConfig.Wg, c, "AWS_DYN_001", CheckIfDynamodbEncrypted)(checkConfig, gt, "AWS_DYN_001")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_DYN_002", CheckIfDynamodbContinuousBackupsEnabled)(checkConfig, gb, "AWS_DYN_002")

	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)

			checkConfig.Wg.Done()

		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
	logger.Logger.Debug("DYN - Checks done")
}
