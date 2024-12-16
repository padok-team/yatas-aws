package iam

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/padok-team/yatas-aws/logger"
)

func GetAllUsers(s aws.Config) []types.User {
	svc := iam.NewFromConfig(s)
	var users []types.User
	input := &iam.ListUsersInput{}
	result, err := svc.ListUsers(context.TODO(), input)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list
		return []types.User{}
	}
	users = append(users, result.Users...)
	for {
		if result.IsTruncated {
			input.Marker = result.Marker
			result, err = svc.ListUsers(context.TODO(), input)
			if err != nil {
				logger.Logger.Error(err.Error())
				// Return an empty list
				return []types.User{}
			}
			users = append(users, result.Users...)
		} else {
			break
		}
	}
	return users
}

func GetAllRoles(s aws.Config) []types.Role {
	svc := iam.NewFromConfig(s)
	var roles []types.Role
	input := &iam.ListRolesInput{}
	result, err := svc.ListRoles(context.TODO(), input)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list
		return []types.Role{}
	}

	roles = append(roles, result.Roles...)
	for {
		if result.IsTruncated {
			input.Marker = result.Marker
			result, err = svc.ListRoles(context.TODO(), input)
			if err != nil {
				logger.Logger.Error(err.Error())
				// Return an empty list
				return []types.Role{}
			}
			roles = append(roles, result.Roles...)
		} else {
			break
		}
	}
	return roles
}

type MFAForUser struct {
	UserName string
	MFAs     []types.MFADevice
}

func GetMfaForUsers(s aws.Config, u []types.User) []MFAForUser {
	svc := iam.NewFromConfig(s)

	var mfaForUsers []MFAForUser
	for _, user := range u {
		input := &iam.ListMFADevicesInput{
			UserName: user.UserName,
		}
		result, err := svc.ListMFADevices(context.TODO(), input)
		if err != nil {
			logger.Logger.Error(err.Error())
			// Return an empty list of mfa devices
			return []MFAForUser{}
		}
		mfaForUsers = append(mfaForUsers, MFAForUser{
			UserName: *user.UserName,
			MFAs:     result.MFADevices,
		})
		for {
			if result.IsTruncated {
				input.Marker = result.Marker
				result, err = svc.ListMFADevices(context.TODO(), input)
				if err != nil {
					logger.Logger.Error(err.Error())
					// Return an empty list of mfa devices
					return []MFAForUser{}
				}
				mfaForUsers = append(mfaForUsers, MFAForUser{
					UserName: *user.UserName,
					MFAs:     result.MFADevices,
				})
			} else {
				break
			}
		}
	}
	return mfaForUsers
}

type AccessKeysForUser struct {
	UserName   string
	AccessKeys []types.AccessKeyMetadata
}

func GetAccessKeysForUsers(s aws.Config, u []types.User) []AccessKeysForUser {
	svc := iam.NewFromConfig(s)

	var accessKeysForUsers []AccessKeysForUser
	for _, user := range u {
		input := &iam.ListAccessKeysInput{
			UserName: user.UserName,
		}
		result, err := svc.ListAccessKeys(context.TODO(), input)
		if err != nil {
			logger.Logger.Error(err.Error())
			// Return an empty list of access keys
			return []AccessKeysForUser{}
		}
		accessKeysForUsers = append(accessKeysForUsers, AccessKeysForUser{
			UserName:   *user.UserName,
			AccessKeys: result.AccessKeyMetadata,
		})
		for {
			if result.IsTruncated {
				input.Marker = result.Marker
				result, err = svc.ListAccessKeys(context.TODO(), input)
				if err != nil {
					logger.Logger.Error(err.Error())
					// Return an empty list of access keys
					return []AccessKeysForUser{}
				}
				accessKeysForUsers = append(accessKeysForUsers, AccessKeysForUser{
					UserName:   *user.UserName,
					AccessKeys: result.AccessKeyMetadata,
				})
			} else {
				break
			}
		}
	}
	return accessKeysForUsers
}

func GetUserPolicies(users []types.User, s aws.Config) []UserPolicies {
	var wgPolicyForUser sync.WaitGroup
	wgPolicyForUser.Add(len(users))
	queue := make(chan UserPolicies, 10)
	for _, user := range users {
		go GetAllPolicyForUser(&wgPolicyForUser, queue, s, user)
	}
	var userPolicies []UserPolicies
	go func() {
		for user := range queue {
			userPolicies = append(userPolicies, user)
			wgPolicyForUser.Done()
		}

	}()
	wgPolicyForUser.Wait()
	return userPolicies
}

type UserToPoliciesElevate struct {
	UserName string
	Policies [][]string
}

func GetUserToPoliciesElevate(userPolicies []UserPolicies) []UserToPoliciesElevate {
	var usersElevatedPolicies []UserToPoliciesElevate
	for _, user := range userPolicies {
		elevation := CheckPolicyForAllowInRequiredPermission(user.Policies, requiredPermissions)
		if elevation != nil {
			usersElevatedPolicies = append(usersElevatedPolicies, UserToPoliciesElevate{
				UserName: user.UserName,
				Policies: elevation,
			})
		}

	}

	return usersElevatedPolicies
}

func GetAllPolicyForUser(wg *sync.WaitGroup, queueCheck chan UserPolicies, s aws.Config, user types.User) {
	var policyList []Policy
	var wgpolicy sync.WaitGroup
	queue := make(chan *string, 100)
	policies := GetPolicyAttachedToUser(s, user)
	wgpolicy.Add(len(policies))
	for _, policy := range policies {
		go GetPolicyDocument(&wgpolicy, queue, s, policy.PolicyArn)

	}
	go func() {
		for t := range queue {
			policyList = append(policyList, JsonDecodePolicyDocument(t))
			wgpolicy.Done()
		}
	}()
	wgpolicy.Wait()
	queueCheck <- UserPolicies{*user.UserName, policyList}
}

func GetPolicyDocument(wg *sync.WaitGroup, queue chan *string, s aws.Config, policyArn *string) {
	policyVersions := GetAllPolicyVersions(s, policyArn)
	SortPolicyVersions(policyVersions)
	input := &iam.GetPolicyVersionInput{
		PolicyArn: policyArn,
		VersionId: policyVersions[0].VersionId,
	}
	svc := iam.NewFromConfig(s)
	result, err := svc.GetPolicyVersion(context.TODO(), input)
	if err != nil {
		logger.Logger.Error(err.Error())
	}
	queue <- result.PolicyVersion.Document
}

func GetPolicyAttachedToUser(s aws.Config, user types.User) []types.AttachedPolicy {
	svc := iam.NewFromConfig(s)
	input := &iam.ListAttachedUserPoliciesInput{
		UserName: user.UserName,
	}
	result, err := svc.ListAttachedUserPolicies(context.TODO(), input)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list of attached policies
		return []types.AttachedPolicy{}
	}
	return result.AttachedPolicies
}

func GetAllPolicyVersions(s aws.Config, policyArn *string) []types.PolicyVersion {
	svc := iam.NewFromConfig(s)
	input := &iam.ListPolicyVersionsInput{
		PolicyArn: policyArn,
	}
	result, err := svc.ListPolicyVersions(context.TODO(), input)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list of policy versions
		return []types.PolicyVersion{}
	}

	return result.Versions
}

func GetRolePolicies(roles []types.Role, s aws.Config) []RolePolicies {
	var wgPolicyForRole sync.WaitGroup
	wgPolicyForRole.Add(len(roles))
	queue := make(chan RolePolicies, 10)
	for _, role := range roles {
		go GetAllPolicyForRole(&wgPolicyForRole, queue, s, role)
	}
	var rolePolicies []RolePolicies
	go func() {
		for role := range queue {
			rolePolicies = append(rolePolicies, role)
			wgPolicyForRole.Done()
		}

	}()
	wgPolicyForRole.Wait()
	return rolePolicies
}

type RoleToPoliciesElevate struct {
	RoleName string
	Policies [][]string
}

func GetRoleToPoliciesElevate(rolePolicies []RolePolicies) []RoleToPoliciesElevate {
	var rolesElevatedPolicies []RoleToPoliciesElevate
	for _, role := range rolePolicies {
		elevation := CheckPolicyForAllowInRequiredPermission(role.Policies, requiredPermissions)
		if elevation != nil {
			rolesElevatedPolicies = append(rolesElevatedPolicies, RoleToPoliciesElevate{
				RoleName: role.RoleName,
				Policies: elevation,
			})
		}

	}

	return rolesElevatedPolicies
}

func GetAllPolicyForRole(wg *sync.WaitGroup, queueCheck chan RolePolicies, s aws.Config, role types.Role) {
	var policyList []Policy
	var wgpolicy sync.WaitGroup
	queue := make(chan *string, 100)
	policies := GetPolicyAttachedToRole(s, role)
	wgpolicy.Add(len(policies))
	for _, policy := range policies {
		go GetPolicyDocument(&wgpolicy, queue, s, policy.PolicyArn)

	}
	go func() {
		for t := range queue {
			policyList = append(policyList, JsonDecodePolicyDocument(t))
			wgpolicy.Done()
		}
	}()
	wgpolicy.Wait()
	queueCheck <- RolePolicies{*role.RoleName, policyList}
}

func GetPolicyAttachedToRole(s aws.Config, role types.Role) []types.AttachedPolicy {
	svc := iam.NewFromConfig(s)
	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: role.RoleName,
	}
	result, err := svc.ListAttachedRolePolicies(context.TODO(), input)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list of attached policies
		return []types.AttachedPolicy{}
	}
	return result.AttachedPolicies
}

func GetPasswordPolicy(s aws.Config) types.PasswordPolicy {
	svc := iam.NewFromConfig(s)

	result, err := svc.GetAccountPasswordPolicy(context.TODO(), &iam.GetAccountPasswordPolicyInput{})
	if err != nil {
		logger.Logger.Error(err.Error())
		return types.PasswordPolicy{}
	}

	return *result.PasswordPolicy
}
