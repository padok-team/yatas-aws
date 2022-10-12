package guardduty

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
)

func GetDetectors(s aws.Config) []string {
	svc := guardduty.NewFromConfig(s)
	input := &guardduty.ListDetectorsInput{}
	result, err := svc.ListDetectors(context.TODO(), input)
	if err != nil {
		fmt.Println(err)
	}
	return result.DetectorIds
}
