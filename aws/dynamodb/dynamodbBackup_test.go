package dynamodb

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stangirard/yatas/config"
)

func TestCheckIfDynamodbContinuousBackupsEnabled(t *testing.T) {
	type args struct {
		checkConfig config.CheckConfig
		dynamodbs   []TableBackups
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfDynamodbEncrypted",
			args: args{
				checkConfig: config.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan config.Check, 1),
				},
				dynamodbs: []TableBackups{
					{
						TableName: "DynamoDB-XXX",
						Backups: types.ContinuousBackupsDescription{
							ContinuousBackupsStatus: types.ContinuousBackupsStatusEnabled,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfDynamodbContinuousBackupsEnabled(tt.args.checkConfig, tt.args.dynamodbs, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckifDynamodbEncrypted() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfDynamodbContinuousBackupsEnabledFail(t *testing.T) {
	type args struct {
		checkConfig config.CheckConfig
		dynamodbs   []TableBackups
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfDynamodbEncrypted",
			args: args{
				checkConfig: config.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan config.Check, 1),
				},
				dynamodbs: []TableBackups{
					{
						TableName: "DynamoDB-XXX",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfDynamodbContinuousBackupsEnabled(tt.args.checkConfig, tt.args.dynamodbs, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckifDynamodbEncrypted() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
