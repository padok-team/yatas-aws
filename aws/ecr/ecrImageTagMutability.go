package ecr

import (
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/stangirard/yatas/config"
)

func CheckIfTagImmutable(checkConfig config.CheckConfig, ecr []types.Repository, testName string) {
	var check config.Check
	check.InitCheck("ECRs tags are immutable", "Check if all ECRs are tag immutable", testName)
	for _, ecr := range ecr {
		if ecr.ImageTagMutability == types.ImageTagMutabilityMutable {
			Message := "ECR " + *ecr.RepositoryName + " is not tag immutable"
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *ecr.RepositoryName}
			check.AddResult(result)
		} else {
			Message := "ECR " + *ecr.RepositoryName + " is tag immutable"
			result := config.Result{Status: "OK", Message: Message, ResourceID: *ecr.RepositoryName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
