package iam

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckNoConsolePasswordForNonHumanUser(checkConfig commons.CheckConfig, consolePasswordForUsers []ConsolePasswordForUser, testName string) {
	var check commons.Check
	check.InitCheck("IAM Non-human users don’t have console password", "Check if non-human users don’t have console password", testName, []string{"Security", "Good Practice", "HDS"})
	for _, user := range consolePasswordForUsers {
		userName := user.UserName

		Message := userName + " has no console password"
		result := commons.Result{Status: "OK", Message: Message, ResourceID: userName}

		if user.HasConsolePassword {
			Message = userName + " has a console password"
			result = commons.Result{Status: "FAIL", Message: Message, ResourceID: userName}
		}
		check.AddResult(result)
	}
	checkConfig.Queue <- check
}
