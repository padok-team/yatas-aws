package lambda

import (
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfLambdaNoErrors(checkConfig commons.CheckConfig, lambdas []types.FunctionConfiguration, testName string) {
	var check commons.Check
	check.InitCheck("Lambdas are not with errors", "Check if all Lambdas are running smoothly", testName, []string{"Security", "Good Practice"})
	for _, lambda := range lambdas {
		if lambda.StateReasonCode != types.StateReasonCodeIdle && lambda.StateReasonCode != "" {
			Message := "Lambda " + *lambda.FunctionName + " is in error with code : " + string(lambda.StateReasonCode)
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		} else {
			Message := "Lambda " + *lambda.FunctionName + " is running smoothly"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
