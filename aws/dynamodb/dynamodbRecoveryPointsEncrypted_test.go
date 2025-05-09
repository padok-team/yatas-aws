package dynamodb

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	backupTypes "github.com/aws/aws-sdk-go-v2/service/backup/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfEncryptionDynamodbRecoveryPointsEnabled(t *testing.T) {
	cases := []struct {
		name                string
		tableRecoveryPoints []TableRecoveryPoints
		expectedStatus      []string
	}{
		{
			name: "Recovery points encrypted",
			tableRecoveryPoints: []TableRecoveryPoints{
				{
					TableName: "test-table",
					RecoveryPoints: []backupTypes.RecoveryPointByResource{
						{
							EncryptionKeyArn: aws.String("arn:aws:kms:region:account:key/test-key"),
							RecoveryPointArn: aws.String("arn:aws:dynamodb:region:account:table/test-table/backup/test-backup"),
						},
					},
				},
			},
			expectedStatus: []string{"OK"},
		},
		{
			name: "Recovery points not encrypted",
			tableRecoveryPoints: []TableRecoveryPoints{
				{
					TableName: "test-table",
					RecoveryPoints: []backupTypes.RecoveryPointByResource{
						{
							EncryptionKeyArn: aws.String(""),
							RecoveryPointArn: aws.String("arn:aws:dynamodb:region:account:table/test-table/backup/test-backup"),
						},
					},
				},
			},
			expectedStatus: []string{"FAIL"},
		},
		{
			name: "Multiple recovery points mixed encryption",
			tableRecoveryPoints: []TableRecoveryPoints{
				{
					TableName: "test-table",
					RecoveryPoints: []backupTypes.RecoveryPointByResource{
						{
							EncryptionKeyArn: aws.String("arn:aws:kms:region:account:key/test-key"),
							RecoveryPointArn: aws.String("arn:aws:dynamodb:region:account:table/test-table/backup/test-backup1"),
						},
						{
							EncryptionKeyArn: aws.String(""),
							RecoveryPointArn: aws.String("arn:aws:dynamodb:region:account:table/test-table/backup/test-backup2"),
						},
					},
				},
			},
			expectedStatus: []string{"OK", "FAIL"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			checkConfig := commons.CheckConfig{Queue: make(chan commons.Check, 1)}
			CheckIfEncryptionDynamodbRecoveryPointsEnabled(checkConfig, c.tableRecoveryPoints, "test")
			check := <-checkConfig.Queue

			if len(check.Results) != len(c.expectedStatus) {
				t.Errorf("Expected %d results, got %d", len(c.expectedStatus), len(check.Results))
			}

			for i, result := range check.Results {
				if result.Status != c.expectedStatus[i] {
					t.Errorf("Expected status %s, got %s", c.expectedStatus[i], result.Status)
				}
			}
		})
	}
}
