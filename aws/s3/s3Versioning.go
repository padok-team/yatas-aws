package s3

import (
	"github.com/stangirard/yatas/config"
)

func CheckIfBucketObjectVersioningEnabled(checkConfig config.CheckConfig, buckets []S3ToVersioning, testName string) {
	var check config.Check
	check.InitCheck("S3 buckets are versioned", "Check if S3 buckets are using object versioning", testName)
	for _, bucket := range buckets {
		if !bucket.Versioning {
			Message := "S3 bucket " + bucket.BucketName + " is not using object versioning"
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + bucket.BucketName + " is using object versioning"
			result := config.Result{Status: "OK", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		}

	}
	checkConfig.Queue <- check
}
