package iam

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/padok-team/yatas-aws/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {
	logger.Logger.Debug("IAM - Checks started")
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check
	users := GetAllUsers(s)
	mfaForUsers := GetMfaForUsers(s, users)
	accessKeysForUsers := GetAccessKeysForUsers(s, users)
	UserToPolicies := GetUserPolicies(users, s)
	UserToPoliciesElevated := GetUserToPoliciesElevate(UserToPolicies)
	roles := GetAllRoles(s)
	RoleToPolicies := GetRolePolicies(roles, s)
	RoleToPoliciesElevated := GetRoleToPoliciesElevate(RoleToPolicies)
	passwordPolicy := GetPasswordPolicy(s)

	go commons.CheckTest(checkConfig.Wg, c, "AWS_IAM_001", CheckIf2FAActivated)(checkConfig, mfaForUsers, "AWS_IAM_001")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_IAM_002", CheckAgeAccessKeyLessThan90Days)(checkConfig, accessKeysForUsers, "AWS_IAM_002")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_IAM_003", CheckIfUserCanElevateRights)(checkConfig, UserToPoliciesElevated, "AWS_IAM_003")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_IAM_004", CheckIfRoleCanElevateRights)(checkConfig, RoleToPoliciesElevated, "AWS_IAM_004")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_IAM_005", CheckIfUserLastPasswordUse120Days)(checkConfig, users, "AWS_IAM_005")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_IAM_006", CheckPasswordPolicy)(checkConfig, passwordPolicy, "AWS_IAM_006")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_IAM_007", CheckNoConsolePasswordForNonHumanUser)(checkConfig, users, "AWS_IAM_007")
	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)

			checkConfig.Wg.Done()

		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
	logger.Logger.Debug("IAM - Checks done")
}
