package documentdb

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/docdb/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func Test_checkIfDocumentDBAuditLogsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		clusters []types.DBCluster
		wantOK   bool
	}{
		{
			name: "audit logs enabled",
			clusters: []types.DBCluster{
				{
					DBClusterIdentifier:           aws.String("my-cluster"),
					DBClusterArn:                  aws.String("arn:aws:rds:us-east-1:123456789012:cluster:my-cluster"),
					EnabledCloudwatchLogsExports:  []string{"audit", "profiler"},
				},
			},
			wantOK: true,
		},
		{
			name: "audit logs not enabled",
			clusters: []types.DBCluster{
				{
					DBClusterIdentifier:          aws.String("my-cluster"),
					DBClusterArn:                 aws.String("arn:aws:rds:us-east-1:123456789012:cluster:my-cluster"),
					EnabledCloudwatchLogsExports: []string{"profiler"},
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
			checkIfDocumentDBAuditLogsEnabled(checkConfig, tt.clusters, "AWS_DOC_004")
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
