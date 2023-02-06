package s3

import (
	"sync"
	"testing"

	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfS3PublicAccessBlockEnabled(t *testing.T) {
	type args struct {
		checkConfig           commons.CheckConfig
		s3toPublicBlockAccess []S3toPublicBlockAccess
		testName              string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "check if s3 public access block enabled",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				s3toPublicBlockAccess: []S3toPublicBlockAccess{
					{
						BucketName: "test",
						Config:     true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfS3PublicAccessBlockEnabled(tt.args.checkConfig, tt.args.s3toPublicBlockAccess, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfS3PublicAccessBlockEnabled() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfS3PublicAccessBlockEnabledFail(t *testing.T) {
	type args struct {
		checkConfig           commons.CheckConfig
		s3toPublicBlockAccess []S3toPublicBlockAccess
		testName              string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "check if s3 public access block enabled",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				s3toPublicBlockAccess: []S3toPublicBlockAccess{
					{
						BucketName: "test",
						Config:     false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfS3PublicAccessBlockEnabled(tt.args.checkConfig, tt.args.s3toPublicBlockAccess, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfS3PublicAccessBlockEnabled() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
