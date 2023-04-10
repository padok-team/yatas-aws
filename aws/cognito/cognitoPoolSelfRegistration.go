package cognito

import (
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfCognitoSelfRegistration(checkConfig commons.CheckConfig, cognitoUserPools []cognitoidentityprovider.DescribeUserPoolOutput, testName string) {
	var check commons.Check
	check.InitCheck("Cognito allows self-registration", "Check if Cognito allows self-registration", testName, []string{"Security", "Good Practice"})
	for _, c := range cognitoUserPools {
		if !c.UserPool.AdminCreateUserConfig.AllowAdminCreateUserOnly {
			Message := "Cognito allows self-registration on " + *c.UserPool.Name
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *c.UserPool.Arn}
			check.AddResult(result)
		} else {
			Message := "Cognito does not allow self-registration on " + *c.UserPool.Name
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *c.UserPool.Arn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
