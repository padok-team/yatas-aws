package acm

import (
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfACMValid(checkConfig config.CheckConfig, certificates []types.CertificateDetail, testName string) {
	var check config.Check
	check.InitCheck("ACM certificates are valid", "Check if certificate is valid", testName)
	for _, certificate := range certificates {
		if certificate.Status == types.CertificateStatusIssued || certificate.Status == types.CertificateStatusInactive {
			Message := "Certificate " + *certificate.CertificateArn + " is valid"
			result := config.Result{Status: "OK", Message: Message, ResourceID: *certificate.CertificateArn}
			check.AddResult(result)
		} else {
			Message := "Certificate " + *certificate.CertificateArn + " is not valid"
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *certificate.CertificateArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
