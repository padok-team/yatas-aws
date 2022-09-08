package eks

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/stangirard/yatas/config"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *config.Config, queue chan []config.Check) {
	var checkConfig config.CheckConfig
	checkConfig.Init(s, c)
	var checks []config.Check
	svc := eks.NewFromConfig(s)
	clusters := GetClusters(svc)
	go config.CheckTest(checkConfig.Wg, c, "AWS_EKS_001", CheckIfLoggingIsEnabled)(checkConfig, clusters, "AWS_EKS_001")
	go config.CheckTest(checkConfig.Wg, c, "AWS_EKS_002", CheckIfEksEndpointPrivate)(checkConfig, clusters, "AWS_EKS_002")
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
