package rds

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas-aws/logger"
)

type RDSGetObjectAPI interface {
	DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
	DescribeDBClusters(ctx context.Context, input *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
	DescribeDBLogFiles(ctx context.Context, input *rds.DescribeDBLogFilesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBLogFilesOutput, error)
	DownloadDBLogFilePortion(ctx context.Context, input *rds.DownloadDBLogFilePortionInput, optFns ...func(*rds.Options)) (*rds.DownloadDBLogFilePortionOutput, error)
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
