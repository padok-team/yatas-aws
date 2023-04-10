package lambda

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/padok-team/yatas-aws/logger"
)

func GetLambdas(s aws.Config) []types.FunctionConfiguration {
	svc := lambda.NewFromConfig(s)
	var lambdas []types.FunctionConfiguration
	input := &lambda.ListFunctionsInput{
		MaxItems: aws.Int32(100),
	}
	result, err := svc.ListFunctions(context.TODO(), input)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list
		return []types.FunctionConfiguration{}
	}
	lambdas = append(lambdas, result.Functions...)
	for {
		if result.NextMarker != nil {
			input.Marker = result.NextMarker
			result, err = svc.ListFunctions(context.TODO(), input)
			if err != nil {
				logger.Logger.Error(err.Error())
				// Return an empty list
				return []types.FunctionConfiguration{}
			}
			lambdas = append(lambdas, result.Functions...)
		} else {
			break
		}
	}
	return lambdas
}

func GetLambdaUrlConfigs(s aws.Config, lambdas []types.FunctionConfiguration) []LambdaUrlConfig {
	svc := lambda.NewFromConfig(s)
	lambdaUrlConfigs := []LambdaUrlConfig{}
	for _, function := range lambdas {
		input := &lambda.ListFunctionUrlConfigsInput{
			FunctionName: function.FunctionName,
		}
		result, err := svc.ListFunctionUrlConfigs(context.TODO(), input)
		if err != nil {
			return nil
		}
		lambdaUrlConfigs = append(lambdaUrlConfigs, LambdaUrlConfig{
			LambdaName: *function.FunctionName,
			LambdaArn:  *function.FunctionArn,
			UrlConfigs: result.FunctionUrlConfigs,
		})
	}
	return lambdaUrlConfigs
}
