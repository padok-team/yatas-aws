package cloudfront

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/stangirard/yatas/plugins/commons"
)

func CheckIfCloudfrontTLS1_2Minimum(checkConfig commons.CheckConfig, d []types.DistributionSummary, testName string) {
	var check commons.Check
	check.InitCheck("Cloudfronts enforce TLS 1.2 at least", "Check if all cloudfront distributions have TLS 1.2 minimum", testName, []string{"Security", "Good Practice"})
	for _, cloudfront := range d {
		if cloudfront.ViewerCertificate != nil && strings.Contains(string(cloudfront.ViewerCertificate.MinimumProtocolVersion), "TLSv1.2") {
			Message := "TLS 1.2 minimum is set on " + *cloudfront.Id
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *cloudfront.Id}
			check.AddResult(result)
		} else {
			Message := "TLS 1.2 minimum is not set on " + *cloudfront.Id
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *cloudfront.Id}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
