package guardduty

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfGuarddutyNoHighFindings(checkConfig commons.CheckConfig, testName string, findings []string) {
	var check commons.Check
	check.InitCheck("GuardDuty has no HIGH severity findings", "Check GuardDuty for HIGH severity findings", testName, []string{"Security", "Good Practice"})

	if len(findings) == 0 {
		Message := "GuardDuty has 0 HIGH severity findings"
		result := commons.Result{Status: "OK", Message: Message}
		check.AddResult(result)
	} else {
		Message := "GuardDuty has at least 1 HIGH severity finding, please perform a review of these findings"
		result := commons.Result{Status: "FAIL", Message: Message}
		check.AddResult(result)
	}

	checkConfig.Queue <- check
}
