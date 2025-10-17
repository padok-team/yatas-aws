package sqs

import (
	"strconv"

	"github.com/padok-team/yatas/plugins/commons"
)

const MAX_RETENTION_PERIOD = 14 * 24 * 60 * 60 // 14 days in seconds

func CheckIfRetentionPeriodIsValid(checkConfig commons.CheckConfig, sqsQueues []SQSToRetentionPeriod, testName string) {
	var check commons.Check
	check.InitCheck("SQS queues have retention period of 14 days maximum", "Check if SQS queues have a retention period configured and not exceeding 14 days", testName, []string{"Security", "Good Practice"})

	for _, sqsQueue := range sqsQueues {
		if sqsQueue.RetentionPeriod == 0 {
			Message := "SQS queue " + sqsQueue.QueueName + " has no retention period configured"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: sqsQueue.QueueUrl}
			check.AddResult(result)
		} else if sqsQueue.RetentionPeriod > MAX_RETENTION_PERIOD {
			Message := "SQS queue " + sqsQueue.QueueName + " has retention period exceeding 14 days (" + formatRetentionPeriod(sqsQueue.RetentionPeriod) + ")"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: sqsQueue.QueueUrl}
			check.AddResult(result)
		} else {
			Message := "SQS queue " + sqsQueue.QueueName + " has valid retention period (" + formatRetentionPeriod(sqsQueue.RetentionPeriod) + ")"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: sqsQueue.QueueUrl}
			check.AddResult(result)
		}
	}

	checkConfig.Queue <- check
}

func formatRetentionPeriod(seconds int) string {
	days := seconds / (24 * 60 * 60)
	if days == 1 {
		return "1 day"
	}
	return strconv.Itoa(days) + " days"
}
