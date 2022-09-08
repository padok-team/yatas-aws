package volumes

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/config"
)

// Main function that runs all the tests and returns the results
func RunChecks(wa *sync.WaitGroup, s aws.Config, c *config.Config, queue chan []config.Check) {

	var checkConfig config.CheckConfig
	checkConfig.Init(s, c)
	var checks []config.Check

	volumes := GetVolumes(s)
	snapshots := GetSnapshots(s)
	couples := couple{volumes, snapshots}

	go config.CheckTest(checkConfig.Wg, c, "AWS_VOL_001", checkIfEncryptionEnabled)(checkConfig, volumes, "AWS_VOL_001")
	go config.CheckTest(checkConfig.Wg, c, "AWS_VOL_002", CheckIfVolumesTypeGP3)(checkConfig, volumes, "AWS_VOL_002")
	go config.CheckTest(checkConfig.Wg, c, "AWS_VOL_003", CheckIfAllVolumesHaveSnapshots)(checkConfig, couples, "AWS_VOL_003")
	go config.CheckTest(checkConfig.Wg, c, "AWS_VOL_004", CheckIfVolumeIsUsed)(checkConfig, volumes, "AWS_VOL_004")

	go config.CheckTest(checkConfig.Wg, c, "AWS_BAK_001", CheckIfAllSnapshotsEncrypted)(checkConfig, snapshots, "AWS_BAK_001")
	go config.CheckTest(checkConfig.Wg, c, "AWS_BAK_002", CheckIfSnapshotYoungerthan24h)(checkConfig, couples, "AWS_BAK_002")

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
