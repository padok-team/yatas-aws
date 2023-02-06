package ecr

import (
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfTagImmutable(checkConfig commons.CheckConfig, ecr []types.Repository, testName string) {
	var check commons.Check
	check.InitCheck("ECRs tags are immutable", "Check if all ECRs are tag immutable", testName, []string{"Security", "Good Practice"})
	for _, ecr := range ecr {
		if ecr.ImageTagMutability == types.ImageTagMutabilityMutable {
			Message := "ECR " + *ecr.RepositoryName + " is not tag immutable"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *ecr.RepositoryName}
			check.AddResult(result)
		} else {
			Message := "ECR " + *ecr.RepositoryName + " is tag immutable"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *ecr.RepositoryName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
