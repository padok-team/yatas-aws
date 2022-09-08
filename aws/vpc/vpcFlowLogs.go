package vpc

import "github.com/stangirard/yatas/config"

func checkIfVPCFLowLogsEnabled(checkConfig config.CheckConfig, VpcFlowLogs []VpcToFlowLogs, testName string) {
	var check config.Check
	check.InitCheck("VPC Flow Logs are activated", "Check if VPC Flow Logs are enabled", testName)
	for _, vpcFlowLog := range VpcFlowLogs {

		if len(vpcFlowLog.FlowLogs) == 0 {
			Message := "VPC Flow Logs are not enabled on " + vpcFlowLog.VpcID
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: vpcFlowLog.VpcID}
			check.AddResult(result)
		} else {
			Message := "VPC Flow Logs are enabled on " + vpcFlowLog.VpcID
			result := config.Result{Status: "OK", Message: Message, ResourceID: vpcFlowLog.VpcID}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
