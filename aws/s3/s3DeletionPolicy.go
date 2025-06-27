package s3

import (
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/padok-team/yatas-aws/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

const MAX_DURATION = int32(90) // Maximum deletion period in days

func isValidRetention(days *int32) bool {
	return days != nil && *days > 0 && *days <= MAX_DURATION
}

func checkIfDeletionPolicyExists(checkConfig commons.CheckConfig, buckets []S3ToLifecycleRules, testName string) {
	var check commons.Check
	check.InitCheck("S3 buckets have a deletion policy", "Check if S3 buckets are using a deletion policy", testName, []string{"Good Practice"})

	for _, bucket := range buckets {
		bucketName := bucket.BucketName
		lifecycleRules := bucket.LifecycleRules
		isVersioningEnabled := bucket.Versioning
		var validRuleIDCurrent *string = nil
		var validRuleIDNonCurrent *string = nil

		Message := "S3 bucket " + bucketName + " is not using a 90 days or less deletion policy"
		result := commons.Result{Status: "FAIL", Message: Message, ResourceID: bucketName}

		logger.Logger.Debug("Checking bucket " + bucketName + " for deletion policy")

		if isVersioningEnabled {
			logger.Logger.Debug("Versioning is enabled for bucket " + bucketName)

			// Check if lifecycle rule action "Permanently delete noncurrent versions of objects" exists
			for _, rule := range lifecycleRules {
				var nonCurrentDays int32 = 0
				if rule.NoncurrentVersionExpiration != nil && rule.NoncurrentVersionExpiration.NoncurrentDays != nil {
					nonCurrentDays = *rule.NoncurrentVersionExpiration.NoncurrentDays
				}

				if rule.Status != types.ExpirationStatusEnabled {
					logger.Logger.Debug("Skipping rule with non-enabled status: " + string(rule.Status))
					continue
				}

				if isValidRetention(&nonCurrentDays) {
					validRuleIDNonCurrent = rule.ID
					continue
				}
			}
		}

		// Check if lifecycle rule action "Expire current versions of objects" exists
		for _, rule := range lifecycleRules {
			var expirationDays int32 = 0
			if rule.Expiration != nil && rule.Expiration.Days != nil {
				expirationDays = *rule.Expiration.Days
			}

			if rule.Status != types.ExpirationStatusEnabled {
				logger.Logger.Debug("Skipping rule with non-enabled status: " + string(rule.Status))
				continue
			}

			if isValidRetention(&expirationDays) {
				validRuleIDCurrent = rule.ID
				continue
			}
		}

		if validRuleIDCurrent != nil && ((isVersioningEnabled && validRuleIDNonCurrent != nil) || !isVersioningEnabled) {
			Message = "S3 bucket " + bucketName + " is using a 90 days or less deletion policy (rule ID: " + *validRuleIDCurrent + ")"
			result = commons.Result{Status: "OK", Message: Message, ResourceID: bucketName}
		}

		check.AddResult(result)
	}
	checkConfig.Queue <- check
}
