package ec2

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/padok-team/yatas-aws/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfAuditLogsEnabledOnBastionInstance(checkConfig commons.CheckConfig, ec2ToIAMPolicies []BastionToIAMPolicies, testName string) {
	var check commons.Check
	check.InitCheck("EC2 bastion instance have audit logs enabled", "Check if EC2 bastion instance have audit logs enabled", testName, []string{"Security", "Good Practice"})

	for _, ec2ToIAM := range ec2ToIAMPolicies {
		// Check if the instance is a bastion by checking if it has bastion in its tag name
		if ec2ToIAM.Instance.Tags != nil {
			for _, tag := range ec2ToIAM.Instance.Tags {
				if *tag.Key == "Name" {
					if strings.Contains(*tag.Value, "bastion") {
						// Check if at least one of the policy name has "SSM" in it
						logger.Logger.Debug("Checking " + fmt.Sprint(len(ec2ToIAM.Policies)) + " policies")
						for _, policy := range ec2ToIAM.Policies {
							// URLDecode the policy document
							decodedPolicy, err := url.QueryUnescape(*policy.Document)
							if err != nil {
								logger.Logger.Error(err.Error())
								Message := "Error while decoding the policy document"
								result := commons.Result{Status: "ERROR", Message: Message, ResourceID: *ec2ToIAM.Instance.InstanceId}
								check.AddResult(result)
								break
							}
							logger.Logger.Debug("Decoded policy document: " + decodedPolicy)
							// If ANY policy has "SSM" in it, the check is OK
							if strings.Contains(decodedPolicy, "ssm") {
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
		checkConfig.Queue <- check
	}
}
