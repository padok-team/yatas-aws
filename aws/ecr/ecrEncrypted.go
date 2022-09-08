package ecr

import (
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfEncrypted(checkConfig config.CheckConfig, ecr []types.Repository, testName string) {
	var check config.Check
	check.InitCheck("ECRs are encrypted", "Check if all ECRs are encrypted", testName)
	for _, ecr := range ecr {
		if ecr.EncryptionConfiguration == nil {
			Message := "ECR " + *ecr.RepositoryName + " is not encrypted"
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *ecr.RepositoryName}
			check.AddResult(result)
		} else {
			Message := "ECR " + *ecr.RepositoryName + " is encrypted"
			result := config.Result{Status: "OK", Message: Message, ResourceID: *ecr.RepositoryName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
