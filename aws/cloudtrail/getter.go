package cloudtrail

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/padok-team/yatas-aws/logger"
)

func GetCloudtrails(s aws.Config) []types.Trail {
	svc := cloudtrail.NewFromConfig(s)
	input := &cloudtrail.DescribeTrailsInput{
		IncludeShadowTrails: aws.Bool(true),
	}
	result, err := svc.DescribeTrails(context.TODO(), input)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list
		return []types.Trail{}
	}
	return result.TrailList
}

type EventSelectorsByLoggingTrail struct {
	TrailName           string
	EventSelectors      []types.EventSelector
	HasInsightSelectors bool
}

func GetEventSelectorsForIsLoggingTrail(s aws.Config, trails []types.Trail) []EventSelectorsByLoggingTrail {
	svc := cloudtrail.NewFromConfig(s)
	var eventSelectorsForIsLoggingTrail []EventSelectorsByLoggingTrail

	for _, trail := range trails {
		status, err := svc.GetTrailStatus(context.TODO(), &cloudtrail.GetTrailStatusInput{
			Name: trail.Name,
		})
		if err != nil || !aws.ToBool(status.IsLogging) {
			continue
		}

		eventSelectors, err := svc.GetEventSelectors(context.TODO(), &cloudtrail.GetEventSelectorsInput{
			TrailName: trail.Name,
		})
		if err != nil {
			continue
		}

		eventSelectorsForIsLoggingTrail = append(eventSelectorsForIsLoggingTrail, EventSelectorsByLoggingTrail{
			TrailName:           aws.ToString(trail.Name),
			EventSelectors:      eventSelectors.EventSelectors,
			HasInsightSelectors: aws.ToBool(trail.HasInsightSelectors),
		})
	}

	return eventSelectorsForIsLoggingTrail
}
