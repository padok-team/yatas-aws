package apigateway

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/stangirard/yatas/plugins/commons"
)

func TestCheckIfTracingEnabled(t *testing.T) {
	type args struct {
		checkConfig commons.CheckConfig
		stages      map[string][]types.Stage
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Check if all stages are tracing enabled",
			args: args{
				checkConfig: commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}},
				stages: map[string][]types.Stage{
					"test-api": {
						{
							TracingEnabled: true,
							StageName:      aws.String("test"),
						},
					},
				},
				testName: "CheckIfTracingEnabled",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfTracingEnabled(tt.args.checkConfig, tt.args.stages, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if len(check.Results) != len(tt.args.stages) {
						t.Errorf("CheckIfTracingEnabled() = %v, want %v", len(check.Results), len(tt.args.stages))
					}
					if check.Status != "OK" {
						t.Errorf("CheckIfTracingEnabled() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
