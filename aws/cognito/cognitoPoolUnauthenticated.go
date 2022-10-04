package cognito

import (
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/stangirard/yatas/plugins/commons"
)

func CheckIfCognitoAllowsUnauthenticated(checkConfig commons.CheckConfig, cognitoPools []cognitoidentity.DescribeIdentityPoolOutput, testName string) {
	var check commons.Check
	check.InitCheck("Cognito allows unauthenticated users", "Check if Cognito allows unauthenticated users", testName, []string{"Security", "Good Practice"})
	for _, c := range cognitoPools {
		if c.AllowUnauthenticatedIdentities {
			Message := "Cognito allows unauthenticated users on " + *c.IdentityPoolName
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *c.IdentityPoolName}
			check.AddResult(result)
		} else {
			Message := "Cognito does not allow unauthenticated users on " + *c.IdentityPoolName
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *c.IdentityPoolName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
