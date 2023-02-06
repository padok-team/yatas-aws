package acm

import (
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfACMValid(checkConfig commons.CheckConfig, certificates []types.CertificateDetail, testName string) {
	var check commons.Check
	check.InitCheck("ACM certificates are valid", "Check if certificate is valid", testName, []string{"Security", "Good Practice"})
	for _, certificate := range certificates {
		if certificate.Status == types.CertificateStatusIssued || certificate.Status == types.CertificateStatusInactive {
			Message := "Certificate " + *certificate.CertificateArn + " is valid"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *certificate.CertificateArn}
			check.AddResult(result)
		} else {
			Message := "Certificate " + *certificate.CertificateArn + " is not valid"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *certificate.CertificateArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
