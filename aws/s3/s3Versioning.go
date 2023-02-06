package s3

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfBucketObjectVersioningEnabled(checkConfig commons.CheckConfig, buckets []S3ToVersioning, testName string) {
	var check commons.Check
	check.InitCheck("S3 buckets are versioned", "Check if S3 buckets are using object versioning", testName, []string{"Security", "Good Practice"})
	for _, bucket := range buckets {
		if !bucket.Versioning {
			Message := "S3 bucket " + bucket.BucketName + " is not using object versioning"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + bucket.BucketName + " is using object versioning"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		}

	}
	checkConfig.Queue <- check
}
