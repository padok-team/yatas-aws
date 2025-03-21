package rds

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas-aws/logger"
)

type RDSGetObjectAPI interface {
	DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
	DescribeDBClusters(ctx context.Context, input *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
	DescribeDBSnapshots(ctx context.Context, input *rds.DescribeDBSnapshotsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBSnapshotsOutput, error)
	DescribeDBClusterSnapshots(ctx context.Context, input *rds.DescribeDBClusterSnapshotsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClusterSnapshotsOutput, error)
}

func GetListRDS(svc RDSGetObjectAPI) []types.DBInstance {
	params := &rds.DescribeDBInstancesInput{}
	var instances []types.DBInstance
	resp, err := svc.DescribeDBInstances(context.TODO(), params)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list of instances
		return []types.DBInstance{}
	}
	instances = append(instances, resp.DBInstances...)
	for {
		if resp.Marker != nil {
			params.Marker = resp.Marker
			resp, err = svc.DescribeDBInstances(context.TODO(), params)
			instances = append(instances, resp.DBInstances...)
			if err != nil {
				logger.Logger.Error(err.Error())
				// Return an empty list of instances
				return []types.DBInstance{}
			}
		} else {
			break
		}
	}

	return instances
}

func GetListDBClusters(svc RDSGetObjectAPI) []types.DBCluster {
	params := &rds.DescribeDBClustersInput{}
	var clusters []types.DBCluster
	resp, err := svc.DescribeDBClusters(context.TODO(), params)
	clusters = append(clusters, resp.DBClusters...)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list of instances
		return []types.DBCluster{}
	}
	for {
		if resp.Marker != nil {
			params.Marker = resp.Marker
			resp, err = svc.DescribeDBClusters(context.TODO(), params)
			clusters = append(clusters, resp.DBClusters...)
			if err != nil {
				logger.Logger.Error(err.Error())
				// Return an empty list of instances
				return []types.DBCluster{}
			}
		} else {
			break
		}
	}

	return clusters
}

func GetListRDSSnapshots(svc RDSGetObjectAPI) []types.DBSnapshot {
	params := &rds.DescribeDBSnapshotsInput{}
	var snapshots []types.DBSnapshot
	resp, err := svc.DescribeDBSnapshots(context.TODO(), params)
	snapshots = append(snapshots, resp.DBSnapshots...)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list of instances
		return []types.DBSnapshot{}
	}
	for {
		if resp.Marker != nil {
			params.Marker = resp.Marker
			resp, err = svc.DescribeDBSnapshots(context.TODO(), params)
			snapshots = append(snapshots, resp.DBSnapshots...)
			if err != nil {
				logger.Logger.Error(err.Error())
				// Return an empty list of instances
				return []types.DBSnapshot{}
			}
		} else {
			break
		}
	}

	return snapshots
}

func GetListDBClusterSnapshots(svc RDSGetObjectAPI) []types.DBClusterSnapshot {
	params := &rds.DescribeDBClusterSnapshotsInput{}
	var snapshots []types.DBClusterSnapshot
	resp, err := svc.DescribeDBClusterSnapshots(context.TODO(), params)
	snapshots = append(snapshots, resp.DBClusterSnapshots...)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list of instances
		return []types.DBClusterSnapshot{}
	}
	for {
		if resp.Marker != nil {
			params.Marker = resp.Marker
			resp, err = svc.DescribeDBClusterSnapshots(context.TODO(), params)
			snapshots = append(snapshots, resp.DBClusterSnapshots...)
			if err != nil {
				logger.Logger.Error(err.Error())
				// Return an empty list of instances
				return []types.DBClusterSnapshot{}
			}
		} else {
			break
		}
	}

	return snapshots
}
