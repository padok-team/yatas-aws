package dynamodb

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/config"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *config.Config, queue chan []config.Check) {

	var checkConfig config.CheckConfig
	checkConfig.Init(s, c)
	var checks []config.Check
	dynamodbs := GetDynamodbs(s)
	gt := GetTables(s, dynamodbs)
	gb := GetContinuousBackups(s, dynamodbs)
	go config.CheckTest(checkConfig.Wg, c, "AWS_DYN_001", CheckIfDynamodbEncrypted)(checkConfig, gt, "AWS_DYN_001")
	go config.CheckTest(checkConfig.Wg, c, "AWS_DYN_002", CheckIfDynamodbContinuousBackupsEnabled)(checkConfig, gb, "AWS_DYN_002")

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
