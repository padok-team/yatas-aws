package lambda

import (
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfLambdaInSecurityGroup(checkConfig config.CheckConfig, lambdas []types.FunctionConfiguration, testName string) {
	var check config.Check
	check.InitCheck("Lambdas are in a security group", "Check if all Lambdas are in a security group", testName)
	for _, lambda := range lambdas {
		if lambda.VpcConfig == nil || lambda.VpcConfig.SecurityGroupIds == nil {
			Message := "Lambda " + *lambda.FunctionName + " is not in a security group"
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		} else {
			Message := "Lambda " + *lambda.FunctionName + " is in a security group"
			result := config.Result{Status: "OK", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
