package lambda

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfLambdaUrlAuth(checkConfig commons.CheckConfig, lambdaUrlConfigs []LambdaUrlConfig, testName string) {
	var check commons.Check
	check.InitCheck("Lambdas has no public URL access", "Check if all Lambdas has no URL AuthType set to None", testName, []string{"Security", "Good Practice"})

	for _, lambda := range lambdaUrlConfigs {
		AuthTypeIsNone := false
		for _, urlConfig := range lambda.UrlConfigs {
			if urlConfig.AuthType == "NONE" {
				AuthTypeIsNone = true
				break
			}
		}

		if AuthTypeIsNone {
			Message := "Lambda " + lambda.LambdaName + " has URL AuthType set to None"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: lambda.LambdaArn}
			check.AddResult(result)
		} else {
			Message := "Lambda " + lambda.LambdaName + " has no URL AuthType set to None"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: lambda.LambdaArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
