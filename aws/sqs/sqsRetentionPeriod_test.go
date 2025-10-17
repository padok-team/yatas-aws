package sqs

import (
	"sync"
	"testing"

	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfRetentionPeriodIsValid(t *testing.T) {
	type args struct {
		checkConfig commons.CheckConfig
		sqsQueues   []SQSToRetentionPeriod
		testName    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Valid retention period - OK case",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				sqsQueues: []SQSToRetentionPeriod{
					{
						QueueName:       "test-queue",
						QueueUrl:        "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
						RetentionPeriod: 86400, // 1 day in seconds
					},
				},
				testName: "AWS_SQS_001",
			},
			want: "OK",
		},
		{
			name: "No retention period configured - FAIL case",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				sqsQueues: []SQSToRetentionPeriod{
					{
						QueueName:       "test-queue-no-retention",
						QueueUrl:        "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue-no-retention",
						RetentionPeriod: 0,
					},
				},
				testName: "AWS_SQS_001",
			},
			want: "FAIL",
		},
		{
			name: "Retention period exceeding 14 days - FAIL case",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				sqsQueues: []SQSToRetentionPeriod{
					{
						QueueName:       "test-queue-long-retention",
						QueueUrl:        "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue-long-retention",
						RetentionPeriod: 8640000, // 100 days in seconds
					},
				},
				testName: "AWS_SQS_001",
			},
			want: "FAIL",
		},
		{
			name: "Retention period exactly 14 days - OK case",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				sqsQueues: []SQSToRetentionPeriod{
					{
						QueueName:       "test-queue-max-retention",
						QueueUrl:        "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue-max-retention",
						RetentionPeriod: 14 * 24 * 60 * 60, // 14 days in seconds
					},
				},
				testName: "AWS_SQS_001",
			},
			want: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfRetentionPeriodIsValid(tt.args.checkConfig, tt.args.sqsQueues, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != tt.want {
						t.Errorf("CheckIfRetentionPeriodIsValid() = %v, want %v", check.Status, tt.want)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfRetentionPeriodIsValidFail(t *testing.T) {
	type args struct {
		checkConfig commons.CheckConfig
		sqsQueues   []SQSToRetentionPeriod
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Multiple queues with issues",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				sqsQueues: []SQSToRetentionPeriod{
					{
						QueueName:       "test-queue-no-retention",
						QueueUrl:        "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue-no-retention",
						RetentionPeriod: 0,
					},
					{
						QueueName:       "test-queue-too-long",
						QueueUrl:        "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue-too-long",
						RetentionPeriod: 10368000, // 120 days in seconds
					},
				},
				testName: "AWS_SQS_001",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfRetentionPeriodIsValid(tt.args.checkConfig, tt.args.sqsQueues, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfRetentionPeriodIsValid() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
