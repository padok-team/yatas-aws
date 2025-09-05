package rds

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func checkIfClusterRestrictedSecurityGroups(checkConfig commons.CheckConfig, clusters []DBClusterWithSGs, testName string) {
	var check commons.Check
	check.InitCheck("Aurora RDS have properly restricted security groups", "Ensure no ingress 0.0.0.0/0 or all ports on SGs for Aurora RDS", testName, []string{"Security", "Best Practice"})

	for _, i := range clusters {
		cluster := i.Cluster
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
			message := "Aurora RDS " + *cluster.DBClusterIdentifier + " has SG opened to 0.0.0.0/0"
			if hasAllPortsOpen {
				message += " and all ports opened"
			}
			check.AddResult(commons.Result{
				Status:     "FAIL",
				Message:    message,
				ResourceID: *cluster.DBClusterArn,
			})
		} else if hasAllPortsOpen {
			check.AddResult(commons.Result{
				Status:     "FAIL",
				Message:    "Aurora RDS " + *cluster.DBClusterIdentifier + " has SG with all ports opened",
				ResourceID: *cluster.DBClusterArn,
			})
		} else {
			check.AddResult(commons.Result{
				Status:     "OK",
				Message:    "Aurora RDS SG is restricted for " + *cluster.DBClusterIdentifier,
				ResourceID: *cluster.DBClusterArn,
			})
		}
	}

	checkConfig.Queue <- check
}
