package vpc

import (
	"strconv"

	"github.com/stangirard/yatas/plugins/commons"
)

func CheckInboudPort(checkConfig commons.CheckConfig, sg2SGRs []SGToSecurityGroupRules, testName string) {
	var check commons.Check
	check.InitCheck("SG with inbound 0.0.0.0/0 must allow only 80 and 443 ports", "Check if SG with 0.0.0.0/0 in inbound allow only 80 or 443", testName, []string{"Security", "Good Practice"})
	for _, sg2sgr := range sg2SGRs {
		error := 0
		for _, sgrule := range sg2sgr.securityGroupRules {
			if *sgrule.IsEgress == false && sgrule.CidrIpv4 != nil && *sgrule.CidrIpv4 == "0.0.0.0/0" {
				if (*sgrule.ToPort != 80 || *sgrule.FromPort != 80) && (*sgrule.ToPort != 443 || *sgrule.FromPort != 443) {
					error++
					port1 := strconv.Itoa(int(*sgrule.FromPort))
					port2 := strconv.Itoa(int(*sgrule.ToPort))
					Message := "SecurityGroup " + *sg2sgr.securityGroup.GroupId + " has forbidden ports range" + port1 + "-" + port2 + " in inbound 0.0.0.0/0" + *sgrule.SecurityGroupRuleId
					result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *sg2sgr.securityGroup.GroupId}
					check.AddResult(result)
				}
			}
		}
		if error == 0 {
			Message := "SecurityGroup " + *sg2sgr.securityGroup.GroupId + " has no forbidden port in inbound 0.0.0.0/0"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *sg2sgr.securityGroup.GroupId}
			check.AddResult(result)
		}

	}
	checkConfig.Queue <- check
}
