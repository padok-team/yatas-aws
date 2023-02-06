package apigateway

import (
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfStagesCloudwatchLogsExist(checkConfig commons.CheckConfig, stages map[string][]types.Stage, testName string) {
	var check commons.Check
	check.InitCheck("ApiGateways logs are sent to Cloudwatch", "Check if all cloudwatch logs are enabled for all stages", testName, []string{"Security", "Good Practice"})
	for apigateway, id := range stages {
		for _, stage := range id {
			if stage.AccessLogSettings != nil && stage.AccessLogSettings.DestinationArn != nil {
				Message := "Cloudwatch logs are enabled on stage " + *stage.StageName + " of ApiGateway " + apigateway
				result := commons.Result{Status: "OK", Message: Message, ResourceID: *stage.StageName}
				check.AddResult(result)
			} else {
				Message := "Cloudwatch logs are not enabled on " + *stage.StageName + " of ApiGateway " + apigateway
				result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *stage.StageName}
				check.AddResult(result)
			}
		}
	}
	checkConfig.Queue <- check
}
