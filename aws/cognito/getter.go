package cognito

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/types"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	ciptypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/padok-team/yatas-aws/logger"
)

func GetCognitoPools(s aws.Config) []types.IdentityPoolShortDescription {
	svc := cognitoidentity.NewFromConfig(s)
	cognitoInput := &cognitoidentity.ListIdentityPoolsInput{
		MaxResults: aws.Int32(50),
	}
	result, err := svc.ListIdentityPools(context.TODO(), cognitoInput)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list of certificates
		return []types.IdentityPoolShortDescription{}
	}
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
			logger.Logger.Error(err.Error())
			// Return an empty list of certificates
			return []cognitoidentity.DescribeIdentityPoolOutput{}
		}
		detailedPools = append(detailedPools, *result)
	}
	return detailedPools
}

func GetCognitoUserPools(s aws.Config) []ciptypes.UserPoolDescriptionType {
	svc := cognitoidentityprovider.NewFromConfig(s)
	logger.Logger.Info("Getting Cognito User Pools")
	cognitoInput := &cognitoidentityprovider.ListUserPoolsInput{
		MaxResults: aws.Int32(50),
	}
	result, err := svc.ListUserPools(context.TODO(), cognitoInput)
	if err != nil {
		logger.Logger.Error(err.Error())
		// Return an empty list of certificates
		return []ciptypes.UserPoolDescriptionType{}
	}
	return result.UserPools
}

func GetDetailedCognitoUserPool(s aws.Config, userPools []ciptypes.UserPoolDescriptionType) []cognitoidentityprovider.DescribeUserPoolOutput {
	svc := cognitoidentityprovider.NewFromConfig(s)
	var detailedUserPools []cognitoidentityprovider.DescribeUserPoolOutput
	for _, userPool := range userPools {
		cognitoInput := &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: userPool.Id,
		}
		result, err := svc.DescribeUserPool(context.TODO(), cognitoInput)
		if err != nil {
			logger.Logger.Error(err.Error())
			// Return an empty list of certificates
			continue
		}
		detailedUserPools = append(detailedUserPools, *result)
	}
	return detailedUserPools
}
