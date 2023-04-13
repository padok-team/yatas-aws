package ec2

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/padok-team/yatas-aws/aws/awschecks"
	"github.com/padok-team/yatas/plugins/commons"
)

type EC2Instance struct {
	Instance types.Instance
}

func (e *EC2Instance) GetID() string {
	return *e.Instance.InstanceId
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check

	svc := ec2.NewFromConfig(s)
	instances := GetEC2s(svc)
	ec2Checks := []awschecks.CheckDefinition{
		{
			Title:          "AWS_EC2_001",
			Description:    "Check if all instances have monitoring enabled",
			Tags:           []string{"Security", "Good Practice"},
			ConditionFn:    Ec2MonitoringEnabledCondition,
			SuccessMessage: "EC2 instance has monitoring enabled",
			FailureMessage: "EC2 instance has no monitoring enabled",
		},
		{
			Title:          "AWS_EC2_002",
			Description:    "Check if all instances have a public IP",
			Tags:           []string{"Security", "Good Practice"},
			ConditionFn:    Ec2PublicIPCondition,
			SuccessMessage: "EC2 instance has no public IP",
			FailureMessage: "EC2 instance has a public IP",
		},
		{
			Title:          "AWS_EC2_003",
			Description:    "Check if instances are running in a Virtual Private Cloud (VPC)",
			Tags:           []string{"Security", "Good Practice"},
			ConditionFn:    Ec2RunningInVPCCondition,
			SuccessMessage: "EC2 instance is running in a VPC",
			FailureMessage: "EC2 instance is not running in a VPC",
		},
	}

	var resources []awschecks.Resource
	for _, instance := range instances {
		resources = append(resources, &EC2Instance{Instance: instance})
	}
	awschecks.AddChecks(&checkConfig, ec2Checks)
	go awschecks.CheckResources(checkConfig, resources, ec2Checks)

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
