package s3

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/padok-team/yatas-aws/logger"
)

func GetListS3(s aws.Config) []types.Bucket {

	svc := s3.NewFromConfig(s)

	params := &s3.ListBucketsInput{}
	resp, err := svc.ListBuckets(context.TODO(), params)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list
		return []types.Bucket{}
	}

	return resp.Buckets
}

func GetListS3NotInRegion(s aws.Config, region string) []types.Bucket {

	svc := s3.NewFromConfig(s)

	params := &s3.ListBucketsInput{}
	resp, err := svc.ListBuckets(context.TODO(), params)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list
		return []types.Bucket{}
	}

	var buckets []types.Bucket
	for _, bucket := range resp.Buckets {
		check, _ := CheckS3Location(s, *bucket.Name, region)
		if !check {
			buckets = append(buckets, bucket)
		}
	}

	return buckets
}

type S3toPublicBlockAccess struct {
	BucketName string
	Config     bool
}

func GetS3ToPublicBlockAccess(s aws.Config, b []types.Bucket) []S3toPublicBlockAccess {

	svc := s3.NewFromConfig(s)

	var s3toPublicBlockAccess []S3toPublicBlockAccess
	for _, bucket := range b {
		params := &s3.GetPublicAccessBlockInput{
			Bucket: aws.String(*bucket.Name),
		}
		resp, err := svc.GetPublicAccessBlock(context.TODO(), params)
		if err != nil {
			if resp != nil && resp.PublicAccessBlockConfiguration != nil && resp.PublicAccessBlockConfiguration.BlockPublicAcls {
				s3toPublicBlockAccess = append(s3toPublicBlockAccess, S3toPublicBlockAccess{*bucket.Name, true})
			} else {
				s3toPublicBlockAccess = append(s3toPublicBlockAccess, S3toPublicBlockAccess{*bucket.Name, false})
			}
		} else {
			s3toPublicBlockAccess = append(s3toPublicBlockAccess, S3toPublicBlockAccess{*bucket.Name, false})
		}
	}

	return s3toPublicBlockAccess
}

type S3ToEncryption struct {
	BucketName string
	Encrypted  bool
}

func GetS3ToEncryption(s aws.Config, b []types.Bucket) []S3ToEncryption {

	svc := s3.NewFromConfig(s)

	var s3toEncryption []S3ToEncryption
	for _, bucket := range b {
		params := &s3.GetBucketEncryptionInput{
			Bucket: aws.String(*bucket.Name),
		}
		_, err := svc.GetBucketEncryption(context.TODO(), params)
		if err != nil && !strings.Contains(err.Error(), "ServerSideEncryptionConfigurationNotFoundError") {
			logger.Logger.Error(err.Error())
		} else if err != nil {
			s3toEncryption = append(s3toEncryption, S3ToEncryption{*bucket.Name, false})
		} else {
			s3toEncryption = append(s3toEncryption, S3ToEncryption{*bucket.Name, true})
		}
	}

	return s3toEncryption
}

type S3ToVersioning struct {
	BucketName string
	Versioning bool
}

func GetS3ToVersioning(s aws.Config, b []types.Bucket) []S3ToVersioning {

	svc := s3.NewFromConfig(s)

	var s3toVersioning []S3ToVersioning
	for _, bucket := range b {
		params := &s3.GetBucketVersioningInput{
			Bucket: aws.String(*bucket.Name),
		}
		resp, err := svc.GetBucketVersioning(context.TODO(), params)
		if err != nil {
			logger.Logger.Error(err.Error())
			// return empty	struct
			return []S3ToVersioning{}
		}
		if resp.Status != types.BucketVersioningStatusEnabled {
			s3toVersioning = append(s3toVersioning, S3ToVersioning{*bucket.Name, false})
		} else {
			s3toVersioning = append(s3toVersioning, S3ToVersioning{*bucket.Name, true})
		}
	}

	return s3toVersioning
}

type S3ToObjectLock struct {
	BucketName string
	ObjectLock bool
}

func GetS3ToObjectLock(s aws.Config, b []types.Bucket) []S3ToObjectLock {

	svc := s3.NewFromConfig(s)

	var s3toObjectLock []S3ToObjectLock
	for _, bucket := range b {
		params := &s3.GetObjectLockConfigurationInput{
			Bucket: aws.String(*bucket.Name),
		}
		resp, err := svc.GetObjectLockConfiguration(context.TODO(), params)
		if err != nil || (resp.ObjectLockConfiguration != nil && resp.ObjectLockConfiguration.ObjectLockEnabled != "Enabled") {
			s3toObjectLock = append(s3toObjectLock, S3ToObjectLock{*bucket.Name, false})
		} else {
			s3toObjectLock = append(s3toObjectLock, S3ToObjectLock{*bucket.Name, true})
		}
	}

	return s3toObjectLock
}

type S3ToReplicationOtherRegion struct {
	BucketName            string
	ReplicatedOtherRegion bool
	OtherRegion           string
}

func GetS3ToReplicationOtherRegion(s aws.Config, b []types.Bucket) []S3ToReplicationOtherRegion {

	svc := s3.NewFromConfig(s)

	var s3toReplicationOtherRegion []S3ToReplicationOtherRegion
	for _, bucket := range b {
		params := &s3.GetBucketReplicationInput{
			Bucket: aws.String(*bucket.Name),
		}
		resp, err := svc.GetBucketReplication(context.TODO(), params)
		if err != nil {
			var ae smithy.APIError
			if errors.As(err, &ae) && ae.ErrorCode() == "ReplicationConfigurationNotFoundError" {
				s3toReplicationOtherRegion = append(s3toReplicationOtherRegion, S3ToReplicationOtherRegion{*bucket.Name, false, ""})
				continue
			}
			logger.Logger.Error(err.Error())
			// return empty	struct
			return []S3ToReplicationOtherRegion{}
		}
		if resp.ReplicationConfiguration == nil {
			s3toReplicationOtherRegion = append(s3toReplicationOtherRegion, S3ToReplicationOtherRegion{*bucket.Name, false, ""})
		} else {
			// Check the region of destination buckets
			found := false
			for _, rule := range resp.ReplicationConfiguration.Rules {
				replicationTarget := strings.TrimPrefix(*rule.Destination.Bucket, "arn:aws:s3:::")
				if ok, otherRegion := CheckS3Location(s, replicationTarget, s.Region); !ok {
					s3toReplicationOtherRegion = append(s3toReplicationOtherRegion, S3ToReplicationOtherRegion{*bucket.Name, true, otherRegion})
					found = true
					break // break the loop if at least one of the replication rule is to other region
				}
			}

			// no destination rule to other region
			if !found {
				s3toReplicationOtherRegion = append(s3toReplicationOtherRegion, S3ToReplicationOtherRegion{*bucket.Name, false, ""})
			}
		}
	}

	return s3toReplicationOtherRegion
}
