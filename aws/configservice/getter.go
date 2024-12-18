package configservice

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/configservice/types"
	"github.com/padok-team/yatas-aws/logger"
)

func GetConfigurationRecorderStatus(s aws.Config) []types.ConfigurationRecorderStatus {
	svc := configservice.NewFromConfig(s)
	result, err := svc.DescribeConfigurationRecorderStatus(context.TODO(), &configservice.DescribeConfigurationRecorderStatusInput{})
	if err != nil {
		logger.Logger.Error(err.Error())
		return []types.ConfigurationRecorderStatus{}
	}

	return result.ConfigurationRecordersStatus
}
