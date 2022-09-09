package vpc

import (
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/plugins/commons"
)

func checkCIDR20(checkConfig commons.CheckConfig, vpcs []types.Vpc, testName string) {
	var check commons.Check
	check.InitCheck("VPC CIDRs are bigger than /20", "Check if VPC CIDR is /20 or bigger", testName)
	for _, vpc := range vpcs {
		cidr := *vpc.CidrBlock
		// split the cidr to / and get the last part as an int
		cidrInt, _ := strconv.Atoi(strings.Split(cidr, "/")[1])
		if cidrInt > 20 {
			Message := "VPC CIDR is not /20 or bigger on " + *vpc.VpcId
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *vpc.VpcId}
			check.AddResult(result)
		} else {
			Message := "VPC CIDR is /20 or bigger on " + *vpc.VpcId
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *vpc.VpcId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
