package s3

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfDeletionPolicyExists(t *testing.T) {
	type args struct {
		checkConfig             commons.CheckConfig
		bucketsToLifecycleRules []S3ToLifecycleRules
		testName                string
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus string
	}{
		{
			name: "S3 buckets has valid deletion policy",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				bucketsToLifecycleRules: []S3ToLifecycleRules{
					{
						BucketName: "test-bucket",
						LifecycleRules: []types.LifecycleRule{
							{Status: "Enabled", Expiration: &types.LifecycleExpiration{Days: aws.Int32(90)}},
						},
					},
				},
				testName: "TestValidRetention",
			},
			expectedStatus: "OK",
		},
		{
			name: "S3 bucket has disabled deletion policy",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				bucketsToLifecycleRules: []S3ToLifecycleRules{
					{
						BucketName: "test-bucket",
						LifecycleRules: []types.LifecycleRule{
							{Status: "Disabled", Expiration: &types.LifecycleExpiration{Days: aws.Int32(90)}},
						},
					},
				},
				testName: "TestNoRetention",
			},
			expectedStatus: "FAIL",
		},
		{
			name: "S3 buckets has invalid deletion policy",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				bucketsToLifecycleRules: []S3ToLifecycleRules{
					{
						BucketName: "test-bucket",
						LifecycleRules: []types.LifecycleRule{
							{Status: "Enabled", Expiration: &types.LifecycleExpiration{Days: aws.Int32(91)}},
						},
					},
				},
				testName: "TestNoRetention",
			},
			expectedStatus: "FAIL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				checkIfDeletionPolicyExists(tt.args.checkConfig, tt.args.bucketsToLifecycleRules, tt.args.testName)
				for check := range tt.args.checkConfig.Queue {
					if check.Status != tt.expectedStatus {
						t.Errorf("Unexpected check result: %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
