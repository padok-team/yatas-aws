package s3

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func checkIfEncryptionEnabled(checkConfig commons.CheckConfig, buckets []S3ToEncryption, testName string) {
	var check commons.Check
	check.InitCheck("S3 are encrypted", "Check if S3 encryption is enabled", testName, []string{"Security", "Good Practice", "HDS"})
	for _, bucket := range buckets {
		if !bucket.Encrypted {
			Message := "S3 bucket " + bucket.BucketName + " is not using encryption"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + bucket.BucketName + " is using encryption"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		}

	}
	checkConfig.Queue <- check
}
