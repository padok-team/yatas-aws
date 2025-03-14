package rds

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfAuditLogsEnabledSuccess(t *testing.T) {
	type args struct {
		checkConfig commons.CheckConfig
		dbInstances []InstanceToLogFiles
		testName    string
	}

	// Mock log file portion to avoid nil pointer in getter
	logPortion := "test log portion"
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test MySQL RDS with audit logs",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				dbInstances: []InstanceToLogFiles{
					{
						Instance: types.DBInstance{
							DBInstanceIdentifier: aws.String("mysql-db"),
							DBInstanceArn:        aws.String("arn:aws:rds:us-west-2:123456789012:db:mysql-db"),
							Engine:               aws.String("mysql"),
						},
						LogFiles: []types.DescribeDBLogFilesDetails{
							{
								LogFileName: aws.String("audit.log"),
								LastWritten: aws.Int64(1640995200000), // 2022-01-01
							},
						},
						RecentLogFilesPortion: logPortion,
					},
				},
				testName: "AWS_RDS_001",
			},
		},
		{
			name: "Test Aurora MySQL with audit logs",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				dbInstances: []InstanceToLogFiles{
					{
						Instance: types.DBInstance{
							DBInstanceIdentifier: aws.String("aurora-mysql-db"),
							DBInstanceArn:        aws.String("arn:aws:rds:us-west-2:123456789012:db:aurora-mysql-db"),
							Engine:               aws.String("aurora-mysql"),
						},
						LogFiles: []types.DescribeDBLogFilesDetails{
							{
								LogFileName: aws.String("audit/audit.log"),
								LastWritten: aws.Int64(1640995200000), // 2022-01-01
							},
						},
						RecentLogFilesPortion: logPortion,
					},
				},
				testName: "AWS_RDS_001",
			},
		},
		{
			name: "Test PostgreSQL RDS with audit logs",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				dbInstances: []InstanceToLogFiles{
					{
						Instance: types.DBInstance{
							DBInstanceIdentifier: aws.String("postgres-db"),
							DBInstanceArn:        aws.String("arn:aws:rds:us-west-2:123456789012:db:postgres-db"),
							Engine:               aws.String("postgres"),
						},
						RecentLogFilesPortion: "AUDIT: SESSION,1,1,WRITE,2024-01-01 00:00:00 UTC",
					},
				},
				testName: "AWS_RDS_001",
			},
		},
		{
			name: "Test Aurora PostgreSQL with audit logs",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				dbInstances: []InstanceToLogFiles{
					{
						Instance: types.DBInstance{
							DBInstanceIdentifier: aws.String("aurora-postgres-db"),
							DBInstanceArn:        aws.String("arn:aws:rds:us-west-2:123456789012:db:aurora-postgres-db"),
							Engine:               aws.String("aurora-postgresql"),
						},
						RecentLogFilesPortion: "AUDIT: SESSION,1,1,DDL,2024-01-01 00:00:00 UTC",
					},
				},
				testName: "AWS_RDS_001",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfAuditLogsEnabled(tt.args.checkConfig, tt.args.dbInstances, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("checkIfAuditLogsEnabled() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfAuditLogsEnabledFail(t *testing.T) {
	type args struct {
		checkConfig commons.CheckConfig
		dbInstances []InstanceToLogFiles
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test MySQL RDS without audit logs",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				dbInstances: []InstanceToLogFiles{
					{
						Instance: types.DBInstance{
							DBInstanceIdentifier: aws.String("mysql-db"),
							DBInstanceArn:        aws.String("arn:aws:rds:us-west-2:123456789012:db:mysql-db"),
							Engine:               aws.String("mysql"),
						},
						LogFiles: []types.DescribeDBLogFilesDetails{
							{
								LogFileName: aws.String("error.log"),
							},
						},
					},
				},
				testName: "AWS_RDS_001",
			},
		},
		{
			name: "Test Aurora MySQL without audit logs",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				dbInstances: []InstanceToLogFiles{
					{
						Instance: types.DBInstance{
							DBInstanceIdentifier: aws.String("aurora-mysql-db"),
							DBInstanceArn:        aws.String("arn:aws:rds:us-west-2:123456789012:db:aurora-mysql-db"),
							Engine:               aws.String("aurora-mysql"),
						},
						LogFiles: []types.DescribeDBLogFilesDetails{
							{
								LogFileName: aws.String("slowquery.log"),
							},
						},
					},
				},
				testName: "AWS_RDS_001",
			},
		},
		{
			name: "Test PostgreSQL RDS without audit logs",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				dbInstances: []InstanceToLogFiles{
					{
						Instance: types.DBInstance{
							DBInstanceIdentifier: aws.String("postgres-db"),
							DBInstanceArn:        aws.String("arn:aws:rds:us-west-2:123456789012:db:postgres-db"),
							Engine:               aws.String("postgres"),
						},
						RecentLogFilesPortion: "LOG: database system is ready to accept connections",
					},
				},
				testName: "AWS_RDS_001",
			},
		},
		{
			name: "Test Aurora PostgreSQL without audit logs",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				dbInstances: []InstanceToLogFiles{
					{
						Instance: types.DBInstance{
							DBInstanceIdentifier: aws.String("aurora-postgres-db"),
							DBInstanceArn:        aws.String("arn:aws:rds:us-west-2:123456789012:db:aurora-postgres-db"),
							Engine:               aws.String("aurora-postgresql"),
						},
						RecentLogFilesPortion: "LOG: checkpoint starting: time",
					},
				},
				testName: "AWS_RDS_001",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfAuditLogsEnabled(tt.args.checkConfig, tt.args.dbInstances, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("checkIfAuditLogsEnabled() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
