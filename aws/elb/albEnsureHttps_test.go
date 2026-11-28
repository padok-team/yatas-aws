package elb

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/padok-team/yatas/plugins/commons"
)

const (
	strongPolicy = "ELBSecurityPolicy-TLS13-1-2-2021-06"
	weakPolicy   = "ELBSecurityPolicy-2016-08"
)

func httpsListener(port int32, policy string) types.Listener {
	return types.Listener{
		Port:           aws.Int32(port),
		Protocol:       types.ProtocolEnumHttps,
		SslPolicy:      aws.String(policy),
		DefaultActions: []types.Action{{Type: types.ActionTypeEnumForward}},
	}
}

func httpListener(port int32, redirect bool) types.Listener {
	action := types.Action{Type: types.ActionTypeEnumForward}
	if redirect {
		action = types.Action{
			Type: types.ActionTypeEnumRedirect,
			RedirectConfig: &types.RedirectActionConfig{
				Protocol:   aws.String("HTTPS"),
				Port:       aws.String("443"),
				Host:       aws.String("#{host}"),
				Path:       aws.String("/#{path}"),
				Query:      aws.String("#{query}"),
				StatusCode: types.RedirectActionStatusCodeEnumHttp301,
			},
		}
	}
	return types.Listener{Port: aws.Int32(port), Protocol: types.ProtocolEnumHttp, DefaultActions: []types.Action{action}}
}

func runElbCheck(fn func(commons.CheckConfig, []LoadBalancerAttributes, string), listeners ...types.Listener) commons.Check {
	cc := commons.CheckConfig{Queue: make(chan commons.Check, 1)}
	fn(cc, []LoadBalancerAttributes{{
		LoadBalancerName: "test-lb",
		LoadBalancerArn:  "arn:aws:elasticloadbalancing:eu-west-3:123456789012:loadbalancer/test-lb",
		Listeners:        listeners,
	}}, "AWS_ELB_00X")
	return <-cc.Queue
}

func TestCheckAlbEnsureHttps(t *testing.T) {
	tests := []struct {
		name      string
		listeners []types.Listener
		want      string
	}{
		{"HTTPS TLS1.2 + HTTP redirect", []types.Listener{httpsListener(443, strongPolicy), httpListener(80, true)}, "OK"},
		{"HTTPS TLS1.2 only", []types.Listener{httpsListener(443, strongPolicy)}, "OK"},
		{"HTTPS on non-standard port", []types.Listener{httpsListener(8443, "ELBSecurityPolicy-FS-1-2-Res-2020-10")}, "OK"},
		{"HTTPS TLS1.3 only policy", []types.Listener{httpsListener(443, "ELBSecurityPolicy-TLS13-1-3-2021-06")}, "OK"},
		{"HTTPS with default policy allowing TLS1.0", []types.Listener{httpsListener(443, weakPolicy)}, "FAIL"},
		{"HTTPS with TLS13-1-1 policy", []types.Listener{httpsListener(443, "ELBSecurityPolicy-TLS13-1-1-2021-06")}, "FAIL"},
		{"no HTTPS listener", []types.Listener{httpListener(80, false)}, "FAIL"},
		{"plain HTTP on port 443", []types.Listener{httpListener(443, false)}, "FAIL"},
		{"HTTP not redirected", []types.Listener{httpsListener(443, strongPolicy), httpListener(80, false)}, "FAIL"},
		{"HTTP on 8080 not redirected", []types.Listener{httpsListener(443, strongPolicy), httpListener(80, true), httpListener(8080, false)}, "FAIL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runElbCheck(CheckAlbEnsureHttps, tt.listeners...); got.Results[0].Status != tt.want {
				t.Errorf("CheckAlbEnsureHttps() = %v, want %s", got.Results[0], tt.want)
			}
		})
	}
}

func TestCheckNlbTLS12Minimum(t *testing.T) {
	tls := func(policy string) types.Listener {
		return types.Listener{Port: aws.Int32(443), Protocol: types.ProtocolEnumTls, SslPolicy: aws.String(policy)}
	}
	tcp := types.Listener{Port: aws.Int32(443), Protocol: types.ProtocolEnumTcp}

	if got := runElbCheck(CheckNlbTLS12Minimum, tls(strongPolicy)); got.Results[0].Status != "OK" {
		t.Errorf("TLS1.2 policy: got %v", got.Results[0])
	}
	if got := runElbCheck(CheckNlbTLS12Minimum, tls(strongPolicy), tls(weakPolicy)); got.Results[0].Status != "FAIL" {
		t.Errorf("weak policy: got %v", got.Results[0])
	}
	if got := runElbCheck(CheckNlbTLS12Minimum, tcp); len(got.Results) != 0 {
		t.Errorf("TCP passthrough should be skipped, got %v", got.Results)
	}
}
