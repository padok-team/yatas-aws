package configservice

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/configservice/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfConfigServiceIsEnabled(checkConfig commons.CheckConfig, testName string, configurationRecorderStatus []types.ConfigurationRecorderStatus) {
	var check commons.Check
	check.InitCheck("AWS Config is enabled in the account", "Check if AWS Config is enabled", testName, []string{"Security", "Good Practice", "HDS"})

	Message := "AWS Config is not enabled"
	result := commons.Result{Status: "FAIL", Message: Message}

	for _, recorderStatus := range configurationRecorderStatus {
		if aws.ToBool(&recorderStatus.Recording) {
			Message = "AWS Config is enabled"
			result = commons.Result{Status: "OK", Message: Message}
			break
		}
	}

	check.AddResult(result)

	checkConfig.Queue <- check
}
