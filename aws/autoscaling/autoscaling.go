package autoscaling

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/stangirard/yatas/config"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *config.Config, queue chan []config.Check) {
	var checkConfig config.CheckConfig
	checkConfig.Init(s, c)
	var checks []config.Check
	svc := autoscaling.NewFromConfig(s)
	groups := GetAutoscalingGroups(svc)

	go config.CheckTest(checkConfig.Wg, c, "AWS_ASG_001", CheckIfDesiredCapacityMaxCapacityBelow80percent)(checkConfig, groups, "AWS_ASG_001")
	go config.CheckTest(checkConfig.Wg, c, "AWS_ASG_002", CheckIfInTwoAvailibilityZones)(checkConfig, groups, "AWS_ASG_002")

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
