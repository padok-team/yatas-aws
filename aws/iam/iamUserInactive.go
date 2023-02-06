package iam

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfUserLastPasswordUse120Days(checkConfig commons.CheckConfig, users []types.User, testName string) {
	var check commons.Check
	check.InitCheck("IAM Users have not used their password for 120 days", "Check if all users have not used their password for 120 days", testName, []string{"Security", "Good Practice"})
	for _, user := range users {
		if user.PasswordLastUsed != nil {
			if time.Since(*user.PasswordLastUsed).Hours() > 120*24 {
				Message := "Password has not been used for more than 120 days on " + *user.UserName
				result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *user.UserName}
				check.AddResult(result)
			} else {
				Message := "Password has been used in the last 120 days on " + *user.UserName
				result := commons.Result{Status: "OK", Message: Message, ResourceID: *user.UserName}
				check.AddResult(result)
			}
		} else {
			Message := "Password has never been used on " + *user.UserName
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *user.UserName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
