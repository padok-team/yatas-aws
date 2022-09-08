package acm

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/stangirard/yatas/config"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *config.Config, queue chan []config.Check) {
	var checkConfig config.CheckConfig
	checkConfig.Init(s, c)
	var checks []config.Check
	svc := acm.NewFromConfig(s)
	certificates := GetCertificates(svc)
	go config.CheckTest(checkConfig.Wg, c, "AWS_ACM_001", CheckIfACMValid)(checkConfig, certificates, "AWS_ACM_001")
	go config.CheckTest(checkConfig.Wg, c, "AWS_ACM_002", CheckIfCertificateExpiresIn90Days)(checkConfig, certificates, "AWS_ACM_002")
	go config.CheckTest(checkConfig.Wg, c, "AWS_ACM_003", CheckIfACMInUse)(checkConfig, certificates, "AWS_ACM_003")
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
