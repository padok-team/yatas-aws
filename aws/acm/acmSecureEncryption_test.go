package acm

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfACMUsesSecureEncryption(t *testing.T) {
	tests := []struct {
		name         string
		certificates []types.CertificateDetail
		expected     string
	}{
		{
			name: "Certificate with secure encryption",
			certificates: []types.CertificateDetail{
				{
					CertificateArn: aws.String("arn:aws:acm:region:account:certificate/secure-cert"),
					KeyAlgorithm:   types.KeyAlgorithmRsa2048,
				},
			},
			expected: "OK",
		},
		{
			name: "Certificate with insecure encryption",
			certificates: []types.CertificateDetail{
				{
					CertificateArn: aws.String("arn:aws:acm:region:account:certificate/insecure-cert"),
					KeyAlgorithm:   types.KeyAlgorithmRsa1024,
				},
			},
			expected: "FAIL",
		},
		{
			name: "Multiple certificates",
			certificates: []types.CertificateDetail{
				{
					CertificateArn: aws.String("arn:aws:acm:region:account:certificate/secure-cert"),
					KeyAlgorithm:   types.KeyAlgorithmRsa2048,
				},
				{
					CertificateArn: aws.String("arn:aws:acm:region:account:certificate/secure-cert"),
					KeyAlgorithm:   types.KeyAlgorithmRsa2048,
				},
				{
					CertificateArn: aws.String("arn:aws:acm:region:account:certificate/insecure-cert"),
					KeyAlgorithm:   types.KeyAlgorithmRsa1024,
				},
			},
			expected: "FAIL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkConfig := commons.CheckConfig{
				Queue: make(chan commons.Check, 1),
			}
			go CheckIfACMUsesSecureEncryption(checkConfig, tt.certificates, tt.name)
			check := <-checkConfig.Queue

			if len(check.Results) != len(tt.certificates) {
				t.Errorf("Expected %d result, got %d", len(tt.certificates), len(check.Results))
			}

			result := check.Results[0].Status

			for _, status := range check.Results {
				if status.Status != "OK" {
					result = status.Status
				}
			}

			if result != tt.expected {
				t.Errorf("Expected status %s, got %s", tt.expected, result)
			}
		})
	}
}
