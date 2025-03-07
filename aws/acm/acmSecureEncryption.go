package acm

import (
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfACMUsesSecureEncryption(checkConfig commons.CheckConfig, certificates []types.CertificateDetail, testName string) {
	var check commons.Check
	check.InitCheck("ACM certificates are using secure encryption", "Check if certificates use secure encryption", testName, []string{"Security", "Good Practice"})

	for _, certificate := range certificates {
		if certificate.KeyAlgorithm != types.KeyAlgorithmRsa1024 {
			Message := "Certificate " + *certificate.CertificateArn + " uses secure encryption"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *certificate.CertificateArn}
			check.AddResult(result)
		} else {
			Message := "Certificate " + *certificate.CertificateArn + " uses insecure encryption"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *certificate.CertificateArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
