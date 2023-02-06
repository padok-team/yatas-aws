package autoscaling

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	checkConfig.Init(s, c)
	var checks []commons.Check
	svc := autoscaling.NewFromConfig(s)
	groups := GetAutoscalingGroups(svc)

	go commons.CheckTest(checkConfig.Wg, c, "AWS_ASG_001", CheckIfDesiredCapacityMaxCapacityBelow80percent)(checkConfig, groups, "AWS_ASG_001")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_ASG_002", CheckIfInTwoAvailibilityZones)(checkConfig, groups, "AWS_ASG_002")

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
