package apigateway

import (
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfStagesProtectedByAcl(checkConfig config.CheckConfig, stages map[string][]types.Stage, testName string) {
	var check config.Check
	check.InitCheck("ApiGateways are protected by an ACL", "Check if all stages are protected by ACL", testName)
	for apigateway, id := range stages {
		for _, stage := range id {
			if stage.WebAclArn != nil && *stage.WebAclArn != "" {
				Message := "Stage " + *stage.StageName + " is protected by ACL" + " of ApiGateway " + apigateway
				result := config.Result{Status: "OK", Message: Message, ResourceID: *stage.StageName}
				check.AddResult(result)
			} else {
				Message := "Stage " + *stage.StageName + " is not protected by ACL" + " of ApiGateway " + apigateway
				result := config.Result{Status: "FAIL", Message: Message, ResourceID: *stage.StageName}
				check.AddResult(result)
			}
		}
	}
	checkConfig.Queue <- check
}
