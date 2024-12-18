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

func GetTrailStatus(s aws.Config, trailList []types.Trail) []cloudtrail.GetTrailStatusOutput {
	svc := cloudtrail.NewFromConfig(s)
	var trailStatusOutput []cloudtrail.GetTrailStatusOutput

	for _, trail := range trailList {
		statusInput := &cloudtrail.GetTrailStatusInput{
			Name: trail.Name,
		}
		result, err := svc.GetTrailStatus(context.TODO(), statusInput)

		if err != nil {
			logger.Logger.Error(err.Error())
			continue
		}

		trailStatusOutput = append(trailStatusOutput, *result)
	}

	return trailStatusOutput
}
