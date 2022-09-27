package loadbalancers

import (
	"github.com/stangirard/yatas/plugins/commons"
)

func CheckIfAccessLogsEnabled(checkConfig commons.CheckConfig, loadBalancers []LoadBalancerAttributes, testName string) {
	var check commons.Check
	check.InitCheck("ELB have access logs enabled", "Check if all load balancers have access logs enabled", testName, []string{"Security", "Good Practice"})
	for _, loadBalancer := range loadBalancers {
		for _, attributes := range loadBalancer.Output.Attributes {

			if *attributes.Key == "access_logs.s3.enabled" && *attributes.Value == "true" {
				Message := "Access logs are enabled on : " + loadBalancer.LoadBalancerName
				result := commons.Result{Status: "OK", Message: Message, ResourceID: loadBalancer.LoadBalancerArn}
				check.AddResult(result)
			} else if *attributes.Key == "access_logs.s3.enabled" && *attributes.Value == "false" {
				Message := "Access logs are not enabled on : " + loadBalancer.LoadBalancerName
				result := commons.Result{Status: "FAIL", Message: Message, ResourceID: loadBalancer.LoadBalancerArn}
				check.AddResult(result)
			} else {
				continue
			}
		}

	}

	checkConfig.Queue <- check
}
