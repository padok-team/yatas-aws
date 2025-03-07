package s3

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func checkIfDeletionPolicyExists(checkConfig commons.CheckConfig, bucketsToLifecycleRules []S3ToLifecycleRules, testName string) {
	var check commons.Check
	check.InitCheck("S3 buckets have a deletion policy", "Check if S3 buckets are using a deletion policy", testName, []string{"Good Practice"})

	const MAX_DURATION = 90 // Maximum retention period in days

	for _, bucketToLifecycleRules := range bucketsToLifecycleRules {
		bucketName := bucketToLifecycleRules.BucketName
		hasValidRetentionPolicy := false
		Message := "S3 bucket " + bucketName + " is not using a 90 days or less retention policy"
		result := commons.Result{Status: "FAIL", Message: Message, ResourceID: bucketName}

		for _, rule := range bucketToLifecycleRules.LifecycleRules {
			if rule.Status == "Enabled" && rule.Expiration.Days != nil && *rule.Expiration.Days <= MAX_DURATION {
				hasValidRetentionPolicy = true
				continue
			}
		}

		if hasValidRetentionPolicy {
			Message = "S3 bucket " + bucketName + " is using a 90 days or less retention policy"
			result = commons.Result{Status: "OK", Message: Message, ResourceID: bucketName}
		}
		check.AddResult(result)
	}
	checkConfig.Queue <- check
}
