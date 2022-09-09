package s3

import (
	"github.com/stangirard/yatas/plugins/commons"
)

func CheckIfBucketInOneZone(checkConfig commons.CheckConfig, buckets BucketAndNotInRegion, testName string) {
	var check commons.Check
	check.InitCheck("S3 buckets are not global but in one zone", "Check if S3 buckets are in one zone", testName)
	for _, bucket := range buckets.Buckets {
		found := false
		for _, region := range buckets.NotInRegion {
			if *bucket.Name == *region.Name {
				Message := "S3 bucket " + *bucket.Name + " is global but should be in " + checkConfig.ConfigAWS.Region
				result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *bucket.Name}
				check.AddResult(result)
				found = true
				break
			}
		}
		if !found {
			Message := "S3 bucket " + *bucket.Name + " is in " + checkConfig.ConfigAWS.Region
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *bucket.Name}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
