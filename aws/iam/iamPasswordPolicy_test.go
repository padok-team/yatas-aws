package iam

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckPasswordPolicy(t *testing.T) {
	tests := []struct {
		name           string
		passwordPolicy types.PasswordPolicy
		expectedStatus string
		expectedErrors []string
	}{
		{
			name: "Strong password policy",
			passwordPolicy: types.PasswordPolicy{
				MinimumPasswordLength:      aws.Int32(12),
				RequireSymbols:             true,
				RequireNumbers:             true,
				RequireUppercaseCharacters: true,
				RequireLowercaseCharacters: true,
				ExpirePasswords:            true,
				MaxPasswordAge:             aws.Int32(90),
				PasswordReusePrevention:    aws.Int32(5),
				AllowUsersToChangePassword: false,
			},
			expectedStatus: "OK",
			expectedErrors: nil,
		},
		{
			name: "Weak password policy",
			passwordPolicy: types.PasswordPolicy{
				MinimumPasswordLength:      aws.Int32(8),
				RequireSymbols:             false,
				RequireNumbers:             false,
				RequireUppercaseCharacters: false,
				RequireLowercaseCharacters: false,
				ExpirePasswords:            false,
				MaxPasswordAge:             aws.Int32(30),
				PasswordReusePrevention:    aws.Int32(2),
				AllowUsersToChangePassword: true,
			},
			expectedStatus: "FAIL",
			expectedErrors: []string{
				"Minimum password length must be at least 12 characters",
				"Password must require at least one symbol",
				"Password must require at least one number",
				"Password must require at least one uppercase character",
				"Password must require at least one lowercase character",
				"Passwords must be set to expire",
				"Maximum password age must be at least 90 days (current: 30 days)",
				"Password reuse prevention must prevent reuse of at least 5 previous passwords (current: 2)",
				"Users should not be allowed to change their passwords",
			},
		},
		{
			name: "Weak password policy on MinimumPasswordLength",
			passwordPolicy: types.PasswordPolicy{
				MinimumPasswordLength:      aws.Int32(8),
				RequireSymbols:             true,
				RequireNumbers:             true,
				RequireUppercaseCharacters: true,
				RequireLowercaseCharacters: true,
				ExpirePasswords:            true,
				MaxPasswordAge:             aws.Int32(90),
				PasswordReusePrevention:    aws.Int32(5),
				AllowUsersToChangePassword: false,
			},
			expectedStatus: "FAIL",
			expectedErrors: []string{
				"Minimum password length must be at least 12 characters",
			},
		},
		{
			name: "Weak password policy on ExpirePasswords",
			passwordPolicy: types.PasswordPolicy{
				MinimumPasswordLength:      aws.Int32(12),
				RequireSymbols:             true,
				RequireNumbers:             true,
				RequireUppercaseCharacters: true,
				RequireLowercaseCharacters: true,
				ExpirePasswords:            false,
				MaxPasswordAge:             aws.Int32(90),
				PasswordReusePrevention:    aws.Int32(5),
				AllowUsersToChangePassword: false,
			},
			expectedStatus: "FAIL",
			expectedErrors: []string{
				"Passwords must be set to expire",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := make(chan commons.Check, 1)
			checkConfig := commons.CheckConfig{
				Queue: queue,
			}

			CheckPasswordPolicy(checkConfig, tt.passwordPolicy, tt.name)

			check := <-queue

			if check.Results[0].Status != tt.expectedStatus {
				t.Errorf("expected status %s, got %s", tt.expectedStatus, check.Results[0].Status)
			}

			if tt.expectedStatus == "FAIL" {
				for _, err := range tt.expectedErrors {
					if !containsError(check.Results[0].Message, err) {
						t.Errorf("expected error message to contain: %s", err)
					}
				}
			}
		})
	}
}

func containsError(message string, error string) bool {
	return strings.Contains(message, error)
}
