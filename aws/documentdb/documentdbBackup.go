package documentdb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/docdb/types"
	"github.com/padok-team/yatas/plugins/commons"
)

// MIN_BACKUP_RETENTION_DAYS is the minimum number of days a backup must be retained.
// A value of 1 ensures that at least one daily backup is always available.
const MIN_BACKUP_RETENTION_DAYS = 30

// checkIfDocumentDBBackupEnabled verifies that automated backups are enabled on
// each DocumentDB cluster by checking that BackupRetentionPeriod >= MIN_BACKUP_RETENTION_DAYS.
// A value of 0 means automated backups are completely disabled.
func checkIfDocumentDBBackupEnabled(checkConfig commons.CheckConfig, clusters []types.DBCluster, testName string) {
	var check commons.Check
	check.InitCheck("DocumentDB clusters have automated backups enabled", "Check if DocumentDB backup retention period is at least 1 day", testName, []string{"Security", "Good Practice"})
	for _, cluster := range clusters {
		if aws.ToInt32(cluster.BackupRetentionPeriod) < MIN_BACKUP_RETENTION_DAYS {
			Message := "DocumentDB automated backup is not enabled on " + *cluster.DBClusterIdentifier
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *cluster.DBClusterArn}
			check.AddResult(result)
		} else {
			Message := "DocumentDB automated backup is enabled on " + *cluster.DBClusterIdentifier
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *cluster.DBClusterArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
