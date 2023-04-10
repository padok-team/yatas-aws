package lambda

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/dlclark/regexp2"
	"github.com/padok-team/yatas/plugins/commons"
)

var secrets_patterns = []*regexp2.Regexp{
	// General
	regexp2.MustCompile("^-----BEGIN (RSA|EC|DSA|GPP) PRIVATE KEY-----$", regexp2.RE2),
	// AWS
	regexp2.MustCompile("(?<![A-Za-z0-9/+=])[A-Za-z0-9/+=]{40}(?![A-Za-z0-9/+=])", regexp2.RE2),           // AWS secret access key
	regexp2.MustCompile("(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}", regexp2.RE2), // AWS access key ID
	regexp2.MustCompile("(\"|')?(AWS|aws|Aws)?_?(SECRET|secret|Secret)?_?(ACCESS|access|Access)?_?(KEY|key|Key)(\"|')?\\s*(:|=>|=)\\s*(\"|')?[A-Za-z0-9/\\+=]{40}(\"|')?", regexp2.RE2),
	regexp2.MustCompile("(\"|')?(AWS|aws|Aws)?_?(ACCOUNT|account|Account)_?(ID|id|Id)?(\"|')?\\s*(:|=>|=)\\s*(\"|')?[0-9]{4}\\-?[0-9]{4}\\-?[0-9]{4}(\"|')?", regexp2.RE2),
}

func string_has_secrets(value string) bool {
	isSecret := false
	for _, pattern := range secrets_patterns {
		if res, _ := pattern.MatchString(value); res {
			isSecret = true
			break
		}
	}
	return isSecret
}

func CheckIfLambdaNoSecrets(checkConfig commons.CheckConfig, lambdas []types.FunctionConfiguration, testName string) {
	var check commons.Check
	check.InitCheck("Lambdas has no hard-coded secrets in environment", "Check if all Lambdas has no secrets as environment variable", testName, []string{"Security", "Good Practice"})

	for _, lambda := range lambdas {
		envSecrets := []string{}
		if lambda.Environment.Error == nil {
			for key, value := range lambda.Environment.Variables {
				if string_has_secrets(value) {
					envSecrets = append(envSecrets, key)
				}
			}
		}

		if len(envSecrets) > 0 {
			Message := "Lambda " + *lambda.FunctionName + " has secrets in environment: " + fmt.Sprint(envSecrets)
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		} else {
			Message := "Lambda " + *lambda.FunctionName + " has no secrets in environment"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: *lambda.FunctionArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
