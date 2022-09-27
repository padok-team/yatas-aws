package s3

import (
	"github.com/stangirard/yatas/plugins/commons"
)

func CheckIfS3PublicAccessBlockEnabled(checkConfig commons.CheckConfig, s3toPublicBlockAccess []S3toPublicBlockAccess, testName string) {
	var check commons.Check
	check.InitCheck("S3 bucket have public access block enabled", "Check if S3 buckets are using Public Access Block", testName, []string{"Security", "Good Practice"})
	for _, bucket := range s3toPublicBlockAccess {
		if !bucket.Config {
			Message := "S3 bucket " + bucket.BucketName + " is not using Public Access Block"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + bucket.BucketName + " is using Public Access Block"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
