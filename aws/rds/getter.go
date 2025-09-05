package rds

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas-aws/logger"
)

type RDSGetObjectAPI interface {
	DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
	DescribeDBClusters(ctx context.Context, input *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
	DescribeDBLogFiles(ctx context.Context, input *rds.DescribeDBLogFilesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBLogFilesOutput, error)
	DownloadDBLogFilePortion(ctx context.Context, input *rds.DownloadDBLogFilePortionInput, optFns ...func(*rds.Options)) (*rds.DownloadDBLogFilePortionOutput, error)
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

func GetListDBLogFiles(svc RDSGetObjectAPI, dbInstanceIdentifier string) []types.DescribeDBLogFilesDetails {
	params := &rds.DescribeDBLogFilesInput{
		DBInstanceIdentifier: &dbInstanceIdentifier,
	}
	var logFiles []types.DescribeDBLogFilesDetails
	resp, err := svc.DescribeDBLogFiles(context.TODO(), params)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list of instances
		return []types.DescribeDBLogFilesDetails{}
	}
	logFiles = append(logFiles, resp.DescribeDBLogFiles...)
	for {
		if resp.Marker != nil {
			params.Marker = resp.Marker
			resp, err = svc.DescribeDBLogFiles(context.TODO(), params)
			logFiles = append(logFiles, resp.DescribeDBLogFiles...)
			if err != nil {
				logger.Logger.Error(err.Error())
				// Return an empty list of instances
				return []types.DescribeDBLogFilesDetails{}
			}
		} else {
			break
		}
	}
	return logFiles
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

type InstanceToLogFiles struct {
	Instance              types.DBInstance
	LogFiles              []types.DescribeDBLogFilesDetails
	RecentLogFilesPortion string
}

func GetInstancesToLogFiles(svc RDSGetObjectAPI, instances []types.DBInstance) []InstanceToLogFiles {
	var wg sync.WaitGroup
	rdsToLogFilesChan := make(chan InstanceToLogFiles, len(instances))
	var rdsToLogFiles []InstanceToLogFiles

	// Rate limit to 5 concurrent requests
	sem := make(chan struct{}, 5)

	for _, instance := range instances {
		wg.Add(1)
		go func(instance types.DBInstance) {
			defer wg.Done()

			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }()

			instanceToLogFile := InstanceToLogFiles{Instance: instance}
			logFiles := GetListDBLogFiles(svc, *instance.DBInstanceIdentifier)
			instanceToLogFile.LogFiles = logFiles

			if *instance.Engine == "postgres" || *instance.Engine == "aurora-postgresql" {
				slices.SortFunc(logFiles, func(a, b types.DescribeDBLogFilesDetails) int {
					return int(*a.LastWritten) - int(*b.LastWritten)
				})
				logFileCount := len(logFiles)
				if logFileCount > 3 {
					logFileCount = 3
				}

				var portionWg sync.WaitGroup
				portions := make(chan string, logFileCount)

				// Inner semaphore for log file portions, limit to 3 concurrent
				portionSem := make(chan struct{}, 3)

				for i := 0; i < logFileCount; i++ {
					portionWg.Add(1)
					go func(i int) {
						defer portionWg.Done()

						// Acquire portion semaphore
						portionSem <- struct{}{}
						defer func() { <-portionSem }()

						// Add small delay between requests
						time.Sleep(100 * time.Millisecond)

						portion := GetLogFilePortion(svc, *instance.DBInstanceIdentifier, *logFiles[i].LogFileName, "")
						portions <- portion
					}(i)
				}

				go func() {
					portionWg.Wait()
					close(portions)
				}()

				for portion := range portions {
					instanceToLogFile.RecentLogFilesPortion += portion
				}
			}
			rdsToLogFilesChan <- instanceToLogFile
		}(instance)
	}

	go func() {
		wg.Wait()
		close(rdsToLogFilesChan)
	}()

	for result := range rdsToLogFilesChan {
		rdsToLogFiles = append(rdsToLogFiles, result)
	}

	return rdsToLogFiles
}

func GetLogFilePortion(svc RDSGetObjectAPI, dbInstanceIdentifier string, logFileName string, marker string) string {
	params := &rds.DownloadDBLogFilePortionInput{
		DBInstanceIdentifier: &dbInstanceIdentifier,
		LogFileName:          &logFileName,
	}
	if marker != "" {
		params.Marker = &marker
	}
	resp, err := svc.DownloadDBLogFilePortion(context.TODO(), params)
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	return *resp.LogFileData
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

type InstanceWithSGs struct {
	Instance       types.DBInstance
	SecurityGroups []ec2Types.SecurityGroup
}

func GetInstancesWithSGs(rdsSvc RDSGetObjectAPI, ec2Svc *ec2.Client) []InstanceWithSGs {
	rdsInstances := GetListRDS(rdsSvc)
	result := make([]InstanceWithSGs, 0, len(rdsInstances))

	for _, inst := range rdsInstances {
		groupIDs := []string{}
		for _, sg := range inst.VpcSecurityGroups {
			if sg.VpcSecurityGroupId != nil {
				groupIDs = append(groupIDs, *sg.VpcSecurityGroupId)
			}
		}

		if len(groupIDs) == 0 {
			result = append(result, InstanceWithSGs{Instance: inst, SecurityGroups: []ec2Types.SecurityGroup{}})
			continue
		}

		resp, err := ec2Svc.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
			GroupIds: groupIDs,
		})
		if err != nil {
			logger.Logger.Error(err.Error())
			continue
		}

		result = append(result, InstanceWithSGs{
			Instance:       inst,
			SecurityGroups: resp.SecurityGroups,
		})
	}

	return result
}

type DBClusterWithSGs struct {
	Cluster        types.DBCluster
	SecurityGroups []ec2Types.SecurityGroup
}

func GetDBClustersWithSGs(rdsSvc RDSGetObjectAPI, ec2Svc *ec2.Client) []DBClusterWithSGs {
	clusters := GetListDBClusters(rdsSvc)
	result := make([]DBClusterWithSGs, 0, len(clusters))

	for _, cluster := range clusters {

		groupIDs := []string{}
		for _, sg := range cluster.VpcSecurityGroups {
			if sg.VpcSecurityGroupId != nil {
				groupIDs = append(groupIDs, *sg.VpcSecurityGroupId)
			}
		}

		if len(groupIDs) == 0 {
			result = append(result, DBClusterWithSGs{Cluster: cluster, SecurityGroups: []ec2Types.SecurityGroup{}})
			continue
		}

		resp, err := ec2Svc.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
			GroupIds: groupIDs,
		})
		if err != nil {
			logger.Logger.Error(err.Error())
			continue
		}

		result = append(result, DBClusterWithSGs{
			Cluster:        cluster,
			SecurityGroups: resp.SecurityGroups,
		})
	}

	return result
}
