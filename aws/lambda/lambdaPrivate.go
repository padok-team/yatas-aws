package lambda

import (
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfLambdaPrivate(checkConfig commons.CheckConfig, lambdas []types.FunctionConfiguration, testName string) {
	var check commons.Check
	check.InitCheck("Lambdas are private", "Check if all Lambdas are private", testName, []string{"Security", "Good Practice"})
	for _, lambda := range lambdas {
		if lambda.VpcConfig == nil {
			Message := "Lambda " + *lambda.FunctionName + " is public"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		} else {
			Message := "Lambda " + *lambda.FunctionName + " is private"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
