package iam

import (
	"time"

	"github.com/stangirard/yatas/config"
)

func CheckAgeAccessKeyLessThan90Days(checkConfig config.CheckConfig, accessKeysForUsers []AccessKeysForUser, testName string) {
	var check config.Check
	check.InitCheck("IAM access key younger than 90 days", "Check if all users have access key less than 90 days", testName)
	for _, accesskeyforuser := range accessKeysForUsers {
		now := time.Now()
		for _, accessKey := range accesskeyforuser.AccessKeys {
			if now.Sub(*accessKey.CreateDate).Hours() > 2160 {
				Message := "Access key " + *accessKey.AccessKeyId + " is older than 90 days on " + accesskeyforuser.UserName
				result := config.Result{Status: "FAIL", Message: Message, ResourceID: accesskeyforuser.UserName}
				check.AddResult(result)

			} else {
				Message := "Access key " + *accessKey.AccessKeyId + " is younger than 90 days on " + accesskeyforuser.UserName
				result := config.Result{Status: "OK", Message: Message, ResourceID: accesskeyforuser.UserName}
				check.AddResult(result)
			}
		}
	}
	checkConfig.Queue <- check
}
