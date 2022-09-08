package cloudfront

import (
	"github.com/stangirard/yatas/config"
)

func CheckIfACLUsed(checkConfig config.CheckConfig, d []SummaryToConfig, testName string) {
	var check config.Check
	check.InitCheck("Cloudfronts are protected by an ACL", "Check if all cloudfront distributions have an ACL used", testName)
	for _, cc := range d {

		if cc.config.WebACLId != nil && *cc.config.WebACLId != "" {
			Message := "ACL is used on " + *cc.summary.Id
			result := config.Result{Status: "OK", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		} else {
			Message := "ACL is not used on " + *cc.summary.Id
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
