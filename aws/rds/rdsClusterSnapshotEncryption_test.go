package rds

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfClusterSnapshotEncryptionEnabled(t *testing.T) {
	cases := []struct {
		name           string
		snapshots      []types.DBClusterSnapshot
		expectedStatus []string
		expectedLength int
	}{
		{
			name: "Snapshot encrypted",
			snapshots: []types.DBClusterSnapshot{
				{
					DBClusterIdentifier:  aws.String("test-cluster"),
					DBClusterSnapshotArn: aws.String("arn:aws:rds:region:account:snapshot:test-snapshot"),
					StorageEncrypted:     aws.Bool(true),
				},
			},
			expectedStatus: []string{"OK"},
			expectedLength: 1,
		},
		{
			name: "Snapshot not encrypted",
			snapshots: []types.DBClusterSnapshot{
				{
					DBClusterIdentifier:  aws.String("test-cluster"),
					DBClusterSnapshotArn: aws.String("arn:aws:rds:region:account:snapshot:test-snapshot"),
					StorageEncrypted:     aws.Bool(false),
				},
			},
			expectedStatus: []string{"FAIL"},
			expectedLength: 1,
		},
		{
			name: "Multiple snapshots mixed encryption",
			snapshots: []types.DBClusterSnapshot{
				{
					DBClusterIdentifier:  aws.String("test-cluster-1"),
					DBClusterSnapshotArn: aws.String("arn:aws:rds:region:account:snapshot:test-snapshot-1"),
					StorageEncrypted:     aws.Bool(true),
				},
				{
					DBClusterIdentifier:  aws.String("test-cluster-2"),
					DBClusterSnapshotArn: aws.String("arn:aws:rds:region:account:snapshot:test-snapshot-2"),
					StorageEncrypted:     aws.Bool(false),
				},
			},
			expectedStatus: []string{"OK", "FAIL"},
			expectedLength: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			checkConfig := commons.CheckConfig{Queue: make(chan commons.Check, 1)}
			checkIfClusterSnapshotEncryptionEnabled(checkConfig, c.snapshots, "test")
			check := <-checkConfig.Queue

			if len(check.Results) != c.expectedLength {
				t.Errorf("Expected %d results, got %d", c.expectedLength, len(check.Results))
			}

			for i, result := range check.Results {
				if result.Status != c.expectedStatus[i] {
					t.Errorf("Expected status %s, got %s", c.expectedStatus[i], result.Status)
				}
			}
		})
	}
}
