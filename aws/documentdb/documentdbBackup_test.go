package documentdb

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/docdb/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func Test_checkIfDocumentDBBackupEnabled(t *testing.T) {
	tests := []struct {
		name     string
		clusters []types.DBCluster
		wantOK   bool
	}{
		{
			name: "backup retention >= 1 day",
			clusters: []types.DBCluster{
				{
					DBClusterIdentifier:   aws.String("my-cluster"),
					DBClusterArn:          aws.String("arn:aws:rds:us-east-1:123456789012:cluster:my-cluster"),
					BackupRetentionPeriod: aws.Int32(1),
				},
			},
			wantOK: true,
		},
		{
			name: "backup disabled (0 days)",
			clusters: []types.DBCluster{
				{
					DBClusterIdentifier:   aws.String("my-cluster"),
					DBClusterArn:          aws.String("arn:aws:rds:us-east-1:123456789012:cluster:my-cluster"),
					BackupRetentionPeriod: aws.Int32(0),
				},
			},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkConfig := commons.CheckConfig{
				Wg:    &sync.WaitGroup{},
				Queue: make(chan commons.Check, 1),
			}
			checkIfDocumentDBBackupEnabled(checkConfig, tt.clusters, "AWS_DOC_002")
			checkConfig.Wg.Add(1)
			go func() {
				for check := range checkConfig.Queue {
					wantStatus := "OK"
					if !tt.wantOK {
						wantStatus = "FAIL"
					}
					if check.Status != wantStatus {
						t.Errorf("got status %v, want %v", check.Status, wantStatus)
					}
					checkConfig.Wg.Done()
				}
			}()
			checkConfig.Wg.Wait()
		})
	}
}
