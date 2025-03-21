package rds

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfRDSSnapshotEncryptionEnabled(t *testing.T) {
	cases := []struct {
		name      string
		snapshots []types.DBSnapshot
		expected  []commons.Result
	}{
		{
			name: "Encrypted snapshot",
			snapshots: []types.DBSnapshot{
				{
					DBInstanceIdentifier: aws.String("db-1"),
					DBSnapshotArn:        aws.String("arn:aws:rds:us-east-1:123456789012:snapshot:db-1"),
					Encrypted:            aws.Bool(true),
				},
			},
			expected: []commons.Result{
				{
					Status:     "OK",
					Message:    "RDS snapshot encryption is enabled on db-1",
					ResourceID: "arn:aws:rds:us-east-1:123456789012:snapshot:db-1",
				},
			},
		},
		{
			name: "Unencrypted snapshot",
			snapshots: []types.DBSnapshot{
				{
					DBInstanceIdentifier: aws.String("db-2"),
					DBSnapshotArn:        aws.String("arn:aws:rds:us-east-1:123456789012:snapshot:db-2"),
					Encrypted:            aws.Bool(false),
				},
			},
			expected: []commons.Result{
				{
					Status:     "FAIL",
					Message:    "RDS snapshot encryption is not enabled on db-2",
					ResourceID: "arn:aws:rds:us-east-1:123456789012:snapshot:db-2",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			checkConfig := commons.CheckConfig{
				Queue: make(chan commons.Check, 1),
			}

			go checkIfRDSSnapshotEncryptionEnabled(checkConfig, c.snapshots, "TEST")

			check := <-checkConfig.Queue
			if len(check.Results) != len(c.expected) {
				t.Errorf("got %d results, expected %d", len(check.Results), len(c.expected))
			}

			for i, result := range check.Results {
				if result.Status != c.expected[i].Status {
					t.Errorf("got status %s, expected %s", result.Status, c.expected[i].Status)
				}
				if result.Message != c.expected[i].Message {
					t.Errorf("got message %s, expected %s", result.Message, c.expected[i].Message)
				}
				if result.ResourceID != c.expected[i].ResourceID {
					t.Errorf("got resource ID %s, expected %s", result.ResourceID, c.expected[i].ResourceID)
				}
			}
		})
	}
}
