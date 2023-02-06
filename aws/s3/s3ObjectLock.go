package s3

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfObjectLockConfigurationEnabled(checkConfig commons.CheckConfig, buckets []S3ToObjectLock, testName string) {
	var check commons.Check
	check.InitCheck("S3 buckets have a retention policy", "Check if S3 buckets are using retention policy", testName, []string{"Security", "Good Practice"})
	for _, bucket := range buckets {
		if !bucket.ObjectLock {
			Message := "S3 bucket " + bucket.BucketName + " is not using retention policy"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + bucket.BucketName + " is using retention policy"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		}

	}
	checkConfig.Queue <- check
}
