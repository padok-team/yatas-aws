package main

import (
	"context"
	"encoding/gob"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/padok-team/yatas-aws/aws/acm"
	"github.com/padok-team/yatas-aws/aws/apigateway"
	"github.com/padok-team/yatas-aws/aws/autoscaling"
	"github.com/padok-team/yatas-aws/aws/cloudfront"
	"github.com/padok-team/yatas-aws/aws/cloudtrail"
	"github.com/padok-team/yatas-aws/aws/cognito"
	"github.com/padok-team/yatas-aws/aws/configservice"
	"github.com/padok-team/yatas-aws/aws/dynamodb"
	"github.com/padok-team/yatas-aws/aws/ec2"
	"github.com/padok-team/yatas-aws/aws/ecr"
	"github.com/padok-team/yatas-aws/aws/eks"
	"github.com/padok-team/yatas-aws/aws/elb"
	"github.com/padok-team/yatas-aws/aws/guardduty"
	"github.com/padok-team/yatas-aws/aws/iam"
	"github.com/padok-team/yatas-aws/aws/lambda"
	"github.com/padok-team/yatas-aws/aws/rds"
	"github.com/padok-team/yatas-aws/aws/s3"
	"github.com/padok-team/yatas-aws/aws/ssm"
	"github.com/padok-team/yatas-aws/aws/volumes"
	"github.com/padok-team/yatas-aws/aws/vpc"
	"github.com/padok-team/yatas-aws/internal"
	"github.com/padok-team/yatas-aws/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

// Create a new session that the SDK will use to load
// credentials from. With either SSO or credentials
func initAuth(a internal.AWS_Account) aws.Config {
	logger.Logger.Debug("Init auth")
	s := initSession(a)
	return s
}

// Create a new session that the SDK will use to load
// credentials from credentials
func createSessionWithCredentials(c internal.AWS_Account) aws.Config {
	var s aws.Config
	var err error

	if c.Profile == "" {
		s, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(c.Region),
			config.WithRetryer(func() aws.Retryer {
				return retry.AddWithMaxAttempts(retry.NewStandard(), 10)
			}),
			config.WithRetryMode(aws.RetryMode(aws.RetryModeAdaptive)),
		)
	} else {
		s, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(c.Region),
			config.WithSharedConfigProfile(c.Profile),
			config.WithRetryer(func() aws.Retryer {
				return retry.AddWithMaxAttempts(retry.NewStandard(), 10)
			}),
			config.WithRetryMode(aws.RetryMode(aws.RetryModeAdaptive)),
		)
	}

	if err != nil {
		logger.Logger.Error(err.Error())
	}

	return s
}

// Create a new session that the SDK will use to load
// credentials from the shared credentials file.
// Usefull for SSO
func createSessionWithSSO(c internal.AWS_Account) aws.Config {

	if c.Profile == "" {
		s, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(c.Region),
			config.WithRetryer(func() aws.Retryer {
				return retry.AddWithMaxAttempts(retry.NewStandard(), 10)
			}),
			config.WithRetryMode(aws.RetryMode(aws.RetryModeAdaptive)),
		)
		if err != nil {
			logger.Logger.Error(err.Error())
		}
		return s
	} else {
		s, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(c.Region),
			config.WithSharedConfigProfile(c.Profile),
			config.WithRetryer(func() aws.Retryer {
				return retry.AddWithMaxAttempts(retry.NewStandard(), 10)
			}),
			config.WithRetryMode(aws.RetryMode(aws.RetryModeAdaptive)),
		)
		if err != nil {
			logger.Logger.Error(err.Error())
		}
		return s

	}

}

// Create a new session that the SDK will use to load
// credentials from. With either SSO or credentials
func initSession(c internal.AWS_Account) aws.Config {

	if c.SSO {
		logger.Logger.Debug("Init SSO")
		return createSessionWithSSO(c)
	} else {
		logger.Logger.Debug("Init Credentials")
		return createSessionWithCredentials(c)
	}
}

// Public Functin used to run the AWS tests
func Run(c *commons.Config, accounts []internal.AWS_Account) ([]commons.Tests, error) {
	logger.Logger.Debug("Run tests")
	var wg sync.WaitGroup
	var queue = make(chan commons.Tests, 10)
	var checks []commons.Tests
	wg.Add(len(accounts))
	for _, account := range accounts {

		go runTestsForAccount(account, c, queue)
	}
	go func() {
		for t := range queue {
			checks = append(checks, t)

			wg.Done()
		}
	}()
	wg.Wait()

	return checks, nil
}

// For each account we run the tests. We use a queue to store the results and a waitgroup to wait for all the tests to be done. This allows to run all tests asynchronously.
func runTestsForAccount(account internal.AWS_Account, c *commons.Config, queue chan commons.Tests) {
	s := initAuth(account)
	checks := initTest(s, c, account)
	queue <- checks
}

// Main function that launched all the test for a given account. If a new category is added, it needs to be added here.
func initTest(s aws.Config, c *commons.Config, a internal.AWS_Account) commons.Tests {

	var checks commons.Tests
	checks.Account = a.Name
	var wg sync.WaitGroup
	queue := make(chan []commons.Check, 100)
	go commons.CheckMacroTest(&wg, c, cognito.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, acm.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, s3.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, volumes.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, rds.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, vpc.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, cloudtrail.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, ecr.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, lambda.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, dynamodb.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, ec2.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, cloudfront.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, apigateway.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, autoscaling.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, elb.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, guardduty.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, iam.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, eks.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, configservice.RunChecks)(&wg, s, c, queue)
	go commons.CheckMacroTest(&wg, c, ssm.RunChecks)(&wg, s, c, queue)

	go func() {
		for t := range queue {

			checks.Checks = append(checks.Checks, t...)

			wg.Done()

		}
	}()
	wg.Wait()

	return checks
}

type YatasPlugin struct {
	logger hclog.Logger
}

func UnmarshalAWS(g *YatasPlugin, c *commons.Config) ([]internal.AWS_Account, error) {
	var accounts []internal.AWS_Account

	for _, r := range c.PluginConfig {
		var tmpAccounts []internal.AWS_Account
		awsFound := false
		for key, value := range r {

			switch key {
			case "pluginName":
				if value == "aws" {
					awsFound = true

				}
			case "accounts":

				for _, v := range value.([]interface{}) {
					var account internal.AWS_Account
					logger.Logger.Debug("🔎")
					logger.Logger.Debug("%v", v)
					for keyaccounts, valueaccounts := range v.(map[string]interface{}) {
						switch keyaccounts {
						case "name":
							account.Name = valueaccounts.(string)
						case "profile":
							account.Profile = valueaccounts.(string)
						case "region":
							account.Region = valueaccounts.(string)
						case "sso":
							account.SSO = valueaccounts.(bool)
						}
					}
					tmpAccounts = append(tmpAccounts, account)

				}

			}
		}
		if awsFound {
			accounts = tmpAccounts
		}

	}
	logger.Logger.Debug("✅")
	logger.Logger.Debug("%v", accounts)
	logger.Logger.Debug("Length of accounts: %d", len(accounts))
	if len(accounts) == 0 {
		logger.Logger.Error("No AWS accounts found in config file")
	}
	return accounts, nil
}

func (g *YatasPlugin) Run(c *commons.Config) []commons.Tests {
	logger.Logger = g.logger
	logger.Logger.Debug("message from YatasPlugin.Run")
	var err error
	var accounts []internal.AWS_Account
	accounts, err = UnmarshalAWS(g, c)
	if err != nil {
		logger.Logger.Error("Error unmarshaling AWS accounts", "error", err)
		return nil
	}
	var checksAll []commons.Tests

	checks, err := runPlugins(c, "aws", accounts)
	if err != nil {
		logger.Logger.Error("Error running plugins", "error", err)
	}
	checksAll = append(checksAll, checks...)
	return checksAll
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Debug,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	yatasPlugin := &YatasPlugin{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"aws": &commons.YatasPlugin{Impl: yatasPlugin},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}

// Run the plugins that are enabled in the config with a switch based on the name of the plugin
func runPlugins(c *commons.Config, _ string, accounts []internal.AWS_Account) ([]commons.Tests, error) {
	var checksAll []commons.Tests

	checksAll, err := Run(c, accounts)
	if err != nil {
		return nil, err
	}

	return checksAll, nil
}
