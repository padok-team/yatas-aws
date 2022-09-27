package eks

import (
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stangirard/yatas/plugins/commons"
	"golang.org/x/exp/slices"
)

func CheckIfEksEndpointPrivate(checkConfig commons.CheckConfig, clusters []types.Cluster, testName string) {
	var check commons.Check
	check.InitCheck("EKS clusters have private endpoint or strict public access", "Check if EKS clusters have private endpoint", testName, []string{"Security", "Good Practice"})
	for _, cluster := range clusters {
		if cluster.ResourcesVpcConfig != nil {
			if cluster.ResourcesVpcConfig.EndpointPublicAccess {
				if ok := slices.Contains(cluster.ResourcesVpcConfig.PublicAccessCidrs, "0.0.0.0/0"); !ok {
					Message := "EKS cluster " + *cluster.Name + " has private endpoint"
					result := commons.Result{Status: "OK", Message: Message, ResourceID: *cluster.Name}
					check.AddResult(result)
				} else {
					Message := "EKS cluster " + *cluster.Name + " has public endpoint"
					result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *cluster.Name}
					check.AddResult(result)
				}
			} else {
				Message := "EKS cluster " + *cluster.Name + " has private endpoint"
				result := commons.Result{Status: "OK", Message: Message, ResourceID: *cluster.Name}
				check.AddResult(result)
			}
		} else {
			Message := "Private endpoint is not enabled for cluster " + *cluster.Name
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *cluster.Name}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
