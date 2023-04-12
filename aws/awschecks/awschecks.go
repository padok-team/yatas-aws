package awschecks

import (
	"fmt"
	"reflect"

	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/padok-team/yatas-aws/logger"
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

func CheckResources(checkConfig commons.CheckConfig, resources []interface{}, checkDefinitions []CheckDefinition) {
	for _, checkDefinition := range checkDefinitions {
		// if !checkConfig.ConfigYatas.CheckExclude(checkDefinition.Title) && checkConfig.ConfigYatas.CheckInclude(checkDefinition.Title) {
		// 	checkConfig.Wg.Add(1)
		// 	logger.Logger.Info("Running check: " + checkDefinition.Title)
		// }
		check := createCheck(checkDefinition)
		for _, resource := range resources {
			result := checkResource(resource, checkDefinition.ConditionFn, checkDefinition.SuccessMessage, checkDefinition.FailureMessage)
			check.AddResult(result)
		}
		checkConfig.Queue <- check
	}
}

func AddChecks(checkConfig *commons.CheckConfig, resources ...[]CheckDefinition) {
	//Print the resources to check
	totalCount := 0
	for _, resourceSlice := range resources {
		totalCount += len(resourceSlice)
	}
	logger.Logger.Info(fmt.Sprintf("Adding %d checks", totalCount))
	checkConfig.Wg.Add(totalCount)
}

func createCheck(checkDefinition CheckDefinition) commons.Check {
	var check commons.Check
	logger.Logger.Info("Creating check: " + checkDefinition.Title)
	check.InitCheck(checkDefinition.Description, checkDefinition.Description, checkDefinition.Title, checkDefinition.Tags)
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
