package eks

import (
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stangirard/yatas/config"
	"golang.org/x/exp/slices"
)

func CheckIfEksEndpointPrivate(checkConfig config.CheckConfig, clusters []types.Cluster, testName string) {
	var check config.Check
	check.InitCheck("EKS clusters have private endpoint or strict public access", "Check if EKS clusters have private endpoint", testName)
	for _, cluster := range clusters {
		if cluster.ResourcesVpcConfig != nil {
			if cluster.ResourcesVpcConfig.EndpointPublicAccess {
				if ok := slices.Contains(cluster.ResourcesVpcConfig.PublicAccessCidrs, "0.0.0.0/0"); !ok {
					Message := "EKS cluster " + *cluster.Name + " has private endpoint"
					result := config.Result{Status: "OK", Message: Message, ResourceID: *cluster.Name}
					check.AddResult(result)
				} else {
					Message := "EKS cluster " + *cluster.Name + " has public endpoint"
					result := config.Result{Status: "FAIL", Message: Message, ResourceID: *cluster.Name}
					check.AddResult(result)
				}
			} else {
				Message := "EKS cluster " + *cluster.Name + " has private endpoint"
				result := config.Result{Status: "OK", Message: Message, ResourceID: *cluster.Name}
				check.AddResult(result)
			}
		} else {
			Message := "Private endpoint is not enabled for cluster " + *cluster.Name
			result := config.Result{Status: "FAIL", Message: Message, ResourceID: *cluster.Name}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
