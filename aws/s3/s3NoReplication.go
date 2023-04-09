package s3

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfBucketNoReplicationOtherRegion(checkConfig commons.CheckConfig, buckets []S3ToReplicationOtherRegion, testName string) {
	var check commons.Check
	check.InitCheck("S3 buckets are not replicated to another region", "Check if S3 buckets are replicated to other region", testName, []string{"Security", "Good Practice"})
	for _, bucket := range buckets {
		if bucket.ReplicatedOtherRegion {
			Message := "S3 bucket " + bucket.BucketName + " is replicated to the " + bucket.OtherRegion + " region"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		} else {
			Message := "S3 bucket " + bucket.BucketName + " is not replicated to another region"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: bucket.BucketName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
