package elb

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/padok-team/yatas/plugins/commons"
)

// sslPolicyEnforcesTLS12 reports whether an ELB predefined security policy sets TLS 1.2 as the minimum version.
// ponytail: matches AWS naming (ELBSecurityPolicy-TLS13-1-2-…, -FS-1-2-…, -TLS-1-2-…, -TLS13-1-3-…).
// TLS13-1-0 / TLS13-1-1 / 2016-08 / FS-2018-06 allow older versions and don't match. Switch to DescribeSSLPolicies if AWS breaks the naming.
func sslPolicyEnforcesTLS12(policy *string) bool {
	p := aws.ToString(policy)
	return strings.Contains(p, "-1-2-") || strings.Contains(p, "-1-3-")
}

func isHttpsRedirect(action types.Action) bool {
	rc := action.RedirectConfig
	return action.Type == types.ActionTypeEnumRedirect && rc != nil &&
		aws.ToString(rc.Protocol) == "HTTPS" &&
		aws.ToString(rc.Host) == "#{host}" &&
		aws.ToString(rc.Path) == "/#{path}" &&
		aws.ToString(rc.Query) == "#{query}" &&
		rc.StatusCode == types.RedirectActionStatusCodeEnumHttp301
}

func CheckAlbEnsureHttps(checkConfig commons.CheckConfig, loadBalancers []LoadBalancerAttributes, testName string) {
	var check commons.Check
	check.InitCheck("ALB only allows HTTPS traffic with TLS 1.2 minimum", "Check if ALB HTTPS listeners enforce TLS 1.2 minimum and HTTP listeners redirect to HTTPS", testName, []string{"Security", "Good Practice", "HDS"})

	for _, loadBalancer := range loadBalancers {
		hasHttps := false
		var issues []string

		for _, listener := range loadBalancer.Listeners {
			port := aws.ToInt32(listener.Port)
			switch listener.Protocol {
			case types.ProtocolEnumHttps:
				hasHttps = true
				if !sslPolicyEnforcesTLS12(listener.SslPolicy) {
					issues = append(issues, fmt.Sprintf("HTTPS listener on port %d uses SSL policy %s which allows TLS < 1.2", port, aws.ToString(listener.SslPolicy)))
				}
			case types.ProtocolEnumHttp:
				redirects := false
				for _, action := range listener.DefaultActions {
					redirects = redirects || isHttpsRedirect(action)
				}
				if !redirects {
					issues = append(issues, fmt.Sprintf("HTTP listener on port %d does not redirect to HTTPS", port))
				}
			}
		}
		if !hasHttps {
			issues = append([]string{"no HTTPS listener"}, issues...)
		}

		if len(issues) == 0 {
			Message := "ALB " + loadBalancer.LoadBalancerName + " only allows HTTPS with TLS 1.2 minimum"
			check.AddResult(commons.Result{Status: "OK", Message: Message, ResourceID: loadBalancer.LoadBalancerArn})
		} else {
			Message := "ALB " + loadBalancer.LoadBalancerName + ": " + strings.Join(issues, "; ")
			check.AddResult(commons.Result{Status: "FAIL", Message: Message, ResourceID: loadBalancer.LoadBalancerArn})
		}
	}

	checkConfig.Queue <- check
}
