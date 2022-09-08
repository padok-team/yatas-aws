package volumes

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/config"
)

func Test_checkIfEncryptionEnabled(t *testing.T) {
	type args struct {
		checkConfig config.CheckConfig
		volumes     []types.Volume
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_checkIfEncryptionEnabled",
			args: args{
				checkConfig: config.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan config.Check, 1),
				},
				volumes: []types.Volume{
					{
						VolumeId:  aws.String("test"),
						Encrypted: aws.Bool(true),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfEncryptionEnabled(tt.args.checkConfig, tt.args.volumes, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfVolumeEncrypted() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func Test_checkIfEncryptionEnabledFail(t *testing.T) {
	type args struct {
		checkConfig config.CheckConfig
		volumes     []types.Volume
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_checkIfEncryptionEnabled",
			args: args{
				checkConfig: config.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan config.Check, 1),
				},
				volumes: []types.Volume{
					{
						VolumeId: aws.String("test"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfEncryptionEnabled(tt.args.checkConfig, tt.args.volumes, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfVolumeEncrypted() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
