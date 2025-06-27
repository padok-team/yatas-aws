package rds

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func checkIfRDSRestrictedSecurityGroups(checkConfig commons.CheckConfig, instances []InstanceWithSGs, testName string) {
	var check commons.Check
	check.InitCheck("RDS have properly restricted security groups", "Ensure no ingress 0.0.0.0/0 or all ports on SGs for RDS", testName, []string{"Security", "Best Practice"})

	for _, i := range instances {
		instance := i.Instance
		sgrs := i.SecurityGroups

		hasPermissiveIngress := false
		hasAllPortsOpen := false

		for _, sg := range sgrs {
			for _, perm := range sg.IpPermissions {
				for _, ipRange := range perm.IpRanges {
					if ipRange.CidrIp != nil && *ipRange.CidrIp == "0.0.0.0/0" {
						hasPermissiveIngress = true
					}
					if (perm.FromPort == nil || *perm.FromPort == 0) && (perm.ToPort == nil || *perm.ToPort >= 65535 || *perm.ToPort == -1) {
						hasAllPortsOpen = true
					}
				}
			}
		}

		if hasPermissiveIngress {
			message := "RDS " + *instance.DBInstanceIdentifier + " has SG opened to 0.0.0.0/0"
			if hasAllPortsOpen {
				message += " and all ports opened"
			}
			check.AddResult(commons.Result{
				Status:     "FAIL",
				Message:    message,
				ResourceID: *instance.DBInstanceArn,
			})
		} else if hasAllPortsOpen {
			check.AddResult(commons.Result{
				Status:     "FAIL",
				Message:    "RDS " + *instance.DBInstanceIdentifier + " has SG with all ports opened",
				ResourceID: *instance.DBInstanceArn,
			})
		} else {
			check.AddResult(commons.Result{
				Status:     "OK",
				Message:    "RDS SG is restricted for " + *instance.DBInstanceIdentifier,
				ResourceID: *instance.DBInstanceArn,
			})
		}
	}

	checkConfig.Queue <- check
}
