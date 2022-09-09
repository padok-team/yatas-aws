package ecr

import (
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/stangirard/yatas/plugins/commons"
)

func CheckIfImageScanningEnabled(checkConfig commons.CheckConfig, ecr []types.Repository, testName string) {
	var check commons.Check
	check.InitCheck("ECRs image are scanned on push", "Check if all ECRs have image scanning enabled", testName)
	for _, ecr := range ecr {
		if !ecr.ImageScanningConfiguration.ScanOnPush {
			Message := "ECR " + *ecr.RepositoryName + " has image scanning disabled"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *ecr.RepositoryName}
			check.AddResult(result)
		} else {
			Message := "ECR " + *ecr.RepositoryName + " has image scanning enabled"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *ecr.RepositoryName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
