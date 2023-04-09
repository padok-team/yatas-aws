package guardduty

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
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
