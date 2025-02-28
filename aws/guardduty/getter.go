package guardduty

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	"github.com/aws/aws-sdk-go-v2/service/guardduty/types"
	"github.com/padok-team/yatas-aws/logger"
)

func GetDetectors(s aws.Config) []string {
	svc := guardduty.NewFromConfig(s)
	input := &guardduty.ListDetectorsInput{}
	result, err := svc.ListDetectors(context.TODO(), input)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list
		return []string{}
	}
	return result.DetectorIds
}

func GetHighFindings(s aws.Config) []string {
	svc := guardduty.NewFromConfig(s)

	detectors := GetDetectors(s)

	input := &guardduty.ListFindingsInput{
		FindingCriteria: &types.FindingCriteria{
			Criterion: map[string]types.Condition{
				"severity": {
					GreaterThanOrEqual: aws.Int64(7), // https://docs.aws.amazon.com/guardduty/latest/ug/guardduty_findings-severity.html#guardduty-finding-severity-level-high
				},
				"service.archived": {
					Equals: []string{"false"},
				},
			},
		},
	}

	var findings []string

	for _, detector := range detectors {
		input.DetectorId = &detector
		result, err := svc.ListFindings(context.TODO(), input)
		if err != nil {
			logger.Logger.Error(err.Error())
			// Return an empty list
			return []string{}
		}

		findings = append(findings, result.FindingIds...)
	}

	return findings
}
