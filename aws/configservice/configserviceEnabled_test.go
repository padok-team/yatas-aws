package configservice

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/configservice/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfConfigServiceIsEnabled(t *testing.T) {
	t.Run("AWS Config is not enabled (empty recorder status)", func(t *testing.T) {
		queue := make(chan commons.Check, 1)
		checkConfig := commons.CheckConfig{Queue: queue}
		configurationRecorderStatus := []types.ConfigurationRecorderStatus{}

		CheckIfConfigServiceIsEnabled(checkConfig, "TestAWSConfigDisabled", configurationRecorderStatus)

		check := <-queue
		if len(check.Results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(check.Results))
		}
		result := check.Results[0]
		if result.Status != "FAIL" {
			t.Errorf("Expected status FAIL, got %s", result.Status)
		}
		if result.Message != "AWS Config is not enabled" {
			t.Errorf("Expected message 'AWS Config is not enabled', got '%s'", result.Message)
		}
	})

	t.Run("AWS Config is enabled (Recording is true)", func(t *testing.T) {
		queue := make(chan commons.Check, 1)
		checkConfig := commons.CheckConfig{Queue: queue}
		configurationRecorderStatus := []types.ConfigurationRecorderStatus{
			{Recording: true},
		}

		CheckIfConfigServiceIsEnabled(checkConfig, "TestAWSConfigEnabled", configurationRecorderStatus)

		check := <-queue
		if len(check.Results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(check.Results))
		}
		result := check.Results[0]
		if result.Status != "OK" {
			t.Errorf("Expected status OK, got %s", result.Status)
		}
		if result.Message != "AWS Config is enabled" {
			t.Errorf("Expected message 'AWS Config is enabled', got '%s'", result.Message)
		}
	})

	t.Run("AWS Config is not enabled (Recording is false)", func(t *testing.T) {
		queue := make(chan commons.Check, 1)
		checkConfig := commons.CheckConfig{Queue: queue}
		configurationRecorderStatus := []types.ConfigurationRecorderStatus{
			{Recording: false},
		}

		CheckIfConfigServiceIsEnabled(checkConfig, "TestAWSConfigDisabledFalse", configurationRecorderStatus)

		check := <-queue
		if len(check.Results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(check.Results))
		}
		result := check.Results[0]
		if result.Status != "FAIL" {
			t.Errorf("Expected status FAIL, got %s", result.Status)
		}
		if result.Message != "AWS Config is not enabled" {
			t.Errorf("Expected message 'AWS Config is not enabled', got '%s'", result.Message)
		}
	})

	t.Run("AWS Config is enabled with at least one Recording as true", func(t *testing.T) {
		queue := make(chan commons.Check, 1)
		checkConfig := commons.CheckConfig{Queue: queue}
		configurationRecorderStatus := []types.ConfigurationRecorderStatus{
			{Recording: true},
			{Recording: false},
		}

		CheckIfConfigServiceIsEnabled(checkConfig, "TestAWSConfigEnabled", configurationRecorderStatus)

		check := <-queue
		if len(check.Results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(check.Results))
		}
		result := check.Results[0]
		if result.Status != "OK" {
			t.Errorf("Expected status OK, got %s", result.Status)
		}
		if result.Message != "AWS Config is enabled" {
			t.Errorf("Expected message 'AWS Config is enabled', got '%s'", result.Message)
		}
	})
}
