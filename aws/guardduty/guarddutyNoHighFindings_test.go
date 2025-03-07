package guardduty_test

import (
	"testing"

	"github.com/padok-team/yatas-aws/aws/guardduty"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfGuarddutyNoHighFindings_NoFindings(t *testing.T) {
	checkConfig := commons.CheckConfig{
		Queue: make(chan commons.Check, 1),
	}
	testName := "TestCheckIfGuarddutyNoHighFindings_NoFindings"
	findings := []string{}

	go guardduty.CheckIfGuarddutyNoHighFindings(checkConfig, testName, findings)
	check := <-checkConfig.Queue

	if len(check.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(check.Results))
	}

	if check.Results[0].Status != "OK" {
		t.Errorf("Expected status OK, got %s", check.Results[0].Status)
	}

	expectedMessage := "GuardDuty has 0 HIGH severity findings"
	if check.Results[0].Message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, check.Results[0].Message)
	}
}

func TestCheckIfGuarddutyNoHighFindings_WithFindings(t *testing.T) {
	checkConfig := commons.CheckConfig{
		Queue: make(chan commons.Check, 1),
	}
	testName := "TestCheckIfGuarddutyNoHighFindings_WithFindings"
	findings := []string{"finding1", "finding2"}

	go guardduty.CheckIfGuarddutyNoHighFindings(checkConfig, testName, findings)
	check := <-checkConfig.Queue

	if len(check.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(check.Results))
	}

	if check.Results[0].Status != "FAIL" {
		t.Errorf("Expected status FAIL, got %s", check.Results[0].Status)
	}

	expectedMessage := "GuardDuty has at least 1 HIGH severity finding, please perform a review of these findings"
	if check.Results[0].Message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, check.Results[0].Message)
	}
}
