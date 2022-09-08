package lambda

import (
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfLambdaPrivate(checkConfig config.CheckConfig, lambdas []types.FunctionConfiguration, testName string) {
	var check config.Check
	check.InitCheck("Lambdas are private", "Check if all Lambdas are private", testName)
	for _, lambda := range lambdas {
		if lambda.VpcConfig == nil {
			Message := "Lambda " + *lambda.FunctionName + " is public"
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		} else {
			Message := "Lambda " + *lambda.FunctionName + " is private"
			result := config.Result{Status: "OK", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
