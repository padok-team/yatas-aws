package s3

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func GetListS3(s aws.Config) []types.Bucket {

	svc := s3.NewFromConfig(s)

	params := &s3.ListBucketsInput{}
	resp, err := svc.ListBuckets(context.TODO(), params)
	if err != nil {
		fmt.Println(err)
	}

	return resp.Buckets
}

func GetListS3NotInRegion(s aws.Config, region string) []types.Bucket {

	svc := s3.NewFromConfig(s)

	params := &s3.ListBucketsInput{}
	resp, err := svc.ListBuckets(context.TODO(), params)
	if err != nil {
		fmt.Println(err)
	}

	var buckets []types.Bucket
	for _, bucket := range resp.Buckets {
		if !CheckS3Location(s, *bucket.Name, region) {
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
			if resp.PublicAccessBlockConfiguration != nil && resp.PublicAccessBlockConfiguration.BlockPublicAcls {
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
			fmt.Println(err)
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
			fmt.Println(err)
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
