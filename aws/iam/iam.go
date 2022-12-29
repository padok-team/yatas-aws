package iam

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {

	var checkConfig commons.CheckConfig
	checkConfig.Init(s, c)
	var checks []commons.Check
	users := GetAllUsers(s)
	mfaForUsers := GetMfaForUsers(s, users)
	accessKeysForUsers := GetAccessKeysForUsers(s, users)
	UserToPolicies := GetUserPolicies(users, s)
	UserToPoliciesElevated := GetUserToPoliciesElevate(UserToPolicies)
	roles := GetAllRoles(s)
	RoleToPolicies := GetRolePolicies(roles, s)
	RoleToPoliciesElevated := GetRoleToPoliciesElevate(RoleToPolicies)

	go commons.CheckTest(checkConfig.Wg, c, "AWS_IAM_001", CheckIf2FAActivated)(checkConfig, mfaForUsers, "AWS_IAM_001")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_IAM_002", CheckAgeAccessKeyLessThan90Days)(checkConfig, accessKeysForUsers, "AWS_IAM_002")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_IAM_003", CheckIfUserCanElevateRights)(checkConfig, UserToPoliciesElevated, "AWS_IAM_003")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_IAM_004", CheckIfRoleCanElevateRights)(checkConfig, RoleToPoliciesElevated, "AWS_IAM_004")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_IAM_005", CheckIfUserLastPasswordUse120Days)(checkConfig, users, "AWS_IAM_005")
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
