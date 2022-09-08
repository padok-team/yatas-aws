package acm

import (
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfACMInUse(checkConfig config.CheckConfig, certificates []types.CertificateDetail, testName string) {
	var check config.Check
	check.InitCheck("ACM certificates are used", "Check if certificate is in use", testName)
	for _, certificate := range certificates {
		if len(certificate.InUseBy) > 0 {
			Message := "Certificate " + *certificate.CertificateArn + " is in use"
			result := config.Result{Status: "OK", Message: Message, ResourceID: *certificate.CertificateArn}
			check.AddResult(result)
		} else {
			Message := "Certificate " + *certificate.CertificateArn + " is not in use"
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *certificate.CertificateArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
