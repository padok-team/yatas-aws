package documentdb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/docdb"
	"github.com/aws/aws-sdk-go-v2/service/docdb/types"
	"github.com/padok-team/yatas-aws/logger"
)

// DocDBGetObjectAPI defines the subset of the DocumentDB SDK client used by this package.
// Defining an interface allows tests to inject a mock client.
type DocDBGetObjectAPI interface {
	DescribeDBClusters(ctx context.Context, input *docdb.DescribeDBClustersInput, optFns ...func(*docdb.Options)) (*docdb.DescribeDBClustersOutput, error)
	DescribeDBClusterSnapshots(ctx context.Context, input *docdb.DescribeDBClusterSnapshotsInput, optFns ...func(*docdb.Options)) (*docdb.DescribeDBClusterSnapshotsOutput, error)
	DescribeDBClusterSnapshotAttributes(ctx context.Context, input *docdb.DescribeDBClusterSnapshotAttributesInput, optFns ...func(*docdb.Options)) (*docdb.DescribeDBClusterSnapshotAttributesOutput, error)
	DescribeDBClusterParameters(ctx context.Context, input *docdb.DescribeDBClusterParametersInput, optFns ...func(*docdb.Options)) (*docdb.DescribeDBClusterParametersOutput, error)
}

// ClusterWithTLSParam groups a DocumentDB cluster with the resolved value of its
// "tls" parameter fetched from the associated cluster parameter group.
type ClusterWithTLSParam struct {
	Cluster  types.DBCluster
	TLSValue string // "enabled" or "disabled"
}

// SnapshotWithPublicFlag groups a DocumentDB cluster snapshot with a boolean
// indicating whether the snapshot is publicly accessible (restore attribute = "all").
type SnapshotWithPublicFlag struct {
	Snapshot types.DBClusterSnapshot
	IsPublic bool
}

// GetDocumentDBClusters returns all DocumentDB clusters in the account.
// Pagination is handled automatically via the Marker field.
func GetDocumentDBClusters(svc DocDBGetObjectAPI) []types.DBCluster {
	var clusters []types.DBCluster
	params := &docdb.DescribeDBClustersInput{}
	resp, err := svc.DescribeDBClusters(context.TODO(), params)
	if err != nil {
		logger.Logger.Error(err.Error())
		return []types.DBCluster{}
	}
	clusters = append(clusters, resp.DBClusters...)
	for {
		if resp.Marker != nil {
			params.Marker = resp.Marker
			resp, err = svc.DescribeDBClusters(context.TODO(), params)
			if err != nil {
				logger.Logger.Error(err.Error())
				return []types.DBCluster{}
			}
			clusters = append(clusters, resp.DBClusters...)
		} else {
			break
		}
	}
	return clusters
}

// GetDocumentDBClusterSnapshots returns all manual snapshots for DocumentDB clusters.
// Only manual snapshots are fetched (SnapshotType = "manual") because automated
// snapshots cannot be made public and do not require a public-access check.
func GetDocumentDBClusterSnapshots(svc DocDBGetObjectAPI) []types.DBClusterSnapshot {
	var snapshots []types.DBClusterSnapshot
	snapshotType := "manual"
	params := &docdb.DescribeDBClusterSnapshotsInput{
		SnapshotType: &snapshotType,
	}
	resp, err := svc.DescribeDBClusterSnapshots(context.TODO(), params)
	if err != nil {
		logger.Logger.Error(err.Error())
		return []types.DBClusterSnapshot{}
	}
	snapshots = append(snapshots, resp.DBClusterSnapshots...)
	for {
		if resp.Marker != nil {
			params.Marker = resp.Marker
			resp, err = svc.DescribeDBClusterSnapshots(context.TODO(), params)
			if err != nil {
				logger.Logger.Error(err.Error())
				return []types.DBClusterSnapshot{}
			}
			snapshots = append(snapshots, resp.DBClusterSnapshots...)
		} else {
			break
		}
	}
	return snapshots
}

// GetClustersWithTLSParam resolves the "tls" parameter for each cluster by querying
// its associated DBClusterParameterGroup. The TLS parameter controls whether
// clients are required to use TLS when connecting to the cluster.
// If the parameter cannot be fetched, the value defaults to "disabled" (fail-safe).
func GetClustersWithTLSParam(svc DocDBGetObjectAPI, clusters []types.DBCluster) []ClusterWithTLSParam {
	var result []ClusterWithTLSParam
	for _, cluster := range clusters {
		tlsValue := "disabled"
		params := &docdb.DescribeDBClusterParametersInput{
			DBClusterParameterGroupName: cluster.DBClusterParameterGroup,
		}
		resp, err := svc.DescribeDBClusterParameters(context.TODO(), params)
		if err != nil {
			logger.Logger.Error(err.Error())
			result = append(result, ClusterWithTLSParam{Cluster: cluster, TLSValue: tlsValue})
			continue
		}
		// Iterate over parameters to find the "tls" entry
		for _, param := range resp.Parameters {
			if param.ParameterName != nil && *param.ParameterName == "tls" && param.ParameterValue != nil {
				tlsValue = *param.ParameterValue
				break
			}
		}
		result = append(result, ClusterWithTLSParam{Cluster: cluster, TLSValue: tlsValue})
	}
	return result
}

// GetSnapshotsWithPublicFlag resolves whether each snapshot is publicly accessible
// by checking its restore attribute via DescribeDBClusterSnapshotAttributes.
// A snapshot is considered public when the "restore" attribute contains the value "all".
// If attributes cannot be fetched, IsPublic defaults to false (fail-safe).
func GetSnapshotsWithPublicFlag(svc DocDBGetObjectAPI, snapshots []types.DBClusterSnapshot) []SnapshotWithPublicFlag {
	var result []SnapshotWithPublicFlag
	for _, snapshot := range snapshots {
		isPublic := false
		resp, err := svc.DescribeDBClusterSnapshotAttributes(context.TODO(), &docdb.DescribeDBClusterSnapshotAttributesInput{
			DBClusterSnapshotIdentifier: snapshot.DBClusterSnapshotIdentifier,
		})
		if err != nil {
			logger.Logger.Error(err.Error())
			result = append(result, SnapshotWithPublicFlag{Snapshot: snapshot, IsPublic: false})
			continue
		}
		if resp.DBClusterSnapshotAttributesResult != nil {
			for _, attr := range resp.DBClusterSnapshotAttributesResult.DBClusterSnapshotAttributes {
				// The "restore" attribute lists the AWS account IDs (or "all") allowed to restore the snapshot
				if attr.AttributeName != nil && *attr.AttributeName == "restore" {
					for _, val := range attr.AttributeValues {
						if val == "all" {
							isPublic = true
							break
						}
					}
				}
			}
		}
		result = append(result, SnapshotWithPublicFlag{Snapshot: snapshot, IsPublic: isPublic})
	}
	return result
}
