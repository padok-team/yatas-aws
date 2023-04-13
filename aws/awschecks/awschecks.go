package awschecks

import (
	"fmt"

	"github.com/padok-team/yatas-aws/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

type CheckFunc func(interface{}) commons.Result

type Resource interface {
	GetID() string
}

type CheckDefinition struct {
	Title          string
	Description    string
	Tags           []string
	ConditionFn    func(Resource) bool
	SuccessMessage string
	FailureMessage string
}

func CheckResources(checkConfig commons.CheckConfig, resources []Resource, checkDefinitions []CheckDefinition) {
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

func checkResource(resource Resource, conditionFn func(Resource) bool, successMessage, failureMessage string) commons.Result {
	if conditionFn(resource) {
		message := successMessage + " - Resource " + resource.GetID()
		return commons.Result{Status: "OK", Message: message, ResourceID: resource.GetID()}
	} else {
		message := failureMessage + " - Resource " + resource.GetID()
		return commons.Result{Status: "FAIL", Message: message, ResourceID: resource.GetID()}
	}
}
