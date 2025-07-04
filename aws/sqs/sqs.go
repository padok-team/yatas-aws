package sqs

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/padok-team/yatas-aws/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {
	logger.Logger.Debug("SQS - Checks started")
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check

	sqsQueues := GetSQSQueues(s)
	sqsRetentionPeriod := GetSQSRetentionPeriod(s, sqsQueues)

	go commons.CheckTest(checkConfig.Wg, c, "AWS_SQS_001", CheckIfRetentionPeriodIsValid)(checkConfig, sqsRetentionPeriod, "AWS_SQS_001")

	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)

			checkConfig.Wg.Done()

		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
	logger.Logger.Debug("SQS - Checks done")
}
