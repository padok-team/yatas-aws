package iam

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckPasswordPolicy(checkConfig commons.CheckConfig, passwordPolicy types.PasswordPolicy, testName string) {
	var check commons.Check
	check.InitCheck("IAM password policy must be defined and enforced", "Check if the password policy is strong", testName, []string{"Security", "Good Practice"})
	var errors []string

	if *passwordPolicy.MinimumPasswordLength < int32(12) {
		errors = append(errors, "Minimum password length must be at least 12 characters")
	}

	if !passwordPolicy.RequireSymbols {
		errors = append(errors, "Password must require at least one symbol")
	}

	if !passwordPolicy.RequireNumbers {
		errors = append(errors, "Password must require at least one number")
	}

	if !passwordPolicy.RequireUppercaseCharacters {
		errors = append(errors, "Password must require at least one uppercase character")
	}

	if !passwordPolicy.RequireLowercaseCharacters {
		errors = append(errors, "Password must require at least one lowercase character")
	}

	if !passwordPolicy.ExpirePasswords {
		errors = append(errors, "Passwords must be set to expire")
	}

	if *passwordPolicy.MaxPasswordAge < int32(90) {
		errors = append(errors, fmt.Sprintf("Maximum password age must be at least 90 days (current: %d days)", *passwordPolicy.MaxPasswordAge))
	}

	if *passwordPolicy.PasswordReusePrevention < int32(5) {
		errors = append(errors, fmt.Sprintf("Password reuse prevention must prevent reuse of at least 5 previous passwords (current: %d)", *passwordPolicy.PasswordReusePrevention))
	}

	if passwordPolicy.AllowUsersToChangePassword {
		errors = append(errors, "Users should not be allowed to change their passwords")
	}

	var result commons.Result
	if len(errors) > 0 {
		result.Status = "FAIL"
		result.Message = fmt.Sprintf("Password policy has %d issues: %v", len(errors), errors)
	} else {
		result.Status = "OK"
		result.Message = "Password policy is strong and meets all requirements"
	}

	check.AddResult(result)
	checkConfig.Queue <- check
}
