package s3

import (
	"sync"
	"testing"

	"github.com/padok-team/yatas/plugins/commons"
)

func Test_checkIfReplicationDisabled(t *testing.T) {
	type args struct {
		checkConfig commons.CheckConfig
		buckets     []S3ToReplicationOtherRegion
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "check if replication disabled",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				buckets: []S3ToReplicationOtherRegion{
					{
						BucketName:            "test",
						ReplicatedOtherRegion: false,
						OtherRegion:           "",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfBucketNoReplicationOtherRegion(tt.args.checkConfig, tt.args.buckets, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfBucketNoReplicationOtherRegion() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func Test_checkIfReplicationDisabledFail(t *testing.T) {
	type args struct {
		checkConfig commons.CheckConfig
		buckets     []S3ToReplicationOtherRegion
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "check if replication disabled",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				buckets: []S3ToReplicationOtherRegion{
					{
						BucketName:            "test",
						ReplicatedOtherRegion: true,
						OtherRegion:           "eu-west-1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfBucketNoReplicationOtherRegion(tt.args.checkConfig, tt.args.buckets, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfBucketNoReplicationOtherRegion() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
