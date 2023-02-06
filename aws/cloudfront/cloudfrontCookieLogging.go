package cloudfront

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfCookieLogginEnabled(checkConfig commons.CheckConfig, d []SummaryToConfig, testName string) {
	var check commons.Check
	check.InitCheck("Cloudfronts are logging Cookies", "Check if all cloudfront distributions have cookies logging enabled", testName, []string{"Security", "Good Practice"})
	for _, cc := range d {
		if cc.config.Logging != nil && *cc.config.Logging.Enabled && cc.config.Logging.IncludeCookies != nil && *cc.config.Logging.IncludeCookies {
			Message := "Cookie logging is enabled on " + *cc.summary.Id
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		} else {
			Message := "Cookie logging is not enabled on " + *cc.summary.Id
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
