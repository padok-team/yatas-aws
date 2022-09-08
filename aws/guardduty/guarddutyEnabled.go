package guardduty

import (
	"github.com/stangirard/yatas/config"
)

func CheckIfGuarddutyEnabled(checkConfig config.CheckConfig, testName string, detectors []string) {
	var check config.Check
	check.InitCheck("GuardDuty is enabled in the account", "Check if GuardDuty is enabled", testName)

	if len(detectors) == 0 {
		Message := "GuardDuty is not enabled"
		result := config.Result{Status: "FAIL", Message: Message}
		check.AddResult(result)
	} else {
		Message := "GuardDuty is enabled"
		result := config.Result{Status: "OK", Message: Message}
		check.AddResult(result)
	}
	checkConfig.Queue <- check
}
