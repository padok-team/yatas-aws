package acm

import (
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/stangirard/yatas/plugins/commons"
)

func CheckIfACMInUse(checkConfig commons.CheckConfig, certificates []types.CertificateDetail, testName string) {
	var check commons.Check
	check.InitCheck("ACM certificates are used", "Check if certificate is in use", testName)
	for _, certificate := range certificates {
		if len(certificate.InUseBy) > 0 {
			Message := "Certificate " + *certificate.CertificateArn + " is in use"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *certificate.CertificateArn}
			check.AddResult(result)
		} else {
			Message := "Certificate " + *certificate.CertificateArn + " is not in use"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *certificate.CertificateArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
