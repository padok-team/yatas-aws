package acm

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfCertificateExpiresIn90Days(checkConfig config.CheckConfig, certificates []types.CertificateDetail, testName string) {
	var check config.Check
	check.InitCheck("ACM certificate expires in more than 90 days", "Check if certificate expires in 90 days", testName)
	for _, certificate := range certificates {
		if certificate.Status == types.CertificateStatusIssued || certificate.Status == types.CertificateStatusInactive {
			if time.Until(*certificate.NotAfter).Hours() > 24*90 {
				Message := "Certificate " + *certificate.CertificateArn + " does not expire in 90 days"
				result := config.Result{Status: "OK", Message: Message, ResourceID: *certificate.CertificateArn}
				check.AddResult(result)
			} else {
				Message := "Certificate " + *certificate.CertificateArn + " expires in 90 days or less"
				result := config.Result{Status: "FAIL", Message: Message, ResourceID: *certificate.CertificateArn}
				check.AddResult(result)
			}
		}
	}
	checkConfig.Queue <- check
}
