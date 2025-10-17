package sqs

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/padok-team/yatas-aws/logger"
)

type SQSToRetentionPeriod struct {
	QueueName       string
	QueueUrl        string
	RetentionPeriod int
}

func GetSQSQueues(s aws.Config) []string {
	svc := sqs.NewFromConfig(s)

	params := &sqs.ListQueuesInput{}
	resp, err := svc.ListQueues(context.TODO(), params)
	if err != nil {
		logger.Logger.Error(err.Error())
		return []string{}
	}

	return resp.QueueUrls
}

func GetSQSRetentionPeriod(s aws.Config, queueUrls []string) []SQSToRetentionPeriod {
	svc := sqs.NewFromConfig(s)
	var sqsRetentionPeriods []SQSToRetentionPeriod

	for _, queueUrl := range queueUrls {
		params := &sqs.GetQueueAttributesInput{
			QueueUrl: aws.String(queueUrl),
			AttributeNames: []types.QueueAttributeName{
				types.QueueAttributeNameMessageRetentionPeriod,
			},
		}

		resp, err := svc.GetQueueAttributes(context.TODO(), params)
		if err != nil {
			logger.Logger.Error(err.Error())
			continue
		}

		retentionPeriodStr, exists := resp.Attributes[string(types.QueueAttributeNameMessageRetentionPeriod)]
		if !exists {
			logger.Logger.Warn("MessageRetentionPeriod not found for queue: " + queueUrl)
			continue
		}

		retentionPeriod, err := strconv.Atoi(retentionPeriodStr)
		if err != nil {
			logger.Logger.Error("Failed to parse MessageRetentionPeriod for queue " + queueUrl + ": " + err.Error())
			continue
		}

		// Extract queue name from URL (last part after the last slash)
		queueName := queueUrl
		if lastSlash := len(queueUrl) - 1; lastSlash >= 0 {
			for i := lastSlash; i >= 0; i-- {
				if queueUrl[i] == '/' {
					queueName = queueUrl[i+1:]
					break
				}
			}
		}

		sqsRetentionPeriods = append(sqsRetentionPeriods, SQSToRetentionPeriod{
			QueueName:       queueName,
			QueueUrl:        queueUrl,
			RetentionPeriod: retentionPeriod,
		})
	}

	return sqsRetentionPeriods
}
