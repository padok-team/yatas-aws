package apigateway

import (
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfStagesCloudwatchLogsExist(checkConfig config.CheckConfig, stages map[string][]types.Stage, testName string) {
	var check config.Check
	check.InitCheck("ApiGateways logs are sent to Cloudwatch", "Check if all cloudwatch logs are enabled for all stages", testName)
	for apigateway, id := range stages {
		for _, stage := range id {
			if stage.AccessLogSettings != nil && stage.AccessLogSettings.DestinationArn != nil {
				Message := "Cloudwatch logs are enabled on stage " + *stage.StageName + " of ApiGateway " + apigateway
				result := config.Result{Status: "OK", Message: Message, ResourceID: *stage.StageName}
				check.AddResult(result)
			} else {
				Message := "Cloudwatch logs are not enabled on " + *stage.StageName + " of ApiGateway " + apigateway
				result := config.Result{Status: "FAIL", Message: Message, ResourceID: *stage.StageName}
				check.AddResult(result)
			}
		}
	}
	checkConfig.Queue <- check
}
