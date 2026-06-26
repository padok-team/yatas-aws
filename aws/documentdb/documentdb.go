// Package documentdb implements security and best-practice checks for AWS DocumentDB clusters.
// Each check follows the same pattern as other services in this plugin:
// resources are fetched via getter.go, and individual check functions are
// dispatched as goroutines through commons.CheckTest.
package documentdb

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/docdb"
	"github.com/padok-team/yatas-aws/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

// RunChecks is the entry point for the DocumentDB check suite.
// It is called by the main plugin dispatcher and runs all checks concurrently.
func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {
	logger.Logger.Debug("DocumentDB - Checks started")
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check

	svc := docdb.NewFromConfig(s)

	// Fetch all resources needed across multiple checks upfront
	clusters := GetDocumentDBClusters(svc)
	snapshots := GetDocumentDBClusterSnapshots(svc)

	// Enriched data structures that require additional API calls per resource
	clustersWithTLS := GetClustersWithTLSParam(svc, clusters)
	snapshotsWithPublicFlag := GetSnapshotsWithPublicFlag(svc, snapshots)

	// Launch each check as a goroutine; results are sent to checkConfig.Queue
	go commons.CheckTest(checkConfig.Wg, c, "AWS_DOC_001", checkIfDocumentDBEncryptionEnabled)(checkConfig, clusters, "AWS_DOC_001")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_DOC_002", checkIfDocumentDBBackupEnabled)(checkConfig, clusters, "AWS_DOC_002")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_DOC_003", checkIfDocumentDBDeletionProtectionEnabled)(checkConfig, clusters, "AWS_DOC_003")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_DOC_004", checkIfDocumentDBAuditLogsEnabled)(checkConfig, clusters, "AWS_DOC_004")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_DOC_005", checkIfDocumentDBSnapshotNotPublic)(checkConfig, snapshotsWithPublicFlag, "AWS_DOC_005")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_DOC_006", checkIfDocumentDBTLSEnabled)(checkConfig, clustersWithTLS, "AWS_DOC_006")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_DOC_007", checkIfDocumentDBSnapshotEncrypted)(checkConfig, snapshots, "AWS_DOC_007")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_DOC_008", checkIfDocumentDBMultiAZEnabled)(checkConfig, clusters, "AWS_DOC_008")

	// Collect results from the queue as checks complete
	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()
	queue <- checks
}
