package cloudtrail

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfCloudtrailIsEnabled(checkConfig commons.CheckConfig, eventSelectorsByTrails []EventSelectorsByLoggingTrail, testName string) {
	var check commons.Check
	check.InitCheck("Cloudtrail is enabled", "Check if Cloudtrail is enabled", testName, []string{"Security", "Good Practice"})
	result := commons.Result{Status: "FAIL", Message: "Cloudtrail is not enabled"}

	var errors []string

	hasManagement := false
	hasData := false
	hasInsights := false

	for _, eventSelectorsByTrail := range eventSelectorsByTrails {
		for _, eventSelector := range eventSelectorsByTrail.EventSelectors {
			if aws.ToBool(eventSelector.IncludeManagementEvents) {
				hasManagement = true
			}
			if len(eventSelector.DataResources) > 0 {
				hasData = true
			}
		}
		if eventSelectorsByTrail.HasInsightSelectors {
			hasInsights = true
		}

		if hasManagement && hasData && hasInsights {
			break
		}
	}

	if !hasManagement {
		errors = append(errors, "CloudTrail does not log management events")
	}
	if !hasData {
		errors = append(errors, "CloudTrail does not log data events")
	}
	if !hasInsights {
		errors = append(errors, "CloudTrail does not log insight events")
	}

	if len(errors) > 0 {
		result.Status = "FAIL"
		result.Message = fmt.Sprintf("CloudTrail configuration has %d issues: %v", len(errors), strings.Join(errors, ", "))
	} else {
		Message := "CloudTrail is enabled with management, data, and insight events"
		result = commons.Result{Status: "OK", Message: Message}
	}

	check.AddResult(result)
	checkConfig.Queue <- check
}
