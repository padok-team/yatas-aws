package cloudfront

import (
	"github.com/stangirard/yatas/config"
)

func CheckIfCookieLogginEnabled(checkConfig config.CheckConfig, d []SummaryToConfig, testName string) {
	var check config.Check
	check.InitCheck("Cloudfronts are logging Cookies", "Check if all cloudfront distributions have cookies logging enabled", testName)
	for _, cc := range d {
		if cc.config.Logging != nil && *cc.config.Logging.Enabled && cc.config.Logging.IncludeCookies != nil && *cc.config.Logging.IncludeCookies {
			Message := "Cookie logging is enabled on " + *cc.summary.Id
			result := config.Result{Status: "OK", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		} else {
			Message := "Cookie logging is not enabled on " + *cc.summary.Id
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
