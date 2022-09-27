package apigateway

import (
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/stangirard/yatas/plugins/commons"
)

func CheckIfTracingEnabled(checkConfig commons.CheckConfig, stages map[string][]types.Stage, testName string) {
	var check commons.Check
	check.InitCheck("ApiGateways have tracing enabled", "Check if all stages are enabled for tracing", testName, []string{"Security", "Good Practice"})
	for apigateway, id := range stages {
		for _, stage := range id {
			if stage.TracingEnabled {
				Message := "Tracing is enabled on stage" + *stage.StageName + " of ApiGateway " + apigateway
				result := commons.Result{Status: "OK", Message: Message, ResourceID: *stage.StageName}
				check.AddResult(result)
			} else {
				Message := "Tracing is not enabled on " + *stage.StageName + " of ApiGateway " + apigateway
				result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *stage.StageName}
				check.AddResult(result)
			}
		}
	}
	checkConfig.Queue <- check
}
