package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/types"
)

func GetCognitoPools(s aws.Config) []types.IdentityPoolShortDescription {
	svc := cognitoidentity.NewFromConfig(s)
	cognitoInput := &cognitoidentity.ListIdentityPoolsInput{
		MaxResults: 50,
	}
	result, err := svc.ListIdentityPools(context.TODO(), cognitoInput)
	if err != nil {
		panic(err)
	}
	fmt.Println("Hello")
	return result.IdentityPools
}

func GetDetailedCognitoPool(s aws.Config, pools []types.IdentityPoolShortDescription) []cognitoidentity.DescribeIdentityPoolOutput {
	svc := cognitoidentity.NewFromConfig(s)
	var detailedPools []cognitoidentity.DescribeIdentityPoolOutput
	for _, pool := range pools {
		cognitoInput := &cognitoidentity.DescribeIdentityPoolInput{
			IdentityPoolId: pool.IdentityPoolId,
		}
		result, err := svc.DescribeIdentityPool(context.TODO(), cognitoInput)
		if err != nil {
			panic(err)
		}
		detailedPools = append(detailedPools, *result)
	}
	return detailedPools
}
