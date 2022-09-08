package s3

import (
	"github.com/stangirard/yatas/config"
)

func CheckIfS3PublicAccessBlockEnabled(checkConfig config.CheckConfig, s3toPublicBlockAccess []S3toPublicBlockAccess, testName string) {
	var check config.Check
	check.InitCheck("S3 bucket have public access block enabled", "Check if S3 buckets are using Public Access Block", testName)
	for _, bucket := range s3toPublicBlockAccess {
		if !bucket.Config {
			Message := "S3 bucket " + bucket.BucketName + " is not using Public Access Block"
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + bucket.BucketName + " is using Public Access Block"
			result := config.Result{Status: "OK", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
