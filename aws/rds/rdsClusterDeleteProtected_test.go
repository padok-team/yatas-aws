package rds

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfClusterDeleteProtectionEnabled(t *testing.T) {
	type args struct {
		checkConfig commons.CheckConfig
		instances   []types.DBCluster
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_checkIfDeleteProtectionEnabled",
			args: args{
				checkConfig: commons.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan commons.Check, 1),
				},
				instances: []types.DBCluster{
					{
						DBClusterIdentifier: aws.String("test"),
						DBClusterArn:        aws.String("arn:aws:rds:us-east-1:123456789012:db:test"),
						StorageEncrypted:    aws.Bool(true),
						DeletionProtection:  aws.Bool(true),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfClusterDeleteProtectionEnabled(tt.args.checkConfig, tt.args.instances, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfDeleteProtected() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfClusterDeleteProtectionEnabledFail(t *testing.T) {
	type args struct {
		checkConfig commons.CheckConfig
		instances   []types.DBCluster
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_checkIfDeleteProtectionEnabled",
			args: args{
				checkConfig: commons.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan commons.Check, 1),
				},
				instances: []types.DBCluster{
					{
						DBClusterIdentifier: aws.String("test"),
						DBClusterArn:        aws.String("arn:aws:rds:us-east-1:123456789012:db:test"),
						StorageEncrypted:    aws.Bool(true),
						DeletionProtection:  aws.Bool(false),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfClusterDeleteProtectionEnabled(tt.args.checkConfig, tt.args.instances, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfDeleteProtected() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
