package ecr

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/padok-team/yatas-aws/logger"
)

func GetECRs(s aws.Config) []types.Repository {
	svc := ecr.NewFromConfig(s)
	var ecrRepositories []types.Repository
	input := &ecr.DescribeRepositoriesInput{
		MaxResults: aws.Int32(100),
	}
	result, err := svc.DescribeRepositories(context.TODO(), input)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list
		return []types.Repository{}
	}
	ecrRepositories = append(ecrRepositories, result.Repositories...)
	for {
		if result.NextToken != nil {
			input.NextToken = result.NextToken
			result, err = svc.DescribeRepositories(context.TODO(), input)
			if err != nil {
				logger.Logger.Error(err.Error())
				// Return an empty list
				return []types.Repository{}
			}
			ecrRepositories = append(ecrRepositories, result.Repositories...)
		} else {
			break
		}
	}

	return ecrRepositories
}
