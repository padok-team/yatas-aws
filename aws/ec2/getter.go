package ec2

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/padok-team/yatas-aws/logger"
)

type EC2GetObjectAPI interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

type IAMGetObjectAPI interface {
	GetRole(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error)
	GetInstanceProfile(ctx context.Context, params *iam.GetInstanceProfileInput, optFns ...func(*iam.Options)) (*iam.GetInstanceProfileOutput, error)
	ListAttachedRolePolicies(ctx context.Context, params *iam.ListAttachedRolePoliciesInput, optFns ...func(*iam.Options)) (*iam.ListAttachedRolePoliciesOutput, error)
	ListPolicyVersions(ctx context.Context, params *iam.ListPolicyVersionsInput, optFns ...func(*iam.Options)) (*iam.ListPolicyVersionsOutput, error)
	GetPolicyVersion(ctx context.Context, params *iam.GetPolicyVersionInput, optFns ...func(*iam.Options)) (*iam.GetPolicyVersionOutput, error)
}

func GetEC2s(svc EC2GetObjectAPI) []ec2Types.Instance {
	input := &ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(context.TODO(), input)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list
		return []ec2Types.Instance{}
	}
	var instances []ec2Types.Instance
	for _, r := range result.Reservations {
		instances = append(instances, r.Instances...)
	}
	for {
		if result.NextToken == nil {
			break
		}
		input.NextToken = result.NextToken
		result, err = svc.DescribeInstances(context.TODO(), input)
		if err != nil {
			logger.Logger.Error(err.Error())
			// Return an empty list of instances
			return []ec2Types.Instance{}
		}
		for _, r := range result.Reservations {
			instances = append(instances, r.Instances...)
		}
	}

	return instances
}

type BastionToIAMPolicies struct {
	Instance ec2Types.Instance
	Policies []iamTypes.PolicyVersion
}

func GetBastionToIAMPolicies(ec2Svc EC2GetObjectAPI, iamSvc IAMGetObjectAPI, instances []ec2Types.Instance) []BastionToIAMPolicies {
	var ec2ToIAMPolicies []BastionToIAMPolicies
	for _, instance := range instances {
		if instance.Tags != nil {
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					if strings.Contains(*tag.Value, "bastion") {
						if instance.IamInstanceProfile != nil {
							if instance.IamInstanceProfile.Arn != nil {
								// Get the IAM role from the instance profile
								instanceProfileName := *instance.IamInstanceProfile.Arn
								getInstanceProfileOutput, err := iamSvc.GetInstanceProfile(context.TODO(), &iam.GetInstanceProfileInput{InstanceProfileName: aws.String(strings.Split(instanceProfileName, "/")[1])})
								if err != nil {
									logger.Logger.Error(err.Error())
									continue
								}
								// Get the policies attached to the role
								role, err := iamSvc.GetRole(context.TODO(), &iam.GetRoleInput{RoleName: aws.String(*getInstanceProfileOutput.InstanceProfile.Roles[0].RoleName)})
								if err != nil {
									logger.Logger.Error(err.Error())
									continue
								}
								logger.Logger.Debug(fmt.Sprintf("Role: %s", *role.Role.RoleName))
								// Get the policies attached to the role
								attachedRolePolicies, err := iamSvc.ListAttachedRolePolicies(context.TODO(), &iam.ListAttachedRolePoliciesInput{RoleName: role.Role.RoleName})
								if err != nil {
									logger.Logger.Error(err.Error())
									continue
								}
								// Get policy details for each policy attached to the role
								var policies []iamTypes.PolicyVersion
								logger.Logger.Debug(fmt.Sprintf("Number of policies attached to the role: %d", len(attachedRolePolicies.AttachedPolicies)))
								for _, policy := range attachedRolePolicies.AttachedPolicies {
									// List all versions of the policy
									policyVersions, err := iamSvc.ListPolicyVersions(context.TODO(), &iam.ListPolicyVersionsInput{PolicyArn: policy.PolicyArn})
									if err != nil {
										logger.Logger.Error(err.Error())
										continue
									}
									// Get the latest version of the policy
									latestPolicyVersion := policyVersions.Versions[len(policyVersions.Versions)-1]
									policyDetail, err := iamSvc.GetPolicyVersion(context.TODO(), &iam.GetPolicyVersionInput{PolicyArn: policy.PolicyArn, VersionId: latestPolicyVersion.VersionId})

									if err != nil {
										logger.Logger.Error(err.Error())
										continue
									}
									policies = append(policies, *policyDetail.PolicyVersion)
								}
								ec2ToIAMPolicies = append(ec2ToIAMPolicies, BastionToIAMPolicies{Instance: instance, Policies: policies})
							}
						}
					}
				}
			}
		}
	}
	return ec2ToIAMPolicies
}
