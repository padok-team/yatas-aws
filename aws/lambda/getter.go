package lambda

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func GetLambdas(s aws.Config) []types.FunctionConfiguration {
	svc := lambda.NewFromConfig(s)
	var lambdas []types.FunctionConfiguration
	input := &lambda.ListFunctionsInput{
		MaxItems: aws.Int32(100),
	}
	result, err := svc.ListFunctions(context.TODO(), input)
	if err != nil {
		fmt.Println(err)
		// Return an empty list
		return []types.FunctionConfiguration{}
	}
	lambdas = append(lambdas, result.Functions...)
	for {
		if result.NextMarker != nil {
			input.Marker = result.NextMarker
			result, err = svc.ListFunctions(context.TODO(), input)
			if err != nil {
				fmt.Println(err)
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
