package awschecks

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/padok-team/yatas/plugins/commons"
)

type CheckFunc func(interface{}) commons.Result

type CheckDefinition struct {
	Title          string
	Description    string
	Tags           []string
	ConditionFn    func(interface{}) bool
	SuccessMessage string
	FailureMessage string
}

func CheckResources(checkConfig commons.CheckConfig, resources []interface{}, testName string, checkDefinitions []CheckDefinition) {
	for _, checkDefinition := range checkDefinitions {
		check := createCheck(testName, checkDefinition)
		for _, resource := range resources {
			result := checkResource(resource, checkDefinition.ConditionFn, checkDefinition.SuccessMessage, checkDefinition.FailureMessage)
			check.AddResult(result)
		}
		checkConfig.Queue <- check
	}
}

func createCheck(testName string, checkDefinition CheckDefinition) commons.Check {
	var check commons.Check
	check.InitCheck(checkDefinition.Title, checkDefinition.Description, testName, checkDefinition.Tags)
	return check
}

func checkResource(resource interface{}, conditionFn func(interface{}) bool, successMessage, failureMessage string) commons.Result {
	if conditionFn(resource) {
		message := successMessage + " - Resource " + getResourceID(resource)
		return commons.Result{Status: "OK", Message: message, ResourceID: getResourceID(resource)}
	} else {
		message := failureMessage + " - Resource " + getResourceID(resource)
		return commons.Result{Status: "FAIL", Message: message, ResourceID: getResourceID(resource)}
	}
}

func getResourceID(resource interface{}) string {
	resourceType := reflect.TypeOf(resource)
	resourceValue := reflect.ValueOf(resource)

	switch resourceType {
	// Add cases for each supported AWS resource type
	case reflect.TypeOf((*ec2Types.Instance)(nil)).Elem():
		return *resourceValue.FieldByName("InstanceId").Interface().(*string)
	// Add more cases as needed
	default:
		fmt.Printf("Unsupported resource type: %s\n", resourceType)
		return ""
	}
}

func Ec2MonitoringEnabledCondition(resource interface{}) bool {
	instance, ok := resource.(*types.Instance)
	if !ok {
		return false
	}
	return instance.Monitoring.State == types.MonitoringStateEnabled
}

func Ec2PublicIPCondition(resource interface{}) bool {
	instance, ok := resource.(*types.Instance)
	if !ok {
		return false
	}
	return instance.PublicIpAddress == nil
}
