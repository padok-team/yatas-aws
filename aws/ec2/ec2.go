package ec2

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/stangirard/yatas/config"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *config.Config, queue chan []config.Check) {

	var checkConfig config.CheckConfig
	checkConfig.Init(s, c)
	var checks []config.Check

	svc := ec2.NewFromConfig(s)
	instances := GetEC2s(svc)
	go config.CheckTest(checkConfig.Wg, c, "AWS_EC2_001", CheckIfEC2PublicIP)(checkConfig, instances, "AWS_EC2_001")
	go config.CheckTest(checkConfig.Wg, c, "AWS_EC2_002", CheckIfMonitoringEnabled)(checkConfig, instances, "AWS_EC2_002")

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
