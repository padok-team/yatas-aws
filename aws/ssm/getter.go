package ssm

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/padok-team/yatas-aws/aws/ec2"
	"github.com/padok-team/yatas-aws/logger"
)

type SSMGetObjectAPI interface {
}

type IAMGetObjectAPI interface {
	GetRole(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error)
	GetRolePolicy(ctx context.Context, params *iam.GetRolePolicyInput, optFns ...func(*iam.Options)) (*iam.GetRolePolicyOutput, error)
	GetInstanceProfile(ctx context.Context, params *iam.GetInstanceProfileInput, optFns ...func(*iam.Options)) (*iam.GetInstanceProfileOutput, error)
	ListRolePolicies(ctx context.Context, params *iam.ListRolePoliciesInput, optFns ...func(*iam.Options)) (*iam.ListRolePoliciesOutput, error)
	ListPolicyVersions(ctx context.Context, params *iam.ListPolicyVersionsInput, optFns ...func(*iam.Options)) (*iam.ListPolicyVersionsOutput, error)
	GetPolicyVersion(ctx context.Context, params *iam.GetPolicyVersionInput, optFns ...func(*iam.Options)) (*iam.GetPolicyVersionOutput, error)
}

type BastionToIAMPolicies struct {
	Instance ec2Types.Instance
	Policies []string
}

func GetBastionToIAMPolicies(ec2Client ec2.EC2GetObjectAPI, iamClient IAMGetObjectAPI, instances []ec2Types.Instance) []BastionToIAMPolicies {
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
								getInstanceProfileOutput, err := iamClient.GetInstanceProfile(context.TODO(), &iam.GetInstanceProfileInput{InstanceProfileName: aws.String(strings.Split(instanceProfileName, "/")[1])})
								if err != nil {
									logger.Logger.Error(err.Error())
									continue
								}
								// Get the policies attached to the role
								role, err := iamClient.GetRole(context.TODO(), &iam.GetRoleInput{RoleName: aws.String(*getInstanceProfileOutput.InstanceProfile.Roles[0].RoleName)})
								if err != nil {
									logger.Logger.Error(err.Error())
									continue
								}
								logger.Logger.Debug(fmt.Sprintf("Role: %s", *role.Role.RoleName))
								// Get the policies attached to the role
								attachedRolePolicies, err := iamClient.ListRolePolicies(context.TODO(), &iam.ListRolePoliciesInput{RoleName: role.Role.RoleName})
								if err != nil {
									logger.Logger.Error(err.Error())
									continue
								}
								// Get policy details for each policy attached to the role
								var policies []string
								logger.Logger.Debug(fmt.Sprintf("Number of policies attached to the role: %d", len(attachedRolePolicies.PolicyNames)))
								for _, policyName := range attachedRolePolicies.PolicyNames {
									// List all versions of the policy
									policy, err := iamClient.GetRolePolicy(context.TODO(), &iam.GetRolePolicyInput{PolicyName: &policyName, RoleName: role.Role.RoleName})
									if err != nil {
										logger.Logger.Error(err.Error())
										continue
									}
									if err != nil {
										logger.Logger.Error(err.Error())
										continue
									}
									policies = append(policies, *policy.PolicyDocument)
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
