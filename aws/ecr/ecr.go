package ecr

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/config"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *config.Config, queue chan []config.Check) {

	var checkConfig config.CheckConfig
	checkConfig.Init(s, c)
	var checks []config.Check
	ecr := GetECRs(s)
	go config.CheckTest(checkConfig.Wg, c, "AWS_ECR_001", CheckIfImageScanningEnabled)(checkConfig, ecr, "AWS_ECR_001")
	go config.CheckTest(checkConfig.Wg, c, "AWS_ECR_002", CheckIfEncrypted)(checkConfig, ecr, "AWS_ECR_002")
	go config.CheckTest(checkConfig.Wg, c, "AWS_ECR_003", CheckIfTagImmutable)(checkConfig, ecr, "AWS_ECR_003")
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
