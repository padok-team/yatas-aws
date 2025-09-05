package vpc

import "github.com/padok-team/yatas/plugins/commons"

func checkIfVPCFLowLogsEnabled(checkConfig commons.CheckConfig, VpcFlowLogs []VpcToFlowLogs, testName string) {
	var check commons.Check
	check.InitCheck("VPC Flow Logs are activated", "Check if VPC Flow Logs are enabled", testName, []string{"Security", "Good Practice", "HDS"})
	for _, vpcFlowLog := range VpcFlowLogs {

		if len(vpcFlowLog.FlowLogs) == 0 {
			Message := "VPC Flow Logs are not enabled on " + vpcFlowLog.VpcID
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: vpcFlowLog.VpcID}
			check.AddResult(result)
		} else {
			Message := "VPC Flow Logs are enabled on " + vpcFlowLog.VpcID
			result := commons.Result{Status: "OK", Message: Message, ResourceID: vpcFlowLog.VpcID}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
