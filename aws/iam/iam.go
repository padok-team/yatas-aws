package iam

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/config"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *config.Config, queue chan []config.Check) {

	var checkConfig config.CheckConfig
	checkConfig.Init(s, c)
	var checks []config.Check
	users := GetAllUsers(s)
	mfaForUsers := GetMfaForUsers(s, users)
	accessKeysForUsers := GetAccessKeysForUsers(s, users)
	UserToPolicies := GetUserPolicies(users, s)
	UserToPoliciesElevated := GetUserToPoliciesElevate(UserToPolicies)

	go config.CheckTest(checkConfig.Wg, c, "AWS_IAM_001", CheckIf2FAActivated)(checkConfig, mfaForUsers, "AWS_IAM_001")
	go config.CheckTest(checkConfig.Wg, c, "AWS_IAM_002", CheckAgeAccessKeyLessThan90Days)(checkConfig, accessKeysForUsers, "AWS_IAM_002")
	go config.CheckTest(checkConfig.Wg, c, "AWS_IAM_003", CheckIfUserCanElevateRights)(checkConfig, UserToPoliciesElevated, "AWS_IAM_003")
	go config.CheckTest(checkConfig.Wg, c, "AWS_IAM_004", CheckIfUserLastPasswordUse120Days)(checkConfig, users, "AWS_IAM_004")
	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)

			checkConfig.Wg.Done()

		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
