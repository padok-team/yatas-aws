package cloudtrail

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfCloudtrailIsEnabled(t *testing.T) {
	t.Run("CloudTrail is not enabled", func(t *testing.T) {
		queue := make(chan commons.Check, 1)
		checkConfig := commons.CheckConfig{Queue: queue}
		trailStatus := []cloudtrail.GetTrailStatusOutput{}

		CheckIfCloudtrailIsEnabled(checkConfig, trailStatus, "TestCloudTrailDisabled")

		check := <-queue
		if len(check.Results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(check.Results))
		}
		result := check.Results[0]
		if result.Status != "FAIL" {
			t.Errorf("Expected status FAIL, got %s", result.Status)
		}
		if result.Message != "Cloudtrail is not enabled" {
			t.Errorf("Expected message 'Cloudtrail is not enabled', got '%s'", result.Message)
		}
	})

	t.Run("CloudTrail is enabled", func(t *testing.T) {
		queue := make(chan commons.Check, 1)
		checkConfig := commons.CheckConfig{Queue: queue}
		trailStatus := []cloudtrail.GetTrailStatusOutput{
			{IsLogging: aws.Bool(true)},
		}

		CheckIfCloudtrailIsEnabled(checkConfig, trailStatus, "TestCloudTrailEnabled")

		check := <-queue
		if len(check.Results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(check.Results))
		}
		result := check.Results[0]
		if result.Status != "OK" {
			t.Errorf("Expected status OK, got %s", result.Status)
		}
		if result.Message != "Cloudtrail is enabled" {
			t.Errorf("Expected message 'Cloudtrail is enabled', got '%s'", result.Message)
		}
	})

	t.Run("CloudTrail is disabled (IsLogging is false)", func(t *testing.T) {
		queue := make(chan commons.Check, 1)
		checkConfig := commons.CheckConfig{Queue: queue}
		trailStatus := []cloudtrail.GetTrailStatusOutput{
			{IsLogging: aws.Bool(false)},
		}

		CheckIfCloudtrailIsEnabled(checkConfig, trailStatus, "TestCloudTrailDisabledFalse")

		check := <-queue
		if len(check.Results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(check.Results))
		}
		result := check.Results[0]
		if result.Status != "FAIL" {
			t.Errorf("Expected status FAIL, got %s", result.Status)
		}
		if result.Message != "Cloudtrail is not enabled" {
			t.Errorf("Expected message 'Cloudtrail is not enabled', got '%s'", result.Message)
		}
	})

	t.Run("CloudTrail is enabled with multiple trail status", func(t *testing.T) {
		queue := make(chan commons.Check, 1)
		checkConfig := commons.CheckConfig{Queue: queue}
		trailStatus := []cloudtrail.GetTrailStatusOutput{
			{IsLogging: aws.Bool(true)},
			{IsLogging: aws.Bool(false)},
		}

		CheckIfCloudtrailIsEnabled(checkConfig, trailStatus, "TestCloudTrailEnabled")

		check := <-queue
		if len(check.Results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(check.Results))
		}
		result := check.Results[0]
		if result.Status != "OK" {
			t.Errorf("Expected status OK, got %s", result.Status)
		}
		if result.Message != "Cloudtrail is enabled" {
			t.Errorf("Expected message 'Cloudtrail is enabled', got '%s'", result.Message)
		}
	})
}
