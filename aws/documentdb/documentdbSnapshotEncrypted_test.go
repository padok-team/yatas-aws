package documentdb

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/docdb/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func Test_checkIfDocumentDBSnapshotEncrypted(t *testing.T) {
	tests := []struct {
		name      string
		snapshots []types.DBClusterSnapshot
		wantOK    bool
	}{
		{
			name: "snapshot encrypted",
			snapshots: []types.DBClusterSnapshot{
				{
					DBClusterSnapshotIdentifier: aws.String("snap-1"),
					DBClusterSnapshotArn:        aws.String("arn:aws:rds:us-east-1:123456789012:cluster-snapshot:snap-1"),
					StorageEncrypted:            aws.Bool(true),
				},
			},
			wantOK: true,
		},
		{
			name: "snapshot not encrypted",
			snapshots: []types.DBClusterSnapshot{
				{
					DBClusterSnapshotIdentifier: aws.String("snap-1"),
					DBClusterSnapshotArn:        aws.String("arn:aws:rds:us-east-1:123456789012:cluster-snapshot:snap-1"),
					StorageEncrypted:            aws.Bool(false),
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
			checkIfDocumentDBSnapshotEncrypted(checkConfig, tt.snapshots, "AWS_DOC_007")
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
