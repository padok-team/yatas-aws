package ssm

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/padok-team/yatas-aws/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

type PolicyDocument struct {
	Version   string
	Statement []struct {
		Sid      string
		Effect   string
		Action   interface{} // can be string or []string
		Resource interface{} // can be string or []string
	}
}

// CheckIfAuditLogsEnabledOnBastionInstance checks if the bastion instance has audit logs enabled by checking if it has the correct permissions and resources
func CheckIfAuditLogsEnabledOnBastionInstance(checkConfig commons.CheckConfig, ec2ToIAMPolicies []BastionToIAMPolicies, testName string) {
	var check commons.Check
	check.InitCheck("EC2 bastion instance have audit logs enabled", "Check if EC2 bastion instance have audit logs enabled (ec2 name tag contains bastion* with a role policy containing s3:PutObject and s3:PutObjectAcl on a bucket with an arn like arn:aws:s3:::ssm-logging*)", testName, []string{"Security", "Good Practice"})

	for _, ec2ToIAM := range ec2ToIAMPolicies {
		// Check if the instance is a bastion by checking if it has bastion in its tag name
		if ec2ToIAM.Instance.Tags != nil {
			for _, tag := range ec2ToIAM.Instance.Tags {
				if *tag.Key == "Name" {
					if strings.Contains(*tag.Value, "bastion") {
						for _, policy := range ec2ToIAM.Policies {
							// URLDecode the policy document
							decodedPolicy, err := url.QueryUnescape(policy)
							if err != nil {
								logger.Logger.Error(err.Error())
								Message := "Error while decoding the policy document"
								result := commons.Result{Status: "ERROR", Message: Message, ResourceID: *ec2ToIAM.Instance.InstanceId}
								check.AddResult(result)
								break
							}
							var policyDoc PolicyDocument
							if err := json.Unmarshal([]byte(decodedPolicy), &policyDoc); err != nil {
								logger.Logger.Error(err.Error())
								Message := "Error while parsing the policy document"
								result := commons.Result{Status: "ERROR", Message: Message, ResourceID: *ec2ToIAM.Instance.InstanceId}
								check.AddResult(result)
								break
							}

							hasRequiredPermissions := false
							for _, statement := range policyDoc.Statement {
								if statement.Effect == "Allow" {
									hasPutObjectPermission := false
									hasPutObjectAclPermission := false
									hasSSMLoggingBucketArn := false

									// Handle Action being either string or []string
									var actions []string
									switch v := statement.Action.(type) {
									case string:
										actions = []string{v}
									case []interface{}:
										actions = make([]string, len(v))
										for i, a := range v {
											actions[i] = a.(string)
										}
									}

									// Handle Resource being either string or []string
									var resources []string
									switch v := statement.Resource.(type) {
									case string:
										resources = []string{v}
									case []interface{}:
										resources = make([]string, len(v))
										for i, r := range v {
											resources[i] = r.(string)
										}
									}

									for _, action := range actions {
										if action == "s3:PutObject" {
											hasPutObjectPermission = true
										}
										if action == "s3:PutObjectAcl" {
											hasPutObjectAclPermission = true
										}
									}

									for _, resource := range resources {
										if strings.Contains(resource, "arn:aws:s3:::ssm-logging") {
											hasSSMLoggingBucketArn = true
										}
									}

									if hasPutObjectPermission && hasPutObjectAclPermission && hasSSMLoggingBucketArn {
										hasRequiredPermissions = true
										break
									}
								}
							}

							if hasRequiredPermissions {
								Message := "EC2 instance " + *ec2ToIAM.Instance.InstanceId + " has audit logs enabled"
								result := commons.Result{Status: "OK", Message: Message, ResourceID: *ec2ToIAM.Instance.InstanceId}
								check.AddResult(result)
								break
							}
						}
					}
				}
			}
		}
		if len(check.Results) == 0 {
			Message := "EC2 instance " + *ec2ToIAM.Instance.InstanceId + " has no audit logs enabled"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *ec2ToIAM.Instance.InstanceId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
