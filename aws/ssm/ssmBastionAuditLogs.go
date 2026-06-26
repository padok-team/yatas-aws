package ssm

import (
	"encoding/json"
	"net/url"
	"slices"
	"strings"

	"github.com/padok-team/yatas-aws/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

type PolicyDocument struct {
	Version   string
	Statement []struct {
		Sid      string
		Effect   string
		Action   any // can be string or []string
		Resource any // can be string or []string
	}
}

func parseStringOrSlice(v any) []string {
	switch val := v.(type) {
	case string:
		return []string{val}
	case []any:
		result := make([]string, len(val))
		for i, item := range val {
			result[i] = item.(string)
		}
		return result
	}
	return nil
}

func isS3Resource(resource string) bool {
	return resource == "*" || strings.HasPrefix(resource, "arn:aws:s3:::")
}

func hasPutObjectOnS3(policyDoc PolicyDocument) bool {
	for _, statement := range policyDoc.Statement {
		if statement.Effect != "Allow" {
			continue
		}
		actions := parseStringOrSlice(statement.Action)
		resources := parseStringOrSlice(statement.Resource)

		if !slices.ContainsFunc(actions, func(a string) bool {
			return a == "s3:PutObject" || a == "s3:*" || a == "*"
		}) {
			continue
		}

		if slices.ContainsFunc(resources, isS3Resource) {
			return true
		}
	}
	return false
}

// CheckIfAuditLogsEnabledOnBastionInstance checks if the bastion instance has audit logs enabled by checking if it has the correct permissions and resources
func CheckIfAuditLogsEnabledOnBastionInstance(checkConfig commons.CheckConfig, ec2ToIAMPolicies []BastionToIAMPolicies, testName string) {
	var check commons.Check
	check.InitCheck("EC2 bastion instance have audit logs enabled", "Check if EC2 bastion instance have audit logs enabled (ec2 name tag contains bastion* with a role policy containing s3:PutObject on any S3 bucket)", testName, []string{"Security", "Good Practice", "HDS"})

	for _, ec2ToIAM := range ec2ToIAMPolicies {
		isBastion := false
		if ec2ToIAM.Instance.Tags != nil {
			for _, tag := range ec2ToIAM.Instance.Tags {
				if *tag.Key == "Name" && strings.Contains(strings.ToLower(*tag.Value), "bastion") {
					isBastion = true
					break
				}
			}
		}

		if !isBastion {
			continue
		}

		instanceHasPermissions := false
		for _, policy := range ec2ToIAM.Policies {
			decodedPolicy, err := url.QueryUnescape(policy)
			if err != nil {
				logger.Logger.Error(err.Error())
				result := commons.Result{Status: "ERROR", Message: "Error while decoding the policy document", ResourceID: *ec2ToIAM.Instance.InstanceId}
				check.AddResult(result)
				break
			}

			var policyDoc PolicyDocument
			if err := json.Unmarshal([]byte(decodedPolicy), &policyDoc); err != nil {
				logger.Logger.Error(err.Error())
				result := commons.Result{Status: "ERROR", Message: "Error while parsing the policy document", ResourceID: *ec2ToIAM.Instance.InstanceId}
				check.AddResult(result)
				break
			}

			if hasPutObjectOnS3(policyDoc) {
				instanceHasPermissions = true
				break
			}
		}

		if instanceHasPermissions {
			result := commons.Result{Status: "OK", Message: "EC2 instance " + *ec2ToIAM.Instance.InstanceId + " has audit logs enabled", ResourceID: *ec2ToIAM.Instance.InstanceId}
			check.AddResult(result)
		} else {
			result := commons.Result{Status: "FAIL", Message: "EC2 instance " + *ec2ToIAM.Instance.InstanceId + " has no audit logs enabled", ResourceID: *ec2ToIAM.Instance.InstanceId}
			check.AddResult(result)
		}
	}

	checkConfig.Queue <- check
}
