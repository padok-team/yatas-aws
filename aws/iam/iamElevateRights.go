package iam

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfUserCanElevateRights(checkConfig commons.CheckConfig, userToPolociesElevated []UserToPoliciesElevate, testName string) {
	var check commons.Check
	check.InitCheck("IAM User can't elevate rights", "Check if  users can elevate rights", testName, []string{"Security", "Good Practice"})
	for _, userPol := range userToPolociesElevated {
		if len(userPol.Policies) > 0 {
			var Message string
			if len(userPol.Policies) > 3 {
				Message = "User " + userPol.UserName + " can elevate rights with " + fmt.Sprint(userPol.Policies[len(userPol.Policies)-3:]) + " only last 3 policies"
			} else {
				Message = "User " + userPol.UserName + " can elevate rights with " + fmt.Sprint(userPol.Policies)
			}
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: userPol.UserName}
			check.AddResult(result)

		} else {
			Message := "User " + userPol.UserName + " cannot elevate rights"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: userPol.UserName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func CheckIfRoleCanElevateRights(checkConfig commons.CheckConfig, roleToPoliciesElevated []RoleToPoliciesElevate, testName string) {
	var check commons.Check
	check.InitCheck("IAM Role can't elevate rights", "Check if roles can elevate rights", testName, []string{"Security", "Good Practice"})
	for _, rolePol := range roleToPoliciesElevated {
		if len(rolePol.Policies) > 0 {
			var Message string
			if len(rolePol.Policies) > 3 {
				Message = "Role " + rolePol.RoleName + " can elevate rights with " + fmt.Sprint(rolePol.Policies[len(rolePol.Policies)-3:]) + " only last 3 policies"
			} else {
				Message = "Role " + rolePol.RoleName + " can elevate rights with " + fmt.Sprint(rolePol.Policies)
			}
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: rolePol.RoleName}
			check.AddResult(result)

		} else {
			Message := "Role " + rolePol.RoleName + " cannot elevate rights"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: rolePol.RoleName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func CheckPolicyForAllowInRequiredPermission(policies []Policy, requiredPermission [][]string) [][]string {
	// Extract all allow statements from policy
	allowStatements := make([]Statement, 0)
	for _, policy := range policies {
		for _, statement := range policy.Statements {
			if statement.Effect == "Allow" {
				allowStatements = append(allowStatements, statement)
			}
		}
	}
	var permissionElevationPossible = [][]string{}
	// Check if any statement is in requiredPermissions
	for _, permissions := range requiredPermissions {
		// Create a map of permissions and false
		permissionMap := make(map[string]bool)
		for _, permission := range permissions {
			permissionMap[permission] = false
		}
		for _, permission := range permissions {
			for _, statement := range allowStatements {
				for _, actions := range statement.Action {
					actions = strings.ReplaceAll(actions, "*", ".*")
					// If regex actions matches permission actions, return true
					found, err := regexp.MatchString(actions, permission)
					if err != nil {
						panic(err)
					}
					if found {
						permissionMap[permission] = true
					}
				}
			}
		}
		// If all permissions are true, return true
		permissionsBool := true
		for _, permission := range permissionMap {
			if !permission {
				permissionsBool = false
			}
		}
		if permissionsBool {
			permissionElevationPossible = append(permissionElevationPossible, permissions)
		}
	}

	return permissionElevationPossible
}
