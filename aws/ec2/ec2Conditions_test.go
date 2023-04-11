package ec2

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Generate unit tests for EC2

func TestEc2MonitoringEnabledCondition(t *testing.T) {
	type args struct {
		resource interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "TestEc2MonitoringEnabledCondition",
			args: args{
				resource: &types.Instance{
					Monitoring: &types.Monitoring{
						State: types.MonitoringStateEnabled,
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ec2MonitoringEnabledCondition(tt.args.resource); got != tt.want {
				t.Errorf("Ec2MonitoringEnabledCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEc2PublicIPCondition(t *testing.T) {
	type args struct {
		resource interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "TestEc2PublicIPCondition",
			args: args{
				resource: &types.Instance{
					PublicIpAddress: nil,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ec2PublicIPCondition(tt.args.resource); got != tt.want {
				t.Errorf("Ec2PublicIPCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEc2RunningInVPCCondition(t *testing.T) {
	type args struct {
		resource interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "TestEc2RunningInVPCCondition",
			args: args{
				resource: &types.Instance{
					VpcId: nil,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ec2RunningInVPCCondition(tt.args.resource); got != tt.want {
				t.Errorf("Ec2RunningInVPCCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}
