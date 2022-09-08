package s3

import (
	"github.com/stangirard/yatas/config"
)

func checkIfEncryptionEnabled(checkConfig config.CheckConfig, buckets []S3ToEncryption, testName string) {
	var check config.Check
	check.InitCheck("S3 are encrypted", "Check if S3 encryption is enabled", testName)
	for _, bucket := range buckets {
		if !bucket.Encrypted {
			Message := "S3 bucket " + bucket.BucketName + " is not using encryption"
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + bucket.BucketName + " is using encryption"
			result := config.Result{Status: "OK", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		}

	}
	checkConfig.Queue <- check
}
