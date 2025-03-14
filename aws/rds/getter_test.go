package rds

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type mockGetRdsAPI func(ctx context.Context, input *rds.DescribeDBInstancesInput) (output *rds.DescribeDBInstancesOutput, err error)

func (m mockGetRdsAPI) DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error) {
	// Return an empty list of RDS instances
	return &rds.DescribeDBInstancesOutput{
		DBInstances: []types.DBInstance{},
	}, nil

}

func (m mockGetRdsAPI) DescribeDBClusters(ctx context.Context, input *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error) {
	// Return an empty list of RDS clusters
	return &rds.DescribeDBClustersOutput{
		DBClusters: []types.DBCluster{},
	}, nil
}

func (m mockGetRdsAPI) DescribeDBLogFiles(ctx context.Context, input *rds.DescribeDBLogFilesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBLogFilesOutput, error) {
	// Return an empty list of RDS log files
	return &rds.DescribeDBLogFilesOutput{
		DescribeDBLogFiles: []types.DescribeDBLogFilesDetails{},
	}, nil
}

func (m mockGetRdsAPI) DownloadDBLogFilePortion(ctx context.Context, input *rds.DownloadDBLogFilePortionInput, optFns ...func(*rds.Options)) (*rds.DownloadDBLogFilePortionOutput, error) {
	// Return an empty list of RDS log files
	return &rds.DownloadDBLogFilePortionOutput{
		LogFileData: aws.String("AUDIT: "),
	}, nil
}

func TestGetListRDS(t *testing.T) {
	tests := []struct {
		name string
		want []types.DBInstance
	}{
		{
			name: "Empty list of RDS instances",
			want: []types.DBInstance{},
		},
	}
	mockGetRdsAPI := mockGetRdsAPI(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetListRDS(mockGetRdsAPI); len(got) != 0 {
				t.Errorf("GetListRDS() = %+v, want %+v", got, tt.want)
			}
			if got := GetListDBClusters(mockGetRdsAPI); len(got) != 0 {
				t.Errorf("GetListDBClusters() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
