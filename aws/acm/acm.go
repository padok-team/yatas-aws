package acm

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check
	svc := acm.NewFromConfig(s)
	certificates := GetCertificates(svc)
	go commons.CheckTest(checkConfig.Wg, c, "AWS_ACM_001", CheckIfACMValid)(checkConfig, certificates, "AWS_ACM_001")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_ACM_002", CheckIfCertificateExpiresIn90Days)(checkConfig, certificates, "AWS_ACM_002")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_ACM_003", CheckIfACMInUse)(checkConfig, certificates, "AWS_ACM_003")
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
