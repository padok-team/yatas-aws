package elb

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckNlbTLS12Minimum(checkConfig commons.CheckConfig, loadBalancers []LoadBalancerAttributes, testName string) {
	var check commons.Check
	check.InitCheck("NLB TLS listeners enforce TLS 1.2 minimum", "Check if NLB TLS listeners use a security policy with TLS 1.2 minimum", testName, []string{"Security", "Good Practice", "HDS"})

	for _, loadBalancer := range loadBalancers {
		hasTLS := false
		var issues []string
		for _, listener := range loadBalancer.Listeners {
			if listener.Protocol != types.ProtocolEnumTls {
				continue
			}
			hasTLS = true
			if !sslPolicyEnforcesTLS12(listener.SslPolicy) {
				issues = append(issues, fmt.Sprintf("TLS listener on port %d uses SSL policy %s which allows TLS < 1.2", aws.ToInt32(listener.Port), aws.ToString(listener.SslPolicy)))
			}
		}
		// ponytail: NLBs without TLS listener are skipped, TCP passthrough (TLS terminated by the target) can't be judged from the AWS API.
		if !hasTLS {
			continue
		}
		if len(issues) == 0 {
			Message := "NLB " + loadBalancer.LoadBalancerName + " TLS listeners enforce TLS 1.2 minimum"
			check.AddResult(commons.Result{Status: "OK", Message: Message, ResourceID: loadBalancer.LoadBalancerArn})
		} else {
			Message := "NLB " + loadBalancer.LoadBalancerName + ": " + strings.Join(issues, "; ")
			check.AddResult(commons.Result{Status: "FAIL", Message: Message, ResourceID: loadBalancer.LoadBalancerArn})
		}
	}

	checkConfig.Queue <- check
}
