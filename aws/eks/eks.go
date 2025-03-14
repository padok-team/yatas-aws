package eks

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/padok-team/yatas-aws/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {
	logger.Logger.Debug("EKS - Checks started")
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check
	svc := eks.NewFromConfig(s)
	clusters := GetClusters(svc)
	go commons.CheckTest(checkConfig.Wg, c, "AWS_EKS_001", CheckIfLoggingIsEnabled)(checkConfig, clusters, "AWS_EKS_001")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_EKS_002", CheckIfEksEndpointPrivate)(checkConfig, clusters, "AWS_EKS_002")
	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)

			checkConfig.Wg.Done()

		}
	}()
	checkConfig.Wg.Wait()
	queue <- checks
	logger.Logger.Debug("EKS - Checks done")
}
