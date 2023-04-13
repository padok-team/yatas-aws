package rds

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas-aws/aws/awschecks"
	"github.com/padok-team/yatas/plugins/commons"
	"github.com/padok-team/yatas/plugins/logger"
)

type RDSInstance struct {
	Instance types.DBInstance
}

func (r *RDSInstance) GetID() string {
	return *r.Instance.DBInstanceIdentifier
}

type RDSCluster struct {
	Cluster types.DBCluster
}

func (c *RDSCluster) GetID() string {
	return *c.Cluster.DBClusterIdentifier
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check

	svc := rds.NewFromConfig(s)
	instances := GetListRDS(svc)
	clusters := GetListDBClusters(svc)

	rdsChecks := []awschecks.CheckDefinition{
		{
			Title:          "AWS_RDS_001",
			Description:    "Check if encryption is enabled",
			Tags:           []string{"Security", "Good Practice"},
			ConditionFn:    CheckIfEncryptionEnabled,
			SuccessMessage: "RDS instance is encrypted",
			FailureMessage: "RDS instance is not encrypted",
		},
		{
			Title:          "AWS_RDS_002",
			Description:    "Check if backup is enabled",
			Tags:           []string{"Security", "Good Practice"},
			ConditionFn:    CheckIfBackupEnabled,
			SuccessMessage: "RDS instance has backup enabled",
			FailureMessage: "RDS instance has no backup enabled",
		},
		{
			Title:          "AWS_RDS_003",
			Description:    "Check if auto upgrade is enabled",
			Tags:           []string{"Security", "Good Practice"},
			ConditionFn:    CheckIfAutoUpgradeEnabled,
			SuccessMessage: "RDS instance has auto upgrade enabled",
			FailureMessage: "RDS instance has no auto upgrade enabled",
		},
		{
			Title:          "AWS_RDS_004",
			Description:    "Check if RDS instance is private",
			Tags:           []string{"Security", "Good Practice"},
			ConditionFn:    CheckIfRDSPrivateEnabled,
			SuccessMessage: "RDS instance is private",
			FailureMessage: "RDS instance is not private",
		},
		{
			Title:          "AWS_RDS_005",
			Description:    "Check if logging is enabled",
			Tags:           []string{"Security", "Good Practice"},
			ConditionFn:    CheckIfLoggingEnabled,
			SuccessMessage: "RDS instance has logging enabled",
			FailureMessage: "RDS instance has no logging enabled",
		},
		{
			Title:          "AWS_RDS_006",
			Description:    "Check if delete protection is enabled",
			Tags:           []string{"Security", "Good Practice"},
			ConditionFn:    CheckIfDeleteProtectionEnabled,
			SuccessMessage: "RDS instance has delete protection enabled",
			FailureMessage: "RDS instance has no delete protection enabled",
		},
	}

	rdsClusterChecks := []awschecks.CheckDefinition{
		{
			Title:          "AWS_RDS_007",
			Description:    "Check if cluster auto upgrade is enabled",
			Tags:           []string{"Security", "Good Practice"},
			ConditionFn:    CheckIfClusterAutoUpgradeEnabled,
			SuccessMessage: "RDS cluster has auto upgrade enabled",
			FailureMessage: "RDS cluster has no auto upgrade enabled",
		},
		{
			Title:          "AWS_RDS_008",
			Description:    "Check if cluster backup is enabled",
			Tags:           []string{"Security", "Good Practice"},
			ConditionFn:    CheckIfClusterBackupEnabled,
			SuccessMessage: "RDS cluster has backup enabled",
			FailureMessage: "RDS cluster has no backup enabled",
		},
		{
			Title:          "AWS_RDS_009",
			Description:    "Check if cluster delete protection is enabled",
			Tags:           []string{"Security", "Good Practice"},
			ConditionFn:    CheckIfClusterDeleteProtectionEnabled,
			SuccessMessage: "RDS cluster has delete protection enabled",
			FailureMessage: "RDS cluster has no delete protection enabled",
		},
		{
			Title:          "AWS_RDS_010",
			Description:    "Check if cluster encryption is enabled",
			Tags:           []string{"Security", "Good Practice"},
			ConditionFn:    CheckIfClusterEncryptionEnabled,
			SuccessMessage: "RDS cluster is encrypted",
			FailureMessage: "RDS cluster is not encrypted",
		},
		{
			Title:          "AWS_RDS_011",
			Description:    "Check if cluster logging is enabled",
			Tags:           []string{"Security", "Good Practice"},
			ConditionFn:    CheckIfClusterLoggingEnabled,
			SuccessMessage: "RDS cluster has logging enabled",
			FailureMessage: "RDS cluster has no logging enabled",
		},
		{
			Title:          "AWS_RDS_012",
			Description:    "Check if cluster RDS instance is private",
			Tags:           []string{"Security", "Good Practice"},
			ConditionFn:    CheckIfClusterRDSPrivateEnabled,
			SuccessMessage: "RDS cluster instance is private",
			FailureMessage: "RDS cluster instance is not private",
		},
	}

	var resources []awschecks.Resource
	for _, instance := range instances {
		resources = append(resources, &RDSInstance{Instance: instance})
	}

	var clusterResources []awschecks.Resource
	for _, cluster := range clusters {
		clusterResources = append(clusterResources, &RDSCluster{Cluster: cluster})
	}

	awschecks.AddChecks(&checkConfig, rdsChecks, rdsClusterChecks)
	go awschecks.CheckResources(checkConfig, resources, rdsChecks)
	go awschecks.CheckResources(checkConfig, clusterResources, rdsClusterChecks)

	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)
			checkConfig.Wg.Done()

		}
	}()

	checkConfig.Wg.Wait()
	logger.Logger().Info("RDS checks done")
	logger.Logger().Info("Lenght of checks: ", len(checks))
	queue <- checks
}

func CheckIfEncryptionEnabled(resource awschecks.Resource) bool {
	instanceResource, ok := resource.(*RDSInstance)
	if !ok {
		return false
	}

	instance := instanceResource.Instance
	return instance.StorageEncrypted
}

func CheckIfBackupEnabled(resource awschecks.Resource) bool {
	instanceResource, ok := resource.(*RDSInstance)
	if !ok {
		return false
	}

	instance := instanceResource.Instance
	return instance.BackupRetentionPeriod > 0
}

func CheckIfAutoUpgradeEnabled(resource awschecks.Resource) bool {
	instanceResource, ok := resource.(*RDSInstance)
	if !ok {
		return false
	}

	instance := instanceResource.Instance
	return instance.AutoMinorVersionUpgrade
}

func CheckIfRDSPrivateEnabled(resource awschecks.Resource) bool {
	instanceResource, ok := resource.(*RDSInstance)
	if !ok {
		return false
	}

	instance := instanceResource.Instance
	return !instance.PubliclyAccessible
}

func CheckIfLoggingEnabled(resource awschecks.Resource) bool {
	instanceResource, ok := resource.(*RDSInstance)
	if !ok {
		return false
	}

	instance := instanceResource.Instance
	return len(instance.EnabledCloudwatchLogsExports) > 0
}

func CheckIfDeleteProtectionEnabled(resource awschecks.Resource) bool {
	instanceResource, ok := resource.(*RDSInstance)
	if !ok {
		return false
	}

	instance := instanceResource.Instance
	return instance.DeletionProtection
}

func CheckIfClusterAutoUpgradeEnabled(resource awschecks.Resource) bool {
	clusterResource, ok := resource.(*RDSCluster)
	if !ok {
		return false
	}

	return clusterResource.Cluster.AutoMinorVersionUpgrade
}

func CheckIfClusterBackupEnabled(resource awschecks.Resource) bool {
	clusterResource, ok := resource.(*RDSCluster)
	if !ok {
		return false
	}

	return *clusterResource.Cluster.BackupRetentionPeriod != 0
}

func CheckIfClusterDeleteProtectionEnabled(resource awschecks.Resource) bool {
	clusterResource, ok := resource.(*RDSCluster)
	if !ok {
		return false
	}

	return *clusterResource.Cluster.DeletionProtection
}

func CheckIfClusterEncryptionEnabled(resource awschecks.Resource) bool {
	clusterResource, ok := resource.(*RDSCluster)
	if !ok {
		return false
	}

	return clusterResource.Cluster.StorageEncrypted
}

func CheckIfClusterLoggingEnabled(resource awschecks.Resource) bool {
	clusterResource, ok := resource.(*RDSCluster)
	if !ok {
		return false
	}

	return len(clusterResource.Cluster.EnabledCloudwatchLogsExports) > 0
}

func CheckIfClusterRDSPrivateEnabled(resource awschecks.Resource) bool {
	clusterResource, ok := resource.(*RDSCluster)
	if !ok {
		return false
	}

	return clusterResource.Cluster.PubliclyAccessible != nil && !*clusterResource.Cluster.PubliclyAccessible
}
