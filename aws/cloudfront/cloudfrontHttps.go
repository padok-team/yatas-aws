package cloudfront

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfHTTPSOnly(checkConfig config.CheckConfig, d []types.DistributionSummary, testName string) {
	var check config.Check
	check.InitCheck("Cloudfronts only allow HTTPS or redirect to HTTPS", "Check if all cloudfront distributions are HTTPS only", testName)
	for _, cloudfront := range d {
		if cloudfront.DefaultCacheBehavior != nil && (cloudfront.DefaultCacheBehavior.ViewerProtocolPolicy == "https-only" || cloudfront.DefaultCacheBehavior.ViewerProtocolPolicy == "redirect-to-https") {
			Message := "Cloudfront distribution is HTTPS only on " + *cloudfront.Id
			result := config.Result{Status: "OK", Message: Message, ResourceID: *cloudfront.Id}
			check.AddResult(result)
		} else {
			Message := "Cloudfront distribution is not HTTPS only on " + *cloudfront.Id
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *cloudfront.Id}
			check.AddResult(result)
		}
	}

	checkConfig.Queue <- check
}
