package s3

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckS3Location(s aws.Config, bucket, region string) bool {

	svc := s3.NewFromConfig(s)

	params := &s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	}
	resp, err := svc.GetBucketLocation(context.TODO(), params)
	if err != nil {

		return false
	}

	if resp.LocationConstraint != "" {
		if string(resp.LocationConstraint) == region {
			return true
		} else {
			return false
		}

	} else {
		return false
	}
}

type BucketAndNotInRegion struct {
	Buckets     []types.Bucket
	NotInRegion []types.Bucket
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {

	var checkConfig commons.CheckConfig
	checkConfig.Init(s, c)
	var checks []commons.Check
	buckets := GetListS3(s)
	bucketsNotInRegion := GetListS3NotInRegion(s, s.Region)
	couple := BucketAndNotInRegion{buckets, bucketsNotInRegion}
	OnlyBucketInRegion := OnlyBucketInRegion(couple)
	S3ToEncryption := GetS3ToEncryption(s, OnlyBucketInRegion)
	S3toPublicBlockAccess := GetS3ToPublicBlockAccess(s, OnlyBucketInRegion)
	S3ToVersioning := GetS3ToVersioning(s, OnlyBucketInRegion)
	S3ToObjectLock := GetS3ToObjectLock(s, OnlyBucketInRegion)

	go commons.CheckTest(checkConfig.Wg, c, "AWS_S3_001", checkIfEncryptionEnabled)(checkConfig, S3ToEncryption, "AWS_S3_001")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_S3_002", CheckIfBucketInOneZone)(checkConfig, couple, "AWS_S3_002")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_S3_003", CheckIfBucketObjectVersioningEnabled)(checkConfig, S3ToVersioning, "AWS_S3_003")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_S3_004", CheckIfObjectLockConfigurationEnabled)(checkConfig, S3ToObjectLock, "AWS_S3_004")
	go commons.CheckTest(checkConfig.Wg, c, "AWS_S3_005", CheckIfS3PublicAccessBlockEnabled)(checkConfig, S3toPublicBlockAccess, "AWS_S3_005")
	// Wait for all the goroutines to finish

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
