package s3

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func Test_checkIfDeletionPolicyExists(t *testing.T) {
	tests := []struct {
		name     string
		bucket   S3ToLifecycleRules
		expected string
	}{
		{
			name: "Non-versioned with valid current expiration",
			bucket: bucket("non-versioned-valid", false, rule(
				"rule1", 90, nil, true,
			)),
			expected: "OK",
		},
		{
			name:     "Non-versioned with no rules",
			bucket:   bucket("non-versioned-none", false),
			expected: "FAIL",
		},
		{
			name: "Versioned with valid current and noncurrent expiration",
			bucket: bucket("versioned-valid-both", true, rule(
				"rule2", 90, intPtr(30), true,
			)),
			expected: "OK",
		},
		{
			name: "Versioned with only current expiration",
			bucket: bucket("versioned-current-only", true, rule(
				"rule3", 90, nil, true,
			)),
			expected: "FAIL",
		},
		{
			name: "Versioned with only noncurrent expiration",
			bucket: bucket("versioned-noncurrent-only", true, rule(
				"rule4", 0, intPtr(30), true,
			)),
			expected: "FAIL",
		},
		{
			name:     "Versioned with no rules",
			bucket:   bucket("versioned-none", true),
			expected: "FAIL",
		},
		{
			name: "Versioned with current > 90 days",
			bucket: bucket("versioned-current-too-long", true, rule(
				"rule5", 120, intPtr(30), true,
			)),
			expected: "FAIL",
		},
		{
			name: "Versioned with noncurrent > 90 days",
			bucket: bucket("versioned-noncurrent-too-long", true, rule(
				"rule6", 60, intPtr(120), true,
			)),
			expected: "FAIL",
		},
		{
			name: "Non-versioned with rule disabled",
			bucket: bucket("non-versioned-disabled", false, rule(
				"rule7", 90, nil, false,
			)),
			expected: "FAIL",
		},
		{
			name: "Versioned with disabled current rule",
			bucket: bucket("versioned-disabled-current", true,
				rule("rule8", 90, nil, false),                 // current disabled
				rule("rule8-noncurrent", 0, intPtr(30), true), // noncurrent enabled
			),
			expected: "FAIL",
		},
		{
			name: "Versioned with both rules enabled",
			bucket: bucket("versioned-enabled-both", true,
				rule("rule9", 90, nil, true),
				rule("rule9-noncurrent", 0, intPtr(30), true),
			),
			expected: "OK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := commons.CheckConfig{
				Queue: make(chan commons.Check, 1),
				Wg:    &sync.WaitGroup{},
			}
			cfg.Wg.Add(1)

			go func() {
				for check := range cfg.Queue {
					for _, result := range check.Results {
						if result.Status != tt.expected {
							t.Errorf("Expected %s, got %s for bucket %s", tt.expected, result.Status, result.ResourceID)
						}
					}
					cfg.Wg.Done()
				}
			}()

			checkIfDeletionPolicyExists(cfg, []S3ToLifecycleRules{tt.bucket}, "test")
			cfg.Wg.Wait()
		})
	}
}

//
// 🔧 Helpers
//

func bucket(name string, versioning bool, rules ...types.LifecycleRule) S3ToLifecycleRules {
	return S3ToLifecycleRules{
		BucketName:     name,
		Versioning:     versioning,
		LifecycleRules: rules,
	}
}

func rule(id string, currentDays int32, noncurrentDays *int32, enabled bool) types.LifecycleRule {
	rule := types.LifecycleRule{
		ID:     strPtr(id),
		Status: types.ExpirationStatusDisabled,
	}

	if enabled {
		rule.Status = types.ExpirationStatusEnabled
	}

	if currentDays > 0 {
		rule.Expiration = &types.LifecycleExpiration{
			Days: int32Ptr(currentDays),
		}
	}

	if noncurrentDays != nil {
		rule.NoncurrentVersionExpiration = &types.NoncurrentVersionExpiration{
			NoncurrentDays: noncurrentDays,
		}
	}

	return rule
}

func int32Ptr(i int32) *int32 { return &i }
func intPtr(i int32) *int32   { return &i }
func strPtr(s string) *string { return &s }
