package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	backupTypes "github.com/aws/aws-sdk-go-v2/service/backup/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/padok-team/yatas-aws/logger"
)

func GetDynamodbs(s aws.Config) []string {
	svc := dynamodb.NewFromConfig(s)
	dynamodbInput := &dynamodb.ListTablesInput{}
	result, err := svc.ListTables(context.TODO(), dynamodbInput)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list
		return []string{}
	}
	return result.TableNames
}

func GetTables(s aws.Config, dynamodbs []string) []*dynamodb.DescribeTableOutput {
	svc := dynamodb.NewFromConfig(s)
	var tables []*dynamodb.DescribeTableOutput
	for _, d := range dynamodbs {
		params := &dynamodb.DescribeTableInput{
			TableName: &d,
		}
		resp, err := svc.DescribeTable(context.TODO(), params)
		if err != nil {
			logger.Logger.Error(err.Error())
			// Return an empty list of certificates
			return []*dynamodb.DescribeTableOutput{}
		}
		tables = append(tables, resp)

	}
	return tables
}

type TableBackups struct {
	TableName string
	Backups   ddbTypes.ContinuousBackupsDescription
}

func GetContinuousBackups(s aws.Config, tables []string) []TableBackups {
	svc := dynamodb.NewFromConfig(s)
	var continuousBackups []TableBackups
	for _, d := range tables {
		params := &dynamodb.DescribeContinuousBackupsInput{
			TableName: &d,
		}
		resp, err := svc.DescribeContinuousBackups(context.TODO(), params)
		if err != nil {
			logger.Logger.Error(err.Error())
			// Return an empty list of certificates
			return []TableBackups{}
		}
		continuousBackups = append(continuousBackups, TableBackups{d, *resp.ContinuousBackupsDescription})
	}
	return continuousBackups
}

type TableRecoveryPoints struct {
	TableName      string
	RecoveryPoints []backupTypes.RecoveryPointByResource
}

func GetTableRecoveryPoints(s aws.Config, tables []*dynamodb.DescribeTableOutput) []TableRecoveryPoints {
	svc := backup.NewFromConfig(s)
	var recoveryPoints []TableRecoveryPoints
	for _, t := range tables {
		params := &backup.ListRecoveryPointsByResourceInput{
			ResourceArn: t.Table.TableArn,
		}
		resp, err := svc.ListRecoveryPointsByResource(context.TODO(), params)
		if err != nil {
			logger.Logger.Error(err.Error())
			// Return an empty list of certificates
			return []TableRecoveryPoints{}
		}
		recoveryPoints = append(recoveryPoints, TableRecoveryPoints{*t.Table.TableName, resp.RecoveryPoints})
	}
	return recoveryPoints
}
