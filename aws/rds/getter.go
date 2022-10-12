package rds

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type RDSGetObjectAPI interface {
	DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
	DescribeDBClusters(ctx context.Context, input *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
}

func GetListRDS(svc RDSGetObjectAPI) []types.DBInstance {

	params := &rds.DescribeDBInstancesInput{}
	var instances []types.DBInstance
	resp, err := svc.DescribeDBInstances(context.TODO(), params)
	instances = append(instances, resp.DBInstances...)
	if err != nil {
		fmt.Println(err)
	}
	for {
		if resp.Marker != nil {
			params.Marker = resp.Marker
			resp, err = svc.DescribeDBInstances(context.TODO(), params)
			instances = append(instances, resp.DBInstances...)
			if err != nil {
				fmt.Println(err)
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
		fmt.Println(err)
	}
	for {
		if resp.Marker != nil {
			params.Marker = resp.Marker
			resp, err = svc.DescribeDBClusters(context.TODO(), params)
			clusters = append(clusters, resp.DBClusters...)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			break
		}
	}

	return clusters
}
