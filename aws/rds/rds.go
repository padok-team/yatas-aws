package rds

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/padok-team/yatas-aws/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {
	logger.Logger.Debug("RDS - Checks started")
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check
	svc := rds.NewFromConfig(s)

	instances := GetListRDS(svc)
	clusters := GetListDBClusters(svc)
	rdsSnapshots := GetListRDSSnapshots(svc)
	auroraSnapshots := GetListDBClusterSnapshots(svc)

	instanceList := []types.DBInstance{}
	instanceList = append(instanceList, instances...)
	instancesToLogFiles := GetInstancesToLogFiles(svc, instanceList)
	ec2Svc := ec2.NewFromConfig(s)
	instancesWithSGs := GetInstancesWithSGs(svc, ec2Svc)

	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_001", checkIfEncryptionEnabled)(checkConfig, instances, "AWS_RDS_001")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_002", checkIfBackupEnabled)(checkConfig, instances, "AWS_RDS_002")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_003", checkIfAutoUpgradeEnabled)(checkConfig, instances, "AWS_RDS_003")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_004", checkIfRDSPrivateEnabled)(checkConfig, instances, "AWS_RDS_004")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_005", CheckIfLoggingEnabled)(checkConfig, instances, "AWS_RDS_005")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_006", CheckIfDeleteProtectionEnabled)(checkConfig, instances, "AWS_RDS_006")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_007", checkIfClusterAutoUpgradeEnabled)(checkConfig, clusters, "AWS_RDS_007")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_008", checkIfClusterBackupEnabled)(checkConfig, clusters, "AWS_RDS_008")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_009", CheckIfClusterDeleteProtectionEnabled)(checkConfig, clusters, "AWS_RDS_009")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_010", checkIfClusterEncryptionEnabled)(checkConfig, clusters, "AWS_RDS_010")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_011", CheckIfClusterLoggingEnabled)(checkConfig, clusters, "AWS_RDS_011")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_012", checkIfClusterRDSPrivateEnabled)(checkConfig, clusters, "AWS_RDS_012")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_013", checkIfAuditLogsEnabled)(checkConfig, instancesToLogFiles, "AWS_RDS_013")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_014", checkIfRDSSnapshotEncryptionEnabled)(checkConfig, rdsSnapshots, "AWS_RDS_014")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_015", checkIfClusterSnapshotEncryptionEnabled)(checkConfig, auroraSnapshots, "AWS_RDS_015")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_RDS_016", checkIfRDSRestrictedSecurityGroups)(checkConfig, instancesWithSGs, "AWS_RDS_016")

	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)

			checkConfig.Wg.Done()

		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
	logger.Logger.Debug("RDS - Checks done")
}
