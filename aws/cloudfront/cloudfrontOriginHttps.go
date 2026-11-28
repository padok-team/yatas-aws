package cloudfront

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/padok-team/yatas/plugins/commons"
)

// viewerHttpsOnly reports whether no cache behavior of the distribution accepts plain HTTP from viewers.
// ponytail: checks all behaviors, not only those routed to a given origin. Precise per-origin mapping if it causes false positives.
func viewerHttpsOnly(d types.DistributionSummary) bool {
	if d.DefaultCacheBehavior == nil || d.DefaultCacheBehavior.ViewerProtocolPolicy == types.ViewerProtocolPolicyAllowAll {
		return false
	}
	if d.CacheBehaviors != nil {
		for _, b := range d.CacheBehaviors.Items {
			if b.ViewerProtocolPolicy == types.ViewerProtocolPolicyAllowAll {
				return false
			}
		}
	}
	return true
}

// CheckIfOriginHTTPSTLS12 checks the CloudFront -> custom origin leg, which goes over the public internet.
// S3 (REST) and VPC origins are skipped: traffic stays on the AWS network.
func CheckIfOriginHTTPSTLS12(checkConfig commons.CheckConfig, d []types.DistributionSummary, testName string) {
	var check commons.Check
	check.InitCheck("Cloudfronts connect to custom origins using HTTPS with TLS 1.2 minimum", "Check if all cloudfront custom origins use HTTPS and do not allow SSLv3, TLSv1 or TLSv1.1", testName, []string{"Security", "Good Practice", "HDS"})
	for _, cloudfront := range d {
		var issues []string
		if cloudfront.Origins != nil {
			for _, origin := range cloudfront.Origins.Items {
				cfg := origin.CustomOriginConfig
				if cfg == nil {
					continue
				}
				id := aws.ToString(origin.Id)
				switch {
				case cfg.OriginProtocolPolicy == types.OriginProtocolPolicyHttpOnly:
					issues = append(issues, "origin "+id+" uses HTTP only")
					continue
				case cfg.OriginProtocolPolicy == types.OriginProtocolPolicyMatchViewer && !viewerHttpsOnly(cloudfront):
					issues = append(issues, "origin "+id+" matches viewer protocol while viewers can use HTTP")
				}
				if cfg.OriginSslProtocols != nil {
					for _, p := range cfg.OriginSslProtocols.Items {
						if p != types.SslProtocolTLSv12 {
							issues = append(issues, "origin "+id+" allows "+string(p))
						}
					}
				}
			}
		}

		if len(issues) == 0 {
			Message := "Custom origins use HTTPS with TLS 1.2 minimum on " + *cloudfront.Id
			check.AddResult(commons.Result{Status: "OK", Message: Message, ResourceID: *cloudfront.Id})
		} else {
			Message := "Custom origins are not HTTPS with TLS 1.2 minimum on " + *cloudfront.Id + ": " + strings.Join(issues, "; ")
			check.AddResult(commons.Result{Status: "FAIL", Message: Message, ResourceID: *cloudfront.Id})
		}
	}
	checkConfig.Queue <- check
}
