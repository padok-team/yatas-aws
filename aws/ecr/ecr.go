package ecr

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/padok-team/yatas-aws/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {
	logger.Logger.Debug("ECR - Checks started")
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check
	ecr := GetECRs(s)
	go commons.CheckTest(checkConfig.Wg, c, "AWS_ECR_001", CheckIfImageScanningEnabled)(checkConfig, ecr, "AWS_ECR_001")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_ECR_002", CheckIfEncrypted)(checkConfig, ecr, "AWS_ECR_002")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_ECR_003", CheckIfTagImmutable)(checkConfig, ecr, "AWS_ECR_003")
	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)

			checkConfig.Wg.Done()

		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
	logger.Logger.Debug("ECR - Checks done")
}
