package iam

import (
	"strings"
	"testing"

	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckNoConsolePasswordForNonHumanUser(t *testing.T) {
	tests := []struct {
		name             string
		consolePasswords []ConsolePasswordForUser
		expectedStatus   []string
		expectedMessages []string
	}{
		{
			name:             "All users have no console password",
			consolePasswords: []ConsolePasswordForUser{{UserName: "user1", HasConsolePassword: false}, {UserName: "user2", HasConsolePassword: false}},
			expectedStatus:   []string{"OK", "OK"},
			expectedMessages: []string{
				"user1 has no console password",
				"user2 has no console password",
			},
		},
		{
			name:             "One user has console password",
			consolePasswords: []ConsolePasswordForUser{{UserName: "user1", HasConsolePassword: false}, {UserName: "user2", HasConsolePassword: true}},
			expectedStatus:   []string{"OK", "FAIL"},
			expectedMessages: []string{
				"user1 has no console password",
				"user2 has a console password",
			},
		},
		{
			name:             "All users have console password",
			consolePasswords: []ConsolePasswordForUser{{UserName: "user1", HasConsolePassword: true}, {UserName: "user2", HasConsolePassword: true}},
			expectedStatus:   []string{"FAIL", "FAIL"},
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

			CheckNoConsolePasswordForNonHumanUser(checkConfig, tt.consolePasswords, tt.name)

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
