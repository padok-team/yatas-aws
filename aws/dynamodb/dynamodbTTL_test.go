package dynamodb

import (
	"sync"
	"testing"

	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfTTLConfiguredAndValid(t *testing.T) {
	type args struct {
		checkConfig commons.CheckConfig
		tables      []TableTTL
		testName    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TTL properly configured - OK case",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				tables: []TableTTL{
					{
						TableName:    "test-table",
						TTLEnabled:   true,
						TTLAttribute: "expiryTime",
					},
				},
				testName: "AWS_DYN_004",
			},
			want: "OK",
		},
		{
			name: "TTL not enabled - FAIL case",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				tables: []TableTTL{
					{
						TableName:    "test-table-no-ttl",
						TTLEnabled:   false,
						TTLAttribute: "",
					},
				},
				testName: "AWS_DYN_004",
			},
			want: "FAIL",
		},
		{
			name: "TTL enabled but no attribute - FAIL case",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				tables: []TableTTL{
					{
						TableName:    "test-table-no-attr",
						TTLEnabled:   true,
						TTLAttribute: "",
					},
				},
				testName: "AWS_DYN_004",
			},
			want: "FAIL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfTTLConfiguredAndValid(tt.args.checkConfig, tt.args.tables, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != tt.want {
						t.Errorf("CheckIfTTLConfiguredAndValid() = %v, want %v", check.Status, tt.want)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfTTLConfiguredAndValidFail(t *testing.T) {
	type args struct {
		checkConfig commons.CheckConfig
		tables      []TableTTL
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Multiple tables with issues",
			args: args{
				checkConfig: commons.CheckConfig{
					Queue: make(chan commons.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				tables: []TableTTL{
					{
						TableName:    "test-table-no-ttl",
						TTLEnabled:   false,
						TTLAttribute: "",
					},
					{
						TableName:    "test-table-no-attr",
						TTLEnabled:   true,
						TTLAttribute: "",
					},
				},
				testName: "AWS_DYN_004",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfTTLConfiguredAndValid(tt.args.checkConfig, tt.args.tables, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfTTLConfiguredAndValid() = %v, want FAIL", check.Status)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
