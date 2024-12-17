package iam

import (
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckNoConsolePasswordForNonHumanUser(t *testing.T) {
	tests := []struct {
		name             string
		users            []types.User
		expectedStatus   []string
		expectedMessages []string
	}{
		{
			name: "All users have no console password",
			users: []types.User{
				{UserName: aws.String("user1"), PasswordLastUsed: nil},
				{UserName: aws.String("user2"), PasswordLastUsed: nil},
			},
			expectedStatus: []string{"OK", "OK"},
			expectedMessages: []string{
				"user1 has no console password",
				"user2 has no console password",
			},
		},
		{
			name: "One user has console password",
			users: []types.User{
				{UserName: aws.String("user1"), PasswordLastUsed: nil},
				{UserName: aws.String("user2"), PasswordLastUsed: aws.Time(time.Now())},
			},
			expectedStatus: []string{"OK", "FAIL"},
			expectedMessages: []string{
				"user1 has no console password",
				"user2 has a console password",
			},
		},
		{
			name: "All users have console password",
			users: []types.User{
				{UserName: aws.String("user1"), PasswordLastUsed: aws.Time(time.Now())},
				{UserName: aws.String("user2"), PasswordLastUsed: aws.Time(time.Now())},
			},
			expectedStatus: []string{"FAIL", "FAIL"},
			expectedMessages: []string{
				"user1 has a console password",
				"user2 has a console password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := make(chan commons.Check, 1)
			checkConfig := commons.CheckConfig{
				Queue: queue,
			}

			CheckNoConsolePasswordForNonHumanUser(checkConfig, tt.users, tt.name)

			check := <-queue

			for i, result := range check.Results {
				if result.Status != tt.expectedStatus[i] {
					t.Errorf("expected status for user %s: %s, got: %s", result.ResourceID, tt.expectedStatus[i], result.Status)
				}
				if !strings.Contains(result.Message, tt.expectedMessages[i]) {
					t.Errorf("expected message for user %s: %s, got: %s", result.ResourceID, tt.expectedMessages[i], result.Message)
				}
			}
		})
	}
}
