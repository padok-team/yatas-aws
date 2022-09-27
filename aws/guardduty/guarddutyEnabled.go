package guardduty

import (
	"github.com/stangirard/yatas/plugins/commons"
)

func CheckIfGuarddutyEnabled(checkConfig commons.CheckConfig, testName string, detectors []string) {
	var check commons.Check
	check.InitCheck("GuardDuty is enabled in the account", "Check if GuardDuty is enabled", testName, []string{"Security", "Good Practice"})

	if len(detectors) == 0 {
		Message := "GuardDuty is not enabled"
		result := commons.Result{Status: "FAIL", Message: Message}
		check.AddResult(result)
	} else {
		Message := "GuardDuty is enabled"
		result := commons.Result{Status: "OK", Message: Message}
		check.AddResult(result)
	}
	checkConfig.Queue <- check
}
