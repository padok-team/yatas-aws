package elb

import (
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckAlbEnsureHttps(checkConfig commons.CheckConfig, loadBalancers []LoadBalancerAttributes, testName string) {
	var check commons.Check
	check.InitCheck("ALB only allows HTTPS traffic", "Check if ALB only allows HTTPS (443) or redirects HTTP (80) to HTTPS", testName, []string{"Security", "Good Practice", "HDS"})

	for _, loadBalancer := range loadBalancers {
		hasHttps := false
		hasHttpRedirect := false
		hasNoHttpListener := true

		for _, listener := range loadBalancer.Listeners {
			// Check for HTTPS listener on port 443
			if *listener.Port == 443 {
				hasHttps = true
			}

			// Check for HTTP listener on port 80 with redirect
			if *listener.Port == 80 && len(listener.DefaultActions) > 0 {
				hasNoHttpListener = false
				action := listener.DefaultActions[0]
				if action.Type == "redirect" && action.RedirectConfig != nil {
					if *action.RedirectConfig.Protocol == "HTTPS" &&
						*action.RedirectConfig.Port == "443" &&
						*action.RedirectConfig.Host == "#{host}" &&
						*action.RedirectConfig.Path == "/#{path}" &&
						*action.RedirectConfig.Query == "#{query}" &&
						action.RedirectConfig.StatusCode == types.RedirectActionStatusCodeEnumHttp301 {
						hasHttpRedirect = true
					}
				}
			}
		}

		if hasHttps && (hasHttpRedirect || hasNoHttpListener) {
			Message := "ALB " + loadBalancer.LoadBalancerName + " has HTTPS enabled and HTTP redirects to HTTPS"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: loadBalancer.LoadBalancerArn}
			check.AddResult(result)
		} else if !hasHttps {
			Message := "ALB " + loadBalancer.LoadBalancerName + " does not have HTTPS listener on port 443"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: loadBalancer.LoadBalancerArn}
			check.AddResult(result)
		} else if !hasHttpRedirect && !hasNoHttpListener {
			Message := "ALB " + loadBalancer.LoadBalancerName + " does not redirect HTTP to HTTPS"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: loadBalancer.LoadBalancerArn}
			check.AddResult(result)
		}
	}

	checkConfig.Queue <- check
}
