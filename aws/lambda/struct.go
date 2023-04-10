package lambda

import "github.com/aws/aws-sdk-go-v2/service/lambda/types"

type LambdaUrlConfig struct {
	LambdaName string
	LambdaArn  string
	UrlConfigs []types.FunctionUrlConfig
}
