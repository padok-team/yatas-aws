package s3

import (
	"github.com/stangirard/yatas/config"
)

func CheckIfObjectLockConfigurationEnabled(checkConfig config.CheckConfig, buckets []S3ToObjectLock, testName string) {
	var check config.Check
	check.InitCheck("S3 buckets have a retention policy", "Check if S3 buckets are using retention policy", testName)
	for _, bucket := range buckets {
		if !bucket.ObjectLock {
			Message := "S3 bucket " + bucket.BucketName + " is not using retention policy"
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + bucket.BucketName + " is using retention policy"
			result := config.Result{Status: "OK", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		}

	}
	checkConfig.Queue <- check
}
