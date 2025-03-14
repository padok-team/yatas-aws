package rds

import (
	"strings"

	"github.com/padok-team/yatas/plugins/commons"
)

func checkIfAuditLogsEnabled(checkConfig commons.CheckConfig, dbInstances []InstanceToLogFiles, testName string) {
	var check commons.Check
	check.InitCheck("RDS / Aurora RDS audit logs should be enabled", "Check if RDS / Aurora RDS audit logs are enabled (supports MySQL, Aurora MySQL, MariaDB, PostgreSQL, Aurora PostgreSQL)", testName, []string{"Security", "Good Practice"})
	// MySQL/MariaDB case: check if log files beginning with audit are present
	for _, dbInstance := range dbInstances {
		if *dbInstance.Instance.Engine == "mysql" || *dbInstance.Instance.Engine == "aurora-mysql" || *dbInstance.Instance.Engine == "mariadb" {
			auditLogFiles := []string{}
			for _, logFile := range dbInstance.LogFiles {
				if strings.HasPrefix(*logFile.LogFileName, "audit") {
					auditLogFiles = append(auditLogFiles, *logFile.LogFileName)
				}
			}
			if len(auditLogFiles) == 0 {
				message := "No audit log files found for instance " + *dbInstance.Instance.DBInstanceIdentifier
				result := commons.Result{Status: "FAIL", Message: message, ResourceID: *dbInstance.Instance.DBInstanceArn}
				check.AddResult(result)
			} else {
				message := "Audit log files found for instance " + *dbInstance.Instance.DBInstanceIdentifier
				result := commons.Result{Status: "OK", Message: message, ResourceID: *dbInstance.Instance.DBInstanceArn}
				check.AddResult(result)
			}
		}

		// PostgreSQL case: check if there is a line with "AUDIT: " in the recent log files portion
		if *dbInstance.Instance.Engine == "postgres" || *dbInstance.Instance.Engine == "aurora-postgresql" {
			if strings.Contains(dbInstance.RecentLogFilesPortion, "AUDIT: ") {
				message := "Audit log files found for instance " + *dbInstance.Instance.DBInstanceIdentifier
				result := commons.Result{Status: "OK", Message: message, ResourceID: *dbInstance.Instance.DBInstanceArn}
				check.AddResult(result)
			} else {
				message := "No audit log files found for instance " + *dbInstance.Instance.DBInstanceIdentifier
				result := commons.Result{Status: "FAIL", Message: message, ResourceID: *dbInstance.Instance.DBInstanceArn}
				check.AddResult(result)
			}
		}
	}

	checkConfig.Queue <- check
}
