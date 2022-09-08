package apigateway

import (
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfTracingEnabled(checkConfig config.CheckConfig, stages map[string][]types.Stage, testName string) {
	var check config.Check
	check.InitCheck("ApiGateways have tracing enabled", "Check if all stages are enabled for tracing", testName)
	for apigateway, id := range stages {
		for _, stage := range id {
			if stage.TracingEnabled {
				Message := "Tracing is enabled on stage" + *stage.StageName + " of ApiGateway " + apigateway
				result := config.Result{Status: "OK", Message: Message, ResourceID: *stage.StageName}
				check.AddResult(result)
			} else {
				Message := "Tracing is not enabled on " + *stage.StageName + " of ApiGateway " + apigateway
				result := config.Result{Status: "FAIL", Message: Message, ResourceID: *stage.StageName}
				check.AddResult(result)
			}
		}
	}
	checkConfig.Queue <- check
}
