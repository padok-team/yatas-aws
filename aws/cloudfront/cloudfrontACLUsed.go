package cloudfront

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfACLUsed(checkConfig commons.CheckConfig, d []SummaryToConfig, testName string) {
	var check commons.Check
	check.InitCheck("Cloudfronts are protected by an ACL", "Check if all cloudfront distributions have an ACL used", testName, []string{"Security", "Good Practice"})
	for _, cc := range d {

		if cc.config.WebACLId != nil && *cc.config.WebACLId != "" {
			Message := "ACL is used on " + *cc.summary.Id
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		} else {
			Message := "ACL is not used on " + *cc.summary.Id
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
