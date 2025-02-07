package cloudtrail

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func createMockEventSelectors(management, data, insights bool) []EventSelectorsByLoggingTrail {
	var eventSelectors []EventSelectorsByLoggingTrail
	eventSelector := EventSelectorsByLoggingTrail{
		HasInsightSelectors: insights,
		EventSelectors: []types.EventSelector{
			{
				IncludeManagementEvents: aws.Bool(management),
				DataResources: func() []types.DataResource {
					if data {
						return []types.DataResource{{Type: aws.String("AWS::S3::Object")}}
					}
					return nil
				}(),
			},
		},
	}
	eventSelectors = append(eventSelectors, eventSelector)
	return eventSelectors
}

func TestCheckIfCloudtrailIsEnabled(t *testing.T) {
	tests := []struct {
		name               string
		eventSelectors     []EventSelectorsByLoggingTrail
		expectedStatus     string
		expectedMessageSub string
	}{
		{
			name:               "All event types are enabled",
			eventSelectors:     createMockEventSelectors(true, true, true),
			expectedStatus:     "OK",
			expectedMessageSub: "CloudTrail is enabled with management, data, and insight events",
		},
		{
			name:               "Missing management events",
			eventSelectors:     createMockEventSelectors(false, true, true),
			expectedStatus:     "FAIL",
			expectedMessageSub: "CloudTrail configuration has 1 issues: CloudTrail does not log management events",
		},
		{
			name:               "Missing data events",
			eventSelectors:     createMockEventSelectors(true, false, true),
			expectedStatus:     "FAIL",
			expectedMessageSub: "CloudTrail configuration has 1 issues: CloudTrail does not log data events",
		},
		{
			name:               "Missing insight events",
			eventSelectors:     createMockEventSelectors(true, true, false),
			expectedStatus:     "FAIL",
			expectedMessageSub: "CloudTrail configuration has 1 issues: CloudTrail does not log insight events",
		},
		{
			name:               "Missing management and insight events",
			eventSelectors:     createMockEventSelectors(false, true, false),
			expectedStatus:     "FAIL",
			expectedMessageSub: "CloudTrail configuration has 2 issues: CloudTrail does not log management events, CloudTrail does not log insight events",
		},
		{
			name:               "All event types are missing",
			eventSelectors:     createMockEventSelectors(false, false, false),
			expectedStatus:     "FAIL",
			expectedMessageSub: "CloudTrail configuration has 3 issues",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkConfig := commons.CheckConfig{Queue: make(chan commons.Check, 1)}
			testName := "TestCloudTrail"
			CheckIfCloudtrailIsEnabled(checkConfig, tt.eventSelectors, testName)

			check := <-checkConfig.Queue
			if len(check.Results) != 1 {
				t.Fatalf("Expected 1 result, got %d", len(check.Results))
			}
			result := check.Results[0]

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, result.Status)
			}

			if tt.expectedMessageSub != "" && !strings.Contains(result.Message, tt.expectedMessageSub) {
				t.Errorf("Expected message to contain %q, got %q", tt.expectedMessageSub, result.Message)
			}
		})
	}
}
