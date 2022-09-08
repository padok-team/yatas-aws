package cloudfront

import (
	"github.com/stangirard/yatas/config"
)

func CheckIfStandardLogginEnabled(checkConfig config.CheckConfig, d []SummaryToConfig, testName string) {
	var check config.Check
	check.InitCheck("Cloudfronts queries are logged", "Check if all cloudfront distributions have standard logging enabled", testName)
	for _, cc := range d {

		if cc.config.Logging != nil && cc.config.Logging.Enabled != nil && *cc.config.Logging.Enabled {
			Message := "Standard logging is enabled on " + *cc.summary.Id
			result := config.Result{Status: "OK", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		} else {
			Message := "Standard logging is not enabled on " + *cc.summary.Id
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *cc.summary.Id}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
