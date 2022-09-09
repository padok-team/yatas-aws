package iam

import (
	"github.com/stangirard/yatas/plugins/commons"
)

func CheckIf2FAActivated(checkConfig commons.CheckConfig, mfaForUsers []MFAForUser, testName string) {
	var check commons.Check
	check.InitCheck("IAM Users have 2FA activated", "Check if all users have 2FA activated", testName)
	for _, mfaForUser := range mfaForUsers {
		if len(mfaForUser.MFAs) == 0 {
			Message := "2FA is not activated on " + mfaForUser.UserName
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: mfaForUser.UserName}
			check.AddResult(result)
		} else {
			Message := "2FA is activated on " + mfaForUser.UserName
			result := commons.Result{Status: "OK", Message: Message, ResourceID: mfaForUser.UserName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
